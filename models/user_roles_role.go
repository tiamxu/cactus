package models

type UserRolesRole struct {
	UserId int `json:"userId"`
	RoleId int `json:"roleId"`
}

func (UserRolesRole) TableName() string {
	return "user_roles_role"
}

// 查询用户角色 ID 列表
func GetUserRoleIDs(userId int) ([]int, error) {
	query := `SELECT roleId FROM user_roles_role WHERE userId = ?`
	var roleIds []int
	err := DB.Select(&roleIds, query, userId)
	if err != nil {
		return nil, err
	}
	return roleIds, err
}
