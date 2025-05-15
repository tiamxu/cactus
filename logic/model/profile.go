package model

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
