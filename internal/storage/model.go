package storage

import "time"

type Entity struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Archived  bool      `db:"archived"`
}
