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

func createSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.CreateSessionReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		appCtx.Logger().WithFields(logrus.Fields(claims)).Info(req)
		if req.Type == model.SingleSessionType && (req.EntityId == 0 || len(req.Members) != 0) {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createSession %d %d %v", req.Type, req.EntityId, req.Members)
			baseDto.ResponseBadRequest(ctx)
			return
		} else if (req.Type == model.GroupSessionType || req.Type == model.SuperGroupSessionType) && req.EntityId == 0 {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createSession %d %d", req.Type, req.EntityId)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		for _, member := range req.Members {
			if member <= 0 {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createSession %d ", member)
				baseDto.ResponseBadRequest(ctx)
				return
			}
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 {
			if requestUid != req.Members[0] {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createSession %d %d", requestUid, req.Members[0])
				baseDto.ResponseBadRequest(ctx)
				return
			}
		}

		if resp, err := l.CreateSession(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("createSession %s", err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("createSession %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func updateSessionType(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.UpdateSessionTypeReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSessionType %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 {
			if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(req.Id, requestUid); err != nil {
				baseDto.ResponseForbidden(ctx)
				return
			} else {
				if sessionUser.Role == model.SessionMember {
					appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSessionType %d", sessionUser.Role)
					baseDto.ResponseForbidden(ctx)
					return
				}
			}
		}

		if err := l.UpdateSessionType(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSessionType %s", err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("updateSessionType %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func updateSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.UpdateSessionReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if req.Mute != nil {
			if *req.Mute != 0 && *req.Mute != 1 {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSession %d", *req.Mute)
				baseDto.ResponseBadRequest(ctx)
				return
			}
		}
		if id, err := strconv.Atoi(ctx.Param("id")); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		} else {
			req.Id = int64(id)
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 {
			if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(req.Id, requestUid); err != nil {
				baseDto.ResponseForbidden(ctx)
				return
			} else {
				if sessionUser.Role == model.SessionMember {
					appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSession %d", sessionUser.Role)
					baseDto.ResponseForbidden(ctx)
					return
				}
			}
		}

		if err := l.UpdateSession(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateSession %s", err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("updateSession %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func deleteSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.DelSessionReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		if id, err := strconv.Atoi(ctx.Param("id")); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		} else {
			req.Id = int64(id)
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 {
			if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(req.Id, requestUid); err != nil {
				baseDto.ResponseForbidden(ctx)
				return
			} else {
				if sessionUser.Role == model.SessionMember {
					appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSession %d", sessionUser.Role)
					baseDto.ResponseForbidden(ctx)
					return
				} else if sessionUser.Type == model.SingleSessionType {
					appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSession %d", sessionUser.Type)
					baseDto.ResponseBadRequest(ctx)
					return
				}
			}
		}

		if err := l.DelSession(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSession %s", err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("deleteSession %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func updateUserSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.UpdateUserSessionReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateUserSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		if req.Status != nil && (*req.Status < 0 || *req.Status > 3) {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateUserSession %d", req.Status)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateUserSession %d %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		if err := l.UpdateUserSession(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("updateUserSession %s", err.Error())
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("updateUserSession %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func getLatestUserSessions(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.QueryLatestUserSessionReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getUserSessions %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getUserSessions %d %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		if resp, err := l.QueryLatestUserSessions(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getUserSessions %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("getUserSessions %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func queryUserSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.QueryUserSessionReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryUserSession %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryUserSession %d %d", requestUid, req.UId)
			baseDto.ResponseForbidden(ctx)
			return
		}

		if resp, err := l.GetUserSessionByEntityId(&req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryUserSession %v %v", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("queryUserSession %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func queryUserSessionBySId(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var (
			uid = ctx.Param("uid")
			sid = ctx.Param("sid")
		)

		iUid, errUId := strconv.ParseInt(uid, 10, 64)
		if errUId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryUserSessionBySId %s", errUId.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		iSid, errSId := strconv.ParseInt(sid, 10, 64)
		if errSId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryUserSessionBySId %s", errSId.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != iUid {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryUserSessionBySId %d %d", requestUid, iUid)
			baseDto.ResponseForbidden(ctx)
			return
		}

		if res, err := l.GetUserSession(iUid, iSid, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryUserSessionBySId %d %d %v", iUid, iSid, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Info("queryUserSessionBySId %d %d %v", iUid, iSid, res)
			baseDto.ResponseSuccess(ctx, res)
		}
	}
}

func getSessionMessages(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var (
			sessionId = ctx.Param("id")
		)
		iSessionId, errSession := strconv.ParseInt(sessionId, 10, 64)
		if errSession != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionMessages %s", errSession.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 {
			if _, err := appCtx.SessionUserModel().FindSessionUser(iSessionId, requestUid); err != nil {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionMessages %s", err.Error())
				baseDto.ResponseForbidden(ctx)
				return
			}
		}
		var req dto.GetSessionMessageReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionMessages %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		req.SId = iSessionId
		if res, err := l.GetSessionMessages(req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionMessages %v %s", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("getSessionMessages %v %v", req, res)
			baseDto.ResponseSuccess(ctx, res)
		}
	}
}

func deleteSessionMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var (
			sessionId = ctx.Param("id")
		)
		iSessionId, errSessionId := strconv.ParseInt(sessionId, 10, 64)
		if errSessionId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSessionMessage %s", errSessionId.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		var req dto.DelSessionMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSessionMessage %s", err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 {
			if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(iSessionId, requestUid); err != nil {
				appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSessionMessage %s", err.Error())
				baseDto.ResponseForbidden(ctx)
				return
			} else {
				if sessionUser.Role != model.SessionOwner {
					appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSessionMessage %d %d %d", sessionUser.UserId, sessionUser.SessionId, sessionUser.Role)
					baseDto.ResponseForbidden(ctx)
					return
				}
			}
		}
		req.SId = iSessionId
		if err := l.DelSessionMessage(&req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteSessionMessage %v %s", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("deleteSessionMessage %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}

func deleteUserSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var (
			uid = ctx.Param("uid")
			sid = ctx.Param("sid")
		)

		iUid, errUId := strconv.ParseInt(uid, 10, 64)
		if errUId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteUserSession %s", errUId.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		iSid, errSId := strconv.ParseInt(sid, 10, 64)
		if errSId != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteUserSession %s", errSId.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(userSdk.UidKey)
		if requestUid > 0 && requestUid != iUid {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteUserSession %d %d", requestUid, iUid)
			baseDto.ResponseForbidden(ctx)
			return
		}
		req := dto.SessionDelUserReq{
			UIds: []int64{iUid},
		}
		if err := l.DelSessionUser(iSid, true, req, claims); err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteUserSession %v %s", req, err)
			baseDto.ResponseInternalServerError(ctx, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("deleteUserSession %v", req)
			baseDto.ResponseSuccess(ctx, nil)
		}
	}
}
