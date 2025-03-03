package seeds

import (
	"log"

	"github.com/kifangamukundi/gm/loan/models"

	"gorm.io/gorm"
)

// Fulls access to all features
// Create, edit and publish blog posts, view entire blog analytics
// Create, edit their own blog posts, view own blog analytics
// Create blog posts but can't publish, no access to blog analytics
// Access blog content, receive emails and newsletters and manage their subscription preferences
// Access ecommerce features, view and purchase products, and manage personal account details
// Manage email collection forms, create and send email campaigns, view and segment email lists
// Oversee customer transactions, manage sales reporting, and monitor product performance in the ecommerce module
// Access customer orders, handle customer inquiries, and view customer profiles
// Access limited sales data for affiliate tracking, manage affiliate links, and view commissions
// Manage own products, view sales and inventory data, and interact with the e-commerce module
// Review and moderate user generated content such as comments, reviews, and forum postss
// View and analyze user behavior, blog performance, email campaign data, and ecomerce statistics

var roleNames = []string{
	"Admin",
	"Officer",
	"Agent",
	"Member",
}

var rolesToAssignPermissions = []string{
	"Admin",
}

func SeedRolesAndAssignPermissions(db *gorm.DB) error {
	var roles []models.Role

	for _, roleName := range roleNames {
		var existingRole models.Role
		if err := db.Where("role_name = ?", roleName).First(&existingRole).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("Error checking if role exists: %v", err)
			return err
		}

		if existingRole.ID == 0 {
			roles = append(roles, models.Role{
				RoleName: roleName,
			})
		}
	}

	if len(roles) > 0 {
		if err := db.CreateInBatches(roles, 100).Error; err != nil {
			log.Printf("Error seeding roles: %v", err)
			return err
		}
		log.Println("Roles seeded successfully.")
	} else {
		log.Println("No new roles to seed.")
	}

	var rolesToAssign []models.Role
	if err := db.Where("role_name IN ?", rolesToAssignPermissions).Find(&rolesToAssign).Error; err != nil {
		log.Printf("Error fetching specified roles: %v", err)
		return err
	}

	var permissions []models.Permission
	if err := db.Find(&permissions).Error; err != nil {
		log.Printf("Error fetching permissions: %v", err)
		return err
	}

	for _, role := range rolesToAssign {
		for _, permission := range permissions {
			var existingRolePermission models.RolePermission
			if err := db.Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).First(&existingRolePermission).Error; err == nil {
				log.Printf("Permission %d already assigned to role %s", permission.ID, role.RoleName)
				continue
			} else if err != gorm.ErrRecordNotFound {
				log.Printf("Error checking existing role_permission: %v", err)
				return err
			}

			rolePermission := models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
			}
			if err := db.Create(&rolePermission).Error; err != nil {
				log.Printf("Error assigning permission to role %v: %v", role.RoleName, err)
				return err
			}
		}
		log.Printf("All permissions assigned to the %s role successfully.", role.RoleName)
	}

	return nil
}
