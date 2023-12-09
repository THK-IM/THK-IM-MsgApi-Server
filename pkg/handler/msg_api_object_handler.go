package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msg-api-server/pkg/app"
	"github.com/thk-im/thk-im-msg-api-server/pkg/errorx"
	"github.com/thk-im/thk-im-msg-api-server/pkg/logic"
)

func getObjectUploadParams(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionObjectLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetUploadParamsReq
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

		res, err := l.GetUploadParams(req)
		if err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, res)
		}
	}
}

func getObjectDownloadUrl(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionObjectLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetDownloadUrlReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(middleware.UidKey)
		req.UId = requestUid

		path, err := l.GetObjectByKey(req)
		if err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			if path != nil {
				dto.Redirect302(ctx, *path)
			} else {
				appCtx.Logger().Warn(err.Error())
				dto.ResponseInternalServerError(ctx, errorx.ErrServerUnknown)
			}
		}
	}
}
