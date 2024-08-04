package model

import "time"

type Log struct {
	Action    string    `db:"action"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}
