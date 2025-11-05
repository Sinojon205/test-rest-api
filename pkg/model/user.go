package model

type User struct {
	Id       int64  `json:"id" db:"id"`
	FullName string `json:"fullName" db:"full_name"`
	Phone    string `json:"phone" db:"phone"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
}

type SignIn struct {
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
}

type UserStatus struct {
	User           User    `json:"user"`
	Points         int64   `json:"points"`
	CreatedTasks   []*Task `json:"createdTask"`
	CompletedTasks []*Task `json:"completedTask"`
}
type UsersWithPoints struct {
	User   User  `json:"user"`
	Points int64 `json:"points"`
}
