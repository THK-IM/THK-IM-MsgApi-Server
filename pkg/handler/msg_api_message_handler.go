package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msg-api-server/pkg/app"
	"github.com/thk-im/thk-im-msg-api-server/pkg/logic"
)

func sendMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SendMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.FUid {
			appCtx.Logger().Warn("permission error")
			dto.ResponseForbidden(ctx)
			return
		}

		if rsp, err := l.SendMessage(req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, rsp)
		}
	}
}

func getUserLatestMessages(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetMessageReq
		if err := ctx.BindQuery(&req); err != nil {
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
		if resp, err := l.GetUserMessages(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, resp)
		}
	}
}

func deleteUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.DeleteMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		if len(req.MessageIds) == 0 && (req.TimeFrom == nil || req.TimeTo == nil) {
			appCtx.Logger().Warn("param time_from or time_to or message_ids error")
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Warn("param uid error")
			dto.ResponseForbidden(ctx)
			return
		}
		if err := l.DeleteUserMessage(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}
