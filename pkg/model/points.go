package model

type Points struct {
	Id     int64 `json:"id" db:"id"`
	UserID int64 `json:"userId" db:"user_id"`
	Points int64 `json:"point" db:"point"`
}
