package primedb

import (
	"fmt"

	"github.com/nixys/nxs-go-appctx-example-restapi/misc"
)

const userTableName = "user"

type User struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

type UserInsertData struct {
	Username string
	Password string
}

type UserUpdateData struct {
	ID       int64   `gorm:"->;column:id"`
	Username *string `gorm:"column:username"`
	Password *string `gorm:"column:password"`
}

func (User) TableName() string {
	return userTableName
}

func (UserInsertData) TableName() string {
	return userTableName
}

func (UserUpdateData) TableName() string {
	return userTableName
}

func (db *DB) UserInsert(user UserInsertData) (User, error) {

	u := User{
		Username: user.Username,
		Password: user.Password,
	}

	r := db.client.
		Create(&u)
	if r.Error != nil {
		return User{}, fmt.Errorf("primedb user insert: %w", r.Error)
	}

	return u, nil
}

func (db *DB) UsersGet() ([]User, error) {

	users := []User{}

	r := db.client.
		Find(&users)
	if r.Error != nil {
		return []User{}, fmt.Errorf("primedb users list: %w", r.Error)
	}

	return users, nil
}

func (db *DB) UserGet(id int64) (User, error) {

	user := User{
		ID: id,
	}

	r := db.client.
		Find(&user)
	if r.Error != nil {
		return User{}, fmt.Errorf("primedb user get: %w", r.Error)
	}

	if r.RowsAffected == 0 {
		return User{}, fmt.Errorf("primedb user get: %w", misc.ErrNotFound)
	}

	return user, nil
}

func (db *DB) UserUpdate(user UserUpdateData) (User, error) {

	u := User{}

	r := db.client.
		Model(UserUpdateData{
			ID: user.ID,
		}).
		Updates(user).
		Find(&u)
	if r.Error != nil {
		return User{}, fmt.Errorf("primedb user update: %w", r.Error)
	}

	if r.RowsAffected == 0 {
		return User{}, fmt.Errorf("primedb user update: %w", misc.ErrNotFound)
	}

	return u, nil
}

func (db *DB) UserDelete(id int64) error {

	r := db.client.
		Delete(User{
			ID: id,
		})
	if r.Error != nil {
		return fmt.Errorf("primedb user delete: %w", r.Error)
	}

	if r.RowsAffected == 0 {
		return fmt.Errorf("primedb user delete: %w", misc.ErrNotFound)
	}

	return nil
}
