package model

import (
	"time"
)

type TaskInput struct {
	TaskName      string `json:"taskName" db:"task_name"`
	CreatorUserID string `json:"creatorUserId" db:"creator_user_id"`
	Points        int64  `json:"point" db:"point"`
}

type Task struct {
	Id             int64     `json:"id" db:"id"`
	TaskName       string    `json:"taskName" db:"task_name"`
	CreatorUserID  string    `json:"creatorUserId" db:"creator_user_id"`
	Points         int64     `json:"point" db:"point"`
	CompilerUserId string    `json:"compilerUserId" db:"compiler_user_id"`
	CompleteDate   time.Time `json:"completeDate" db:"complete_date"`
	IsComplete     bool      `json:"isComplete" db:"is_complete"`
}

type CompleteTask struct {
	Id             int64     `json:"id" db:"id"`
	CompilerUserId string    `json:"compilerUserId" db:"compiler_user_id"`
	CompleteDate   time.Time `json:"completeDate" db:"complete_date"`
	IsComplete     bool      `json:"isComplete" db:"is_complete"`
}
