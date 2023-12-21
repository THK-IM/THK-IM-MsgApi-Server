package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/event"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/logic"
)

func pushMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.PushMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("pushMessage %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if req.Type > event.SignalExtended {
			if rsp, err := l.PushMessage(req); err != nil {
				appCtx.Logger().Errorf("pushMessage %v", err)
				baseDto.ResponseInternalServerError(ctx, err)
			} else {
				appCtx.Logger().Errorf("pushMessage %v %v", req, rsp)
				baseDto.ResponseSuccess(ctx, rsp)
			}
		} else {
			appCtx.Logger().Errorf("pushMessage %v", req)
			baseDto.ResponseBadRequest(ctx)
			return
		}
	}
}

func sendSessionMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SendMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("sendSystemMessage %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if rsp, err := l.SendMessage(req); err != nil {
			appCtx.Logger().Errorf("sendSystemMessage %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("sendSystemMessage %v %v", req, rsp)
			baseDto.ResponseSuccess(ctx, rsp)
		}
	}
}

func sendSystemMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SendSysMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("sendSystemMessage %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if rsp, err := l.SendSysMessage(req); err != nil {
			appCtx.Logger().Errorf("sendSystemMessage %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("sendSystemMessage %v %v", req, rsp)
			baseDto.ResponseSuccess(ctx, rsp)
		}
	}
}
