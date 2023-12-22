package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/logic"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func sendMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.SendMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("sendMessage %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.FUid {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("sendMessage %d %d", requestUid, req.FUid)
			baseDto.ResponseForbidden(ctx)
			return
		}

		if rsp, err := l.SendMessage(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("sendMessage %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("sendMessage %d %d", req.FUid, req.SId)
			baseDto.ResponseSuccess(ctx, rsp)
		}
	}
}

func getUserLatestMessages(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.GetMessageReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getUserLatestMessages %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getUserLatestMessages %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}
		if resp, err := l.GetUserMessages(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getUserLatestMessages %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("getUserLatestMessages: %d, %d, %d, %d", req.CTime, req.UId, req.Count, req.Offset)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func deleteUserMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.DeleteMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Error("deleteUserMessage", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if len(req.MessageIds) == 0 && (req.TimeFrom == nil || req.TimeTo == nil) {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteUserMessage %v, %d, %d", req.MessageIds, req.TimeFrom, req.TimeTo)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteUserMessage %d, %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}
		if err := l.DeleteUserMessage(&req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteUserMessage %v %s", req, err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("deleteUserMessage %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}
