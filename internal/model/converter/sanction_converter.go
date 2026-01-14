package converter

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
)

func SanctionToResponse(sanction *entity.Sanction) *model.SanctionResponse {
	return &model.SanctionResponse{
		ID:          sanction.ID,
		CompanyID:   sanction.CompanyID,
		Description: sanction.Description,
		Note:        sanction.Note,
		Name:        sanction.Name,
		CreatedAt:   sanction.CreatedAt.Local().String(),
		UpdatedAt:   sanction.UpdatedAt.Local().String(),
	}
}

func EmSancToResponse(emSanc *entity.EmployeeSanction) *model.EmSancResponse {
	return &model.EmSancResponse{
		ID:          emSanc.ID,
		EmployeeID:  emSanc.EmployeeID,
		SanctionID:  emSanc.SanctionID,
		Reason:      emSanc.Reason,
		StartDate:   emSanc.StartDate,
		EndDate:     emSanc.EndDate,
		CompanyID:   emSanc.CompanyID,
		Employee:    *EmployeeToResponse(&emSanc.Employee),
		DocumentUrl: emSanc.DocumentUrl,
		Sanction:    *SanctionToResponse(&emSanc.Sanction),
		CreatedAt:   emSanc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   emSanc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
