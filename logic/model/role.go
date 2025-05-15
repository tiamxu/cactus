package model

type Role struct {
	ID     int    `db:"id" json:"id"`
	Code   string `db:"code" json:"code"`
	Name   string `db:"name" json:"name"`
	Enable bool   `db:"enable" json:"enable"`
}

func (Role) TableName() string {
	return "role"
}
