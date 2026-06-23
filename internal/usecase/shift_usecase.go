package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ShiftUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	ShiftRepository *repository.ShiftRepository
	ShiftDayRepo    *repository.ShiftDayRepository
}

func NewShiftUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	shiftRepository *repository.ShiftRepository,
	shiftDayRepo *repository.ShiftDayRepository,
) *ShiftUseCase {
	return &ShiftUseCase{
		DB:              db,
		Log:             log,
		Validate:        validate,
		ShiftRepository: shiftRepository,
		ShiftDayRepo:    shiftDayRepo,
	}
}

// Create Shift
func (c *ShiftUseCase) Create(ctx context.Context, request *model.CreateShiftRequest) (*model.ShiftResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	shift := &entity.Shift{
		CompanyID:     request.CompanyID,
		Name:          request.Name,
		LateTolerance: request.LateTolerance,
	}

	if err := c.ShiftRepository.Create(tx, shift); err != nil {
		c.Log.WithError(err).Error("Failed to create shift")
		return nil, fiber.ErrInternalServerError
	}

	if len(request.ShiftDayRequests) > 0 {
		if len(request.ShiftDayRequests) != 7 {
			c.Log.Error("Invalid shift days count: must provide exactly 7 weekdays")
			return nil, fiber.ErrBadRequest
		}

		var weekdayMask uint8
		for _, shiftDayRequest := range request.ShiftDayRequests {
			bit := uint8(1 << (shiftDayRequest.Weekday - 1))
			if weekdayMask&bit != 0 {
				c.Log.Error("Invalid shift days: duplicate weekday")
				return nil, fiber.ErrBadRequest
			}
			weekdayMask |= bit
		}
		if weekdayMask != 0b1111111 {
			c.Log.Error("Invalid shift days: weekdays must cover 1 to 7")
			return nil, fiber.ErrBadRequest
		}

		shiftDays := make([]entity.ShiftDay, 0, len(request.ShiftDayRequests))
		for _, shiftDayRequest := range request.ShiftDayRequests {
			checkIn, _ := lib.ParseTimeToUnixMilli(shiftDayRequest.CheckIn)

			checkOut, _ := lib.ParseTimeToUnixMilli(shiftDayRequest.CheckOut)

			breakStart, _ := lib.ParseTimeToUnixMilli(shiftDayRequest.BreakStart)

			breakEnd, _ := lib.ParseTimeToUnixMilli(shiftDayRequest.BreakEnd)

			shiftDays = append(shiftDays, entity.ShiftDay{
				ShiftID:         shift.ID,
				Weekday:         shiftDayRequest.Weekday,
				DayType:         shiftDayRequest.DayType,
				CheckIn:         checkIn,
				CheckOut:        checkOut,
				BreakStart:      breakStart,
				BreakEnd:        breakEnd,
				MaxBreakMinutes: shiftDayRequest.MaxBreakMinutes,
			})
		}

		if err := c.ShiftDayRepo.CreateBatch(tx, shiftDays); err != nil {
			c.Log.WithError(err).Error("Failed to create shift days")
			return nil, fiber.ErrInternalServerError
		}
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return &model.ShiftResponse{
		ID:            shift.ID,
		CompanyID:     shift.CompanyID,
		Name:          shift.Name,
		LateTolerance: shift.LateTolerance,
		CreatedAt:     shift.CreatedAt,
		UpdatedAt:     shift.UpdatedAt,
	}, nil
}

func (c *ShiftUseCase) AssignEmployee(ctx context.Context, request *model.AssignEmployeeToShiftRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return fiber.ErrBadRequest
	}

	var shiftTotal int64
	shiftTotal, err := c.ShiftRepository.CountByIDAndCompanyID(tx, request.ShiftID, request.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to check shift existence")
		return fiber.ErrInternalServerError
	}
	if shiftTotal == 0 {
		c.Log.Error("Shift not found")
		return fiber.ErrBadRequest
	}

	employeeTotal, err := c.ShiftRepository.CountEmployeeByIDAndCompanyID(tx, request.EmployeeID, request.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to check employee existence")
		return fiber.ErrInternalServerError
	}
	if employeeTotal == 0 {
		c.Log.Error("Employee not found")
		return fiber.ErrBadRequest
	}

	if err := c.ShiftRepository.DeleteEmployeeShiftsByEmployeeID(tx, request.EmployeeID); err != nil {
		c.Log.WithError(err).Error("Failed to remove previous employee shift")
		return fiber.ErrInternalServerError
	}

	if err := c.ShiftRepository.AssignEmployeeToShift(tx, request.EmployeeID, request.ShiftID); err != nil {
		c.Log.WithError(err).Error("Failed to assign employee to shift")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *ShiftUseCase) Search(ctx context.Context, request *model.SearchShiftRequest) ([]model.ShiftResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, 0, fiber.ErrBadRequest
	}
	shifts, total, err := c.ShiftRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to search shifts")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.ShiftResponse, 0, len(shifts))

	for _, shift := range shifts {
		shifDays := make([]model.ShiftDayResponse, len(shift.ShiftDays))
		for i, shifDay := range shift.ShiftDays {
			shifDays[i] = *model.ShiftDayToResponse(&shifDay)
		}

		responses = append(responses, model.ShiftResponse{
			ID:            shift.ID,
			CompanyID:     shift.CompanyID,
			Name:          shift.Name,
			LateTolerance: shift.LateTolerance,
			ShiftDays:     shifDays,
			CreatedAt:     shift.CreatedAt,
			UpdatedAt:     shift.UpdatedAt,
		})
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

func (c *ShiftUseCase) DetailShift(ctx context.Context, request *model.DetailShifRequest) (*model.ShiftResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// cek user exist
	count, err := c.ShiftRepository.CountByIDAndCompanyID(tx, request.ShiftID, request.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to count user by email")
		return nil, fiber.ErrInternalServerError
	}

	if count == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Shift not found.")
	}

	var shift = new(entity.Shift)

	err = c.ShiftRepository.FindByIdAndCompany(tx, shift, request.ShiftID, request.CompanyID, "Employees", "ShiftDays")
	if err != nil {
		c.Log.WithError(err).Error("Failed to find ")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.ShiftToResponse(shift), nil
}

func (c *ShiftUseCase) Update(ctx context.Context, shiftID string, companyID string, request *model.UpdateShiftRequest) (*model.ShiftResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	shift := new(entity.Shift)
	if err := c.ShiftRepository.FindByIdAndCompany(tx, shift, shiftID, companyID, "Employees", "ShiftDays"); err != nil {
		c.Log.WithError(err).Error("Shift not found")
		return nil, fiber.ErrNotFound
	}

	if request.Name != nil {
		name := strings.TrimSpace(*request.Name)
		if name == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		shift.Name = name
	}

	if request.LateTolerance != nil {
		shift.LateTolerance = *request.LateTolerance
	}

	if request.ShiftDayRequests != nil {
		if len(*request.ShiftDayRequests) != 7 {
			c.Log.Error("Invalid shift days count: must provide exactly 7 weekdays")
			return nil, fiber.ErrBadRequest
		}

		var weekdayMask uint8
		shiftDays := make([]entity.ShiftDay, 0, len(*request.ShiftDayRequests))
		for _, shiftDayRequest := range *request.ShiftDayRequests {
			bit := uint8(1 << (shiftDayRequest.Weekday - 1))
			if weekdayMask&bit != 0 {
				c.Log.Error("Invalid shift days: duplicate weekday")
				return nil, fiber.ErrBadRequest
			}
			weekdayMask |= bit

			checkIn, _ := lib.ParseTimeToUnixMilli(shiftDayRequest.CheckIn)
			checkOut, _ := lib.ParseTimeToUnixMilli(shiftDayRequest.CheckOut)
			breakStart, _ := lib.ParseTimeToUnixMilli(shiftDayRequest.BreakStart)
			breakEnd, _ := lib.ParseTimeToUnixMilli(shiftDayRequest.BreakEnd)

			shiftDays = append(shiftDays, entity.ShiftDay{
				ShiftID:         shift.ID,
				Weekday:         shiftDayRequest.Weekday,
				DayType:         shiftDayRequest.DayType,
				CheckIn:         checkIn,
				CheckOut:        checkOut,
				BreakStart:      breakStart,
				BreakEnd:        breakEnd,
				MaxBreakMinutes: shiftDayRequest.MaxBreakMinutes,
			})
		}

		if weekdayMask != 0b1111111 {
			c.Log.Error("Invalid shift days: weekdays must cover 1 to 7")
			return nil, fiber.ErrBadRequest
		}

		if err := c.ShiftDayRepo.DeleteByShiftID(tx, shift.ID); err != nil {
			c.Log.WithError(err).Error("Failed to delete previous shift days")
			return nil, fiber.ErrInternalServerError
		}

		if err := c.ShiftDayRepo.CreateBatch(tx, shiftDays); err != nil {
			c.Log.WithError(err).Error("Failed to create shift days")
			return nil, fiber.ErrInternalServerError
		}

		shift.ShiftDays = shiftDays
	}

	shift.UpdatedAt = time.Now().UnixMilli()

	if err := c.ShiftRepository.Update(tx, shift); err != nil {
		c.Log.WithError(err).Error("Failed to update shift")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.ShiftToResponse(shift), nil
}

func (c *ShiftUseCase) BulkAssignEmployees(ctx context.Context, request *model.BulkAssignEmployeesToShiftRequest) (*model.ShiftResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	count, err := c.ShiftRepository.CountByIDAndCompanyID(tx, request.ShiftID, request.CompanyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to check shift")
		return nil, fiber.ErrInternalServerError
	}
	if count == 0 {
		return nil, fiber.ErrNotFound
	}

	if err := c.ShiftRepository.DeleteEmployeesByShiftID(tx, request.ShiftID); err != nil {
		c.Log.WithError(err).Error("Failed to remove existing employees")
		return nil, fiber.ErrInternalServerError
	}

	if len(request.EmployeeIDs) > 0 {
		if err := c.ShiftRepository.BulkAssignEmployees(tx, request.EmployeeIDs, request.ShiftID); err != nil {
			c.Log.WithError(err).Error("Failed to bulk assign employees")
			return nil, fiber.ErrInternalServerError
		}
	}

	var shift = new(entity.Shift)
	if err := c.ShiftRepository.FindByIdAndCompany(tx, shift, request.ShiftID, request.CompanyID, "Employees", "ShiftDays"); err != nil {
		c.Log.WithError(err).Error("Failed to reload shift")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.ShiftToResponse(shift), nil
}

func (c *ShiftUseCase) Delete(ctx context.Context, shiftID string, companyID string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	shift := new(entity.Shift)
	if err := c.ShiftRepository.FindByIdAndCompany(tx, shift, shiftID, companyID); err != nil {
		c.Log.WithError(err).Error("Shift not found")
		return fiber.ErrNotFound
	}

	if err := c.ShiftRepository.Delete(tx, shift); err != nil {
		c.Log.WithError(err).Error("Failed to delete shift")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}
