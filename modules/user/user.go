package user

import (
	"fmt"

	"github.com/nixys/nxs-go-appctx-example-restapi/ds/primedb"
	"github.com/nixys/nxs-go-appctx-example-restapi/misc"
)

type User struct {
	db primedb.DB
}

type Data struct {
	ID       int64
	Username string
	Password string
}

type CreateData struct {
	Username string
}

type UpdateData struct {
	ID       int64
	Username *string
	Password *string
}

type InitSettings struct {
	DB primedb.DB
}

func Init(s InitSettings) *User {
	return &User{
		db: s.DB,
	}
}

func (usr *User) Create(d CreateData) (Data, error) {

	u, err := usr.db.UserInsert(primedb.UserInsertData{
		Username: d.Username,
		Password: misc.TokenGen(15),
	})
	if err != nil {
		return Data{}, fmt.Errorf("user create: %w", err)
	}

	return Data{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

func (usr *User) Get(id int64) (Data, error) {

	d, err := usr.db.UserGet(id)
	if err != nil {
		return Data{}, fmt.Errorf("user get: %w", err)
	}

	return Data{
		ID:       d.ID,
		Username: d.Username,
		Password: d.Password,
	}, nil
}

func (usr *User) GetAll() ([]Data, error) {

	var users []Data

	ds, err := usr.db.UsersGet()
	if err != nil {
		return nil, fmt.Errorf("users get all: %w", err)
	}

	for _, d := range ds {
		users = append(
			users,
			Data{
				ID:       d.ID,
				Username: d.Username,
				Password: d.Password,
			},
		)
	}

	return users, nil
}

func (usr *User) Update(d UpdateData) (Data, error) {

	if d.ID <= 0 {
		return Data{}, fmt.Errorf("user update: %w", misc.ErrNotFound)
	}

	u, err := usr.db.UserUpdate(primedb.UserUpdateData{
		ID:       d.ID,
		Username: d.Username,
		Password: d.Password,
	})
	if err != nil {
		return Data{}, fmt.Errorf("user update: %w", err)
	}

	return Data{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}, nil
}

func (usr *User) Delete(id int64) error {

	if err := usr.db.UserDelete(id); err != nil {
		return fmt.Errorf("user delete: %w", err)
	}

	return nil
}
