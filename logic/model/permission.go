package model

type Permission struct {
	ID          int          `db:"id" json:"id"`
	Name        string       `db:"name" json:"name"`
	Code        string       `db:"code" json:"code"`
	Type        string       `db:"type" json:"type"`
	ParentId    *int         `db:"parentId" json:"parentId"`
	Path        string       `db:"path" json:"path"`
	Redirect    string       `db:"redirect" json:"redirect"`
	Icon        string       `db:"icon" json:"icon"`
	Component   string       `db:"component" json:"component"`
	Layout      string       `db:"layout" json:"layout"`
	KeepAlive   int          `db:"keepAlive" json:"keepAlive"`
	Method      string       `db:"method" json:"method"`
	Description string       `db:"description" json:"description"`
	Show        int          `db:"show" json:"show"`
	Enable      int          `db:"enable" json:"enable"`
	Order       int          `db:"order" json:"order"`
	Children    []Permission `db:"children" json:"children"`
}

func (Permission) TableName() string {
	return "permission"
}
