package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/logic"
)

func sendMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SendMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("sendMessage %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.FUid {
			appCtx.Logger().Errorf("sendMessage %d %d", requestUid, req.FUid)
			baseDto.ResponseForbidden(ctx)
			return
		}

		if rsp, err := l.SendMessage(req); err != nil {
			appCtx.Logger().Errorf("sendMessage %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("sendMessage %d %d", req.FUid, req.SId)
			baseDto.ResponseSuccess(ctx, rsp)
		}
	}
}

func getUserLatestMessages(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetMessageReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().Errorf("getUserLatestMessages %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Errorf("getUserLatestMessages %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}
		if resp, err := l.GetUserMessages(req); err != nil {
			appCtx.Logger().Errorf("getUserLatestMessages %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("getUserLatestMessages: %d, %d, %d, %d", req.CTime, req.UId, req.Count, req.Offset)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func deleteUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.DeleteMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Error("deleteUserMessage", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if len(req.MessageIds) == 0 && (req.TimeFrom == nil || req.TimeTo == nil) {
			appCtx.Logger().Errorf("deleteUserMessage %v, %d, %d", req.MessageIds, req.TimeFrom, req.TimeTo)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Errorf("deleteUserMessage %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}
		if err := l.DeleteUserMessage(&req); err != nil {
			appCtx.Logger().Errorf("deleteUserMessage %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("deleteUserMessage %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}
