package usecase

import (
	"context"
	"errors"
	"fmt"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttendanceUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	AttendanceRepository *repository.AttendanceRepository
	LocationRepository   *repository.OfficeLocationRepository
	ShiftRepository      *repository.ShiftRepository
	ShiftDayRepo         *repository.ShiftDayRepository
	AttendanceLogRepo    *repository.AttendanceLogRepository
	FaceServiceURL       string
}

func NewAttendanceUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	attendanceRepository *repository.AttendanceRepository,
	locationRepository *repository.OfficeLocationRepository,
	shiftRepository *repository.ShiftRepository,
	shiftDayRepo *repository.ShiftDayRepository,
	attendanceLogRepo *repository.AttendanceLogRepository,
	faceServiceURL string,
) *AttendanceUseCase {
	return &AttendanceUseCase{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		AttendanceRepository: attendanceRepository,
		LocationRepository:   locationRepository,
		ShiftRepository:      shiftRepository,
		ShiftDayRepo:         shiftDayRepo,
		AttendanceLogRepo:    attendanceLogRepo,
		FaceServiceURL:       faceServiceURL,
	}
}

func (c *AttendanceUseCase) validateLocation(tx *gorm.DB, employeeID string, lat, lng float64) (distance float64, err error) {
	locations, err := c.LocationRepository.GetByEmployeeID(tx, employeeID)
	if err != nil {
		return 0, fiber.ErrInternalServerError
	}

	for _, location := range locations {
		locLat, err := strconv.ParseFloat(location.Lat, 64)
		if err != nil {
			continue
		}
		locLng, err := strconv.ParseFloat(location.Lng, 64)
		if err != nil {
			continue
		}
		d := lib.DistanceMeter(lat, lng, locLat, locLng)
		if d <= float64(location.Radius) {
			return d, nil
		}
	}

	return 0, fiber.NewError(400, "Anda diluar jangkauan")
}

func (c *AttendanceUseCase) Search(ctx context.Context, request *model.SearchAttendanceRequest) ([]model.AttendanceResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	attendances, total, err := c.AttendanceRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting attendances")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.AttendanceResponse, len(attendances))
	for i, attendance := range attendances {
		responses[i] = *model.AttendandeToResponse(&attendance)
	}

	return responses, total, nil
}

func (c *AttendanceUseCase) Detail(ctx context.Context, requestID string, companyID string) (*model.AttendanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	attendance := new(entity.Attendance)
	if err := c.AttendanceRepository.FindByIdAndCompany(tx, attendance, requestID, companyID, "Employee", "Employee.EmployeeContract", "Employee.EmployeeContract.Position"); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Attendance tidak ditemukan")
		}
		c.Log.WithError(err).Error("Failed to find attendance")
		return nil, fiber.ErrInternalServerError
	}

	logs, err := c.AttendanceLogRepo.FindByAttendanceID(tx, attendance.ID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to find attendance logs")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	response := model.AttendandeToResponse(attendance)
	response.Logs = make([]model.AttendanceLogResponse, len(logs))
	for i, log := range logs {
		response.Logs[i] = *model.AttendanceLogToResponse(&log)
	}

	return response, nil
}

func (c *AttendanceUseCase) Update(ctx context.Context, requestID string, companyID string, request *model.UpdateAttendanceRequest) (*model.AttendanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	attendance := new(entity.Attendance)
	if err := c.AttendanceRepository.FindByIdAndCompany(tx, attendance, requestID, companyID, "Employee", "Employee.EmployeeContract", "Employee.EmployeeContract.Position"); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Attendance tidak ditemukan")
		}
		c.Log.WithError(err).Error("Failed to find attendance")
		return nil, fiber.ErrInternalServerError
	}

	if request.Date != nil {
		attendance.Date = *request.Date
	}
	if request.CheckInTime != nil {
		attendance.CheckInTime = *request.CheckInTime
	}
	if request.CheckOutTime != nil {
		attendance.CheckOutTime = *request.CheckOutTime
	}
	if request.TotalWorkMinutes != nil {
		attendance.TotalWorkMinutes = *request.TotalWorkMinutes
	}
	if request.TotalBreakMinutes != nil {
		attendance.TotalBreakMinutes = *request.TotalBreakMinutes
	}
	if request.Status != nil {
		attendance.Status = *request.Status
	}
	if request.IsApproved != nil {
		attendance.IsApproved = *request.IsApproved
	}

	if err := c.AttendanceRepository.Update(tx, attendance); err != nil {
		c.Log.WithError(err).Error("Failed to update attendance")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.AttendandeToResponse(attendance), nil
}

func (c *AttendanceUseCase) Delete(ctx context.Context, requestID string, companyID string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	attendance := new(entity.Attendance)
	if err := c.AttendanceRepository.FindByIdAndCompany(tx, attendance, requestID, companyID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Attendance tidak ditemukan")
		}
		c.Log.WithError(err).Error("Failed to find attendance")
		return fiber.ErrInternalServerError
	}

	if err := c.AttendanceLogRepo.DeleteByAttendanceID(tx, attendance.ID); err != nil {
		c.Log.WithError(err).Error("Failed to delete attendance logs")
		return fiber.ErrInternalServerError
	}

	if err := c.AttendanceRepository.Delete(tx, attendance); err != nil {
		c.Log.WithError(err).Error("Failed to delete attendance")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *AttendanceUseCase) CheckIn(ctx context.Context, request *model.CheckInAttendanceRequest) (*model.AttendanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.EmployeeID == "" {
		c.Log.Error("Bad Request")
		return nil, fiber.NewError(400, "User tidak bisa melakukan check-in karena bukan karyawan")
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	now := time.Now()
	nowMilli := now.UnixMilli()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()

	shifts, err := c.ShiftRepository.FindByEmployeeID(tx, request.EmployeeID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Shift tidak ditemukan")
	}
	if len(shifts) == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Shift tidak ditemukan")
	}
	shift := shifts[0]

	var existingAttendance entity.Attendance
	err = c.AttendanceRepository.FindByEmployeeIDAndDate(tx, &existingAttendance, request.EmployeeID, nowMilli)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fiber.ErrInternalServerError
	}
	if existingAttendance.ID != "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Sudah check-in hari ini")
	}

	locations, err := c.LocationRepository.GetByEmployeeID(tx, request.EmployeeID)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	fmt.Println(locations)

	isInRange := false
	locationDistance := 0.0
	for _, location := range locations {
		lat, err := strconv.ParseFloat(location.Lat, 64)
		if err != nil {
			continue
		}
		lng, err := strconv.ParseFloat(location.Lng, 64)
		if err != nil {
			continue
		}
		distance := lib.DistanceMeter(request.Lat, request.Lng, lat, lng)
		if distance <= float64(location.Radius) {
			isInRange = true
			locationDistance = distance
			break
		}
	}
	if !isInRange {
		return nil, fiber.NewError(400, "Anda diluar jangkauan")
	}

	attendance := &entity.Attendance{
		CompanyID:   request.CompanyID,
		EmployeeID:  request.EmployeeID,
		Date:        startOfDay,
		CheckInTime: nowMilli,
		Status:      "HADIR",
	}

	// faceResult, err := lib.VerifyFace(c.FaceServiceURL+"/verify-face", request.EmployeeID, request.FaceImageUrl)
	// if err != nil {
	// 	return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal memverifikasi wajah")
	// }

	attendanceLog := &entity.AttendanceLog{
		Type:               "CHECK_IN",
		Time:               nowMilli,
		Lat:                request.Lat,
		Lng:                request.Lng,
		LocationDistance:   locationDistance,
		IsLocationVerified: true,
		// IsFaceVerified:     faceResult.IsVerified,
		// FaceConfidence:     faceResult.Confidence,
		IsFaceVerified: false,
		FaceConfidence: 0,
		FaceImageURL:   request.FaceImageUrl,
		DeviceInfo:     request.DeviceInfo,
	}

	weekday := int(now.Weekday())
	if weekday == int(time.Sunday) {
		weekday = 7
	}

	var shiftDay entity.ShiftDay
	if err := c.ShiftDayRepo.FindByShiftIDAndWeekday(tx, &shiftDay, shift.ID, weekday); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Shift hari ini tidak ditemukan")
		}
		return nil, fiber.ErrInternalServerError
	}

	checkInDateType := time.UnixMilli(shiftDay.CheckIn)
	startTimeToday := time.Date(
		now.Year(), now.Month(), now.Day(),
		checkInDateType.Hour(), checkInDateType.Minute(), checkInDateType.Second(), checkInDateType.Nanosecond(),
		now.Location(),
	)

	if now.After(startTimeToday.Add(time.Duration(shift.LateTolerance) * time.Minute)) {
		attendance.Status = "TERLAMBAT"
	}

	if err := c.AttendanceRepository.Create(tx, attendance); err != nil {
		c.Log.WithError(err).Error("Failed to create attendance")
		return nil, fiber.ErrInternalServerError
	}

	attendanceLog.AttendanceID = attendance.ID

	if err := c.AttendanceLogRepo.Create(tx, attendanceLog); err != nil {
		c.Log.WithError(err).Error("Failed to create attendance log")
		return nil, fiber.ErrInternalServerError
	}

	if err := c.AttendanceRepository.FindById(tx, attendance, attendance.ID, "Employee", "Employee.EmployeeContract", "Employee.EmployeeContract.Position"); err != nil {
		c.Log.WithError(err).Error("Failed to reload attendance with employee")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.AttendandeToResponse(attendance), nil
}

func (c *AttendanceUseCase) CheckOut(ctx context.Context, request *model.CheckInAttendanceRequest) (*model.AttendanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.EmployeeID == "" {
		return nil, fiber.NewError(400, "User tidak bisa melakukan check-out karena bukan karyawan")
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	now := time.Now()
	nowMilli := now.UnixMilli()

	var attendance entity.Attendance
	err := c.AttendanceRepository.FindByEmployeeIDAndDate(tx, &attendance, request.EmployeeID, nowMilli)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Belum check-in hari ini")
		}
		return nil, fiber.ErrInternalServerError
	}

	if attendance.CheckOutTime != 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Sudah check-out hari ini")
	}

	distance, err := c.validateLocation(tx, request.EmployeeID, request.Lat, request.Lng)
	if err != nil {
		return nil, err
	}

	// faceResult, err := lib.VerifyFace(c.FaceServiceURL+"/verify-face", request.EmployeeID, request.FaceImageUrl)
	// if err != nil {
	// 	return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal memverifikasi wajah")
	// }

	attendance.CheckOutTime = nowMilli

	checkInTime := time.UnixMilli(attendance.CheckInTime)
	totalWorkMinutes := int(now.Sub(checkInTime).Minutes()) - attendance.TotalBreakMinutes
	if totalWorkMinutes < 0 {
		totalWorkMinutes = 0
	}
	attendance.TotalWorkMinutes = totalWorkMinutes

	if err := c.AttendanceRepository.Update(tx, &attendance); err != nil {
		c.Log.WithError(err).Error("Failed to update attendance")
		return nil, fiber.ErrInternalServerError
	}

	attendanceLog := &entity.AttendanceLog{
		AttendanceID:       attendance.ID,
		Type:               "CHECK_OUT",
		Time:               nowMilli,
		Lat:                request.Lat,
		Lng:                request.Lng,
		LocationDistance:   distance,
		IsLocationVerified: true,
		// IsFaceVerified:     faceResult.IsVerified,
		// FaceConfidence:     faceResult.Confidence,
		IsFaceVerified: false,
		FaceConfidence: 0,
		FaceImageURL:   request.FaceImageUrl,
		DeviceInfo:     request.DeviceInfo,
	}

	if err := c.AttendanceLogRepo.Create(tx, attendanceLog); err != nil {
		c.Log.WithError(err).Error("Failed to create attendance log")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.AttendandeToResponse(&attendance), nil
}

func (c *AttendanceUseCase) BreakIn(ctx context.Context, request *model.CheckInAttendanceRequest) (*model.AttendanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.EmployeeID == "" {
		return nil, fiber.NewError(400, "User tidak bisa melakukan break-in karena bukan karyawan")
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	nowMilli := time.Now().UnixMilli()

	var attendance entity.Attendance
	if err := c.AttendanceRepository.FindByEmployeeIDAndDate(tx, &attendance, request.EmployeeID, nowMilli); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Belum check-in hari ini")
		}
		return nil, fiber.ErrInternalServerError
	}

	if attendance.CheckOutTime != 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Sudah check-out, tidak bisa break-in")
	}

	breakInCount, err := c.AttendanceLogRepo.CountByAttendanceIDAndType(tx, attendance.ID, "BREAK_IN")
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	breakOutCount, err := c.AttendanceLogRepo.CountByAttendanceIDAndType(tx, attendance.ID, "BREAK_OUT")
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	if breakInCount > breakOutCount {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Sedang dalam break")
	}

	distance, err := c.validateLocation(tx, request.EmployeeID, request.Lat, request.Lng)
	if err != nil {
		return nil, err
	}

	attendanceLog := &entity.AttendanceLog{
		AttendanceID:       attendance.ID,
		Type:               "BREAK_IN",
		Time:               nowMilli,
		Lat:                request.Lat,
		Lng:                request.Lng,
		LocationDistance:   distance,
		IsLocationVerified: true,
		IsFaceVerified:     false,
		FaceConfidence:     0,
		FaceImageURL:       request.FaceImageUrl,
		DeviceInfo:         request.DeviceInfo,
	}

	if err := c.AttendanceLogRepo.Create(tx, attendanceLog); err != nil {
		c.Log.WithError(err).Error("Failed to create attendance log")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.AttendandeToResponse(&attendance), nil
}

func (c *AttendanceUseCase) BreakOut(ctx context.Context, request *model.CheckInAttendanceRequest) (*model.AttendanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.EmployeeID == "" {
		return nil, fiber.NewError(400, "User tidak bisa melakukan break-out karena bukan karyawan")
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	now := time.Now()
	nowMilli := now.UnixMilli()

	var attendance entity.Attendance
	if err := c.AttendanceRepository.FindByEmployeeIDAndDate(tx, &attendance, request.EmployeeID, nowMilli); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Belum check-in hari ini")
		}
		return nil, fiber.ErrInternalServerError
	}

	if attendance.CheckOutTime != 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Sudah check-out, tidak bisa break-out")
	}

	breakInCount, err := c.AttendanceLogRepo.CountByAttendanceIDAndType(tx, attendance.ID, "BREAK_IN")
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	breakOutCount, err := c.AttendanceLogRepo.CountByAttendanceIDAndType(tx, attendance.ID, "BREAK_OUT")
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	if breakInCount == breakOutCount {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Belum break-in")
	}

	var lastBreakIn entity.AttendanceLog
	if err := c.AttendanceLogRepo.FindLastByAttendanceIDAndType(tx, &lastBreakIn, attendance.ID, "BREAK_IN"); err != nil {
		return nil, fiber.ErrInternalServerError
	}

	distance, err := c.validateLocation(tx, request.EmployeeID, request.Lat, request.Lng)
	if err != nil {
		return nil, err
	}

	breakDuration := int(now.Sub(time.UnixMilli(lastBreakIn.Time)).Minutes())
	attendance.TotalBreakMinutes += breakDuration

	if err := c.AttendanceRepository.Update(tx, &attendance); err != nil {
		c.Log.WithError(err).Error("Failed to update attendance")
		return nil, fiber.ErrInternalServerError
	}

	attendanceLog := &entity.AttendanceLog{
		AttendanceID:       attendance.ID,
		Type:               "BREAK_OUT",
		Time:               nowMilli,
		Lat:                request.Lat,
		Lng:                request.Lng,
		LocationDistance:   distance,
		IsLocationVerified: true,
		IsFaceVerified:     false,
		FaceConfidence:     0,
		FaceImageURL:       request.FaceImageUrl,
		DeviceInfo:         request.DeviceInfo,
	}

	if err := c.AttendanceLogRepo.Create(tx, attendanceLog); err != nil {
		c.Log.WithError(err).Error("Failed to create attendance log")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.AttendandeToResponse(&attendance), nil
}
