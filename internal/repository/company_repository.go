package repository

import (
	"hr-sas/internal/entity"

	"github.com/sirupsen/logrus"
)

type CompanyRepository struct {
	Repository[entity.Company]
	Log *logrus.Logger
}

func NewCompanyRepository(log *logrus.Logger) *CompanyRepository {
	return &CompanyRepository{
		Log: log,
	}
}
