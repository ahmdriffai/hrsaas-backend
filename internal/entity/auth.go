package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            string `gorm:"column:id;primaryKey"`
	Name          string `gorm:"column:name;not null"`
	Email         string `gorm:"column:email;uniqueIndex;not null"`
	EmailVerified bool   `gorm:"column:email_verified;default:false"`
	Image         *string
	CompanyID     string `gorm:"column:company_id"`
	Role          string `gorm:"column:role;not null"`
	Password      string `gorm:"column:password;not null"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Sessions  []Session `gorm:"constraint:OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return nil
}

type Session struct {
	ID        string    `gorm:"column:id;primaryKey"`
	ExpiredAt time.Time `gorm:"column:expired_at"`
	Token     string    `gorm:"column:token;uniqueIndex"`
	IPAddress *string   `gorm:"column:ip_address"`
	UserAgent *string   `gorm:"column:user_agent"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Session) TableName() string {
	return "sessions"
}

func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.NewString()
	return nil
}
