package db

import "database/sql"

type Repo struct {
	*Queries
	DB *sql.DB
}

func NewRepo(dtb *sql.DB) *Repo {
	return &Repo{
		DB:      dtb,
		Queries: New(dtb),
	}
}
