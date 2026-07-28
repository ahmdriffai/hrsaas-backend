package usecase

import (
	"context"
	"errors"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/pkg"
	"hr-sas/internal/repository"
	"io"
	"mime/multipart"
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
	EmployeeRepository   *repository.EmployeeRepository
	UserRepository       *repository.UserRepository
	UploadUseCase        *UploadUseCase
	S3Client             *pkg.S3Client
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
	employeeRepository *repository.EmployeeRepository,
	userRepository *repository.UserRepository,
	uploadUseCase *UploadUseCase,
	s3Client *pkg.S3Client,
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
		EmployeeRepository:   employeeRepository,
		UserRepository:       userRepository,
		UploadUseCase:        uploadUseCase,
		S3Client:             s3Client,
		FaceServiceURL:       faceServiceURL,
	}
}

func (c *AttendanceUseCase) RegisterFace(ctx context.Context, request *model.RegisterFaceRequest) (*model.RegisterFaceResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	employee := new(entity.Employee)
	if err := c.EmployeeRepository.FindByIdAndCompany(c.DB.WithContext(ctx), employee, request.EmployeeID, request.CompanyID); err != nil {
		c.Log.WithError(err).Error("Failed to find employee by ID")
		return nil, fiber.ErrNotFound
	}

	// get image

	// s3 usage
	image, err := c.S3Client.GetObjectBytes(request.ObjectKey, true)
	if err != nil {
		return nil, err
	}

	//
	// image, err := readMultipart(request.File)
	// if err != nil {
	// 	c.Log.WithError(err).Error("Failed to read face image")
	// 	return nil, fiber.ErrBadRequest
	// }

	// image, err = resizeImage(image, 1080) // compress dulu
	// if err != nil {
	// 	c.Log.WithError(err).Error("Failed to resize image")
	// 	return nil, fiber.ErrBadRequest
	// }

	if err := lib.RegisterFace(c.FaceServiceURL+"/register", employee.ID, request.ObjectKey, image); err != nil {
		c.Log.WithError(err).Error("Failed to register face")
		return nil, fiber.NewError(fiber.StatusBadGateway, "Gagal mendaftarkan wajah")
	}

	if err != nil {
		c.Log.WithError(err).Error("Failed to upload face image")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal mengunggah wajah")
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(c.DB.WithContext(ctx), user, employee.UserID); err != nil {
		c.Log.WithError(err).Error("Failed to find user")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal menemukan data user")
	}

	imageUrl := c.S3Client.GetPublicURL(request.ObjectKey)
	user.Image = &imageUrl

	if err := c.UserRepository.Update(c.DB.WithContext(ctx), user); err != nil {
		c.Log.WithError(err).Error("Failed to update user image")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal memperbarui data user")
	}

	return &model.RegisterFaceResponse{
		EmployeeID:   employee.ID,
		FaceImageURL: imageUrl,
	}, nil
}

func (c *AttendanceUseCase) FaceStatus(ctx context.Context, id, companyID string) (*model.FaceStatusResponse, error) {
	employee := new(entity.Employee)
	if err := c.EmployeeRepository.FindByIdAndCompany(c.DB.WithContext(ctx), employee, id, companyID); err != nil {
		c.Log.WithError(err).Error("Failed to find employee by ID")
		return nil, fiber.ErrNotFound
	}

	result := lib.CheckFaceExistence(c.FaceServiceURL+"/check-exists", employee.ID)
	return &model.FaceStatusResponse{
		EmployeeID: employee.ID,
		Registered: result.Registered,
	}, nil
}

func (c *AttendanceUseCase) DeleteFace(ctx context.Context, id, companyID string) error {
	employee := new(entity.Employee)
	if err := c.EmployeeRepository.FindByIdAndCompany(c.DB.WithContext(ctx), employee, id, companyID); err != nil {
		c.Log.WithError(err).Error("Failed to find employee by ID")
		return fiber.ErrNotFound
	}

	if err := lib.DeleteFace(c.FaceServiceURL+"/delete", employee.ID); err != nil {
		c.Log.WithError(err).Error("Failed to delete face")
		return fiber.NewError(fiber.StatusBadGateway, "Gagal menghapus wajah")
	}

	return nil
}

func (c *AttendanceUseCase) uploadFace(ctx context.Context, file *multipart.FileHeader) (string, error) {
	uploaded, err := c.UploadUseCase.Upload(ctx, &model.UploadRequest{File: file})
	if err != nil {
		c.Log.WithError(err).Error("Failed to upload face image")
		return "", err
	}
	return uploaded.Url, nil
}

func readMultipart(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	return io.ReadAll(src)
}

func (c *AttendanceUseCase) verifyAndStoreFace(ctx context.Context, employeeID string, file *multipart.FileHeader) (string, *lib.FaceRecognizeResponse, error) {
	image, err := readMultipart(file)
	if err != nil {
		c.Log.WithError(err).Error("Failed to read face image")
		return "", nil, fiber.ErrBadRequest
	}

	result, err := lib.RecognizeFace(c.FaceServiceURL+"/recognize", employeeID, file.Filename, image)
	if err != nil {
		c.Log.WithError(err).Error("Failed to recognize face")
		return "", nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal memverifikasi wajah !")
	}

	if !result.Match {
		return "", result, nil
	}

	uploadURL, err := c.uploadFace(ctx, file)
	if err != nil {
		c.Log.WithError(err).Error("Failed to upload face image")
		return "", nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan wajah !")
	}

	return uploadURL, result, nil
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

func (c *AttendanceUseCase) SearchLog(ctx context.Context, request *model.SearchAttendanceLogRequest) ([]model.AttendanceLogResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	logs, total, err := c.AttendanceLogRepo.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting attendance logs")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.AttendanceLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = *model.AttendanceLogToResponse(&log)
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
		locationDistance = distance
		if distance <= float64(location.Radius) {
			isInRange = true
		} else {
			isInRange = false
			break
		}
	}

	isApproved := isInRange

	attendance := &entity.Attendance{
		CompanyID:   request.CompanyID,
		EmployeeID:  request.EmployeeID,
		Date:        startOfDay,
		CheckInTime: nowMilli,
		Status:      "HADIR",
	}

	faceImageURL, faceResult, err := c.verifyAndStoreFace(ctx, request.EmployeeID, request.File)
	if err != nil {
		return nil, err
	}
	if !faceResult.Match {
		return nil, fiber.NewError(fiber.StatusBadRequest, faceResult.Message)
	}

	attendanceLog := &entity.AttendanceLog{
		Type:               "CHECK_IN",
		Time:               nowMilli,
		Lat:                request.Lat,
		Lng:                request.Lng,
		LocationDistance:   locationDistance,
		IsLocationVerified: isInRange,
		IsFaceVerified:     faceResult.Match,
		// FaceConfidence:     0,
		FaceImageURL: faceImageURL,
		DeviceInfo:   request.DeviceInfo,
		IsApproved:   isApproved,
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

	// if !request.IsAllowed {
	// 	return nil, fiber.NewError(400, "Anda tidak diizinkan melakukan check-out disini")
	// }

	isInRange := false
	locationDistance := 0.0
	locations, err := c.LocationRepository.GetByEmployeeID(tx, request.EmployeeID)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

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
		locationDistance = distance

		if distance <= float64(location.Radius) {
			isInRange = true
		} else {
			isInRange = false
		}
	}

	isApproved := isInRange

	faceImageURL, faceResult, err := c.verifyAndStoreFace(ctx, request.EmployeeID, request.File)
	if err != nil {
		return nil, err
	}
	if !faceResult.Match {
		return nil, fiber.NewError(fiber.StatusBadRequest, faceResult.Message)
	}

	attendance.CheckOutTime = nowMilli

	checkInTime := time.UnixMilli(attendance.CheckInTime)
	totalWorkMinutes := max(int(now.Sub(checkInTime).Minutes())-attendance.TotalBreakMinutes, 0)
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
		LocationDistance:   locationDistance,
		IsLocationVerified: isInRange,
		IsFaceVerified:     faceResult.Match,
		FaceConfidence:     0, // Python /recognize belum mengembalikan confidence
		FaceImageURL:       faceImageURL,
		DeviceInfo:         request.DeviceInfo,
		IsApproved:         isApproved,
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

	// if !request.IsAllowed {
	// 	return nil, fiber.NewError(400, "Anda tidak diizinkan melakukan break-in disini")
	// }

	isInRange := false
	locationDistance := 0.0
	locations, err := c.LocationRepository.GetByEmployeeID(tx, request.EmployeeID)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

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

	isApproved := isInRange

	faceImageURL, err := c.uploadFace(ctx, request.File)
	if err != nil {
		return nil, err
	}

	attendanceLog := &entity.AttendanceLog{
		AttendanceID:       attendance.ID,
		Type:               "BREAK_IN",
		Time:               nowMilli,
		Lat:                request.Lat,
		Lng:                request.Lng,
		LocationDistance:   locationDistance,
		IsLocationVerified: isInRange,
		IsFaceVerified:     false,
		FaceConfidence:     0,
		FaceImageURL:       faceImageURL,
		DeviceInfo:         request.DeviceInfo,
		IsApproved:         isApproved,
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

	// if !request.IsAllowed {
	// 	return nil, fiber.NewError(400, "Anda tidak diizinkan melakukan break-out disini")
	// }

	isInRange := false
	locationDistance := 0.0
	locations, err := c.LocationRepository.GetByEmployeeID(tx, request.EmployeeID)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

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

	isApproved := isInRange

	breakDuration := int(now.Sub(time.UnixMilli(lastBreakIn.Time)).Minutes())
	attendance.TotalBreakMinutes += breakDuration

	if err := c.AttendanceRepository.Update(tx, &attendance); err != nil {
		c.Log.WithError(err).Error("Failed to update attendance")
		return nil, fiber.ErrInternalServerError
	}

	faceImageURL, err := c.uploadFace(ctx, request.File)
	if err != nil {
		return nil, err
	}

	attendanceLog := &entity.AttendanceLog{
		AttendanceID:       attendance.ID,
		Type:               "BREAK_OUT",
		Time:               nowMilli,
		Lat:                request.Lat,
		Lng:                request.Lng,
		LocationDistance:   locationDistance,
		IsLocationVerified: isInRange,
		IsFaceVerified:     false,
		FaceConfidence:     0,
		FaceImageURL:       faceImageURL,
		DeviceInfo:         request.DeviceInfo,
		IsApproved:         isApproved,
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
