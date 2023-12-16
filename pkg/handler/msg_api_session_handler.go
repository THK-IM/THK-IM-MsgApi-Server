package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msg-api-server/pkg/app"
	"github.com/thk-im/thk-im-msg-api-server/pkg/dto"
	"github.com/thk-im/thk-im-msg-api-server/pkg/logic"
	"github.com/thk-im/thk-im-msg-api-server/pkg/model"
	"strconv"
)

func createSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.CreateSessionReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("createSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		appCtx.Logger().Info(req)
		if req.Type == model.SingleSessionType && (req.EntityId == 0 || len(req.Members) != 0) {
			appCtx.Logger().Errorf("createSession %d %d %v", req.Type, req.EntityId, req.Members)
			baseDto.ResponseBadRequest(ctx)
			return
		} else if (req.Type == model.GroupSessionType || req.Type == model.SuperGroupSessionType) && req.EntityId == 0 {
			appCtx.Logger().Errorf("createSession %d %d", req.Type, req.EntityId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		for _, member := range req.Members {
			if member <= 0 {
				appCtx.Logger().Errorf("createSession %d ", member)
				baseDto.ResponseBadRequest(ctx)
				return
			}
		}

		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 {
			if requestUid != req.Members[0] {
				appCtx.Logger().Errorf("createSession %d %d", requestUid, req.Members[0])
				baseDto.ResponseBadRequest(ctx)
				return
			}
		}

		if resp, err := l.CreateSession(req); err != nil {
			appCtx.Logger().Errorf("createSession %s", err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("updateSession %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func updateSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.UpdateSessionReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Errorf("updateSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if req.Mute != nil {
			if *req.Mute != 0 && *req.Mute != 1 {
				appCtx.Logger().Errorf("updateSession %d", *req.Mute)
				baseDto.ResponseBadRequest(ctx)
				return
			}
		}
		if id, err := strconv.Atoi(ctx.Param("id")); err != nil {
			appCtx.Logger().Errorf("updateSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		} else {
			req.Id = int64(id)
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 {
			if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(req.Id, requestUid); err != nil {
				baseDto.ResponseForbidden(ctx)
				return
			} else {
				if sessionUser.Role == model.SessionMember {
					appCtx.Logger().Errorf("updateSession %d", sessionUser.Role)
					baseDto.ResponseForbidden(ctx)
					return
				}
			}
		}

		if err := l.UpdateSession(req); err != nil {
			appCtx.Logger().Errorf("updateSession %s", err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("updateSession %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func updateUserSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.UpdateUserSessionReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Errorf("updateUserSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if req.Status != nil && (*req.Status < 0 || *req.Status > 3) {
			appCtx.Logger().Errorf("updateUserSession %d", req.Status)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Errorf("updateUserSession %d %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		if err := l.UpdateUserSession(req); err != nil {
			appCtx.Logger().Errorf("updateUserSession %s", err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("updateUserSession %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func getUserSessions(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetUserSessionsReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			appCtx.Logger().Errorf("getUserSessions %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Errorf("getUserSessions %d %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		if resp, err := l.GetUserSessions(req); err != nil {
			appCtx.Logger().Errorf("getUserSessions %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("getUserSessions %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func getUserSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var (
			uid = ctx.Param("uid")
			sid = ctx.Param("sid")
		)

		iUid, errUId := strconv.ParseInt(uid, 10, 64)
		if errUId != nil {
			appCtx.Logger().Errorf("getUserSession %s", errUId.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		iSid, errSId := strconv.ParseInt(sid, 10, 64)
		if errSId != nil {
			appCtx.Logger().Errorf("getUserSession %s", errSId.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != iUid {
			appCtx.Logger().Errorf("getUserSession %d %d", requestUid, iUid)
			baseDto.ResponseForbidden(ctx)
			return
		}

		if res, err := l.GetUserSession(iUid, iSid); err != nil {
			appCtx.Logger().Errorf("getUserSession %d %d %v", iUid, iSid, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Info("getUserSession %d %d %v", iUid, iSid, res)
			baseDto.ResponseSuccess(ctx, res)
		}
	}
}

func getSessionMessages(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var (
			sessionId = ctx.Param("id")
		)
		iSessionId, errSession := strconv.ParseInt(sessionId, 10, 64)
		if errSession != nil {
			appCtx.Logger().Errorf("getSessionMessages %s", errSession.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 {
			if _, err := appCtx.SessionUserModel().FindSessionUser(iSessionId, requestUid); err != nil {
				appCtx.Logger().Errorf("getSessionMessages %s", err.Error())
				baseDto.ResponseForbidden(ctx)
				return
			}
		}
		var req dto.GetSessionMessageReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().Errorf("getSessionMessages %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		req.SId = iSessionId
		if res, err := l.GetSessionMessages(req); err != nil {
			appCtx.Logger().Errorf("getSessionMessages %v %s", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Errorf("getSessionMessages %v %v", req, res)
			baseDto.ResponseSuccess(ctx, res)
		}
	}
}

func deleteSessionMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var (
			sessionId = ctx.Param("id")
		)
		iSessionId, errSessionId := strconv.ParseInt(sessionId, 10, 64)
		if errSessionId != nil {
			appCtx.Logger().Errorf("deleteSessionMessage %s", errSessionId.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		var req dto.DelSessionMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Errorf("deleteSessionMessage %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 {
			if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(iSessionId, requestUid); err != nil {
				appCtx.Logger().Errorf("deleteSessionMessage %s", err.Error())
				baseDto.ResponseForbidden(ctx)
				return
			} else {
				if sessionUser.Role != model.SessionOwner {
					appCtx.Logger().Errorf("deleteSessionMessage %d %d %d", sessionUser.UserId, sessionUser.SessionId, sessionUser.Role)
					baseDto.ResponseForbidden(ctx)
					return
				}
			}
		}
		req.SId = iSessionId
		if err := l.DelSessionMessage(&req); err != nil {
			appCtx.Logger().Errorf("deleteSessionMessage %v %s", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Infof("deleteSessionMessage %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func deleteUserSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var (
			uid = ctx.Param("uid")
			sid = ctx.Param("sid")
		)

		iUid, errUId := strconv.ParseInt(uid, 10, 64)
		if errUId != nil {
			appCtx.Logger().Errorf("deleteUserSession %s", errUId.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		iSid, errSId := strconv.ParseInt(sid, 10, 64)
		if errSId != nil {
			appCtx.Logger().Errorf("deleteUserSession %s", errSId.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(middleware.UidKey)
		if requestUid > 0 && requestUid != iUid {
			appCtx.Logger().Errorf("deleteUserSession %d %d", requestUid, iUid)
			baseDto.ResponseForbidden(ctx)
			return
		}
		req := dto.SessionDelUserReq{
			UIds: []int64{iUid},
		}
		if err := l.DelSessionUser(iSid, true, req); err != nil {
			appCtx.Logger().Errorf("deleteUserSession %v %s", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().Errorf("deleteUserSession %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}
