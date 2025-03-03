package seeds

import (
	"log"

	"github.com/kifangamukundi/gm/loan/models"

	"gorm.io/gorm"
)

// assign these to agent field_overview, create_member, view_members, edit_members, delete_member, create_loan, view_loans, edit_loan, delete_loan, office_overview
var permissionNames = []string{
	"create_permission",
	"create_role", "data_collection_overview",
	"delete_media", "delete_permission", "delete_role", "edit_permission",
	"edit_role", "edit_user", "security_overview", "upload_media",
	"view_permissions", "view_roles", "view_users",
	"field_overview", "create_agent", "view_agents", "edit_agent", "delete_agent",
	"create_group", "view_groups", "edit_group", "delete_group",
	"create_officer", "view_officers", "edit_officer", "delete_officer",
	"create_member", "view_members", "edit_member", "delete_member",
	"create_loan", "view_loans", "edit_loan", "delete_loan",
	"office_overview",
}

func SeedPermissions(db *gorm.DB) error {
	var permissionRecords []models.Permission

	for _, permissionName := range permissionNames {
		var existingPermission models.Permission
		if err := db.Where("permission_name = ?", permissionName).First(&existingPermission).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("Error checking if permission exists: %v", err)
			return err
		}

		if existingPermission.ID == 0 {
			permissionRecords = append(permissionRecords, models.Permission{PermissionName: permissionName})
		}
	}

	if len(permissionRecords) > 0 {
		if err := db.CreateInBatches(permissionRecords, 100).Error; err != nil {
			log.Printf("Error seeding permissions: %v", err)
			return err
		}
	}

	log.Println("Permissions seeded successfully.")
	return nil
}
