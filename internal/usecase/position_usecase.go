package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PositionUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	PositionRepository *repository.PositionRepository
}

func NewPositionUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	positionRepository *repository.PositionRepository,
) *PositionUseCase {
	return &PositionUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		PositionRepository: positionRepository,
	}
}

/*
List Position Usecase
*/
func (c *PositionUseCase) Create(ctx context.Context, request *model.CreatePositionRequest) (*model.PositionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	position := &entity.Position{
		Name:       request.Name,
		CompanyID:  request.CompanyID,
		IsApprover: request.IsApprover,
	}

	if request.ParentID != nil {
		position.ParentID = request.ParentID
	}

	if err := c.PositionRepository.Create(tx, position); err != nil {
		c.Log.WithError(err).Error("Failed to create position")
		return nil, fiber.ErrInternalServerError
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return &model.PositionResponse{
		ID:         position.ID,
		Name:       position.Name,
		CompanyID:  position.CompanyID,
		ParentID:   position.ParentID,
		IsApprover: position.IsApprover,
	}, nil

}

/*
List Position Usecase
*/
func (c *PositionUseCase) Search(
	ctx context.Context,
	request *model.SearchPositionRequest,
) ([]model.PositionResponse, int64, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var result []model.PositionResponse

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	positions, total, err := c.PositionRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting sanctions")
		return nil, 0, fiber.ErrInternalServerError
	}

	posMap := make(map[string]entity.Position)
	for _, p := range positions {
		posMap[p.ID] = p
	}

	// 🔥 BUILD UNTUK SEMUA POSITION
	for _, p := range positions {
		result = append(result, c.buildTree(p, posMap))
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return result, total, nil
}

func (c *PositionUseCase) Detail(ctx context.Context, id string) (*model.PositionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	position, err := c.PositionRepository.FindByID(tx, id, "Parent")
	if err != nil {
		c.Log.WithError(err).Error("Position not found")
		return nil, fiber.ErrNotFound
	}

	response := model.PositionToResponse(position)
	if position.Parent != nil {
		response.Parent = model.PositionToResponse(position.Parent)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
}

func (c *PositionUseCase) Update(ctx context.Context, id string, request *model.UpdatePositionRequest) (*model.PositionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	position, err := c.PositionRepository.FindByID(tx, id, "Parent")
	if err != nil {
		c.Log.WithError(err).Error("Position not found")
		return nil, fiber.ErrNotFound
	}

	if request.Name != nil {
		name := strings.TrimSpace(*request.Name)
		if name == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		position.Name = name
	}

	if request.ParentID != nil {
		if *request.ParentID == id {
			return nil, fiber.NewError(fiber.StatusBadRequest, "parent_id cannot be the same as position id")
		}
		if strings.TrimSpace(*request.ParentID) == "" {
			position.ParentID = nil
		} else {
			position.ParentID = request.ParentID
		}
	}

	if request.IsApprover != nil {
		position.IsApprover = *request.IsApprover
	}

	position.UpdatedAt = time.Now().UnixMilli()

	if err := c.PositionRepository.Update(tx, position); err != nil {
		c.Log.WithError(err).Error("Failed to update position")
		return nil, fiber.ErrInternalServerError
	}

	reloaded, err := c.PositionRepository.FindByID(tx, id, "Parent")
	if err != nil {
		c.Log.WithError(err).Error("Failed to reload position")
		return nil, fiber.ErrInternalServerError
	}

	response := model.PositionToResponse(reloaded)
	if reloaded.Parent != nil {
		response.Parent = model.PositionToResponse(reloaded.Parent)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
}

func (c *PositionUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	position, err := c.PositionRepository.FindByID(tx, id)
	if err != nil {
		c.Log.WithError(err).Error("Position not found")
		return fiber.ErrNotFound
	}

	if err := c.PositionRepository.Delete(tx, position); err != nil {
		c.Log.WithError(err).Error("Failed to delete position")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *PositionUseCase) buildTree(
	pos entity.Position,
	posMap map[string]entity.Position,
) model.PositionResponse {

	response := model.PositionResponse{
		ID:         pos.ID,
		CompanyID:  pos.CompanyID,
		Name:       pos.Name,
		ParentID:   pos.ParentID,
		IsApprover: pos.IsApprover,
	}

	if pos.ParentID != nil {
		parent := posMap[*pos.ParentID]
		parentResponse := c.buildTree(parent, posMap)
		response.Parent = &parentResponse
	}

	return response
}
