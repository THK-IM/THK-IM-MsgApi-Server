package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msg-api-server/pkg/app"
	"github.com/thk-im/thk-im-msg-api-server/pkg/logic"
)

func ackUserMessages(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.AckUserMessagesReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		if len(req.MsgIds) == 0 {
			appCtx.Logger().Warn("param msgIds error")
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Warn("param uid error")
			dto.ResponseForbidden(ctx)
			return
		}
		if err := l.AckUserMessages(req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func readUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.ReadUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Warn("param uid error")
			dto.ResponseForbidden(ctx)
			return
		}
		if err := l.ReadUserMessages(req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func revokeUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.RevokeUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Warn("param uid error")
			dto.ResponseForbidden(ctx)
			return
		}
		if err := l.RevokeUserMessage(req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func reeditUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.ReeditUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Warn("param uid error")
			dto.ResponseForbidden(ctx)
			return
		}
		if err := l.ReeditUserMessage(req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func forwardUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.ForwardUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.FUid {
			appCtx.Logger().Warn("param uid error")
			dto.ResponseForbidden(ctx)
			return
		}

		// 鉴权
		su, errSu := appCtx.SessionUserModel().FindSessionUser(req.ForwardSId, req.FUid)
		if errSu != nil {
			appCtx.Logger().Warn("session user err: ", errSu)
			dto.ResponseForbidden(ctx)
			return
		}
		if su.UserId <= 0 {
			appCtx.Logger().Warn("session user not existed")
			dto.ResponseForbidden(ctx)
			return
		}

		if resp, err := l.ForwardUserMessages(req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, resp)
		}
	}
}
