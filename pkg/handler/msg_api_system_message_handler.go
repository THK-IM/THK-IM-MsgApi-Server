package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/event"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/logic"
)

func pushMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.PushMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("pushMessage %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if req.Type > event.SignalExtended {
			if rsp, err := l.PushMessage(req, claims); err != nil {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("pushMessage %v", err)
				baseDto.ResponseInternalServerError(ctx, err)
			} else {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("pushMessage %v %v", req, rsp)
				baseDto.ResponseSuccess(ctx, rsp)
			}
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("pushMessage %v", req)
			baseDto.ResponseBadRequest(ctx)
			return
		}
	}
}

func sendSessionMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.SendMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("sendSystemMessage %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if rsp, err := l.SendMessage(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("sendSystemMessage %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("sendSystemMessage %v %v", req, rsp)
			baseDto.ResponseSuccess(ctx, rsp)
		}
	}
}

func sendSystemMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.SendSysMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("sendSystemMessage %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if rsp, err := l.SendSysMessage(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("sendSystemMessage %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("sendSystemMessage %v %v", req, rsp)
			baseDto.ResponseSuccess(ctx, rsp)
		}
	}
}
