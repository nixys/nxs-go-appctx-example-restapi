package mysql

import "example/misc"

const usersTableName = "users"

// User contains data from `user` table
type User struct {
	ID       int64  `json:"id" gorm:"primaryKey"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// UserInsertData contains data to create new user record
type UserInsertData struct {
	Username string
	Password string
}

// UserUpdateData contains data to update user record
type UserUpdateData struct {
	ID       int64   `json:"id" gorm:"->"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (User) TableName() string {
	return usersTableName
}

func (UserInsertData) TableName() string {
	return usersTableName
}

func (UserUpdateData) TableName() string {
	return usersTableName
}

// UserInsert inserts new user record
func (m *MySQL) UserInsert(user UserInsertData) (User, error) {

	u := User{
		Username: user.Username,
		Password: user.Password,
	}

	r := m.client.
		Create(&u)
	if r.Error != nil {
		return User{}, r.Error
	}

	return u, nil
}

// UsersGet gets all users
func (m *MySQL) UsersGet() ([]User, error) {

	users := []User{}

	r := m.client.
		Find(&users)
	if r.Error != nil {
		return []User{}, r.Error
	}

	return users, nil
}

// UserGet gets specified users
func (m *MySQL) UserGet(id int64) (User, error) {

	user := User{
		ID: id,
	}

	r := m.client.
		Find(&user)
	if r.Error != nil {
		return User{}, r.Error
	}

	if r.RowsAffected == 0 {
		return User{}, misc.ErrNotFound
	}

	return user, nil
}

// UsersCount gets users count
func (m *MySQL) UsersCount() (int64, error) {

	var count int64

	r := m.client.
		Model(User{}).
		Count(&count)
	if r.Error != nil {
		return 0, r.Error
	}

	return count, nil
}

// UserUpdate updates user record
func (m *MySQL) UserUpdate(user UserUpdateData) (User, error) {

	u := User{}

	r := m.client.
		Model(UserUpdateData{
			ID: user.ID,
		}).
		Updates(user).
		Find(&u)
	if r.Error != nil {
		return User{}, r.Error
	}

	if r.RowsAffected == 0 {
		return User{}, misc.ErrNotFound
	}

	return u, nil
}

func (m *MySQL) UserDelete(id int64) error {

	r := m.client.
		Delete(User{
			ID: id,
		})
	if r.Error != nil {
		return r.Error
	}

	if r.RowsAffected == 0 {
		return misc.ErrNotFound
	}

	return nil
}
