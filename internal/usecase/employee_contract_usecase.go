package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strings"

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
		EmployeeID:       request.EmployeeID,
		ContractType:     request.ContractType,
		StartDate:        startDate,
		EndDate:          endDate,
		DivisionID:       request.DivisionID,
		PositionID:       request.PositionID,
		EmploymentStatus: request.EmploymentStatus,
		Salary:           request.Salary,
		IsActive:         true,
	}

	if err := c.Repo.DeactivateActiveByEmployee(tx, request.EmployeeID); err != nil {
		c.Log.WithError(err).Error("Failed to deactivate previous active employee contracts")
		return nil, fiber.ErrInternalServerError
	}

	if err := c.Repo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create employee contract")
		return nil, fiber.ErrInternalServerError
	}

	// Auto-create time off balances for quota-based types.
	// periodYear := time.UnixMilli(startDate).UTC().Year()
	// types, err := c.TimeOffTypeRepo.ListQuotaBased(tx)
	// if err != nil {
	// 	c.Log.WithError(err).Error("Failed to list time off types")
	// 	return nil, fiber.ErrInternalServerError
	// }
	// for _, t := range types {
	// 	if _, err := c.TimeOffBalanceRepo.FindByEmployeeTypeYear(tx, request.EmployeeID, t.ID, periodYear); err == nil {
	// 		continue
	// 	} else if err != gorm.ErrRecordNotFound {
	// 		c.Log.WithError(err).Error("Failed to check time off balance")
	// 		return nil, fiber.ErrInternalServerError
	// 	}

	// 	balance := &entity.Time_Off_Balance{
	// 		EmployeeId:    request.EmployeeID,
	// 		TimeOffTypeId: t.ID,
	// 		PeriodYear:    periodYear,
	// 		EntitledDays:  t.DefaultQuotaDays,
	// 		UsedDays:      0,
	// 		RemainingDays: t.DefaultQuotaDays,
	// 	}
	// 	if err := c.TimeOffBalanceRepo.Create(tx, balance); err != nil {
	// 		c.Log.WithError(err).Error("Failed to create time off balance")
	// 		return nil, fiber.ErrInternalServerError
	// 	}
	// }

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}
	return &model.EmployeeContractResponse{
		ID:               item.ID,
		EmployeeID:       item.EmployeeID,
		ContractType:     item.ContractType,
		StartDate:        item.StartDate,
		EndDate:          item.EndDate,
		DivisionID:       item.DivisionID,
		PositionID:       item.PositionID,
		Salary:           item.Salary,
		IsActive:         item.IsActive,
		EmploymentStatus: item.EmploymentStatus,
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

	items, total, err := c.Repo.List(tx, request, true)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list employee contracts")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.EmployeeContractResponse, len(items))
	for i, item := range items {
		responses[i] = *model.EmployeeContractToResponse(&item)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

func (c *EmployeeContractUseCase) Detail(ctx context.Context, id string) (*model.EmployeeContractResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.Repo.FindByID(tx, id, true)
	if err != nil {
		c.Log.WithError(err).Error("Employee contract not found")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeContractToResponse(item), nil
}

func (c *EmployeeContractUseCase) Update(ctx context.Context, id string, request *model.UpdateEmployeeContractRequest) (*model.EmployeeContractResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	item, err := c.Repo.FindByID(tx, id, true)
	if err != nil {
		c.Log.WithError(err).Error("Employee contract not found")
		return nil, fiber.ErrNotFound
	}

	if request.ContractType != nil {
		contractType := strings.TrimSpace(*request.ContractType)
		if contractType == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "contract_type cannot be empty")
		}
		item.ContractType = contractType
	}
	if request.StartDate != nil {
		startDate, err := lib.ParseDateToUnixMilli(strings.TrimSpace(*request.StartDate))
		if err != nil {
			return nil, fiber.ErrBadRequest
		}
		item.StartDate = startDate
	}
	if request.EndDate != nil {
		if strings.TrimSpace(*request.EndDate) == "" {
			item.EndDate = nil
		} else {
			endDate, err := lib.ParseDateToUnixMilli(strings.TrimSpace(*request.EndDate))
			if err != nil {
				return nil, fiber.ErrBadRequest
			}
			item.EndDate = &endDate
		}
	}
	if request.DivisionID != nil {
		item.DivisionID = strings.TrimSpace(*request.DivisionID)
	}
	if request.PositionID != nil {
		item.PositionID = strings.TrimSpace(*request.PositionID)
	}
	if request.EmploymentStatus != nil {
		employeeStatus := strings.TrimSpace(*request.EmploymentStatus)
		if employeeStatus == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "employee_status cannot be empty")
		}
		item.EmploymentStatus = employeeStatus
	}
	if request.Salary != nil {
		item.Salary = *request.Salary
	}
	if request.IsActive != nil {
		if *request.IsActive {
			if err := c.Repo.DeactivateActiveByEmployee(tx, item.EmployeeID); err != nil {
				c.Log.WithError(err).Error("Failed to deactivate previous active employee contracts")
				return nil, fiber.ErrInternalServerError
			}
		}
		item.IsActive = *request.IsActive
	}

	if item.EndDate != nil && item.StartDate > *item.EndDate {
		return nil, fiber.NewError(fiber.StatusBadRequest, "start_date cannot be greater than end_date")
	}

	if err := c.Repo.Update(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to update employee contract")
		return nil, fiber.ErrInternalServerError
	}

	reloaded, err := c.Repo.FindByID(tx, id, true)
	if err != nil {
		c.Log.WithError(err).Error("Failed to reload employee contract")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeContractToResponse(reloaded), nil
}

func (c *EmployeeContractUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item, err := c.Repo.FindByID(tx, id, false)
	if err != nil {
		c.Log.WithError(err).Error("Employee contract not found")
		return fiber.ErrNotFound
	}

	if err := c.Repo.Delete(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to delete employee contract")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}
