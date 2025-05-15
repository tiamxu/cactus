package repo

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
