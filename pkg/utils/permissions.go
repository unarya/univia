package utils

// Permissions defines a set of permission constants used across the system.
// Convention: "allow_<action>_<resource>"
var Permissions = map[string]string{
	// ===== Auth & User =====
	"ALLOW_REGISTER_USER":  "allow_register_user",
	"ALLOW_LOGIN_USER":     "allow_login_user",
	"ALLOW_GET_USER":       "allow_get_user",
	"ALLOW_UPDATE_USER":    "allow_update_user",
	"ALLOW_RESET_PASSWORD": "allow_reset_password",

	// ===== Role & Permission Management =====
	"ALLOW_LIST_ROLES":         "allow_list_roles",
	"ALLOW_CREATE_ROLE":        "allow_create_role",
	"ALLOW_ASSIGN_PERMISSIONS": "allow_assign_permissions",
	"ALLOW_DELETE_ROLE":        "allow_delete_role",

	"ALLOW_LIST_PERMISSIONS":  "allow_list_permissions",
	"ALLOW_CREATE_PERMISSION": "allow_create_permission",

	// ===== Posts / Content =====
	"ALLOW_LIST_POSTS":  "allow_list_posts",
	"ALLOW_CREATE_POST": "allow_create_post",
	"ALLOW_UPDATE_POST": "allow_update_post",
	"ALLOW_DELETE_POST": "allow_delete_post",
	"ALLOW_LIKE_POST":   "allow_like_post",
	"ALLOW_UNDO_LIKE":   "allow_undo_like_post",

	// ===== Team / Workspace =====
	"ALLOW_CREATE_TEAM":      "allow_create_team",
	"ALLOW_DELETE_TEAM":      "allow_delete_team",
	"ALLOW_INVITE_MEMBER":    "allow_invite_member",
	"ALLOW_REMOVE_MEMBER":    "allow_remove_member",
	"ALLOW_ASSIGN_TEAM_ROLE": "allow_assign_team_role",

	// ===== Notifications =====
	"ALLOW_LIST_NOTIFICATIONS":       "allow_list_notifications",
	"ALLOW_SEEN_SINGLE_NOTIFICATION": "allow_seen_single_notification",
	"ALLOW_SEEN_ALL_NOTIFICATIONS":   "allow_seen_all_notifications",

	// ===== Billing =====
	"ALLOW_VIEW_BILLING":        "allow_view_billing",
	"ALLOW_UPDATE_PAYMENT":      "allow_update_payment",
	"ALLOW_CANCEL_SUBSCRIPTION": "allow_cancel_subscription",

	// ===== System / Infrastructure =====
	"ALLOW_HEALTH_CHECK":  "allow_health_check",
	"ALLOW_VIEW_LOGS":     "allow_view_logs",
	"ALLOW_MANAGE_SERVER": "allow_manage_server",
}
