package models

import "github.com/jmoiron/sqlx"

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
