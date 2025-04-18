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

// 查询用户信息
func GetUserByID(id int) (*User, error) {
	query := `
        SELECT id, username, password, enable, createTime, updateTime 
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

	return users, nil
}

// ExistsByUsername 检查用户名是否存在
func ExistsByUsername(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)`
	err := DB.QueryRow(query, username).Scan(&exists)
	return exists, err
}

func GetUsersByCondition(gender, enable, username string, limit, offset int) ([]User, int64, error) {
	var users []User
	var total int64

	// 构建基础查询
	query := "SELECT * FROM user WHERE 1=1"
	var params []interface{}

	// 添加筛选条件
	if enable != "" {
		query += " AND enable = ?"
		params = append(params, enable)
	}

	// 获取总数
	countQuery := "SELECT COUNT(*) FROM user WHERE 1=1"
	if enable != "" {
		countQuery += " AND enable = ?"
	}
	err := DB.Get(&total, countQuery, params...)
	if err != nil {
		return nil, 0, err
	}

	// 添加分页
	query += " LIMIT ? OFFSET ?"
	params = append(params, limit, offset)

	// 执行查询
	err = DB.Select(&users, query, params...)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetPasswordHash 获取用户当前密码哈希
func GetPasswordHash(uid int) (string, error) {
	var currentPasswordHash string
	err := DB.Get(&currentPasswordHash,
		"SELECT password FROM user WHERE id = ?",
		uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("用户不存在")
		}
		return "", fmt.Errorf("查询密码失败: %w", err)
	}
	return currentPasswordHash, nil
}

// UpdatePassword 更新用户密码
func UpdatePassword(uid int, newHash string) error {
	_, err := DB.Exec("UPDATE user SET password = ? WHERE id = ?", newHash, uid)
	if err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}
	return nil
}
