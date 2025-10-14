package mysql

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	Permissions "github.com/unarya/univia/internal/api/modules/permission/models"
	Posts "github.com/unarya/univia/internal/api/modules/post/models"
	Profiles "github.com/unarya/univia/internal/api/modules/profile/models"
	Roles "github.com/unarya/univia/internal/api/modules/role/models"
	Users "github.com/unarya/univia/internal/api/modules/user/models"
	"github.com/unarya/univia/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase initializes and migrates the mysql.
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
		log.Fatalf("Failed to connect to mysql: %v", err)
	}

	DB = database
	fmt.Println("Connected to mysql!")
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
		log.Printf("Failed to get generic mysql object: %v", err)
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

// SeedData initializes the mysql with essential data including:
// - Post categories for content organization
// - System and team roles for access control hierarchy
// - Permissions for granular access management
// - Role-permission mappings based on organizational structure
// - Default admin user with profile for initial system access
//
// This function is idempotent and uses FirstOrCreate to prevent duplicate entries.
// It should be called during application initialization or via migration scripts.
func SeedData(db *gorm.DB) error {
	// Seed in order of dependencies
	if err := seedCategories(db); err != nil {
		return err
	}

	roles, err := seedRoles(db)
	if err != nil {
		return err
	}

	permissions, err := seedPermissions(db)
	if err != nil {
		return err
	}

	if err := assignRolePermissions(db, roles, permissions); err != nil {
		return err
	}

	if err := seedDefaultAdminUser(db, roles["super_admin"]); err != nil {
		return err
	}

	fmt.Println("✅ Seed data completed successfully:")
	fmt.Println("   - Categories created")
	fmt.Printf("   - %d roles initialized\n", len(roles))
	fmt.Printf("   - %d permissions configured\n", len(permissions))
	fmt.Println("   - Role-permission mappings established")
	fmt.Println("   - Default admin user with profile created")
	fmt.Println("⚠️  SECURITY WARNING: Change default admin credentials immediately!")

	return nil
}

// seedCategories creates post categories for content organization
func seedCategories(db *gorm.DB) error {
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
		if err := db.FirstOrCreate(&cat, Posts.Category{Name: name}).Error; err != nil {
			return fmt.Errorf("failed to create category %s: %w", name, err)
		}
	}

	return nil
}

// seedRoles creates system and team roles, returns a map of role names to role objects
func seedRoles(db *gorm.DB) (map[string]Roles.Role, error) {
	systemRoles := []string{"super_admin", "general_admin", "billing_manager", "support"}
	teamRoles := []string{"team_owner", "team_admin", "team_editor", "team_viewer", "team_guest"}
	allRoles := append(systemRoles, teamRoles...)

	roleMap := make(map[string]Roles.Role)

	for _, roleName := range allRoles {
		var role Roles.Role
		// Always use WHERE + FirstOrCreate to ensure no duplicate
		if err := db.
			Where("name = ?", roleName).
			Attrs(Roles.Role{
				ID:   uuid.New(),
				Name: roleName,
			}).
			FirstOrCreate(&role).Error; err != nil {
			return nil, fmt.Errorf("failed to create or fetch role %s: %w", roleName, err)
		}

		// Ensure valid ID
		if role.ID == uuid.Nil {
			if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
				return nil, fmt.Errorf("failed to reload role %s: %w", roleName, err)
			}
		}

		roleMap[roleName] = role
	}

	return roleMap, nil
}

// seedPermissions creates all system permissions, returns a slice of permission objects
func seedPermissions(db *gorm.DB) ([]Permissions.Permission, error) {
	var permEntities []Permissions.Permission

	for _, permName := range utils.Permissions {
		var perm Permissions.Permission

		if err := db.
			Where("name = ?", permName).
			Attrs(Permissions.Permission{
				ID:   uuid.New(),
				Name: permName,
			}).
			FirstOrCreate(&perm).Error; err != nil {
			return nil, fmt.Errorf("failed to create or fetch permission %s: %w", permName, err)
		}

		if perm.ID == uuid.Nil {
			if err := db.Where("name = ?", permName).First(&perm).Error; err != nil {
				return nil, fmt.Errorf("failed to reload permission %s: %w", permName, err)
			}
		}

		permEntities = append(permEntities, perm)
	}

	return permEntities, nil
}

// assignRolePermissions assigns permissions to roles based on access control requirements
func assignRolePermissions(db *gorm.DB, roles map[string]Roles.Role, permissions []Permissions.Permission) error {
	// Define permission sets for each role
	permissionSets := map[string][]string{
		"super_admin":     getAllPermissionNames(permissions),
		"general_admin":   getGeneralAdminPermissions(),
		"billing_manager": getBillingManagerPermissions(),
		"support":         getSupportPermissions(),
		"team_owner":      getTeamOwnerPermissions(),
		"team_admin":      getTeamAdminPermissions(),
		"team_editor":     getTeamEditorPermissions(),
		"team_viewer":     getTeamViewerPermissions(),
		"team_guest":      getTeamGuestPermissions(),
	}

	// Build map of permission name → ID
	permMap := make(map[string]uuid.UUID)
	for _, perm := range permissions {
		permMap[perm.Name] = perm.ID
	}

	// Assign permissions to roles
	for roleName, permNames := range permissionSets {
		role, exists := roles[roleName]
		if !exists || role.ID == uuid.Nil {
			return fmt.Errorf("invalid role %s (missing ID)", roleName)
		}

		for _, permName := range permNames {
			permID, exists := permMap[permName]
			if !exists || permID == uuid.Nil {
				continue
			}

			rp := Roles.RolePermission{
				RoleID:       role.ID,
				PermissionID: permID,
			}

			if err := db.
				Where("role_id = ? AND permission_id = ?", role.ID, permID).
				FirstOrCreate(&rp).Error; err != nil {
				return fmt.Errorf("failed to assign permission %s to role %s: %w", permName, roleName, err)
			}
		}
	}

	return nil
}

// seedDefaultAdminUser creates the default super admin account with profile
// Reads plain password from ADMIN_PASSWORD env variable (default: "admin123")
// Hashes the password using bcrypt before storing in mysql
func seedDefaultAdminUser(db *gorm.DB, superAdminRole Roles.Role) error {
	// Load admin configuration from environment
	config := loadAdminConfig()

	// Hash the plain password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.PlainPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	// Create admin user
	adminUser := Users.User{
		Username:    config.Username,
		Email:       config.Email,
		PhoneNumber: uint64(config.PhoneNumber),
		Password:    string(hashedPassword),
		Status:      true,
		RoleID:      superAdminRole.ID,
	}

	if err := db.Where(Users.User{Email: adminUser.Email}).FirstOrCreate(&adminUser).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	var adminID uuid.UUID
	var user Users.User

	if err := db.Model(&Users.User{}).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find admin user: %w", err)
	}
	adminID = user.ID
	// Create admin profile with configured values
	if err := createAdminProfile(db, adminID, config); err != nil {
		return err
	}

	// Log admin info (mask sensitive data)
	fmt.Printf("   - Admin user: %s (%s)\n", config.Username, maskEmail(config.Email))
	if config.IsDefaultPassword {
		fmt.Println("   ⚠️  Using default password 'admin123' - CHANGE IMMEDIATELY in production!")
	}

	return nil
}

// AdminConfig holds configuration for the default admin account
type AdminConfig struct {
	Username          string
	Email             string
	PhoneNumber       int
	PlainPassword     string
	IsDefaultPassword bool
	ProfilePic        string
	CoverPhoto        string
	BackgroundColor   string
	Location          string
	Bio               string
	Interests         []string
	SocialLinks       map[string]string
}

// loadAdminConfig loads admin configuration from environment variables with defaults
func loadAdminConfig() AdminConfig {
	plainPassword := getEnvOrDefault("ADMIN_PASSWORD", "admin123")
	isDefaultPassword := getEnv("ADMIN_PASSWORD") == ""

	config := AdminConfig{
		Username:          getEnvOrDefault("ADMIN_USERNAME", "admin"),
		Email:             getEnvOrDefault("ADMIN_EMAIL", "ties.node@outlook.com"),
		PhoneNumber:       getEnvAsIntOrDefault("ADMIN_PHONE", 773598329),
		PlainPassword:     plainPassword,
		IsDefaultPassword: isDefaultPassword,
		ProfilePic:        getEnvOrDefault("ADMIN_PROFILE_PIC", "https://ui-avatars.com/api/?name=Super+Admin&size=256&background=7b2cbf&color=fff"),
		CoverPhoto:        getEnvOrDefault("ADMIN_COVER_PHOTO", "https://images.unsplash.com/photo-1557683316-973673baf926?w=1200&h=400&fit=crop"),
		BackgroundColor:   getEnvOrDefault("ADMIN_BG_COLOR", "#7b2cbf"),
		Location:          getEnvOrDefault("ADMIN_LOCATION", "Headquarters"),
		Bio:               getEnvOrDefault("ADMIN_BIO", "Super Admin of the system. Full access to all features."),
		Interests: []string{
			"System Administration",
			"Security",
			"Scaling",
			"DevOps",
		},
		SocialLinks: map[string]string{
			"github":  getEnvOrDefault("ADMIN_GITHUB", "https://github.com/superadmin"),
			"twitter": getEnvOrDefault("ADMIN_TWITTER", "https://twitter.com/superadmin"),
		},
	}

	return config
}

// createAdminProfile creates the profile for the admin user
func createAdminProfile(db *gorm.DB, userID uuid.UUID, config AdminConfig) error {
	interests, err := json.Marshal(config.Interests)
	if err != nil {
		return fmt.Errorf("failed to marshal interests: %w", err)
	}

	socialLinks, err := json.Marshal(config.SocialLinks)
	if err != nil {
		return fmt.Errorf("failed to marshal social links: %w", err)
	}

	adminProfile := Profiles.Profile{
		UserID:          userID,
		ProfilePic:      config.ProfilePic,
		CoverPhoto:      config.CoverPhoto,
		BackgroundColor: config.BackgroundColor,
		Gender:          "Other",
		Location:        config.Location,
		Bio:             config.Bio,
		Interests:       datatypes.JSON(interests),
		SocialLinks:     datatypes.JSON(socialLinks),
	}

	if err := db.Where(Profiles.Profile{UserID: userID}).FirstOrCreate(&adminProfile).Error; err != nil {
		return fmt.Errorf("failed to create profile for super admin: %w", err)
	}

	return nil
}

// =================== Permission Set Helpers ===================

// getAllPermissionNames returns all permission names
func getAllPermissionNames(permissions []Permissions.Permission) []string {
	names := make([]string, len(permissions))
	for i, perm := range permissions {
		names[i] = perm.Name
	}
	return names
}

// getGeneralAdminPermissions returns permissions for general admin
// Excludes sensitive billing and orchestrator management
func getGeneralAdminPermissions() []string {
	excluded := map[string]bool{
		utils.Permissions["ALLOW_VIEW_BILLING"]:        true,
		utils.Permissions["ALLOW_UPDATE_PAYMENT"]:      true,
		utils.Permissions["ALLOW_CANCEL_SUBSCRIPTION"]: true,
		utils.Permissions["ALLOW_MANAGE_SERVER"]:       true,
	}

	var perms []string
	for _, permName := range utils.Permissions {
		if !excluded[permName] {
			perms = append(perms, permName)
		}
	}
	return perms
}

// getBillingManagerPermissions returns permissions for billing manager
func getBillingManagerPermissions() []string {
	return []string{
		utils.Permissions["ALLOW_VIEW_BILLING"],
		utils.Permissions["ALLOW_UPDATE_PAYMENT"],
		utils.Permissions["ALLOW_CANCEL_SUBSCRIPTION"],
		utils.Permissions["ALLOW_GET_USER"],
		utils.Permissions["ALLOW_VIEW_POST"],
	}
}

// getSupportPermissions returns permissions for support staff
func getSupportPermissions() []string {
	return []string{
		utils.Permissions["ALLOW_GET_USER"],
		utils.Permissions["ALLOW_VIEW_POST"],
		utils.Permissions["ALLOW_VIEW_NOTIFICATION"],
		utils.Permissions["ALLOW_CREATE_POST"],
	}
}

// getTeamOwnerPermissions returns permissions for team owner
// Full control including team deletion
func getTeamOwnerPermissions() []string {
	return []string{
		utils.Permissions["ALLOW_CREATE_POST"],
		utils.Permissions["ALLOW_UPDATE_POST"],
		utils.Permissions["ALLOW_DELETE_POST"],
		utils.Permissions["ALLOW_VIEW_POST"],
		utils.Permissions["ALLOW_MANAGE_TEAM"],
		utils.Permissions["ALLOW_DELETE_TEAM"],
		utils.Permissions["ALLOW_INVITE_MEMBER"],
		utils.Permissions["ALLOW_REMOVE_MEMBER"],
		utils.Permissions["ALLOW_VIEW_NOTIFICATION"],
		utils.Permissions["ALLOW_SEND_NOTIFICATION"],
	}
}

// getTeamAdminPermissions returns permissions for team admin
// Same as owner but cannot delete team
func getTeamAdminPermissions() []string {
	return []string{
		utils.Permissions["ALLOW_CREATE_POST"],
		utils.Permissions["ALLOW_UPDATE_POST"],
		utils.Permissions["ALLOW_DELETE_POST"],
		utils.Permissions["ALLOW_VIEW_POST"],
		utils.Permissions["ALLOW_MANAGE_TEAM"],
		utils.Permissions["ALLOW_INVITE_MEMBER"],
		utils.Permissions["ALLOW_REMOVE_MEMBER"],
		utils.Permissions["ALLOW_VIEW_NOTIFICATION"],
		utils.Permissions["ALLOW_SEND_NOTIFICATION"],
	}
}

// getTeamEditorPermissions returns permissions for team editor
// Content management focused
func getTeamEditorPermissions() []string {
	return []string{
		utils.Permissions["ALLOW_CREATE_POST"],
		utils.Permissions["ALLOW_UPDATE_POST"],
		utils.Permissions["ALLOW_DELETE_POST"],
		utils.Permissions["ALLOW_VIEW_POST"],
		utils.Permissions["ALLOW_VIEW_NOTIFICATION"],
	}
}

// getTeamViewerPermissions returns permissions for team viewer
// Read-only access
func getTeamViewerPermissions() []string {
	return []string{
		utils.Permissions["ALLOW_VIEW_POST"],
		utils.Permissions["ALLOW_VIEW_NOTIFICATION"],
	}
}

// getTeamGuestPermissions returns permissions for team guest
// Minimal public access
func getTeamGuestPermissions() []string {
	return []string{
		utils.Permissions["ALLOW_VIEW_POST"],
	}
}

// =================== Environment Helpers ===================

// getEnv retrieves an environment variable value
func getEnv(key string) string {
	return os.Getenv(key)
}

// getEnvOrDefault retrieves an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault retrieves an environment variable as int or returns a default value
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// maskEmail masks an email address for secure logging
// Example: user@example.com -> u***r@example.com
func maskEmail(email string) string {
	if len(email) < 3 {
		return "***"
	}

	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atIndex = i
			break
		}
	}

	if atIndex <= 0 {
		return email[:1] + "***"
	}

	// Show first and last character before @
	if atIndex == 1 {
		return email[:1] + "***" + email[atIndex:]
	}

	return email[:1] + "***" + email[atIndex-1:]
}
