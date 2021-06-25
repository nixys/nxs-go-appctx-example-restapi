package mysql

import (
	"database/sql"
)

// User contains data from `user` table
type User struct {
	ID       int64
	Username string
	Password string
}

// UserInsertData contains data to create new user record
type UserInsertData struct {
	Username string
	Password string
}

// UserInsert inserts new user record
func (m *MySQL) UserInsert(user UserInsertData) (User, error) {

	s, err := m.client.Prepare("INSERT INTO `users` (`username`, `password`) VALUES(?, ?)")
	if err != nil {
		return User{}, err
	}
	defer s.Close()

	r, err := s.Exec(
		user.Username,
		user.Password,
	)
	if err != nil {
		return User{}, err
	}

	id, _ := r.LastInsertId()

	return User{
		ID:       id,
		Username: user.Username,
		Password: user.Password,
	}, nil
}

// UsersGet gets all users
func (m *MySQL) UsersGet() ([]User, error) {

	type table struct {
		ID       sql.NullInt64  `db:"id"`
		Username sql.NullString `db:"username"`
		Password sql.NullString `db:"password"`
	}

	var users []User

	t := []table{}

	err := m.client.Select(&t, "SELECT `id`, `username`,`password` FROM `users`")
	if err != nil {
		return users, err
	}

	for _, u := range t {
		users = append(users, User{
			ID:       u.ID.Int64,
			Username: u.Username.String,
			Password: u.Password.String,
		})
	}

	return users, nil
}
