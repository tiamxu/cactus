package model

type UserRolesRole struct {
	UserId int `json:"userId"`
	RoleId int `json:"roleId"`
}

func (UserRolesRole) TableName() string {
	return "user_roles_role"
}
