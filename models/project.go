package models

import "time"

type Project struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	DeleteAt    time.Time `db:"deleted_at" json:"deleted_at"`
	Tags        []string  `db:"tags"`
}

func Create1() (*Project, error) {
	project := &Project{}
	query := `
	INSERT INTO project (id, name, description,  status, created_at, updated_at, tags)
	VALUES (:id, :name, :description, :status, :created_at, :updated_at, :tags)`
	_, err := DB.NamedExec(query, project)
	if err != nil {

	}
	return project, nil
}
