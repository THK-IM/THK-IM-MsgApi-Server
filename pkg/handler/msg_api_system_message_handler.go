package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/event"
	"github.com/thk-im/thk-im-msg-api-server/pkg/app"
	"github.com/thk-im/thk-im-msg-api-server/pkg/logic"
)

func pushExtendedMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.PushMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		if req.Type > event.SignalExtended {
			if rsp, err := l.PushMessage(req); err != nil {
				appCtx.Logger().Warn(err.Error())
				dto.ResponseInternalServerError(ctx, err)
			} else {
				dto.ResponseSuccess(ctx, rsp)
			}
		} else {
			appCtx.Logger().Warn("param type error")
			dto.ResponseBadRequest(ctx)
			return
		}
	}
}

func sendSystemMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SendMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
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
