package models

type RolePermissionsPermission struct {
	RoleId       int `json:"roleId"`
	PermissionId int `json:"permissionId"`
}

func (RolePermissionsPermission) TableName() string {
	return "role_permissions_permission"
}
