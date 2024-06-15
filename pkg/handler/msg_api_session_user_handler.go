package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/logic"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
	"strconv"
)

func getLatestSessionUsers(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.QuerySessionUsersReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getLatestSessionUsers %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if req.Count <= 0 {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getLatestSessionUsers %v", req)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if req.Role != nil && (*req.Role > model.SessionOwner || *req.Role < model.SessionMember) {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getLatestSessionUsers %v", req)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getLatestSessionUsers %v", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		req.SId = sessionId

		requestUid := ctx.GetInt64(userSdk.UidKey)
		appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("getLatestSessionUsers check permission %d %d ", requestUid, sessionId)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkReadPermission(appCtx, requestUid, sessionId, claims); !hasPermission {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getLatestSessionUsers %d %d ", requestUid, sessionId)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}
		if resp, err := l.QueryLatestSessionUsers(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getLatestSessionUsers %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("getLatestSessionUsers %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func getSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionUser %v", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		userId, errUserId := strconv.ParseInt(ctx.Param("uid"), 10, 64)
		if errUserId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionUser %v", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkReadPermission(appCtx, requestUid, sessionId, claims); !hasPermission {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionUser %d %d ", requestUid, sessionId)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}
		if resp, err := l.QuerySessionUser(sessionId, userId, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionUser %d %d %v", sessionId, userId, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("getSessionUser %d %d %v", sessionId, userId, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func getSessionUserCount(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionUserCount %v", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkReadPermission(appCtx, requestUid, sessionId, claims); !hasPermission {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionUserCount %d %d ", requestUid, sessionId)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}
		if resp, err := l.QuerySessionUserCount(sessionId, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionUserCount %d %d", sessionId, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("getSessionUserCount %d %d", sessionId, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func addSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.SessionAddUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("addSessionUser %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if req.Role > model.SessionOwner || req.Role < model.SessionMember {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("addSessionUser %v", req)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("addSessionUser %v %v", req, errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds, claims); !hasPermission {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("addSessionUser %v", req)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}

		if e := l.AddSessionUser(sessionId, req, claims); e != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("addSessionUser %v %v", req, e)
			baseDto.ResponseInternalServerError(ctx, e)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("addSessionUser %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func deleteSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.SessionDelUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSessionUser %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSessionUser %s", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds, claims); !hasPermission {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSessionUser %d %d %v", requestUid, sessionId, req)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}

		if e := l.DelSessionUser(sessionId, true, req, claims); e != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSessionUser %d %d %v %v", requestUid, sessionId, req, e)
			baseDto.ResponseInternalServerError(ctx, e)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("deleteSessionUser %d %d %v", requestUid, sessionId, req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func updateSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.SessionUserUpdateReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSessionUser %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if req.Mute != nil && *req.Mute != 0 && *req.Mute != 1 {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSessionUser %d", *req.Mute)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if req.Role != nil && (*req.Role > model.SessionOwner || *req.Role < model.SessionMember) {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSessionUser %v", req)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSessionUser %v", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		req.SId = sessionId

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds, claims); !hasPermission {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSessionUser %d %d %v", requestUid, sessionId, req)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}

		if e := l.UpdateSessionUser(req, claims); e != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSessionUser %d %d %v %v", requestUid, sessionId, req, e)
			baseDto.ResponseInternalServerError(ctx, e)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("updateSessionUser %d %d %v", requestUid, sessionId, req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func checkReadPermission(appCtx *app.Context, uId, sessionId int64, claims baseDto.ThkClaims) bool {
	sessionUser, err := appCtx.SessionUserModel().FindSessionUser(sessionId, uId)
	if err != nil {
		appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("checkReadPermission %d %d %v", uId, sessionId, err)
		return false
	}
	if sessionUser.UserId <= 0 {
		appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("checkReadPermission %d", sessionUser.UserId)
		return false
	}
	return true
}

func checkPermission(appCtx *app.Context, uId, sessionId int64, oprUIds []int64, claims baseDto.ThkClaims) bool {
	sessionUser, err := appCtx.SessionUserModel().FindSessionUser(sessionId, uId)
	if err != nil {
		appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("checkReadPermission %d %d %v %v", uId, sessionId, oprUIds, err)
		return false
	}
	if sessionUser.UserId <= 0 {
		appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("checkReadPermission %d", sessionUser.UserId)
		return false
	}
	if sessionUser.Role <= model.SessionAdmin {
		return false
	}
	sessionUsers, errSessionUser := appCtx.SessionUserModel().FindSessionUsers(sessionId, oprUIds)
	if errSessionUser != nil {
		return false
	}
	for _, su := range sessionUsers {
		if su.Role >= sessionUser.Role {
			return false
		}
	}
	return true
}
