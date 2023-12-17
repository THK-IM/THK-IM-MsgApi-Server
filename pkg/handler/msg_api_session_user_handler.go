package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/logic"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"strconv"
)

func getSessionUsers(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.QuerySessionUsersReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			appCtx.Logger().Errorf("getSessionUsers %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if req.Count <= 0 {
			appCtx.Logger().Errorf("getSessionUsers %v", req)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if req.Role != nil && (*req.Role > model.SessionOwner || *req.Role < model.SessionMember) {
			appCtx.Logger().Errorf("getSessionUsers %v", req)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().Errorf("getSessionUsers %v", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		req.SId = sessionId

		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkReadPermission(appCtx, requestUid, sessionId); !hasPermission {
				appCtx.Logger().Errorf("getSessionUsers %d %d ", requestUid, sessionId)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}
		if resp, err := l.QuerySessionUsers(req); err != nil {
			appCtx.Logger().Errorf("getSessionUsers %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("getSessionUsers %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func getSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().Errorf("getSessionUser %v", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		userId, errUserId := strconv.ParseInt(ctx.Param("uid"), 10, 64)
		if errUserId != nil {
			appCtx.Logger().Errorf("getSessionUser %v", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkReadPermission(appCtx, requestUid, sessionId); !hasPermission {
				appCtx.Logger().Errorf("getSessionUser %d %d ", requestUid, sessionId)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}
		if resp, err := l.QuerySessionUser(sessionId, userId); err != nil {
			appCtx.Logger().Errorf("getSessionUser %d %d %v", sessionId, userId, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("getSessionUser %d %d %v", sessionId, userId, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func addSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SessionAddUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Errorf("addSessionUser %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if req.Role > model.SessionOwner || req.Role < model.SessionMember {
			appCtx.Logger().Errorf("addSessionUser %v", req)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().Errorf("addSessionUser %v %v", req, errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds); !hasPermission {
				appCtx.Logger().Errorf("addSessionUser %v", req)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}

		if e := l.AddSessionUser(sessionId, req); e != nil {
			appCtx.Logger().Errorf("addSessionUser %v %v", req, e)
			baseDto.ResponseInternalServerError(ctx, e)
		} else {
			appCtx.Logger().Infof("addSessionUser %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func deleteSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SessionDelUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Errorf("deleteSessionUser %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().Errorf("deleteSessionUser %s", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds); !hasPermission {
				appCtx.Logger().Errorf("deleteSessionUser %d %d %v", requestUid, sessionId, req)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}

		if e := l.DelSessionUser(sessionId, true, req); e != nil {
			appCtx.Logger().Errorf("deleteSessionUser %d %d %v %v", requestUid, sessionId, req, e)
			baseDto.ResponseInternalServerError(ctx, e)
		} else {
			appCtx.Logger().Infof("deleteSessionUser %d %d %v", requestUid, sessionId, req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func updateSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SessionUserUpdateReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Errorf("updateSessionUser %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if req.Mute != nil && *req.Mute != 0 && *req.Mute != 1 {
			appCtx.Logger().Errorf("updateSessionUser %d", *req.Mute)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if req.Role != nil && (*req.Role > model.SessionOwner || *req.Role < model.SessionMember) {
			appCtx.Logger().Errorf("updateSessionUser %v", req)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().Errorf("updateSessionUser %v", errSessionId)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		req.SId = sessionId

		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds); !hasPermission {
				appCtx.Logger().Errorf("updateSessionUser %d %d %v", requestUid, sessionId, req)
				baseDto.ResponseForbidden(ctx)
				return
			}
		}

		if e := l.UpdateSessionUser(req); e != nil {
			appCtx.Logger().Errorf("updateSessionUser %d %d %v %v", requestUid, sessionId, req, e)
			baseDto.ResponseInternalServerError(ctx, e)
		} else {
			appCtx.Logger().Infof("updateSessionUser %d %d %v", requestUid, sessionId, req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func checkReadPermission(appCtx *app.Context, uId, sessionId int64) bool {
	if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(sessionId, uId); err != nil {
		appCtx.Logger().Infof("checkReadPermission %d %d %v", uId, sessionId, err)
		return false
	} else {
		if sessionUser.UserId > 0 {
			return true
		}
	}
	return true
}

func checkPermission(appCtx *app.Context, uId, sessionId int64, oprUIds []int64) bool {
	if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(sessionId, uId); err != nil {
		appCtx.Logger().Infof("checkReadPermission %d %d %v %v", uId, sessionId, oprUIds, err)
		return false
	} else {
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
	}
	return true
}
