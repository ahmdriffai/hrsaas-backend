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

type TimeOffUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	Validate              *validator.Validate
	TimeOffRequestRepo    *repository.TimeOffRequestRepository
	TimeOffTypeRepo       *repository.TimeOffTypeRepository
	TimeOffBalanceRepo    *repository.TimeOffBalanceRepository
	TimeOffApprovalRepo   *repository.TimeOffApprovalRepository
	TimeOffAttachmentRepo *repository.TimeOffAttachmentRepository
}

func NewTimeOffUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	timeOffRequestRepo *repository.TimeOffRequestRepository,
	timeOffTypeRepo *repository.TimeOffTypeRepository,
	timeOffBalanceRepo *repository.TimeOffBalanceRepository,
	timeOffApprovalRepo *repository.TimeOffApprovalRepository,
	timeOffAttachmentRepo *repository.TimeOffAttachmentRepository,
) *TimeOffUseCase {
	return &TimeOffUseCase{
		DB:                    db,
		Log:                   log,
		Validate:              validate,
		TimeOffRequestRepo:    timeOffRequestRepo,
		TimeOffTypeRepo:       timeOffTypeRepo,
		TimeOffBalanceRepo:    timeOffBalanceRepo,
		TimeOffApprovalRepo:   timeOffApprovalRepo,
		TimeOffAttachmentRepo: timeOffAttachmentRepo,
	}
}

func (c *TimeOffUseCase) CreateRequest(ctx context.Context, employeeID string, request *model.CreateTimeOffRequest) (*model.TimeOffRequestResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	item := &entity.Time_Off_Requests{
		EmployeeId:    employeeID,
		TimeOffTypeId: request.TimeOffTypeID,
		StartDate:     mustParseEpoch(request.StartDate),
		EndDate:       mustParseEpoch(request.EndDate),
		RequestReason: &request.RequestReason,
		RequestStatus: "PENDING",
		CreatedAt:     nowEpoch(),
	}

	if err := c.TimeOffRequestRepo.Create(tx, item); err != nil {
		c.Log.WithError(err).Error("Failed to create time off request")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return &model.TimeOffRequestResponse{
		ID:            item.ID,
		EmployeeID:    item.EmployeeId,
		TimeOffTypeID: item.TimeOffTypeId,
		StartDate:     request.StartDate,
		EndDate:       request.EndDate,
		RequestedDays: request.RequestedDays,
		RequestReason: request.RequestReason,
		RequestStatus: item.RequestStatus,
		CreatedAt:     time.UnixMilli(item.CreatedAt),
	}, nil
}

func (c *TimeOffUseCase) ListRequests(ctx context.Context, request *model.SearchTimeOffRequest) ([]model.TimeOffRequestResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.TimeOffRequestRepo.List(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off requests")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffRequestResponse, len(items))
	for i, item := range items {
		requestReason := ""
		if item.RequestReason != nil {
			requestReason = *item.RequestReason
		}

		responses[i] = model.TimeOffRequestResponse{
			ID:            item.ID,
			EmployeeID:    item.EmployeeId,
			TimeOffTypeID: item.TimeOffTypeId,
			StartDate:     epochToDateString(item.StartDate),
			EndDate:       epochToDateString(item.EndDate),
			RequestedDays: 0,
			RequestReason: requestReason,
			RequestStatus: item.RequestStatus,
			CreatedAt:     time.UnixMilli(item.CreatedAt),
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

func (c *TimeOffUseCase) ListTypes(ctx context.Context) ([]model.TimeOffTypeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	items, err := c.TimeOffTypeRepo.List(tx)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off types")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffTypeResponse, len(items))
	for i, item := range items {
		responses[i] = model.TimeOffTypeResponse{
			ID:               item.ID,
			Name:             item.Name,
			Category:         item.Category,
			IsQuotaBased:     item.IsQuotaBased,
			DefaultQuotaDays: float64(item.DefaultQuotaDays),
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return responses, nil
}

func (c *TimeOffUseCase) ListBalances(ctx context.Context, employeeID string, request *model.SearchTimeOffBalanceRequest) ([]model.TimeOffBalanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate search query")
		return nil, fiber.ErrBadRequest
	}

	items, err := c.TimeOffBalanceRepo.ListByEmployee(tx, employeeID, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to list time off balances")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.TimeOffBalanceResponse, len(items))
	for i, item := range items {
		responses[i] = model.TimeOffBalanceResponse{
			ID:            item.ID,
			EmployeeID:    item.EmployeeId,
			TimeOffTypeID: item.TimeOffTypeId,
			PeriodYear:    item.PeriodYear,
			EntitledDays:  float64(item.EntitledDays),
			UsedDays:      float64(item.UsedDays),
			RemainingDays: float64(item.RemainingDays),
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return responses, nil
}

// TODO: Replace with a proper date parsing strategy shared across the codebase.
func mustParseEpoch(date string) int64 {
	parsed, _ := lib.ParseDateToUnixMilli(date)
	return parsed
}

func nowEpoch() int64 {
	return time.Now().UnixMilli()
}

func epochToDateString(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.UnixMilli(ts).Format("2006-01-02")
}
