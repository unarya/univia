package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	Permissions "github.com/deva-labs/univia/api/gin/src/modules/permission/models"
	Posts "github.com/deva-labs/univia/api/gin/src/modules/post/models"
	Profiles "github.com/deva-labs/univia/api/gin/src/modules/profile/models"
	Roles "github.com/deva-labs/univia/api/gin/src/modules/role/models"
	Users "github.com/deva-labs/univia/api/gin/src/modules/user/models"
	"github.com/deva-labs/univia/api/gin/src/utils"
	"gorm.io/datatypes"
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

// SeedData initializes the database with essential data including:
// - Post categories for content organization
// - System and team roles for access control hierarchy
// - Permissions for granular access management
// - Role-permission mappings based on organizational structure
// - Default admin user for initial system access
//
// This function is idempotent and uses FirstOrCreate to prevent duplicate entries.
// It should be called during application initialization or via migration scripts.
func SeedData(db *gorm.DB) error {
	// ---- Seed Post Categories ----
	// These categories help organize and classify user-generated content
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

	// ---- Define Role Hierarchy ----
	// System-level roles: Control platform-wide operations and configurations
	systemRoles := []string{"super_admin", "general_admin", "billing_manager", "support"}

	// Team-level roles: Control team-specific resources and collaboration
	teamRoles := []string{"team_owner", "team_admin", "team_editor", "team_viewer", "team_guest"}

	allRoles := append(systemRoles, teamRoles...)

	// Create all roles in database
	var roleEntities []Roles.Role
	for _, roleName := range allRoles {
		var role Roles.Role
		if err := db.FirstOrCreate(&role, Roles.Role{Name: roleName}).Error; err != nil {
			return fmt.Errorf("failed to create role %s: %w", roleName, err)
		}
		roleEntities = append(roleEntities, role)
	}

	// ---- Seed Permissions ----
	// Permissions define atomic actions that can be performed in the system
	var permEntities []Permissions.Permission
	for _, permName := range utils.Permissions {
		var perm Permissions.Permission
		if err := db.FirstOrCreate(&perm, Permissions.Permission{Name: permName}).Error; err != nil {
			return fmt.Errorf("failed to create permission %s: %w", permName, err)
		}
		permEntities = append(permEntities, perm)
	}

	// ---- Assign Permissions to System Roles ----

	// Super Admin: Full system access with all permissions
	// This role has unrestricted access to all platform features and configurations
	var superAdminRole Roles.Role
	if err := db.First(&superAdminRole, "name = ?", "super_admin").Error; err != nil {
		return fmt.Errorf("failed to find super_admin role: %w", err)
	}

	for _, perm := range permEntities {
		rp := Roles.RolePermission{
			RoleID:       superAdminRole.ID,
			PermissionID: perm.ID,
		}
		if err := db.FirstOrCreate(&rp, Roles.RolePermission{
			RoleID:       rp.RoleID,
			PermissionID: rp.PermissionID,
		}).Error; err != nil {
			return fmt.Errorf("failed to assign permission to super_admin: %w", err)
		}
	}

	// General Admin: Administrative access excluding sensitive billing and server operations
	// Can manage users, content, and general platform operations
	var generalAdminRole Roles.Role
	if err := db.First(&generalAdminRole, "name = ?", "general_admin").Error; err != nil {
		return fmt.Errorf("failed to find general_admin role: %w", err)
	}

	// Excluded permissions for general admin (sensitive operations only)
	excludedPermissions := map[string]bool{
		utils.Permissions["ALLOW_VIEW_BILLING"]:        true,
		utils.Permissions["ALLOW_UPDATE_PAYMENT"]:      true,
		utils.Permissions["ALLOW_CANCEL_SUBSCRIPTION"]: true,
		utils.Permissions["ALLOW_MANAGE_SERVER"]:       true,
	}

	for _, perm := range permEntities {
		if excludedPermissions[perm.Name] {
			continue
		}
		rp := Roles.RolePermission{
			RoleID:       generalAdminRole.ID,
			PermissionID: perm.ID,
		}
		if err := db.FirstOrCreate(&rp, rp).Error; err != nil {
			return fmt.Errorf("failed to assign permission to general_admin: %w", err)
		}
	}

	// Billing Manager: Focused access to billing, payment, and subscription management
	var billingManagerRole Roles.Role
	if err := db.First(&billingManagerRole, "name = ?", "billing_manager").Error; err != nil {
		return fmt.Errorf("failed to find billing_manager role: %w", err)
	}

	billingPermissions := map[string]bool{
		utils.Permissions["ALLOW_VIEW_BILLING"]:        true,
		utils.Permissions["ALLOW_UPDATE_PAYMENT"]:      true,
		utils.Permissions["ALLOW_CANCEL_SUBSCRIPTION"]: true,
	}

	for _, perm := range permEntities {
		if billingPermissions[perm.Name] {
			rp := Roles.RolePermission{
				RoleID:       billingManagerRole.ID,
				PermissionID: perm.ID,
			}
			if err := db.FirstOrCreate(&rp, rp).Error; err != nil {
				return fmt.Errorf("failed to assign permission to billing_manager: %w", err)
			}
		}
	}

	// ---- Assign Permissions to Team Roles ----

	// Team Owner: Full control over team resources including posts, members, and notifications
	// Can delete the team and manage all team settings
	var teamOwnerRole Roles.Role
	if err := db.First(&teamOwnerRole, "name = ?", "team_owner").Error; err != nil {
		return fmt.Errorf("failed to find team_owner role: %w", err)
	}

	teamOwnerPermissions := map[string]bool{
		utils.Permissions["ALLOW_CREATE_POST"]:       true,
		utils.Permissions["ALLOW_UPDATE_POST"]:       true,
		utils.Permissions["ALLOW_DELETE_POST"]:       true,
		utils.Permissions["ALLOW_VIEW_POST"]:         true,
		utils.Permissions["ALLOW_MANAGE_TEAM"]:       true,
		utils.Permissions["ALLOW_DELETE_TEAM"]:       true,
		utils.Permissions["ALLOW_INVITE_MEMBER"]:     true,
		utils.Permissions["ALLOW_REMOVE_MEMBER"]:     true,
		utils.Permissions["ALLOW_VIEW_NOTIFICATION"]: true,
		utils.Permissions["ALLOW_SEND_NOTIFICATION"]: true,
	}

	for _, perm := range permEntities {
		if teamOwnerPermissions[perm.Name] {
			rp := Roles.RolePermission{
				RoleID:       teamOwnerRole.ID,
				PermissionID: perm.ID,
			}
			if err := db.FirstOrCreate(&rp, rp).Error; err != nil {
				return fmt.Errorf("failed to assign permission to team_owner: %w", err)
			}
		}
	}

	// Team Admin: Same as owner but cannot delete the team
	// Can manage team operations and members but lacks destructive team actions
	var teamAdminRole Roles.Role
	if err := db.First(&teamAdminRole, "name = ?", "team_admin").Error; err != nil {
		return fmt.Errorf("failed to find team_admin role: %w", err)
	}

	teamAdminPermissions := map[string]bool{
		utils.Permissions["ALLOW_CREATE_POST"]:       true,
		utils.Permissions["ALLOW_UPDATE_POST"]:       true,
		utils.Permissions["ALLOW_DELETE_POST"]:       true,
		utils.Permissions["ALLOW_VIEW_POST"]:         true,
		utils.Permissions["ALLOW_MANAGE_TEAM"]:       true,
		utils.Permissions["ALLOW_INVITE_MEMBER"]:     true,
		utils.Permissions["ALLOW_REMOVE_MEMBER"]:     true,
		utils.Permissions["ALLOW_VIEW_NOTIFICATION"]: true,
		utils.Permissions["ALLOW_SEND_NOTIFICATION"]: true,
	}

	for _, perm := range permEntities {
		if teamAdminPermissions[perm.Name] {
			rp := Roles.RolePermission{
				RoleID:       teamAdminRole.ID,
				PermissionID: perm.ID,
			}
			if err := db.FirstOrCreate(&rp, rp).Error; err != nil {
				return fmt.Errorf("failed to assign permission to team_admin: %w", err)
			}
		}
	}

	// Team Editor: Can create, update, and delete posts; view notifications
	// Focused on content management without administrative privileges
	var teamEditorRole Roles.Role
	if err := db.First(&teamEditorRole, "name = ?", "team_editor").Error; err != nil {
		return fmt.Errorf("failed to find team_editor role: %w", err)
	}

	teamEditorPermissions := map[string]bool{
		utils.Permissions["ALLOW_CREATE_POST"]:       true,
		utils.Permissions["ALLOW_UPDATE_POST"]:       true,
		utils.Permissions["ALLOW_DELETE_POST"]:       true,
		utils.Permissions["ALLOW_VIEW_POST"]:         true,
		utils.Permissions["ALLOW_VIEW_NOTIFICATION"]: true,
	}

	for _, perm := range permEntities {
		if teamEditorPermissions[perm.Name] {
			rp := Roles.RolePermission{
				RoleID:       teamEditorRole.ID,
				PermissionID: perm.ID,
			}
			if err := db.FirstOrCreate(&rp, rp).Error; err != nil {
				return fmt.Errorf("failed to assign permission to team_editor: %w", err)
			}
		}
	}

	// Team Viewer: Read-only access to posts and notifications
	// Can view team content but cannot make modifications
	var teamViewerRole Roles.Role
	if err := db.First(&teamViewerRole, "name = ?", "team_viewer").Error; err != nil {
		return fmt.Errorf("failed to find team_viewer role: %w", err)
	}

	teamViewerPermissions := map[string]bool{
		utils.Permissions["ALLOW_VIEW_POST"]:         true,
		utils.Permissions["ALLOW_VIEW_NOTIFICATION"]: true,
	}

	for _, perm := range permEntities {
		if teamViewerPermissions[perm.Name] {
			rp := Roles.RolePermission{
				RoleID:       teamViewerRole.ID,
				PermissionID: perm.ID,
			}
			if err := db.FirstOrCreate(&rp, rp).Error; err != nil {
				return fmt.Errorf("failed to assign permission to team_viewer: %w", err)
			}
		}
	}

	// Team Guest: Minimal access to public posts only
	// Most restricted role, typically for external collaborators
	var teamGuestRole Roles.Role
	if err := db.First(&teamGuestRole, "name = ?", "team_guest").Error; err != nil {
		return fmt.Errorf("failed to find team_guest role: %w", err)
	}

	teamGuestPermissions := map[string]bool{
		utils.Permissions["ALLOW_VIEW_POST"]: true,
	}

	for _, perm := range permEntities {
		if teamGuestPermissions[perm.Name] {
			rp := Roles.RolePermission{
				RoleID:       teamGuestRole.ID,
				PermissionID: perm.ID,
			}
			if err := db.FirstOrCreate(&rp, rp).Error; err != nil {
				return fmt.Errorf("failed to assign permission to team_guest: %w", err)
			}
		}
	}

	// ---- Seed Default Admin User ----
	// Creates the initial super admin account for system bootstrap
	// Default credentials:
	//   - Username: admin
	//   - Email: ties.node@outlook.com
	//   - Password: admin123 (pre-hashed with bcrypt)
	// IMPORTANT: Change these credentials immediately after first login in production!
	adminUser := Users.User{
		Username:    "admin",
		Email:       "ties.node@outlook.com",
		PhoneNumber: 773598329,
		Password:    "$2a$12$OqiRAY9.CA7pj1zK4p42wuq0d63xf0l/ZXD7uQMDrBWU4.uGvdt12", // admin123
		Status:      true,
		RoleID:      superAdminRole.ID,
	}

	if err := db.FirstOrCreate(&adminUser, Users.User{Email: "ties.node@outlook.com"}).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	// ---- Seed Default Profile for Super Admin ----
	// Creates a profile linked to the default super admin account
	// This ensures the admin has a complete identity in the system
	interests, _ := json.Marshal([]string{"System Administration", "Security", "Scaling"})
	socialLinks, _ := json.Marshal([]string{"https://github.com/superadmin", "https://twitter.com/superadmin"})
	adminProfile := Profiles.Profile{
		UserID:          adminUser.ID,
		ProfilePic:      "https://example.com/images/admin-profile.png",
		CoverPhoto:      "https://example.com/images/admin-cover.jpg",
		BackgroundColor: "#7b2cbf",
		Gender:          "Other",
		Location:        "Headquarters",
		Bio:             "Super Admin of the system. Full access to all features.",
		Interests:       datatypes.JSON(interests),
		SocialLinks:     datatypes.JSON(socialLinks),
	}

	if err := db.FirstOrCreate(&adminProfile, Profiles.Profile{UserID: adminUser.ID}).Error; err != nil {
		return fmt.Errorf("failed to create profile for super admin: %w", err)
	}
	fmt.Println("✅ Seed data completed successfully:")
	fmt.Printf("   - %d categories created\n", len(categories))
	fmt.Printf("   - %d roles initialized\n", len(allRoles))
	fmt.Printf("   - %d permissions configured\n", len(permEntities))
	fmt.Println("   - Role-permission mappings established")
	fmt.Println("   - Default admin user created")
	fmt.Println("⚠️  SECURITY WARNING: Change default admin credentials immediately!")

	return nil
}
