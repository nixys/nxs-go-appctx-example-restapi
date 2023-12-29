package endpoints

import (
	"context"

	"github.com/nixys/nxs-go-appctx-example-restapi/modules/user"
)

type userTx struct {
	User userDataTx `json:"user"`
}

type usersTx struct {
	Users []userDataTx `json:"users"`
}

type userDataTx struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserData(c context.Context, rd interface{}, message string) interface{} {

	if len(message) > 0 {
		return nil
	}

	user, b := rd.(user.Data)
	if b == false {
		return nil
	}

	return userTx{
		User: userDataTx{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
		},
	}
}

func UsersData(c context.Context, rd interface{}, message string) interface{} {

	if len(message) > 0 {
		return nil
	}

	users, b := rd.([]user.Data)
	if b == false {
		return nil
	}

	return usersTx{
		Users: func() []userDataTx {
			us := []userDataTx{}
			for _, u := range users {
				us = append(
					us,
					userDataTx{
						ID:       u.ID,
						Username: u.Username,
						Password: u.Password,
					},
				)
			}
			return us
		}(),
	}
}
