package usecase

import (
	"context"
	"fmt"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"

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
		Name:      request.Name,
		CompanyID: request.CompanyID,
	}

	if request.ParentID != nil {
		position.ParentID = request.ParentID
	}

	if err := c.PositionRepository.Create(tx, position); err != nil {
		c.Log.WithError(err).Error("Failed to create position")
		return nil, fiber.ErrInternalServerError
	}

	fmt.Println("Tesssss ====")

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return &model.PositionResponse{
		Name:      position.Name,
		CompanyID: position.CompanyID,
		ParentID:  position.ParentID,
	}, nil

}

/*
List Position Usecase
*/
func (c *PositionUseCase) Search(
	ctx context.Context,
	request *model.SeachPositionRequest,
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

func (c *PositionUseCase) buildTree(
	pos entity.Position,
	posMap map[string]entity.Position,
) model.PositionResponse {

	response := model.PositionResponse{
		ID:        pos.ID,
		CompanyID: pos.CompanyID,
		Name:      pos.Name,
	}

	if pos.ParentID != nil {
		parent := posMap[*pos.ParentID]
		parentResponse := c.buildTree(parent, posMap)
		response.Parent = &parentResponse
	}

	return response
}
