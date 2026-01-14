package converter

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
)

func EmployeeToResponse(employee *entity.Employee) *model.EmployeeResponse {
	return &model.EmployeeResponse{
		ID:             employee.ID,
		CompanyID:      employee.CompanyID,
		UserID:         employee.UserID,
		Fullname:       employee.Fullname,
		BirthPlace:     employee.BirthPlace,
		BirthDate:      employee.BirthDate,
		BlodType:       employee.BlodType,
		MaritalStatus:  employee.MaritalStatus,
		Religion:       employee.Religion,
		Phone:          employee.Phone,
		Timezone:       employee.Timezone,
		EmployeeNumber: employee.EmployeeNumber,
	}
}
