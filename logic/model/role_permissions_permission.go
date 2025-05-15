package model

type RolePermissionsPermission struct {
	RoleId       int `db:"roleId" json:"roleId"`
	PermissionId int `db:"permissionId" json:"permissionId"`
}

func (RolePermissionsPermission) TableName() string {
	return "role_permissions_permission"
}
