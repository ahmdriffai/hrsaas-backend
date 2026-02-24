package usecase

import (
	"context"
	"hr-sas/internal/entity"
	"hr-sas/internal/lib"

	"hr-sas/internal/model"
	"hr-sas/internal/model/converter"
	"hr-sas/internal/repository"
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
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, sessionRepository *repository.SessionRepository, companyRepository *repository.CompanyRepository) *UserUseCase {
	return &UserUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
		CompanyRepository: companyRepository,
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
		c.Log.Warnf("Failed find user by token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	// Check expiry
	if session.ExpiredAt.Before(time.Now()) {
		if err := c.SessionRepository.Delete(tx, session); err != nil {
			c.Log.WithError(err).Error("Failed to delete session by user id")
			return nil, fiber.ErrInternalServerError
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Session expired")
	}

	// find user
	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, session.UserID, "Employee"); err != nil {
		c.Log.Warnf("Failed find user by token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if user.CompanyID == "" {
		return nil, fiber.NewError(fiber.StatusForbidden, "User not associated with any company")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
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

	// create user
	user := &entity.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  string(passwordHash),
		Role:      "USER",
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

	return converter.UserToResponse(user), nil
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
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by email : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "email or password not match")
	}

	// find session by user id
	session := new(entity.Session)
	totalSession, err := c.SessionRepository.CountByUserId(tx, user.ID)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	if totalSession > 3 {
		return nil, fiber.NewError(fiber.StatusConflict, "User already logged in")
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed compare password : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "email or password not match")
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
		ExpiredAt: time.Now().Add(24 * time.Hour),
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
		User:  *converter.UserToResponse(user),
		Token: token,
	}, nil
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
