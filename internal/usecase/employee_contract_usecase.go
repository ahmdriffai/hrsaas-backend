package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeContractUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	Repo               *repository.EmployeeContractRepository
	TimeOffTypeRepo    *repository.TimeOffTypeRepository
	TimeOffBalanceRepo *repository.TimeOffBalanceRepository
}

func NewEmployeeContractUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	repo *repository.EmployeeContractRepository,
	timeOffTypeRepo *repository.TimeOffTypeRepository,
	timeOffBalanceRepo *repository.TimeOffBalanceRepository,
) *EmployeeContractUseCase {
	return &EmployeeContractUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		Repo:               repo,
		TimeOffTypeRepo:    timeOffTypeRepo,
		TimeOffBalanceRepo: timeOffBalanceRepo,
	}
}

// TODO: Add business validation (only one active contract per employee).
func (c *EmployeeContractUseCase) Create(ctx context.Context, request *model.CreateEmployeeContractRequest) (*model.EmployeeContractResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	startDate, err := lib.ParseDateToUnixMilli(request.StartDate)
	if err != nil {
		return nil, fiber.ErrBadRequest
	}

	var endDate *int64
	if request.EndDate != nil && *request.EndDate != "" {
		parsed, err := lib.ParseDateToUnixMilli(*request.EndDate)
		if err != nil {
			return nil, fiber.ErrBadRequest
		}
		endDate = &parsed
	}

	item := &entity.EmployeeContract{
		EmployeeID:   request.EmployeeID,
		ContractType: request.ContractType,
		StartDate:    startDate,
		EndDate:      endDate,
		DivisionID:   request.DivisionID,
		PositionID:   request.PositionID,
		Salary:       request.Salary,
	}

	if err := c.Repo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create employee contract")
		return nil, fiber.ErrInternalServerError
	}

	// Auto-create time off balances for quota-based types.
	periodYear := time.UnixMilli(startDate).UTC().Year()
	types, err := c.TimeOffTypeRepo.ListQuotaBased(tx)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off types")
		return nil, fiber.ErrInternalServerError
	}
	for _, t := range types {
		if _, err := c.TimeOffBalanceRepo.FindByEmployeeTypeYear(tx, request.EmployeeID, t.ID, periodYear); err == nil {
			continue
		} else if err != nil && err != gorm.ErrRecordNotFound {
			c.Log.WithError(err).Error("Failed to check time off balance")
			return nil, fiber.ErrInternalServerError
		}

		balance := &entity.Time_Off_Balance{
			EmployeeId:    request.EmployeeID,
			TimeOffTypeId: t.ID,
			PeriodYear:    periodYear,
			EntitledDays:  t.DefaultQuotaDays,
			UsedDays:      0,
			RemainingDays: t.DefaultQuotaDays,
		}
		if err := c.TimeOffBalanceRepo.Create(tx, balance); err != nil {
			c.Log.WithError(err).Error("Failed to create time off balance")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	var endDateStr *string
	if endDate != nil {
		str := time.UnixMilli(*endDate).Format("2006-01-02")
		endDateStr = &str
	}

	return &model.EmployeeContractResponse{
		ID:           item.ID,
		EmployeeID:   item.EmployeeID,
		ContractType: item.ContractType,
		StartDate:    time.UnixMilli(item.StartDate).Format("2006-01-02"),
		EndDate:      endDateStr,
		DivisionID:   item.DivisionID,
		PositionID:   item.PositionID,
		Salary:       item.Salary,
	}, nil
}

// TODO: Add authorization filtering for non-admin users.
func (c *EmployeeContractUseCase) List(ctx context.Context, request *model.SearchEmployeeContractRequest) ([]model.EmployeeContractResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.Repo.List(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list employee contracts")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.EmployeeContractResponse, len(items))
	for i, item := range items {
		var endDateStr *string
		if item.EndDate != nil {
			str := time.UnixMilli(*item.EndDate).Format("2006-01-02")
			endDateStr = &str
		}

		responses[i] = model.EmployeeContractResponse{
			ID:           item.ID,
			EmployeeID:   item.EmployeeID,
			ContractType: item.ContractType,
			StartDate:    time.UnixMilli(item.StartDate).Format("2006-01-02"),
			EndDate:      endDateStr,
			DivisionID:   item.DivisionID,
			PositionID:   item.PositionID,
			Salary:       item.Salary,
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}
