package user

import (
	"errors"

	"example/ctx"
	"example/db/mysql"
	"example/misc"

	appctx "github.com/nixys/nxs-go-appctx/v2"
)

// Record defines user data
type Record struct {
	ID       int64
	Username string
	Password string
}

// CreateData defines data to create new user record
type CreateData struct {
	Username string
}

var (
	// ErrUsernameEmpty defines an error when username is empty
	ErrUsernameEmpty = errors.New("empty username")
)

// GetAll retrives all user records
func GetAll(appCtx *appctx.AppContext) ([]Record, error) {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	usrs, err := cc.MySQL.UsersGet()
	if err != nil {
		return []Record{}, err
	}

	users := []Record{}
	for _, u := range usrs {
		users = append(users, Record{
			ID:       u.ID,
			Username: u.Username,
			Password: u.Password,
		})
	}

	return users, nil
}

// Create creates new user
func Create(appCtx *appctx.AppContext, user CreateData) (Record, error) {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	// Check username is not empty
	if len(user.Username) == 0 {
		return Record{}, ErrUsernameEmpty
	}

	// Make DB request to create new user record
	usr, err := cc.MySQL.UserInsert(
		mysql.UserInsertData{
			Username: user.Username,
			Password: misc.TokenGen(15),
		})
	if err != nil {
		return Record{}, err
	}

	return Record{
		ID:       usr.ID,
		Username: usr.Username,
		Password: usr.Password,
	}, nil
}

// Count counts users within the DB
func Count(appCtx *appctx.AppContext) (int64, error) {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	usrs, err := cc.MySQL.UsersGet()
	if err != nil {
		return 0, err
	}

	return int64(len(usrs)), nil
}
