package models

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Profile struct {
	ID       int    `db:"id" json:"id"`
	Gender   int    `db:"gender" json:"gender"`
	Avatar   string `db:"avatar" json:"avatar"`
	Address  string `db:"address" json:"address"`
	Email    string `db:"email" json:"email"`
	UserId   int    `db:"userId" json:"userId"`
	NickName string `db:"nickName" json:"nickName"`
}

func (Profile) TableName() string {
	return "profile"
}

// 查询用户详情
func GetProfileByUserID(userId int) (*Profile, error) {
	query := `SELECT id, gender, avatar, address, email, userId, nickName FROM profile WHERE userId = ?`
	profile := &Profile{}
	err := DB.Get(profile, query, userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func GetProfilesByUserIDs(userIds []int) ([]Profile, error) {
	query := `SELECT * FROM profile WHERE userId IN (?)`

	query, args, err := sqlx.In(query, userIds)
	if err != nil {
		return nil, err
	}

	var profiles []Profile
	err = DB.Select(&profiles, query, args...)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func GetProfilesByUserIds(userIds []int) ([]Profile, error) {
	if len(userIds) == 0 {
		return nil, nil

	}
	// query := `SELECT * FROM profiles WHERE userId IN (?)`
	// query, args, _ := sqlx.In(query, userIds)
	// query = DB.Rebind(query)
	query, args, err := sqlx.In("SELECT * FROM profile WHERE userId IN (?)", userIds)
	if err != nil {
		return nil, err
	}
	var profiles []Profile
	err = DB.Select(&profiles, query, args...)
	if err != nil {
		return nil, err
	}
	return profiles, err
}

func GetProfilesByCondition(gender, enable, username string, pageNo, pageSize int) ([]*Profile, int64, error) {
	query := "SELECT p.* FROM profile p WHERE 1=1"
	// countQuery := "SELECT COUNT(*) FROM profile"
	var args []interface{}
	var total int64

	if gender != "" {
		whereCause := " AND gender = ?"
		query = query + whereCause
		// countQuery = countQuery + whereCause
		args = append(args, gender)
	}
	if enable != "" {
		whereCause := " AND p.userId IN (SELECT id FROM user WHERE enable = ?)"
		query = query + whereCause
		// countQuery = countQuery + whereCause
		args = append(args, enable)
	}
	if username != "" {
		whereCause := " AND nickName LIKE ?"
		query = query + whereCause
		// countQuery = countQuery + whereCause
		args = append(args, "%"+username+"%")
	}
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS t"
	fmt.Println("countQuery", countQuery)
	err := DB.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, errors.New("查询总数失败")
	}

	// 分页查询
	pageQuery := query + " LIMIT ? OFFSET ?"
	pageArgs := append(args, pageSize, (pageNo-1)*pageSize)
	fmt.Printf("Executing SQL: %s\nWith args: %v\n", pageQuery, pageArgs)
	var profileList []*Profile
	err = DB.Select(&profileList, pageQuery, pageArgs...)
	if err != nil {
		return nil, 0, errors.New("查询用户资料失败")
	}
	fmt.Println("profileList", profileList)
	return profileList, total, nil
}

func UpdateProfileByWhere(p Profile) error {
	query := `
        UPDATE profile
        SET gender = :gender, 
            address = :address, 
            email = :email, 
            nickName = :nickName
        WHERE id = :id
    `

	result, err := DB.NamedExec(query, map[string]interface{}{
		"id":       p.ID,
		"gender":   p.Gender,
		"address":  p.Address,
		"email":    p.Email,
		"nickName": p.NickName,
	})
	if err != nil {
		return err

	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("无法获取受影响行数")
	}

	if rowsAffected == 0 {
		return errors.New("未找到要更新的记录")
	}
	return nil
}
