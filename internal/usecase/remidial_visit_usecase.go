package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RemidialVisitUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	Validate                *validator.Validate
	RemidialVisitRepository *repository.RemidialVisitRepository
	EmployeeRepository      *repository.EmployeeRepository
}

func NewRemidialVisitUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	remidialVisitRepo *repository.RemidialVisitRepository,
	employeeRepo *repository.EmployeeRepository,
) *RemidialVisitUseCase {
	return &RemidialVisitUseCase{
		DB:                      db,
		Log:                     log,
		Validate:                validate,
		RemidialVisitRepository: remidialVisitRepo,
		EmployeeRepository:      employeeRepo,
	}
}

func (c *RemidialVisitUseCase) SearchNasabah(ctx context.Context, request *model.SearchNasabahRequest) ([]model.NasabahData, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	data, err := c.RemidialVisitRepository.SearchNasabah(ctx, request)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *RemidialVisitUseCase) Create(ctx context.Context, request *model.CreateRemidialVisitRequest) (*model.RemidialVisitResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	visit := &entity.RemidialVisit{
		CompanyID:                 request.CompanyID,
		EmployeeID:                request.EmployeeID,
		NasabahID:                 request.NasabahID,
		NasabahName:               request.NasabahName,
		NoPjm:                     request.NoPjm,
		LoanType:                  request.LoanType,
		Unit:                      request.Unit,
		Collectibility:            request.Collectibility,
		LoanLimit:                 request.LoanLimit,
		OutstandingBalance:        request.OutstandingBalance,
		OverduePrincipal:          request.OverduePrincipal,
		OverdueInterest:           request.OverdueInterest,
		OverdueTotal:              request.OverdueTotal,
		OverduePrincipalFrequency: request.OverduePrincipalFrequency,
		OverdueInterestFrequency:  request.OverdueInterestFrequency,
		OverduePrincipalDays:      request.OverduePrincipalDays,
		OverdueInterestDays:       request.OverdueInterestDays,
		LoanStatus:                request.LoanStatus,
		TotalPaid:                 request.TotalPaid,
		Note:                      request.Note,
	}

	if err := c.RemidialVisitRepository.Create(tx, visit); err != nil {
		c.Log.WithError(err).Error("Failed to create remidial visit")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return c.toResponse(ctx, visit), nil
}

func (c *RemidialVisitUseCase) List(ctx context.Context, request *model.SearchRemidialVisitRequest) ([]model.RemidialVisitResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.RemidialVisitRepository.List(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get list of remidial visits")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.RemidialVisitResponse, len(items))
	for i, item := range items {
		responses[i] = *c.toResponse(ctx, &item)
	}

	return responses, total, nil
}

func (c *RemidialVisitUseCase) toResponse(ctx context.Context, visit *entity.RemidialVisit) *model.RemidialVisitResponse {
	var employeeName string
	employee := new(entity.Employee)
	if err := c.EmployeeRepository.FindByIdAndCompany(c.DB.WithContext(ctx), employee, visit.EmployeeID, visit.CompanyID); err == nil {
		employeeName = employee.Fullname
	}

	return &model.RemidialVisitResponse{
		ID:                        visit.ID,
		CompanyID:                 visit.CompanyID,
		EmployeeID:                visit.EmployeeID,
		EmployeeName:              employeeName,
		NasabahID:                 visit.NasabahID,
		NasabahName:               visit.NasabahName,
		NoPjm:                     visit.NoPjm,
		LoanType:                  visit.LoanType,
		Unit:                      visit.Unit,
		Collectibility:            visit.Collectibility,
		LoanLimit:                 visit.LoanLimit,
		OutstandingBalance:        visit.OutstandingBalance,
		OverduePrincipal:          visit.OverduePrincipal,
		OverdueInterest:           visit.OverdueInterest,
		OverdueTotal:              visit.OverdueTotal,
		OverduePrincipalFrequency: visit.OverduePrincipalFrequency,
		OverdueInterestFrequency:  visit.OverdueInterestFrequency,
		OverduePrincipalDays:      visit.OverduePrincipalDays,
		OverdueInterestDays:       visit.OverdueInterestDays,
		LoanStatus:                visit.LoanStatus,
		TotalPaid:                 visit.TotalPaid,
		Note:                      visit.Note,
		CreatedAt:                 visit.CreatedAt,
	}
}
