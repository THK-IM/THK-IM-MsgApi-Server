package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/logic"
)

func updateUserOnlineStatus(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.PostUserOnlineReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("updateUserOnlineStatus %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if err := l.UpdateUserOnlineStatus(&req); err != nil {
			appCtx.Logger().Errorf("updateUserOnlineStatus %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("updateUserOnlineStatus %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func queryUserOnlineStatus(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.QueryUsersOnlineStatusReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Errorf("queryUsersOnlineStatus %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if res, err := l.GetUsersOnlineStatus(req.UIds); err != nil {
			appCtx.Logger().Errorf("queryUsersOnlineStatus %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("queryUsersOnlineStatus %v %v", req, res)
			baseDto.ResponseSuccess(ctx, res)
		}
	}
}

func kickOffUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.KickUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Errorf("kickOffUser %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if err := l.KickUser(&req); err != nil {
			appCtx.Logger().Errorf("kickOffUser %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("kickOffUser %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}

	}
}
