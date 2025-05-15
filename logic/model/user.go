package model

import (
	"time"
)

type User struct {
	ID         int       `db:"id" json:"id"`
	Username   string    `db:"username" json:"username"`
	Password   string    `db:"password" json:"password"`
	Enable     bool      `db:"enable" json:"enable"`
	CreateTime time.Time `db:"createTime" json:"createTime"`
	UpdateTime time.Time `db:"updateTime" json:"updateTime"`
}

func (u *User) TableName() string {
	return "user"
}
