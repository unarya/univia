package config

import (
	"fmt"
	"log"
	"os"
	Permissions "univia/src/modules/permission/models"
	Posts "univia/src/modules/post/models"
	Roles "univia/src/modules/role/models"
	Users "univia/src/modules/user/models"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase initializes and migrates the database.
func ConnectDatabase() *gorm.DB {
	// Retrieve connection info from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Format MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Connect to MySQL
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable SQL query logs
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = database
	fmt.Println("Connected to database!")
	if err := SeedData(DB); err != nil {
		log.Fatalf("Failed to seed: %v", err)
	}

	return DB
}

func CheckConnection() bool {
	if DB == nil {
		return false
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Failed to get generic database object: %v", err)
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		log.Printf("Database ping failed: %v", err)
		return false
	}

	var result int
	if err := DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		log.Printf("Test query failed: %v", err)
		return false
	}
	return result == 1
}

func SeedData(db *gorm.DB) error {
	// ---- Seed Categories ----
	categories := []string{
		"Social Media Trends", "Anime Fan Communities", "Music Production Tips",
		"Programming Tutorials", "Digital Art & Design", "K-Pop Culture",
		"Web Development", "Cosplay Inspiration", "Game Development",
		"AI & Machine Learning", "Manga Discussions", "Indie Music Artists",
		"Cybersecurity Corner", "Live Streaming Tips", "Anime Reviews",
		"Hip-Hop Culture", "Mobile App Development", "Viral Content Analysis",
		"Tech Gadgets Talk", "Songwriting & Composition",
	}

	for _, name := range categories {
		var cat Posts.Category
		db.FirstOrCreate(&cat, Posts.Category{
			Name: name,
		})
	}

	// ---- Seed Roles ----
	adminRole := Roles.Role{Name: "admin"}
	userRole := Roles.Role{Name: "user"}
	db.FirstOrCreate(&adminRole, Roles.Role{Name: "admin"})
	db.FirstOrCreate(&userRole, Roles.Role{Name: "user"})

	// Reselect ids of roles
	var adminRoleID uuid.UUID
	if err := db.Model(&Roles.Role{}).
		Where("name = ?", adminRole.Name).
		Select("id", &adminRoleID).Error; err != nil {
	}
	// ---- Seed Admin User ----
	adminUser := Users.User{
		Username:    "admin",
		Email:       "ties.node@outlook.com",
		PhoneNumber: 773598329,
		Password:    "$2a$12$OqiRAY9.CA7pj1zK4p42wuq0d63xf0l/ZXD7uQMDrBWU4.uGvdt12",
		Status:      true,
		RoleID:      adminRoleID,
	}
	db.FirstOrCreate(&adminUser, Users.User{Email: "ties.node@outlook.com"})

	// ---- Seed Permissions ----
	permNames := []string{
		"allow_create_role",
		"allow_create_permission",
		"allow_assign_permissions",
		"allow_list_roles",
		"allow_list_permissions",
	}

	var permissions []Permissions.Permission
	for _, name := range permNames {
		var perm Permissions.Permission
		db.FirstOrCreate(&perm, Permissions.Permission{Name: name})
		if err := db.Where("name = ?", perm.Name).First(&perm).Error; err != nil {
			return err
		}
		permissions = append(permissions, perm)
	}

	// ---- Seed Role-Permissions (assign all perms to admin) ----
	for _, perm := range permissions {
		rp := Roles.RolePermission{
			RoleID:       adminRoleID,
			PermissionID: perm.ID,
		}
		db.FirstOrCreate(&rp, Roles.RolePermission{
			RoleID:       adminRole.ID,
			PermissionID: perm.ID,
		})
	}

	fmt.Println("✅ Seed data completed with UUID models")
	return nil
}
