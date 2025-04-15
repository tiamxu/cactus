package models

import "database/sql"

type Profile struct {
	ID       int    `json:"id"`
	Gender   int    `json:"gender"`
	Avatar   string `json:"avatar"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	UserId   int    `json:"userId"`
	NickName string `json:"nickName"`
}

func (Profile) TableName() string {
	return "profile"
}

func GetProfileByUserID(userId int) (*Profile, error) {
	query := `
        SELECT id, gender,avatar,address,email, userId, nickName 
        FROM profile 
        WHERE userId = ?`

	p := &Profile{}
	err := DB.QueryRow(query, userId).Scan(
		&p.ID,
		&p.Gender,
		&p.Avatar,
		&p.Address,
		&p.Email,
		&p.UserId,
		&p.NickName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 没有找到记录不算错误
		}
		return nil, err
	}
	return p, nil
}
