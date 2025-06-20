package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/tiamxu/cactus/logic/model"
)

// 创建环境
func CreateEnvironment(ctx context.Context, env *model.Environment) error {
	query := `INSERT INTO environments 
		(name, code, description, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, NOW(), NOW())`

	result, err := DB.ExecContext(ctx, query,
		env.Name, env.Code, env.Description, env.Status)
	if err != nil {
		return fmt.Errorf("创建环境失败: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取环境ID失败: %v", err)
	}
	env.ID = int(id)
	return nil
}

// 根据ID获取环境
func GetEnvironmentByID(ctx context.Context, id int) (*model.Environment, error) {
	var env model.Environment
	query := "SELECT * FROM environments WHERE id = ?"
	err := DB.GetContext(ctx, &env, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("环境不存在")
		}
		return nil, fmt.Errorf("查询环境失败: %v", err)
	}
	return &env, nil
}

// 获取环境列表
func ListEnvironments(ctx context.Context, name string, status *int, page, pageSize int) ([]*model.Environment, int64, error) {
	var envs []*model.Environment
	var total int64

	// 构建查询条件
	var where []string
	var args []interface{}

	if name != "" {
		where = append(where, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}
	if status != nil {
		where = append(where, "status = ?")
		args = append(args, *status)
	}

	// 计算总数
	countQuery := "SELECT COUNT(*) FROM environments"
	if len(where) > 0 {
		countQuery += " WHERE " + strings.Join(where, " AND ")
	}
	err := DB.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询环境总数失败: %v", err)
	}

	// 查询数据
	query := "SELECT * FROM environments"
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)

	err = DB.SelectContext(ctx, &envs, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询环境列表失败: %v", err)
	}

	return envs, total, nil
}

// 更新环境
func UpdateEnvironment(ctx context.Context, env *model.Environment) error {
	query := `UPDATE environments SET 
		name = ?, code = ?, description = ?, status = ?, updated_at = NOW() 
		WHERE id = ?`

	_, err := DB.ExecContext(ctx, query,
		env.Name, env.Code, env.Description, env.Status, env.ID)
	if err != nil {
		return fmt.Errorf("更新环境失败: %v", err)
	}
	return nil
}

// 删除
func DeleteEnvironment(ctx context.Context, id int) error {
	query := "DELETE FROM environments WHERE id = ?"
	_, err := DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("删除环境失败: %v", err)
	}
	return nil
}
