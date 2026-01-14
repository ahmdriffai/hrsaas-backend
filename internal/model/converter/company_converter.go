package converter

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
)

func CompanyToResponse(company *entity.Company) *model.CompanyResponse {
	return &model.CompanyResponse{
		ID:      company.ID,
		Name:    company.Name,
		LogoUrl: company.LogoUrl,
		BussinessField: company.BussinessField,
		Address:        company.Address,
		Province:       company.Province,
		City:           company.City,
		District:       company.District,
		Village:        company.Village,
		ZipCode:        company.ZipCode,
		PhoneNumber:    company.PhoneNumber,
		FaxNumber:      company.FaxNumber,
		Email:          company.Email,
		Website:        company.Website,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}
}
