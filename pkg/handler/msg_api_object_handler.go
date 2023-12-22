package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseErrorx "github.com/thk-im/thk-im-base-server/errorx"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/logic"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func getObjectUploadParams(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionObjectLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.GetUploadParamsReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getObjectUploadParams %v %s", req, err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getObjectUploadParams %v %d", req, requestUid)
			baseDto.ResponseForbidden(ctx)
			return
		}

		res, err := l.GetUploadParams(req, claims)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getObjectUploadParams %v %d", req, requestUid)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("getObjectUploadParams %v %v", req, res)
			baseDto.ResponseSuccess(ctx, res)
		}
	}
}

func getObjectDownloadUrl(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionObjectLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.GetDownloadUrlReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getObjectDownloadUrl %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		req.UId = requestUid

		path, err := l.GetObjectByKey(req, claims)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getObjectDownloadUrl %v", err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			if path != nil {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("getObjectDownloadUrl %s", *path)
				baseDto.Redirect302(ctx, *path)
			} else {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getObjectDownloadUrl %s", "path is nil")
				baseDto.ResponseInternalServerError(ctx, baseErrorx.ErrInternalServerError)
			}
		}
	}
}
