package repository

import (
	"context"
	"database/sql"
	"errors"
	"test-rest-api/pkg/model"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type Repository struct {
	db *sqlx.DB
}

func (r *Repository) AddRefferer(id int64) error {
	query := "INSERT INTO task(task_name, creator_user_id, point,compiler_user_id, complete_date, is_complete) VALUES($1,$2,$3,$4,$5,$6)"

	tx, err := r.db.Beginx()
	if err != nil {
		err = tx.Rollback()
		return err
	}
	_, err = tx.Exec(query, "Adding refferer", id, 1, id, time.Now().Format(time.RFC1123), true)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = "SELECT * FROM points WHERE user_id=$1"
	var p model.Points
	err = tx.Get(&p, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		query = "INSERT INTO points(user_id, point) VALUES($1,$2)"
		_, err = tx.Exec(query, id, 1)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else if err != nil {
		err = tx.Rollback()
		return err
	} else {
		p.Points += 1
		query = "UPDATE points SET user_id=$1, point=$2 WHERE id=$3"
		_, err = tx.Exec(query, p.UserID, p.Points, p.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(a model.User) (int64, error) {
	query := "INSERT INTO users(full_name,email,phone,password) VALUES($1,$2,$3,$4)  RETURNING id"
	lastInsertId := int64(0)
	err := r.db.QueryRow(query, a.FullName, a.Email, a.Phone, a.Password).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (r *Repository) AddTask(a *model.TaskInput) error {
	query := "INSERT INTO task(task_name, creator_user_id, point) VALUES($1,$2,$3)"

	_, err := r.db.Exec(query, a.TaskName, a.CreatorUserID, a.Points)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ComplitTask(userId int64, inputTask *model.CompleteTask) error {
	query := "SELECT * from tasks WHERE id=$1"

	tx, err := r.db.Beginx()
	if err != nil {
		err = tx.Rollback()
		return err
	}
	var t model.Task
	err = tx.Get(&t, query, inputTask.Id)
	if err != nil {
		err = tx.Rollback()
		return err
	}

	query = "UPDATE tasks SET compiler_user_id=$1, complete_date=$2, is_complete=$3 WHERE id=$4"
	_, err = tx.Exec(query, inputTask.CompilerUserId, inputTask.CompleteDate, true, inputTask.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = "SELECT * FROM points WHERE user_id=$1"
	var p model.Points
	err = tx.Get(&p, query, inputTask.CompilerUserId)
	if errors.Is(err, sql.ErrNoRows) {
		query = "INSERT INTO points(user_id, point) VALUES($1,$2)"
		_, err = tx.Exec(query, inputTask.CompilerUserId, t.Points)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else if err != nil {
		err = tx.Rollback()
		return err
	} else {
		p.Points += t.Points
		query = "UPDATE points SET user_id=$1, point=$2 WHERE id=$3"
		_, err = tx.Exec(query, p.UserID, p.Points, p.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) RemoveTask(id int64) error {
	query := "SELECT *  FROM tasks WHERE id=$1"

	tx, err := r.db.Beginx()
	if err != nil {
		err = tx.Rollback()
		return err
	}
	var t model.Task
	err = tx.Get(&t, query, id)
	if err != nil {
		err = tx.Rollback()
		return err
	}

	query = "DELETE  from tasks where id=$1"
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = "SELECT * FROM points WHERE user_id=$1"
	var p model.Points
	err = tx.Get(&p, query, id)
	if err != nil {
		err = tx.Rollback()
		return err
	} else {
		p.Points -= t.Points
		query = "UPDATE points SET user_id=$1, point=$2 WHERE id=$3"
		_, err = tx.Exec(query, p.UserID, p.Points, p.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) UpdateUser(user model.User) (int64, error) {
	query := "UPDATE users  SET full_name=$1,email=$2, password=$3 phone=$4 WHERE id=$5"

	res, err := r.db.ExecContext(getContext(), query, user.FullName, user.Email, user.Password, user.Phone, user.Id)
	if err != nil {
		return 0, err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (r *Repository) GetUserByEmail(email string) (*model.User, error) {
	query := "SELECT * FROM users WHERE email=$1"
	var user model.User
	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *Repository) GetUser(email, password string) (*model.User, error) {
	query := "SELECT * FROM users WHERE email=$1 and password = $2"
	var user model.User
	err := r.db.Get(&user, query, email, password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserStatus(id int64) (*model.UserStatus, error) {
	var res model.UserStatus
	user, err := r.getUserById(id)
	if err != nil {
		return nil, err
	}
	res.User = *user

	query := "SELECT * FROM tasks WHERE creator_user_id=$1"
	var createdTasks []*model.Task
	err = r.db.Get(&createdTasks, query, id)
	if err != nil {
		return nil, err
	}
	res.CreatedTasks = createdTasks

	query = "SELECT * FROM tasks WHERE compiler_user_id=$1"
	var compiledTasks []*model.Task
	err = r.db.Get(&compiledTasks, query, id)
	if err != nil {
		return nil, err
	}
	res.CreatedTasks = compiledTasks

	query = "SELECT * FROM tasks WHERE compiler_user_id=$1"
	var point model.Points
	err = r.db.Get(&point, query, id)
	if err != nil {
		return nil, err
	}
	res.Points = point.Points

	return &res, nil

}

func (r *Repository) GetLeaders() ([]*model.UsersWithPoints, error) {

	//first way
	type userWithPoints struct {
		CompilerUserId int64 `json:"compilerUserId" db:"compiler_user_id"`
		SumPoints      int64 `json:"sumPoints" db:"sum_points"`
	}
	query := "SELECT compiler_user_id,sum(points) as sum_points FROM tasks WHERE compiler_user_id=$1 GROUP BY compiler_user_id ORDER BY DESC sum_points Limit 10"

	var uwp []*userWithPoints
	err := r.db.Get(&uwp, query)
	if err != nil {
		return nil, err
	}
	n := len(uwp)
	userIds := make([]int64, n, n)
	for i, u := range uwp {
		userIds[i] = u.CompilerUserId
	}

	query = "SELECT * FROM users WHERE id in($1)"
	var users []*model.User
	err = r.db.Get(&users, query, userIds)
	if err != nil {
		return nil, err
	}
	var res = make([]*model.UsersWithPoints, n, n)
	for i, u := range users {
		res[i].User = *u
		w, ok := lo.Find(uwp, func(wp *userWithPoints) bool {
			return wp.CompilerUserId == u.Id
		})
		if ok {
			res[i].Points = w.SumPoints
		}
	}
	return res, nil

	// second way
	// query := "SELECT * FROM points  ORDER BY DESC points Limit 10"

	// var uwp []*model.Points
	// err := r.db.Get(&uwp, query)
	// if err != nil {
	// 	return nil, err
	// }
	// n := len(uwp)
	// userIds := make([]int64, n, n)
	// for i, u := range uwp {
	// 	userIds[i] = u.UserID
	// }

	// query = "SELECT * FROM users WHERE id in($1)"
	// var users []*model.User
	// err = r.db.Get(&users, query, userIds)
	// if err != nil {
	// 	return nil, err
	// }
	// var res = make([]*model.UsersWithPoints, n, n)
	// for i, u := range users {
	// 	res[i].User = *u
	// 	w, ok := lo.Find(uwp, func(wp *model.Points) bool {
	// 		return wp.UserID == u.Id
	// 	})
	// 	if ok {
	// 		res[i].Points = w.Points
	// 	}
	// }
	// return res, nil
}

func (r *Repository) getUserById(id int64) (*model.User, error) {
	query := "SELECT * FROM users WHERE id=$1"
	var user model.User
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}
