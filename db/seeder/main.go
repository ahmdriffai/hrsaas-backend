package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ── DB config ─────────────────────────────────────────────────────────────────
const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "hr_saas"
	dbSSL      = "disable"
	dbTZ       = "Asia/Jakarta"
)

// ── Minimal entity structs (mirror entity package) ────────────────────────────
type Company struct {
	ID        string `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (Company) TableName() string { return "companies" }

type Division struct {
	ID        string `gorm:"column:id;primaryKey"`
	CompanyID string `gorm:"column:company_id"`
	Name      string `gorm:"column:name"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (Division) TableName() string { return "divisions" }

type Position struct {
	ID         string  `gorm:"column:id;primaryKey"`
	CompanyID  string  `gorm:"column:company_id"`
	Name       string  `gorm:"column:name"`
	ParentID   *string `gorm:"column:parent_id"`
	IsApprover bool    `gorm:"column:is_approver"`
	CreatedAt  int64   `gorm:"column:created_at"`
	UpdatedAt  int64   `gorm:"column:updated_at"`
}

func (Position) TableName() string { return "positions" }

type Role struct {
	ID        string `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name"`
	CompanyID string `gorm:"column:company_id"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (Role) TableName() string { return "roles" }

type User struct {
	ID            string  `gorm:"column:id;primaryKey"`
	Name          string  `gorm:"column:name"`
	Email         string  `gorm:"column:email"`
	Password      string  `gorm:"column:password"`
	EmailVerified bool    `gorm:"column:email_verified"`
	Image         *string `gorm:"column:image"`
	CompanyID     string  `gorm:"column:company_id"`
	CreatedAt     int64   `gorm:"column:created_at"`
	UpdatedAt     int64   `gorm:"column:updated_at"`
}

func (User) TableName() string { return "users" }

type UserRole struct {
	UserID string `gorm:"column:user_id;primaryKey"`
	RoleID string `gorm:"column:role_id;primaryKey"`
}

func (UserRole) TableName() string { return "user_roles" }

type Permission struct {
	ID        string `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (Permission) TableName() string { return "permissions" }

type RolePermission struct {
	RoleID       string `gorm:"column:role_id;primaryKey"`
	PermissionID string `gorm:"column:permission_id;primaryKey"`
}

func (RolePermission) TableName() string { return "role_permissions" }

type Employee struct {
	ID             string `gorm:"column:id;primaryKey"`
	CompanyID      string `gorm:"column:company_id"`
	UserID         string `gorm:"column:user_id"`
	EmployeeNumber string `gorm:"column:employee_number"`
	Fullname       string `gorm:"column:fullname"`
	Gender         string `gorm:"column:gender"`
	BirthPlace     string `gorm:"column:birth_place"`
	BirthDate      int64  `gorm:"column:birth_date"`
	BlodType       string `gorm:"column:blood_type"`
	MaritalStatus  string `gorm:"column:marital_status"`
	Religion       string `gorm:"column:religion"`
	Phone          string `gorm:"column:phone"`
	Timezone       string `gorm:"column:timezone"`
	CreatedAt      int64  `gorm:"column:created_at"`
	UpdatedAt      int64  `gorm:"column:updated_at"`
}

func (Employee) TableName() string { return "employees" }

type EmployeeContract struct {
	ID           string  `gorm:"column:id;primaryKey"`
	EmployeeID   string  `gorm:"column:employee_id"`
	ContractType string  `gorm:"column:contract_type"`
	StartDate    int64   `gorm:"column:start_date"`
	EndDate      *int64  `gorm:"column:end_date"`
	DivisionID   string  `gorm:"column:division_id"`
	PositionID   string  `gorm:"column:position_id"`
	Salary       float64 `gorm:"column:salary"`
	IsActive     bool    `gorm:"column:is_active"`
}

func (EmployeeContract) TableName() string { return "employee_contracts" }

// ── Helpers ───────────────────────────────────────────────────────────────────
func now() int64 { return time.Now().UnixMilli() }

func newID() string { return uuid.NewString() }

func hashPassword(plain string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt error: %v", err)
	}
	return string(b)
}

func ptr[T any](v T) *T { return &v }

// ── Seed data ─────────────────────────────────────────────────────────────────
type seedEmployee struct {
	name     string
	email    string
	empNum   string
	gender   string
	position string // key into positionMap
	salary   float64
	approver bool // whether this position is an approver
}

var employees = []seedEmployee{
	{
		name: "Budi Santoso", email: "direktur.utama@company.com",
		empNum: "EMP-001", gender: "Laki-laki",
		position: "Direktur Utama", salary: 25_000_000, approver: true,
	},
	{
		name: "Siti Rahayu", email: "direktur.operasional@company.com",
		empNum: "EMP-002", gender: "Perempuan",
		position: "Direktur Operasional", salary: 20_000_000, approver: true,
	},
	{
		name: "Andi Wijaya", email: "kadiv.operasional@company.com",
		empNum: "EMP-003", gender: "Laki-laki",
		position: "Kadiv Operasional", salary: 15_000_000, approver: true,
	},
	{
		name: "Dewi Kusuma", email: "kabag.operasional@company.com",
		empNum: "EMP-004", gender: "Perempuan",
		position: "Kabag Operasional", salary: 10_000_000, approver: false,
	},
	{
		name: "Reza Firmansyah", email: "staff.operasional@company.com",
		empNum: "EMP-005", gender: "Laki-laki",
		position: "Staff Operasional", salary: 6_000_000, approver: false,
	},
}

// ── Main ──────────────────────────────────────────────────────────────────────
func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSL, dbTZ,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	// ── 1. Company ─────────────────────────────────────────────────────────────
	var company Company
	if err := db.Where("name = ?", "PT Contoh Perusahaan").First(&company).Error; err != nil {
		company = Company{
			ID: newID(), Name: "PT Contoh Perusahaan",
			CreatedAt: now(), UpdatedAt: now(),
		}
		db.Create(&company)
		fmt.Printf("Created company: %s\n", company.Name)
	} else {
		fmt.Printf("Using existing company: %s (%s)\n", company.Name, company.ID)
	}

	// ── 2. Division: Operasional ───────────────────────────────────────────────
	var division Division
	if err := db.Where("name = ? AND company_id = ?", "Operasional", company.ID).First(&division).Error; err != nil {
		division = Division{
			ID: newID(), CompanyID: company.ID, Name: "Operasional",
			CreatedAt: now(), UpdatedAt: now(),
		}
		db.Create(&division)
		fmt.Printf("Created division: %s\n", division.Name)
	} else {
		fmt.Printf("Using existing division: %s (%s)\n", division.Name, division.ID)
	}

	// ── 3. Positions (hierarchical) ────────────────────────────────────────────
	// Order matters: parent must exist before child
	type posSpec struct {
		name      string
		parentKey string // empty = root
		approver  bool
	}
	posSpecs := []posSpec{
		{"Direktur Utama", "", true},
		{"Direktur Operasional", "Direktur Utama", true},
		{"Kadiv Operasional", "Direktur Operasional", true},
		{"Kabag Operasional", "Kadiv Operasional", false},
		{"Staff Operasional", "Kabag Operasional", false},
	}

	positionMap := map[string]string{} // name -> id

	for _, spec := range posSpecs {
		var pos Position
		if err := db.Where("name = ? AND company_id = ?", spec.name, company.ID).First(&pos).Error; err != nil {
			var parentID *string
			if spec.parentKey != "" {
				parentID = ptr(positionMap[spec.parentKey])
			}
			pos = Position{
				ID: newID(), CompanyID: company.ID, Name: spec.name,
				ParentID: parentID, IsApprover: spec.approver,
				CreatedAt: now(), UpdatedAt: now(),
			}
			db.Create(&pos)
			fmt.Printf("Created position: %s\n", pos.Name)
		} else {
			fmt.Printf("Using existing position: %s (%s)\n", pos.Name, pos.ID)
		}
		positionMap[spec.name] = pos.ID
	}

	// ── 4. Roles: Admin & Staff ───────────────────────────────────────────────
	ensureRole := func(name string) Role {
		var r Role
		if err := db.Where("name = ? AND company_id = ?", name, company.ID).First(&r).Error; err != nil {
			r = Role{ID: newID(), Name: name, CompanyID: company.ID, CreatedAt: now(), UpdatedAt: now()}
			db.Create(&r)
			fmt.Printf("Created role: %s\n", r.Name)
		} else {
			fmt.Printf("Using existing role: %s (%s)\n", r.Name, r.ID)
		}
		return r
	}
	adminRole := ensureRole("ADMIN")
	staffRole := ensureRole("STAFF")

	// ── 5. Permissions → Admin role ────────────────────────────────────────────
	permissionNames := []string{
		"USERS", "EMPLOYEES", "EMPLOYEE_CONTRACTS", "DIVISIONS",
		"SANCTIONS", "EMPLOYEE_SANCTIONS", "POSITIONS", "OFFICE_LOCATIONS",
		"SHIFTS", "TIME_OFF_REQUESTS", "TIME_OFF_TYPES", "TIME_OFF_BALANCES",
		"PERMISSIONS", "ROLES", "EMPLOYEE_DOCUMENTS", "VISITS", "HOLIDAYS",
	}

	for _, pname := range permissionNames {
		var perm Permission
		if err := db.Where("name = ?", pname).First(&perm).Error; err != nil {
			perm = Permission{ID: newID(), Name: pname, CreatedAt: now(), UpdatedAt: now()}
			db.Create(&perm)
			fmt.Printf("Created permission: %s\n", perm.Name)
		}

		var rp RolePermission
		if err := db.Where("role_id = ? AND permission_id = ?", adminRole.ID, perm.ID).First(&rp).Error; err != nil {
			db.Create(&RolePermission{RoleID: adminRole.ID, PermissionID: perm.ID})
		}
	}
	fmt.Printf("Assigned %d permissions to role Admin\n", len(permissionNames))

	// ── 6. Users, Employees, Contracts ────────────────────────────────────────
	defaultPassword := hashPassword("Password123!")
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	contractStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	for i, emp := range employees {
		// User
		var user User
		if err := db.Where("email = ?", emp.email).First(&user).Error; err != nil {
			user = User{
				ID: newID(), Name: emp.name, Email: emp.email,
				Password: defaultPassword, EmailVerified: true,
				CompanyID: company.ID, CreatedAt: now(), UpdatedAt: now(),
			}
			db.Create(&user)
			fmt.Printf("[%d] Created user: %s\n", i+1, user.Email)
		} else {
			fmt.Printf("[%d] Using existing user: %s (%s)\n", i+1, user.Email, user.ID)
		}

		// User ↔ Role
		var ur UserRole
		if err := db.Where("user_id = ? AND role_id = ?", user.ID, staffRole.ID).First(&ur).Error; err != nil {
			db.Create(&UserRole{UserID: user.ID, RoleID: staffRole.ID})
		}

		// Employee
		var employee Employee
		if err := db.Where("user_id = ?", user.ID).First(&employee).Error; err != nil {
			employee = Employee{
				ID: newID(), CompanyID: company.ID, UserID: user.ID,
				EmployeeNumber: emp.empNum, Fullname: emp.name,
				Gender: emp.gender, BirthPlace: "Jakarta",
				BirthDate: birthDate, BlodType: "O",
				MaritalStatus: "single", Religion: "Islam",
				Phone:     fmt.Sprintf("08100000000%d", i+1),
				Timezone:  "Asia/Jakarta",
				CreatedAt: now(), UpdatedAt: now(),
			}
			db.Create(&employee)
			fmt.Printf("[%d] Created employee: %s (%s)\n", i+1, employee.Fullname, employee.EmployeeNumber)
		} else {
			fmt.Printf("[%d] Using existing employee: %s (%s)\n", i+1, employee.Fullname, employee.EmployeeNumber)
		}

		// Contract
		var contract EmployeeContract
		if err := db.Where("employee_id = ? AND position_id = ?", employee.ID, positionMap[emp.position]).First(&contract).Error; err != nil {
			contract = EmployeeContract{
				ID: newID(), EmployeeID: employee.ID,
				ContractType: "PKWTT",
				StartDate:    contractStart,
				EndDate:      nil,
				DivisionID:   division.ID,
				PositionID:   positionMap[emp.position],
				Salary:       emp.salary,
				IsActive:     true,
			}
			db.Create(&contract)
			fmt.Printf("[%d] Created contract: %s – %s (Rp %.0f)\n", i+1, emp.name, emp.position, emp.salary)
		} else {
			fmt.Printf("[%d] Using existing contract for: %s\n", i+1, emp.name)
		}
	}

	fmt.Println("\nSeeding selesai!")
	fmt.Println("Default password semua karyawan: Password123!")
}
