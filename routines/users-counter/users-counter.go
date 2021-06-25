package userscounter

import (
	"context"
	"example/ctx"
	"example/modules/user"
	"time"

	appctx "github.com/nixys/nxs-go-appctx/v2"

	"github.com/sirupsen/logrus"
)

// Runtime executes the routine
func Runtime(cr context.Context, appCtx *appctx.AppContext, crc chan interface{}) {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	timer := time.NewTimer(time.Duration(cc.Conf.CounterInterval) * time.Second)

	if usersCounterIterate(appCtx) != nil {
		appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
		return
	}

	for {
		select {
		case <-timer.C:
			// Do the some actions
			if usersCounterIterate(appCtx) != nil {
				appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
				return
			}
			timer.Reset(time.Duration(cc.Conf.CounterInterval) * time.Second)
		case <-cr.Done():
			// Program termination.
			appCtx.Log().Info("counter routine done")
			return
		case <-crc:
			// Updated context application data.
			// Set the new one in current goroutine.
			appCtx.Log().Info("counter routine reload")
		}
	}
}

func usersCounterIterate(appCtx *appctx.AppContext) error {

	users, err := user.Count(appCtx)
	if err != nil {
		appCtx.Log().Errorf("users count error: %v", err)
		return err
	}

	appCtx.Log().WithFields(logrus.Fields{
		"amount of users": users,
	}).Info("user counter iteration successfully finished")

	return nil
}
