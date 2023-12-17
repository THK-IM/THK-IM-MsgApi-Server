package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/logic"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func ackUserMessages(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.AckUserMessagesReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("ackUserMessages %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if len(req.MsgIds) == 0 {
			appCtx.Logger().Errorf("ackUserMessages %v", req.MsgIds)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Errorf("ackUserMessages %d %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}
		if err := l.AckUserMessages(req); err != nil {
			appCtx.Logger().Errorf("ackUserMessages %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("ackUserMessages %d, %d, %v", req.UId, req.SId, req.MsgIds)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func readUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.ReadUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("readUserMessage %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Errorf("readUserMessage %d %d", requestUid, req.UId)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if err := l.ReadUserMessages(req); err != nil {
			appCtx.Logger().Errorf("readUserMessage %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("readUserMessage %d, %d, %v", req.UId, req.SId, req.MsgIds)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func revokeUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.RevokeUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("revokeUserMessage %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Errorf("revokeUserMessage %d %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}
		if err := l.RevokeUserMessage(req); err != nil {
			appCtx.Logger().Errorf("revokeUserMessage %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("revokeUserMessage %d, %d, %d", req.UId, req.SId, req.MsgId)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func reeditUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.ReeditUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("reeditUserMessage %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Errorf("reeditUserMessage %d %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}
		if err := l.ReeditUserMessage(req); err != nil {
			appCtx.Logger().Errorf("reeditUserMessage %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("reeditUserMessage %d, %d, %d", req.UId, req.SId, req.MsgId)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func forwardUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.ForwardUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("forwardUserMessage %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.FUid {
			appCtx.Logger().Errorf("forwardUserMessage %d %d", requestUid, req.FUid)
			baseDto.ResponseForbidden(ctx)
			return
		}

		// 鉴权
		su, errSu := appCtx.SessionUserModel().FindSessionUser(req.ForwardSId, req.FUid)
		if errSu != nil {
			appCtx.Logger().Errorf("forwardUserMessage %s", errSu.Error())
			baseDto.ResponseForbidden(ctx)
			return
		}
		if su.UserId <= 0 {
			appCtx.Logger().Errorf("forwardUserMessage %d", su.UserId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		if resp, err := l.ForwardUserMessages(req); err != nil {
			appCtx.Logger().Errorf("forwardUserMessage %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("reeditUserMessage %d, %v, %v", req.ForwardSId, req.ForwardFromUIds, req.ForwardClientIds)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}
