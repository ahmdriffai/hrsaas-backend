package usecase

import (
	"context"
	"fmt"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"
	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type EmployeeUseCase struct {
	DB                         *gorm.DB
	Log                        *logrus.Logger
	Validate                   *validator.Validate
	EmployeeRepository         *repository.EmployeeRepository
	UserRepository             *repository.UserRepository
	EmployeeContractRepository *repository.EmployeeContractRepository
	PositionRepository         *repository.PositionRepository
	DivisionRepository         *repository.DivisionRepository
}

func NewEmployeeUseCase(db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	employeeRepository *repository.EmployeeRepository,
	userRepository *repository.UserRepository,
	employeeContractRepository *repository.EmployeeContractRepository,
	positionRepository *repository.PositionRepository,
	divisionRepository *repository.DivisionRepository,
) *EmployeeUseCase {
	return &EmployeeUseCase{
		DB:                         db,
		Log:                        log,
		Validate:                   validate,
		EmployeeRepository:         employeeRepository,
		UserRepository:             userRepository,
		EmployeeContractRepository: employeeContractRepository,
		PositionRepository:         positionRepository,
		DivisionRepository:         divisionRepository,
	}
}

// Create Employee Usecase
func (c *EmployeeUseCase) Create(ctx context.Context, request *model.CreateEmployeeRequest) (*model.EmployeeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	if err := c.ValidateUniqueEmployeeAndEmail(tx, request.EmployeeNumber, request.CompanyID, request.Email, ""); err != nil {
		return nil, err
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.WithError(err).Error("Failed to hash password")
		return nil, fiber.ErrInternalServerError
	}

	// create user entity
	user := &entity.User{
		Name:      request.Fullname,
		Email:     request.Email,
		Password:  string(passwordHash),
		CompanyID: request.CompanyID,
	}

	// create user in database
	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.WithError(err).Error("Failed to create user")
		return nil, fiber.ErrInternalServerError
	}

	birthDate, _ := lib.ParseDateToUnixMilli(request.BirthDate)

	// create employee entity
	employee := &entity.Employee{
		Fullname:       request.Fullname,
		Gender:         request.Gender,
		BirthPlace:     request.BirthPlace,
		BirthDate:      birthDate,
		IdentityNumber: request.IdentityNumber,
		BlodType:       request.BlodType,
		MaritalStatus:  request.MaritalStatus,
		Religion:       request.Religion,
		Phone:          request.Phone,
		Address:        request.Address,
		City:           request.City,
		Timezone:       request.Timezone,
		CompanyID:      request.CompanyID,
		EmployeeNumber: request.EmployeeNumber,
		UserID:         user.ID,
	}

	if err := c.EmployeeRepository.Create(tx, employee); err != nil {
		c.Log.WithError(err).Error("Failed to create employee")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeToResponse(employee), nil
}

// Search Employee
func (c *EmployeeUseCase) Search(ctx context.Context, request *model.SearchEmployeeRequest) ([]model.EmployeeResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	employees, total, err := c.EmployeeRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting employee")
		return nil, 0, fiber.ErrInternalServerError
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.EmployeeResponse, len(employees))
	for i, employee := range employees {
		responses[i] = *model.EmployeeToResponse(&employee)
	}

	return responses, total, nil
}

// Import Excel Employee
func (c *EmployeeUseCase) ImportExcel(ctx context.Context, request *model.ImportExcelEmployeeRequest) (int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return 0, fiber.ErrBadRequest
	}

	// buka file
	file, err := request.File.Open()
	if err != nil {
		c.Log.WithError(err).Error("Failed to open Excel file")
		return 0, fiber.ErrInternalServerError
	}
	defer file.Close()

	// baca excel file

	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		c.Log.WithError(err).Error("Failed to open Excel file")
		return 0, fiber.ErrInternalServerError
	}

	// ambil sheet pertama
	sheetName := excelFile.GetSheetName(0)

	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		c.Log.WithError(err).Error("Failed to read Excel rows")
		return 0, fiber.ErrInternalServerError
	}

	totalData := int64(len(rows) - 1) // kurangi header

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		if len(row) < 14 {
			c.Log.WithError(fmt.Errorf("row %d has insufficient columns", i+1)).Error("Invalid row format")
			return 0, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Row %d has insufficient columns", i+1))
		}

		if err := c.ValidateUniqueEmployeeAndEmail(tx, row[2], request.CompanyID, row[10], row[0]); err != nil {
			return 0, err
		}

		// hash password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(row[2]), bcrypt.DefaultCost)
		if err != nil {
			c.Log.WithError(err).Error("Failed to hash password")
			return 0, fiber.ErrInternalServerError
		}

		user := &entity.User{
			Name:      row[1],
			Email:     row[10],
			Password:  string(passwordHash),
			CompanyID: request.CompanyID,
		}

		if err := c.UserRepository.Create(tx, user); err != nil {
			c.Log.WithError(err).Errorf("Failed to create user on row %d", i+2)
			return 0, fiber.ErrInternalServerError
		}

		birthDate, _ := lib.ParseDateToUnixMilli(row[4])

		employee := &entity.Employee{
			Fullname:       row[1],
			BirthPlace:     row[3],
			BirthDate:      birthDate,
			BlodType:       row[6],
			MaritalStatus:  row[7],
			Religion:       row[8],
			Phone:          row[9],
			Timezone:       "UTC",
			CompanyID:      request.CompanyID,
			UserID:         user.ID,
			EmployeeNumber: row[2],
		}

		if err := c.EmployeeRepository.Create(tx, employee); err != nil {
			c.Log.WithError(err).Error("Failed to create employee")
			return 0, fiber.ErrInternalServerError
		}

		startDate, _ := lib.ParseDateToUnixMilli(row[14])
		salary, _ := strconv.ParseFloat(row[15], 64)

		position := new(entity.Position)
		err = c.PositionRepository.FindByName(tx, position, row[13])
		if err != nil {
			c.Log.WithError(err).Error("Failed to find position")
			return 0, fiber.ErrNotFound
		}

		division := new(entity.Division)
		err = c.DivisionRepository.FindByName(tx, division, row[12])
		if err != nil {
			c.Log.WithError(err).Error("Failed to find division")
			return 0, fiber.ErrNotFound
		}

		employeeContract := &entity.EmployeeContract{
			EmployeeID:   employee.ID,
			StartDate:    startDate,
			ContractType: row[11],
			Salary:       salary,
			PositionID:   position.ID,
			DivisionID:   division.ID,
		}

		if err := c.EmployeeContractRepository.Create(tx, employeeContract); err != nil {
			c.Log.WithError(err).Error("Failed to create employee contract")
			return 0, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return 0, fiber.ErrInternalServerError
	}

	return totalData, nil
}

// Validate unique employee number and email
func (c *EmployeeUseCase) ValidateUniqueEmployeeAndEmail(
	tx *gorm.DB,
	employeeNumber string,
	companyID string,
	email string,
	row string,
) error {

	// validate employee number
	total, err := c.EmployeeRepository.CountByEmployeeNumberAndCompanyID(tx, employeeNumber, companyID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to count employee by number and company")
		return fiber.ErrInternalServerError
	}

	if total > 0 {
		return fiber.NewError(fiber.StatusConflict, "Employee number already in use on row "+row)
	}

	// validate email
	count, err := c.UserRepository.CountByEmail(tx, email)
	if err != nil {
		c.Log.WithError(err).Error("Failed to count user by email ")
		return fiber.ErrInternalServerError
	}

	if count > 0 {
		return fiber.NewError(fiber.StatusConflict, "Email already registered on row "+row)
	}

	return nil
}

// Detail employee
func (c *EmployeeUseCase) Detail(ctx context.Context, request *model.DetailEmployeeRequest) (*model.EmployeeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	employee := new(entity.Employee)
	err := c.EmployeeRepository.FindByIdAndCompany(
		tx, employee, request.ID, request.CompanyID,
		"EmployeeContract", "EmployeeContract.Position",
		"EmployeeContract.Division", "EmployeeDocuments",
		"User", "User.Roles",
	)
	if err != nil {
		c.Log.WithError(err).Error("Failed to find employee by ID")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeToResponse(employee), nil
}

func (c *EmployeeUseCase) Update(ctx context.Context, companyID string, request *model.UpdateEmployeeRequest) (*model.EmployeeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	employee := new(entity.Employee)
	if err := c.EmployeeRepository.FindByIdAndCompany(tx, employee, request.ID, companyID, "EmployeeContract", "EmployeeContract.Position", "EmployeeContract.Division", "EmployeeDocuments", "User", "User.Roles"); err != nil {
		c.Log.WithError(err).Error("Failed to find employee by ID")
		return nil, fiber.ErrNotFound
	}

	if request.EmployeeNumber != nil {
		employeeNumber := strings.TrimSpace(*request.EmployeeNumber)
		if employeeNumber == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "employee_number cannot be empty")
		}

		var count int64
		if err := tx.Model(&entity.Employee{}).
			Where("employee_number = ?", employeeNumber).
			Where("company_id = ?", companyID).
			Where("id <> ?", employee.ID).
			Count(&count).Error; err != nil {
			c.Log.WithError(err).Error("Failed to count employee by number and company")
			return nil, fiber.ErrInternalServerError
		}
		if count > 0 {
			return nil, fiber.NewError(fiber.StatusConflict, "Employee number already in use")
		}

		employee.EmployeeNumber = employeeNumber
	}

	if request.Fullname != nil {
		fullname := strings.TrimSpace(*request.Fullname)
		if fullname == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Nama Lengkap tidak boleh kosong")
		}
		employee.Fullname = fullname
		employee.User.Name = fullname
	}

	if request.Gender != nil {
		gender := strings.TrimSpace(*request.Gender)
		if gender == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Gender Harus di isi")
		}
		employee.Gender = gender
	}

	if request.BirthPlace != nil {
		birthPlace := strings.TrimSpace(*request.BirthPlace)
		if birthPlace == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Tempat lahir tidak boleh kosong")
		}
		employee.BirthPlace = birthPlace
	}

	if request.BirthDate != nil {
		birthDate := strings.TrimSpace(*request.BirthDate)
		if birthDate == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Tanggal lahir tidak boleh kosong")
		}

		parsedBirthDate, err := lib.ParseDateToUnixMilli(birthDate)
		if err != nil {
			c.Log.WithError(err).Error("Failed to parse birth date")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Format tanggal lahir tidak valid")
		}
		employee.BirthDate = parsedBirthDate
	}

	if request.IdentityNumber != nil {
		identityNumber := strings.TrimSpace(*request.IdentityNumber)
		if identityNumber == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Nomor identitas tidak boleh kosong")
		}
		employee.IdentityNumber = identityNumber
	}

	if request.BlodType != nil {
		employee.BlodType = strings.TrimSpace(*request.BlodType)
	}

	if request.MaritalStatus != nil {
		maritalStatus := strings.TrimSpace(*request.MaritalStatus)
		if maritalStatus == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Status pernikahan tidak boleh kosong")
		}
		employee.MaritalStatus = maritalStatus
	}

	if request.Religion != nil {
		religion := strings.TrimSpace(*request.Religion)
		if religion == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Agama tidak boleh kosong")
		}
		employee.Religion = religion
	}

	if request.Phone != nil {
		phone := strings.TrimSpace(*request.Phone)
		if phone == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Nomor telepon tidak boleh kosong")
		}
		employee.Phone = phone
	}

	if request.Address != nil {
		address := strings.TrimSpace(*request.Address)
		if address == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Alamat tidak boleh kosong")
		}
		employee.Address = address
	}

	if request.City != nil {
		city := strings.TrimSpace(*request.City)
		if city == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Kota tidak boleh kosong")
		}
		employee.City = city
	}

	if request.Timezone != nil {
		timezone := strings.TrimSpace(*request.Timezone)
		if timezone == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Timezone tidak boleh kosong")
		}
		employee.Timezone = timezone
	}

	if request.Email != nil {
		email := strings.TrimSpace(*request.Email)
		if email == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Email tidak boleh kosong")
		}

		count, err := c.UserRepository.CountByEmailExcludeID(tx, email, employee.UserID)
		if err != nil {
			c.Log.WithError(err).Error("Failed to count user by email")
			return nil, fiber.ErrInternalServerError
		}
		if count > 0 {
			return nil, fiber.NewError(fiber.StatusConflict, "Email sudah terdaftar")
		}

		employee.User.Email = email
	}

	nowMilli := time.Now().UnixMilli()
	employee.UpdatedAt = nowMilli
	employee.User.UpdatedAt = nowMilli

	if err := c.UserRepository.Update(tx, &employee.User); err != nil {
		c.Log.WithError(err).Error("Gagal memperbarui user employee")
		return nil, fiber.ErrInternalServerError
	}

	if err := c.EmployeeRepository.Update(tx, employee); err != nil {
		c.Log.WithError(err).Error("Gagal memperbarui employee")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Gagal menyelesaikan transaksi")
		return nil, fiber.ErrInternalServerError
	}

	return model.EmployeeToResponse(employee), nil
}

func (c *EmployeeUseCase) Delete(ctx context.Context, id string, companyID string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	employee := new(entity.Employee)
	if err := c.EmployeeRepository.FindByIdAndCompany(tx, employee, id, companyID, "User"); err != nil {
		c.Log.WithError(err).Error("Failed to find employee by ID")
		return fiber.ErrNotFound
	}

	if err := c.EmployeeRepository.Delete(tx, employee); err != nil {
		c.Log.WithError(err).Error("Failed to delete employee")
		if strings.Contains(strings.ToLower(err.Error()), "foreign key") || strings.Contains(err.Error(), "SQLSTATE 23503") {
			return fiber.NewError(fiber.StatusConflict, "Employee cannot be deleted because it is still used by other records")
		}
		return fiber.ErrInternalServerError
	}

	if err := c.UserRepository.Delete(tx, &employee.User); err != nil {
		c.Log.WithError(err).Error("Failed to delete employee user")
		if strings.Contains(strings.ToLower(err.Error()), "foreign key") || strings.Contains(err.Error(), "SQLSTATE 23503") {
			return fiber.NewError(fiber.StatusConflict, "Employee user cannot be deleted because it is still used by other records")
		}
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}
