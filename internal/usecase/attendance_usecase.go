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
	AttendanceLogRepo    *repository.AttendanceLogRepository
}

func NewAttendanceUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	attendanceRepository *repository.AttendanceRepository,
	locationRepository *repository.OfficeLocationRepository,
	shiftRepository *repository.ShiftRepository,
	attendanceLogRepo *repository.AttendanceLogRepository,
) *AttendanceUseCase {
	return &AttendanceUseCase{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		AttendanceRepository: attendanceRepository,
		LocationRepository:   locationRepository,
		ShiftRepository:      shiftRepository,
		AttendanceLogRepo:    attendanceLogRepo,
	}
}

// Implement Attendance Use Case Check-In,
func (c *AttendanceUseCase) CheckIn(ctx context.Context, request *model.CheckInAttendanceRequest) (*model.AttendanceResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// cek user adalah employee
	if request.EmployeeID == "" {
		c.Log.Error("Bad Request")
		return nil, fiber.NewError(400, "User tidak bisa melakukan check-in karena bukan karyawan")
	}

	// validate request body
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	// init time now
	now := time.Now()

	fmt.Println(request.EmployeeID)
	// TODO: Ambil shift karyawan (misalnya dari employee.shift_id → shifts)
	// → Dapatkan start_time, late_tolerance
	shifts, err := c.ShiftRepository.FindByEmployeeID(tx, request.EmployeeID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Shift tidak ditemukan")
	}
	if len(shifts) == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Shift tidak ditemukan")
	}
	shift := shifts[0]

	// TODO: Cek apakah sudah ada attendance hari ini?
	// → Query attendance WHERE employee_id & date = today
	var existingAttendance entity.Attendance
	err = c.AttendanceRepository.FindByEmployeeIDAndDate(tx, &existingAttendance, request.EmployeeID, now)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fiber.ErrInternalServerError
	}
	// → Jika ada → return error "Sudah check-in hari ini
	if existingAttendance.ID != "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Sudah check-in hari ini")
	}

	// TODO: Validasi lokasi
	// - Ambil office_locations yang terdaftar oleh employee
	locations, err := c.LocationRepository.GetByEmployeeID(tx, request.EmployeeID)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	fmt.Println(locations)

	// - Hitung jarak (lat,lng user ke kantor)
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
	// Jika lokasi valid, lanjutkan proses
	if !isInRange {
		return nil, fiber.NewError(400, "Anda diluar jangkauan")
	}

	// attendance object
	attendance := &entity.Attendance{
		CompanyID:   request.CompanyID,
		EmployeeID:  request.EmployeeID,
		Date:        now,
		CheckInTime: now,
		Status:      "HADIR", // Default status
	}

	// attendance log object
	attendanceLog := &entity.AttendanceLog{
		Type:               "CHECK_IN",
		Time:               now,
		Lat:                request.Lat,
		Lng:                request.Lng,
		LocationDistance:   locationDistance,
		IsLocationVerified: true,
		IsFaceVerified:     false,
		FaceConfidence:     0,
		FaceImageURL:       request.FaceImageUrl,
		DeviceInfo:         request.DeviceInfo,
	}

	// TODO: Validasi wajah
	// - Kirim face_image ke Face Recognition Service
	//  - Dapat confidence score
	//  - is_face_verified = confidence >= threshold (misal 0.75)

	// TODO: Tentukan status
	//  - Jika now() > shift.start_time + late_tolerance → status = TERLAMBAT
	//  - Else → status = HADIR
	if now.After(shift.StartTime.Add(time.Duration(shift.LateTolerance) * time.Minute)) {
		// Set status to LATE
		attendance.Status = "TERLAMBAT"
	}

	// TODO: Simpan ke attendance:
	if err := c.AttendanceRepository.Create(tx, attendance); err != nil {
		c.Log.WithError(err).Error("Failed to create attendance")
		return nil, fiber.ErrInternalServerError
	}

	attendanceLog.AttendanceID = attendance.ID

	// TODO: Simpan ke attendance_logs:
	if err := c.AttendanceLogRepo.Create(tx, attendanceLog); err != nil {
		c.Log.WithError(err).Error("Failed to create attendance log")
		return nil, fiber.ErrInternalServerError
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	// return success
	return nil, nil
}

func (c *AttendanceUseCase) CheckOut(ctx context.Context, request *model.CheckInAttendanceRequest) (*model.AttendanceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.EmployeeID == "" {
		c.Log.Error("Bad Request")
		return nil, fiber.NewError(400, "User tidak bisa melakukan check-out karena bukan karyawan")
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	now := time.Now()

	var attendance entity.Attendance
	err := c.AttendanceRepository.FindByEmployeeIDAndDate(tx, &attendance, request.EmployeeID, now)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Belum check-in hari ini")
		}
		return nil, fiber.ErrInternalServerError
	}

	if !attendance.CheckOutTime.IsZero() {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Sudah check-out hari ini")
	}

	locations, err := c.LocationRepository.GetByEmployeeID(tx, request.EmployeeID)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	isInRange := false
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
			break
		}
	}
	if !isInRange {
		return nil, fiber.NewError(400, "Anda diluar jangkauan")
	}

	attendance.CheckOutTime = now
	totalWorkMinutes := int(now.Sub(attendance.CheckInTime).Minutes()) - attendance.TotalBreakMinutes
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
		Time:               now,
		Lat:                request.Lat,
		Lng:                request.Lng,
		LocationDistance:   0,
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

// TODO: Implement Attendance Use Case Break-In
// TODO: Implement Attendance Use Case Break-Out
