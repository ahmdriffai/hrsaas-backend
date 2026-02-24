package repository

import (
	"hr-sas/internal/entity"

	"github.com/sirupsen/logrus"
)

type AttendanceLogRepository struct {
	Repository[entity.AttendanceLog]
	Log *logrus.Logger
}

func NewAttendanceLogRepository(log *logrus.Logger) *AttendanceLogRepository {
	return &AttendanceLogRepository{
		Log: log,
	}
}
