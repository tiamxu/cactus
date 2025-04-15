package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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

// 数据库操作方法
func GetUserByUsername(username string) (*User, error) {
	query := `
		SELECT 
			id, username, password, enable
		FROM user 
		WHERE username = ?
		LIMIT 1`

	var user User
	err := DB.Get(&user, query, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

// GetByID 根据ID获取用户
func GetUserByID(id int) (*User, error) {
	query := `
        SELECT id, username, enable, createTime, updateTime 
        FROM user 
        WHERE id = ?`

	user := &User{}
	err := DB.Get(user, query, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUsersByIDs(ids []int) ([]User, error) {
	query := `SELECT * FROM user WHERE id IN (?)`

	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}

	var users []User
	err = DB.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func CreateUser(user *User) error {
	query := `
		INSERT INTO users (
			username, password, email, 
			created_at, updated_at, status
		) VALUES (
			:username, :password, :email, 
			NOW(), NOW(), :status
		)`

	result, err := DB.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("create user failed: %w", err)
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)
	return nil
}

func ListUsers() ([]User, error) {
	var users []User
	err := DB.Select(&users, `
		SELECT id, username, email 
		FROM users 
		WHERE deleted_at IS NULL`)
	return users, err
}

// Create 创建用户
func Create(user *User) error {
	query := `
        INSERT INTO users 
        (username, password, email, status)
        VALUES (?, ?, ?, ?)`

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	// if err != nil {
	// 	return err
	// }

	result, err := DB.Exec(query,
		user.Username,
		user.Password,
		user.Enable,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)
	return nil
}

// Update 更新用户信息
func Update(user *User) error {
	query := `
        UPDATE user 
        SET username=?,status=?, updated_at=NOW() 
        WHERE id=? AND deleted_at IS NULL`

	_, err := DB.Exec(query,
		user.Username,
		user.Enable,
		user.ID,
	)
	return err
}

// Delete 软删除用户
func Delete(id uint) error {
	query := `UPDATE users SET deleted_at=NOW() WHERE id=?`
	_, err := DB.Exec(query, id)
	return err
}

// List 用户列表
func List(page, pageSize int) ([]*User, error) {
	query := `
        SELECT id, username, status, created_at, updated_at 
        FROM users 
        WHERE deleted_at IS NULL 
        LIMIT ? OFFSET ?`

	rows, err := DB.Query(query, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Enable,
			&user.CreateTime,
			&user.UpdateTime,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// ExistsByUsername 检查用户名是否存在
func ExistsByUsername(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)`
	err := DB.QueryRow(query, username).Scan(&exists)
	return exists, err
}
