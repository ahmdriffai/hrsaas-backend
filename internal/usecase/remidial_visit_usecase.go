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

func (c *RemidialVisitUseCase) SearchNasabah(ctx context.Context, request *model.SearchNasabahRequest) ([]model.SearchNasabahResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	data, err := c.RemidialVisitRepository.SearchNasabah(ctx, request)
	if err != nil {
		return nil, err
	}

	var responses []model.SearchNasabahResponse
	for _, nasabah := range data {
		responses = append(responses, model.SearchNasabahResponse{
			BranchCode:                nasabah.KodeCabang,
			NoPjm:                     nasabah.NoPjm,
			NasabahID:                 nasabah.NasabahID,
			NasabahName:               nasabah.Nama,
			NIK:                       nasabah.NIK,
			PlaceOfBirth:              nasabah.TmpLahir,
			DateOfBirth:               nasabah.TglLahir,
			Address:                   nasabah.Alamat,
			Phone:                     nasabah.Phone,
			Email:                     nasabah.Email,
			LoanType:                  nasabah.JnsPjm,
			Unit:                      nasabah.Unit,
			Collectibility:            nasabah.Col,
			LoanLimit:                 nasabah.Plafond,
			OutstandingBalance:        nasabah.Saldo,
			OverduePrincipal:          nasabah.TunggakanPokok,
			OverdueInterest:           nasabah.TunggakanBunga,
			OverdueTotal:              nasabah.TunggakanTotal,
			OverduePrincipalFrequency: nasabah.TunggakanPokokFrek,
			OverdueInterestFrequency:  nasabah.TunggakanBungaFrek,
			OverduePrincipalDays:      nasabah.TunggakanPokokHari,
			OverdueInterestDays:       nasabah.TunggakanBungaHari,
			LoanStatus:                nasabah.StatusPinjaman,
		})
	}

	return responses, nil
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
		Img_Url:                   request.ImgUrl,
		Latitude:                  request.Lat,
		Longitude:                 request.Lng,
		NasabahID:                 request.Pinjaman.NasabahID,
		NasabahName:               request.Pinjaman.NasabahName,
		NoPjm:                     request.Pinjaman.NoPjm,
		LoanType:                  request.Pinjaman.LoanType,
		Unit:                      request.Pinjaman.Unit,
		Collectibility:            request.Pinjaman.Collectibility,
		LoanLimit:                 request.Pinjaman.LoanLimit,
		OutstandingBalance:        request.Pinjaman.OutstandingBalance,
		OverduePrincipal:          request.Pinjaman.OverduePrincipal,
		OverdueInterest:           request.Pinjaman.OverdueInterest,
		OverdueTotal:              request.Pinjaman.OverdueTotal,
		OverduePrincipalFrequency: request.Pinjaman.OverduePrincipalFrequency,
		OverdueInterestFrequency:  request.Pinjaman.OverdueInterestFrequency,
		OverduePrincipalDays:      request.Pinjaman.OverduePrincipalDays,
		OverdueInterestDays:       request.Pinjaman.OverdueInterestDays,
		LoanStatus:                request.Pinjaman.LoanStatus,
		TotalPaid:                 request.TotalPaid,
		Commitment:                request.Commitment,
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
		ID:           visit.ID,
		CompanyID:    visit.CompanyID,
		EmployeeID:   visit.EmployeeID,
		EmployeeName: employeeName,
		ImgUrl:       visit.Img_Url,
		Lat:          visit.Latitude,
		Lng:          visit.Longitude,
		Pinjaman: model.DetailPinjamanResponse{
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
		},
		TotalPaid:  visit.TotalPaid,
		Commitment: visit.Commitment,
		CreatedAt:  visit.CreatedAt,
	}
}

func (c *RemidialVisitUseCase) Update(ctx context.Context, request *model.UpdateRemidialVisitRequest) (*model.RemidialVisitResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	visit := new(entity.RemidialVisit)
	if err := c.RemidialVisitRepository.FindById(tx, visit, request.ID); err != nil {
		c.Log.WithError(err).Error("Remidial visit not found")
		return nil, fiber.ErrNotFound
	}

	visit.Img_Url = request.ImgUrl
	visit.Latitude = request.Lat
	visit.Longitude = request.Lng
	visit.TotalPaid = request.TotalPaid
	visit.Commitment = request.Commitment

	if err := c.RemidialVisitRepository.Update(tx, visit); err != nil {
		c.Log.WithError(err).Error("Failed to update remidial visit")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return c.toResponse(ctx, visit), nil
}

func (c *RemidialVisitUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	visit := new(entity.RemidialVisit)
	if err := c.RemidialVisitRepository.FindById(tx, visit, id); err != nil {
		c.Log.WithError(err).Error("Kunjungan tidak ditemukan")
		return fiber.ErrNotFound
	}

	if err := c.RemidialVisitRepository.Delete(tx, visit); err != nil {
		c.Log.WithError(err).Error("Gagal menghapus kunjungan")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Gagal menyelesaikan transaksi")
		return fiber.ErrInternalServerError
	}

	return nil
}
