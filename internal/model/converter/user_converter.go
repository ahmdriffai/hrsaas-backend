package converter

import (
	"hr-sas/internal/entity"
	"hr-sas/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	var employeeResponse *model.EmployeeResponse

	if user.Employee != nil {
		employeeResponse = EmployeeToResponse(user.Employee)
	}

	return &model.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Image:         user.Image,
		Role:          user.Role,
		CompanyID:     user.CompanyID,
		EmailVerified: user.EmailVerified,
		Employee:      employeeResponse,
	}
}
