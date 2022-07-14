package user

import (
	"example/ctx"
	"example/ds/mysql"
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

// UpdateData defines data to update user record
type UpdateData struct {
	ID       int64
	Username *string
	Password *string
}

// Create creates new user
func Create(appCtx *appctx.AppContext, user CreateData) (Record, error) {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	// Check username is not empty
	if len(user.Username) == 0 {
		return Record{}, misc.ErrUsernameEmpty
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

// Get retrives specified user records
func Get(appCtx *appctx.AppContext, id int64) (Record, error) {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	usr, err := cc.MySQL.UserGet(id)
	if err != nil {
		return Record{}, err
	}

	return Record{
		ID:       usr.ID,
		Username: usr.Username,
		Password: usr.Password,
	}, nil
}

// Update updates user
func Update(appCtx *appctx.AppContext, user UpdateData) (Record, error) {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	// Check id specified
	if user.ID == 0 {
		return Record{}, misc.ErrIDEmpty
	}

	// Make DB request to update user record
	usr, err := cc.MySQL.UserUpdate(mysql.UserUpdateData{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
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

// Delete deletes user
func Delete(appCtx *appctx.AppContext, id int64) error {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	// Check id specified
	if id == 0 {
		return misc.ErrIDEmpty
	}

	// Make DB request to delete user record
	return cc.MySQL.UserDelete(id)
}

// Count counts users within the DB
func Count(appCtx *appctx.AppContext) (int64, error) {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	count, err := cc.MySQL.UsersCount()
	if err != nil {
		return 0, err
	}

	return count, nil
}
