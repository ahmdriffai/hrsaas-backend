package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"

	"hr-sas/internal/model"
	"hr-sas/internal/repository"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	UserRepository    *repository.UserRepository
	SessionRepository *repository.SessionRepository
	CompanyRepository *repository.CompanyRepository
	RoleRepository    *repository.RoleRepository
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, sessionRepository *repository.SessionRepository, companyRepository *repository.CompanyRepository, roleRepository *repository.RoleRepository) *UserUseCase {
	return &UserUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
		CompanyRepository: companyRepository,
		RoleRepository:    roleRepository,
	}
}

/*
Verify User
*/
func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// find session
	session := new(entity.Session)
	if err := c.SessionRepository.FindByToken(tx, session, request.Token); err != nil {
		c.Log.Warnf("Gagal menemukan user by token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	expiredAt := time.Unix(session.CreatedAt, 0)

	// Check expiry
	if expiredAt.Before(time.Now()) {
		if err := c.SessionRepository.Delete(tx, session); err != nil {
			c.Log.WithError(err).Error("Gagal menghapus session by user id")
			return nil, fiber.ErrInternalServerError
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Session expired")
	}

	// find user
	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, session.UserID, "Roles", "Employee"); err != nil {
		c.Log.Warnf("Gagal menemukan user by token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if user.CompanyID == "" {
		return nil, fiber.NewError(fiber.StatusForbidden, "User tidak terasosiasi dengan perusahaan manapun")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Gagal menyelesaikan transaksi : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return model.UserToResponse(user), nil
}

/*
Register User
*/
func (c *UserUseCase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// cek user exist
	count, err := c.UserRepository.CountByEmail(tx, request.Email)
	if err != nil {
		c.Log.WithError(err).Error("Failed to count user by email")
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 {
		return nil, fiber.NewError(fiber.StatusConflict, "Email already registered")
	}

	// create company
	company := &entity.Company{
		Name: request.CompanyName,
	}

	if err := c.CompanyRepository.Create(tx, company); err != nil {
		c.Log.WithError(err).Error("Failed to create company")
		return nil, fiber.ErrInternalServerError
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.WithError(err).Error("Failed to hash password")
		return nil, fiber.ErrInternalServerError
	}

	adminRole := &entity.Role{Name: "ADMIN"}
	if err := c.RoleRepository.Create(tx, adminRole); err != nil {
		c.Log.WithError(err).Error("Failed to create admin role")
		return nil, fiber.ErrInternalServerError
	}

	// create user
	user := &entity.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  string(passwordHash),
		Roles:     []entity.Role{*adminRole},
		CompanyID: company.ID,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.WithError(err).Error("Failed to create user")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return model.UserToResponse(user), nil
}

/*
Login User
*/
func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.LoginUserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// find user by email
	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email, "Roles"); err != nil {
		c.Log.Warnf("Gagal menemukan user by email : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "email dan password tidak valid")
	}

	// find session by user id
	session := new(entity.Session)
	totalSession, err := c.SessionRepository.CountByUserId(tx, user.ID)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	if totalSession > 10000 {
		return nil, fiber.NewError(fiber.StatusConflict, "User sudah login di perangkat lain")
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Password tidak valid : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "email dan password tidak valid")
	}

	// create token
	token, err := lib.GenerateToken(32)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	// create session
	session = &entity.Session{
		UserID:    user.ID,
		Token:     token,
		IPAddress: &request.Ip,
		UserAgent: &request.UserAgent,
		ExpiredAt: time.Now().Add(24 * time.Hour).UnixMilli(),
	}

	if err := c.SessionRepository.Create(tx, session); err != nil {
		c.Log.WithError(err).Error("Failed to create session")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.LoginUserResponse{
		User:  *model.UserToResponse(user),
		Token: token,
	}, nil
}

func (c *UserUseCase) List(ctx context.Context, request *model.SearchUserRequest) ([]model.UserResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Gagal memvalidasi query pencarian")
		return nil, 0, fiber.ErrBadRequest
	}

	users, total, err := c.UserRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("Gagal mencari user")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.UserResponse, len(users))
	for i := range users {
		responses[i] = *model.UserToResponse(&users[i])
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Gagal menyelesaikan transaksi")
		return nil, 0, fiber.ErrInternalServerError
	}

	return responses, total, nil
}

func (c *UserUseCase) Detail(ctx context.Context, id string) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, id, "Roles", "Employee"); err != nil {
		c.Log.WithError(err).Error("User tidak ditemukan")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Gagal menyelesaikan transaksi")
		return nil, fiber.ErrInternalServerError
	}

	return model.UserToResponse(user), nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID, "Roles", "Employee"); err != nil {
		c.Log.WithError(err).Error("User tidak ditemukan")
		return nil, fiber.ErrNotFound
	}

	if request.Name != nil {
		name := strings.TrimSpace(*request.Name)
		if name == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Nama tidak boleh kosong")
		}
		user.Name = name
	}

	if request.Email != nil {
		email := strings.TrimSpace(*request.Email)
		if email == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Email tidak boleh kosong")
		}

		count, err := c.UserRepository.CountByEmailExcludeID(tx, email, user.ID)
		if err != nil {
			c.Log.WithError(err).Error("Gagal menghitung user by email")
			return nil, fiber.ErrInternalServerError
		}
		if count > 0 {
			return nil, fiber.NewError(fiber.StatusConflict, "Email sudah terdaftar")
		}

		user.Email = email
	}

	if request.Image != nil {
		user.Image = request.Image
	}

	if request.CompanyID != nil {
		user.CompanyID = *request.CompanyID
	}

	if request.EmailVerified != nil {
		user.EmailVerified = *request.EmailVerified
	}

	user.UpdatedAt = time.Now().UnixMilli()

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.WithError(err).Error("Failed to update user")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Gagal menyelesaikan transaksi")
		return nil, fiber.ErrInternalServerError
	}

	return model.UserToResponse(user), nil
}

func (c *UserUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, id); err != nil {
		c.Log.WithError(err).Error("User tidak ditemukan")
		return fiber.ErrNotFound
	}

	if err := c.UserRepository.Delete(tx, user); err != nil {
		c.Log.WithError(err).Error("Gagal menghapus user")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Gagal menyelesaikan transaksi")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *UserUseCase) ChangePassword(ctx context.Context, userID string, request *model.ChangePasswordRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, userID); err != nil {
		c.Log.WithError(err).Error("User tidak ditemukan")
		return fiber.ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.CurrentPassword)); err != nil {
		c.Log.WithError(err).Warn("Password saat ini tidak cocok")
		return fiber.NewError(fiber.StatusBadRequest, "Password saat ini salah")
	}

	if request.CurrentPassword == request.NewPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Password baru harus berbeda dari password saat ini")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Log.WithError(err).Error("Failed to hash new password")
		return fiber.ErrInternalServerError
	}

	user.Password = string(passwordHash)
	user.UpdatedAt = time.Now().UnixMilli()

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.WithError(err).Error("Gagal memperbarui password user")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Gagal menyelesaikan transaksi")
		return fiber.ErrInternalServerError
	}

	return nil
}

/*
Assign Roles to User
*/
func (c *UserUseCase) AssignRoles(ctx context.Context, request *model.AssignRoleRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.UserID, "Roles"); err != nil {
		c.Log.WithError(err).Error("User not found")
		return nil, fiber.ErrNotFound
	}

	roles, err := c.RoleRepository.FindByIDs(tx, request.Roles)
	if err != nil {
		c.Log.WithError(err).Error("Failed to find roles")
		return nil, fiber.ErrInternalServerError
	}

	if err := c.UserRepository.AssignRoles(tx, user, roles); err != nil {
		c.Log.WithError(err).Error("Failed to assign roles to user")
		return nil, fiber.ErrInternalServerError
	}

	user.Roles = append(user.Roles, roles...)

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.UserToResponse(user), nil
}

/*
Remove Roles from User
*/
func (c *UserUseCase) RemoveRoles(ctx context.Context, request *model.RemoveRoleRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.UserID, "Roles"); err != nil {
		c.Log.WithError(err).Error("User not found")
		return nil, fiber.ErrNotFound
	}

	roles, err := c.RoleRepository.FindByIDs(tx, request.Roles)
	if err != nil {
		c.Log.WithError(err).Error("Failed to find roles")
		return nil, fiber.ErrInternalServerError
	}

	if err := c.UserRepository.RemoveRoles(tx, user, roles); err != nil {
		c.Log.WithError(err).Error("Failed to remove roles from user")
		return nil, fiber.ErrInternalServerError
	}

	removeSet := make(map[string]bool, len(roles))
	for _, r := range roles {
		removeSet[r.ID] = true
	}
	remaining := user.Roles[:0]
	for _, r := range user.Roles {
		if !removeSet[r.ID] {
			remaining = append(remaining, r)
		}
	}
	user.Roles = remaining

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return model.UserToResponse(user), nil
}

/*
Logout User
*/
func (c *UserUseCase) Logout(ctx context.Context, userId string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// delete session by user id
	if err := c.SessionRepository.DeleteByUserId(tx, userId); err != nil {
		c.Log.WithError(err).Error("Failed to delete session by user id")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
