package models

type RolePermissionsPermission struct {
	RoleId       int `db:"roleId" json:"roleId"`
	PermissionId int `db:"permissionId" json:"permissionId"`
}

func (RolePermissionsPermission) TableName() string {
	return "role_permissions_permission"
}

func GetPermissionsIdsByWhere(roleId int) ([]int, error) {
	var perIdList []int
	err := DB.Select(&perIdList,
		"SELECT permissionId FROM role_permissions_permission WHERE roleId = ?",
		roleId)
	if err != nil {
		return nil, err
	}
	return perIdList, nil
}
