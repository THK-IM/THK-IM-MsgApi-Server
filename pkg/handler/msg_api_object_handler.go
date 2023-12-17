package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseErrorx "github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/logic"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func getObjectUploadParams(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionObjectLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetUploadParamsReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().Errorf("getObjectUploadParams %v %s", req, err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Errorf("getObjectUploadParams %v %d", req, requestUid)
			baseDto.ResponseForbidden(ctx)
			return
		}

		res, err := l.GetUploadParams(req)
		if err != nil {
			appCtx.Logger().Errorf("getObjectUploadParams %v %d", req, requestUid)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("getObjectUploadParams %v %v", req, res)
			baseDto.ResponseSuccess(ctx, res)
		}
	}
}

func getObjectDownloadUrl(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionObjectLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetDownloadUrlReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().Errorf("getObjectDownloadUrl %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		req.UId = requestUid

		path, err := l.GetObjectByKey(req)
		if err != nil {
			appCtx.Logger().Errorf("getObjectDownloadUrl %v", err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			if path != nil {
				appCtx.Logger().Infof("getObjectDownloadUrl %s", *path)
				baseDto.Redirect302(ctx, *path)
			} else {
				appCtx.Logger().Errorf("getObjectDownloadUrl %s", "path is nil")
				baseDto.ResponseInternalServerError(ctx, baseErrorx.ErrInternalServerError)
			}
		}
	}
}
