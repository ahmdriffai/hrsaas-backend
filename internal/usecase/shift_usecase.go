package usecase

import (
	"context"
	"fmt"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
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

	parseShiftTime := func(value string) (time.Time, error) {
		if value == "" {
			return time.Time{}, nil
		}

		t, err := time.Parse("15:04", value)
		if err == nil {
			return t, nil
		}

		t, err = time.Parse("15:04:05", value)
		if err != nil {
			return time.Time{}, err
		}

		return t, nil
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

		shiftDays := make([]entity.ShiftDays, 0, len(request.ShiftDayRequests))
		for _, shiftDayRequest := range request.ShiftDayRequests {
			checkIn, err := parseShiftTime(shiftDayRequest.CheckIn)
			if err != nil {
				c.Log.WithError(err).Error("Failed to parse check_in")
				return nil, fiber.ErrBadRequest
			}

			checkOut, err := parseShiftTime(shiftDayRequest.CheckOut)
			if err != nil {
				c.Log.WithError(err).Error("Failed to parse check_out")
				return nil, fiber.ErrBadRequest
			}

			breakStart, err := parseShiftTime(shiftDayRequest.BreakStart)
			if err != nil {
				c.Log.WithError(err).Error("Failed to parse break_start")
				return nil, fiber.ErrBadRequest
			}

			breakEnd, err := parseShiftTime(shiftDayRequest.BreakEnd)
			if err != nil {
				c.Log.WithError(err).Error("Failed to parse break_end")
				return nil, fiber.ErrBadRequest
			}

			fmt.Println(shift.ID)

			shiftDays = append(shiftDays, entity.ShiftDays{
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
		CreatedAt:     shift.CreatedAt.Format(time.DateTime),
		UpdatedAt:     shift.UpdatedAt.Format(time.DateTime),
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

	shiftDaysByShiftID := make(map[string][]model.ShiftDayResponse, len(shifts))
	if len(shifts) > 0 {
		shiftIDs := make([]string, 0, len(shifts))
		for _, shift := range shifts {
			shiftIDs = append(shiftIDs, shift.ID)
		}

		type shiftDayRow struct {
			ShiftID         string `gorm:"column:shift_id"`
			Weekday         int    `gorm:"column:weekday"`
			DayType         string `gorm:"column:day_type"`
			CheckIn         string `gorm:"column:check_in"`
			CheckOut        string `gorm:"column:check_out"`
			BreakStart      string `gorm:"column:break_start"`
			BreakEnd        string `gorm:"column:break_end"`
			MaxBreakMinutes int    `gorm:"column:max_break_minutes"`
		}

		var shiftDays []shiftDayRow
		if err := tx.
			Table("shift_days").
			Select("shift_id, weekday, day_type, check_in, check_out, break_start, break_end, max_break_minutes").
			Where("shift_id IN ?", shiftIDs).
			Order("shift_id ASC, weekday ASC").
			Find(&shiftDays).Error; err != nil {
			c.Log.WithError(err).Error("Failed to search shift days")
			return nil, 0, fiber.ErrInternalServerError
		}

		for _, day := range shiftDays {
			shiftDaysByShiftID[day.ShiftID] = append(shiftDaysByShiftID[day.ShiftID], model.ShiftDayResponse{
				Weekday:         day.Weekday,
				DayType:         day.DayType,
				CheckIn:         day.CheckIn,
				CheckOut:        day.CheckOut,
				BreakStart:      day.BreakStart,
				BreakEnd:        day.BreakEnd,
				MaxBreakMinutes: day.MaxBreakMinutes,
			})
		}
	}

	responses := make([]model.ShiftResponse, 0, len(shifts))
	for _, shift := range shifts {
		responses = append(responses, model.ShiftResponse{
			ID:            shift.ID,
			CompanyID:     shift.CompanyID,
			Name:          shift.Name,
			LateTolerance: shift.LateTolerance,
			ShiftDays:     shiftDaysByShiftID[shift.ID],
			CreatedAt:     shift.CreatedAt.Format(time.DateTime),
			UpdatedAt:     shift.UpdatedAt.Format(time.DateTime),
		})
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}
