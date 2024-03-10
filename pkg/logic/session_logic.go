package logic

import (
	"fmt"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseErrorx "github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
)

type SessionLogic struct {
	appCtx *app.Context
}

func NewSessionLogic(appCtx *app.Context) SessionLogic {
	return SessionLogic{
		appCtx: appCtx,
	}
}

func (l *SessionLogic) CreateSession(req dto.CreateSessionReq, claims baseDto.ThkClaims) (*dto.CreateSessionRes, error) {
	lockKey := fmt.Sprintf(sessionCreateLockKey, l.appCtx.Config().Name, req.UId, req.EntityId)
	locker := l.appCtx.NewLocker(lockKey, 1000, 1000)
	success, lockErr := locker.Lock()
	if lockErr != nil || !success {
		return nil, baseErrorx.ErrServerBusy
	}
	defer func() {
		if success, lockErr = locker.Release(); lockErr != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("release locker success: %t, error: %s", success, lockErr.Error())
		}
	}()

	if req.Type == model.SingleSessionType {
		if len(req.Members) > 0 {
			return nil, baseErrorx.ErrParamsError
		}
		userSession, err := l.appCtx.UserSessionModel().FindUserSessionByEntityId(req.UId, req.EntityId, req.Type, true)
		if err != nil {
			return nil, err
		}
		// 如果session已经存在
		if userSession.UserId > 0 {
			// 如果单聊被删除，需要恢复
			if userSession.Deleted == 1 {
				session, errSession := l.appCtx.SessionModel().FindSession(userSession.SessionId)
				if errSession != nil {
					return nil, errSession
				}
				userSessions, errUserSessions := l.appCtx.SessionUserModel().AddUser(
					session, []int64{userSession.EntityId}, []int64{userSession.UserId},
					[]int{model.SessionOwner}, []string{""}, []string{""}, 2,
				)
				if errUserSessions != nil {
					return nil, errUserSessions
				}
				userSession = userSessions[0]
			}
			return &dto.CreateSessionRes{
				SId:      userSession.SessionId,
				EntityId: userSession.EntityId,
				ParentId: userSession.ParentId,
				Type:     userSession.Type,
				Name:     userSession.Name,
				Remark:   userSession.Remark,
				Role:     userSession.Role,
				Mute:     userSession.Mute,
				CTime:    userSession.CreateTime,
				MTime:    userSession.UpdateTime,
				Status:   userSession.Status,
				Top:      userSession.Top,
				IsNew:    false,
			}, nil
		}
	} else {
		userSession, err := l.appCtx.UserSessionModel().FindUserSessionByEntityId(req.UId, req.EntityId, req.Type, true)
		if err != nil {
			return nil, err
		}
		// 如果session已经存在
		if userSession.UserId > 0 {
			// 群删除后不能恢复
			if userSession.Deleted == 1 {
				return nil, errorx.ErrSessionAlreadyDeleted
			}
			return &dto.CreateSessionRes{
				SId:      userSession.SessionId,
				EntityId: userSession.EntityId,
				ParentId: userSession.ParentId,
				Type:     userSession.Type,
				Name:     userSession.Name,
				Remark:   userSession.Remark,
				Role:     userSession.Role,
				Mute:     userSession.Mute,
				CTime:    userSession.CreateTime,
				MTime:    userSession.UpdateTime,
				Status:   userSession.Status,
				Top:      userSession.Top,
				IsNew:    false,
			}, nil
		}
	}
	return l.createNewSession(req)
}

func (l *SessionLogic) createNewSession(req dto.CreateSessionReq) (*dto.CreateSessionRes, error) {
	session, err := l.appCtx.SessionModel().CreateEmptySession(req.Type, req.ExtData, req.Name, req.Remark)
	if err != nil {
		return nil, err
	}
	var userSession *model.UserSession
	if req.Type == model.SingleSessionType {
		entityIds := []int64{req.EntityId, req.UId}
		uIds := []int64{req.UId, req.EntityId}
		roles := []int{model.SessionOwner, model.SessionOwner}
		noteNames := []string{"", ""}
		noteAvatars := []string{"", ""}
		if userSessions, errUserSessions := l.appCtx.SessionUserModel().AddUser(session, entityIds, uIds, roles, noteNames, noteAvatars, 2); err != nil {
			return nil, errUserSessions
		} else {
			userSession = userSessions[0]
		}
	} else {
		if req.EntityId <= 0 {
			err = baseErrorx.ErrParamsError
			return nil, err
		}
		members := make([]int64, 0)
		entityIds := make([]int64, 0)
		noteNames := make([]string, 0)
		noteAvatars := make([]string, 0)
		roles := make([]int, 0)
		// 插入自己的角色和entity_id
		members = append(members, req.UId)
		noteNames = append(noteNames, req.UserNoteName)
		noteAvatars = append(noteAvatars, req.UserNoteAvatar)
		entityIds = append(entityIds, req.EntityId)
		roles = append(roles, model.SessionOwner)
		// 插入群成员的角色和entity_id
		for i, m := range req.Members {
			members = append(members, m)
			if i < len(req.MemberNames) {
				noteNames = append(noteNames, req.MemberNames[i])
			} else {
				noteNames = append(noteNames, "")
			}
			if i < len(req.MemberAvatars) {
				noteAvatars = append(noteAvatars, req.MemberAvatars[i])
			} else {
				noteAvatars = append(noteAvatars, "")
			}
			noteAvatars = append(noteAvatars, req.UserNoteAvatar)
			entityIds = append(entityIds, req.EntityId)
			roles = append(roles, model.SessionMember)
		}
		maxMember := l.appCtx.Config().IM.MaxGroupMember
		if req.Type == model.SuperGroupSessionType {
			maxMember = l.appCtx.Config().IM.MaxGroupMember
		}
		if userSessions, errUserSessions := l.appCtx.SessionUserModel().AddUser(session, entityIds, members, roles, noteNames, noteAvatars, maxMember); err != nil {
			return nil, errUserSessions
		} else {
			userSession = userSessions[0]
		}
	}

	res := &dto.CreateSessionRes{
		SId:      userSession.SessionId,
		EntityId: userSession.EntityId,
		ParentId: userSession.ParentId,
		Type:     userSession.Type,
		Name:     userSession.Name,
		Remark:   userSession.Remark,
		Role:     userSession.Role,
		Mute:     userSession.Mute,
		Status:   userSession.Status,
		Top:      userSession.Top,
		CTime:    userSession.CreateTime,
		MTime:    userSession.UpdateTime,
		IsNew:    true,
	}
	return res, nil
}

func (l *SessionLogic) UpdateSession(req dto.UpdateSessionReq, claims baseDto.ThkClaims) error {
	lockKey := fmt.Sprintf(sessionUpdateLockKey, l.appCtx.Config().Name, req.Id)
	locker := l.appCtx.NewLocker(lockKey, 1000, 1000)
	success, lockErr := locker.Lock()
	if lockErr != nil || !success {
		return baseErrorx.ErrServerBusy
	}
	defer func() {
		if success, lockErr = locker.Release(); lockErr != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("release locker success: %t, error: %s", success, lockErr.Error())
		}
	}()
	err := l.appCtx.SessionModel().UpdateSession(req.Id, req.Name, req.Remark, req.Mute, req.ExtData)
	if err != nil {
		return err
	}
	sessionUsers, errSessionUsers := l.appCtx.SessionUserModel().FindAllSessionUsers(req.Id)
	if errSessionUsers != nil {
		return errSessionUsers
	}
	uIds := make([]int64, 0)
	for _, su := range sessionUsers {
		uIds = append(uIds, su.UserId)
	}
	var mute *string
	if req.Mute == nil {
		mute = nil
	} else if *req.Mute == 0 {
		sql := "mute & (mute ^ 1)"
		mute = &sql
	} else if *req.Mute == 1 {
		sql := "mute | 1"
		mute = &sql
	}
	return l.appCtx.UserSessionModel().UpdateUserSession(uIds, req.Id, req.Name, req.Remark, mute, req.ExtData, nil, nil, nil, nil, nil)
}

func (l *SessionLogic) UpdateSessionType(req dto.UpdateSessionTypeReq, claims baseDto.ThkClaims) error {
	lockKey := fmt.Sprintf(sessionUpdateLockKey, l.appCtx.Config().Name, req.Id)
	locker := l.appCtx.NewLocker(lockKey, 1000, 1000)
	success, lockErr := locker.Lock()
	if lockErr != nil || !success {
		return baseErrorx.ErrServerBusy
	}
	defer func() {
		if success, lockErr = locker.Release(); lockErr != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("release locker success: %t, error: %s", success, lockErr.Error())
		}
	}()
	err := l.appCtx.SessionModel().UpdateSessionType(req.Id, req.Type)
	if err != nil {
		return err
	}
	sessionUsers, errSessionUsers := l.appCtx.SessionUserModel().FindAllSessionUsers(req.Id)
	if errSessionUsers != nil {
		return errSessionUsers
	}
	uIds := make([]int64, 0)
	for _, su := range sessionUsers {
		uIds = append(uIds, su.UserId)
	}
	err = l.appCtx.UserSessionModel().UpdateUserSessionType(uIds, req.Id, req.Type)
	if err == nil {
		err = l.appCtx.SessionUserModel().UpdateType(req.Id, req.Type)
	}
	return err
}

func (l *SessionLogic) DelSession(req dto.DelSessionReq, claims baseDto.ThkClaims) error {
	err := l.appCtx.SessionUserModel().DelSession(req.Id)
	if err != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("DelSession %v  %s", req, err)
	}
	return err
}

func (l *SessionLogic) UpdateUserSession(req dto.UpdateUserSessionReq, claims baseDto.ThkClaims) (err error) {
	lockKey := fmt.Sprintf(userSessionUpdateLockKey, l.appCtx.Config().Name, req.UId, req.SId)
	locker := l.appCtx.NewLocker(lockKey, 1000, 1000)
	success, lockErr := locker.Lock()
	if lockErr != nil || !success {
		return baseErrorx.ErrServerBusy
	}
	defer func() {
		if success, lockErr = locker.Release(); lockErr != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("release locker success: %t, error: %s", success, lockErr.Error())
		}
	}()
	err = l.appCtx.UserSessionModel().UpdateUserSession([]int64{req.UId}, req.SId, nil,
		nil, nil, nil, req.NoteName, req.Top, req.Status, nil, req.ParentId,
	)
	if err == nil {
		err = l.appCtx.SessionUserModel().UpdateUser(req.SId, []int64{req.UId}, nil, req.Status, req.NoteName, req.NoteAvatar, nil)
	} else {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("UpdateUserSession, %v %v", req, err)
	}
	return
}

func (l *SessionLogic) QueryLatestUserSessions(req dto.QueryLatestUserSessionReq, claims baseDto.ThkClaims) (*dto.QueryLatestUserSessionsRes, error) {
	userSessions, err := l.appCtx.UserSessionModel().QueryLatestUserSessions(req.UId, req.MTime, req.Offset, req.Count, req.Types)
	if err != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("QueryLatestUserSessions, %v %v", req, err)
		return nil, err
	}
	dtoUserSessions := make([]*dto.UserSession, 0)
	for _, userSession := range userSessions {
		dtoUserSession := l.convUserSession(userSession)
		dtoUserSessions = append(dtoUserSessions, dtoUserSession)
	}
	return &dto.QueryLatestUserSessionsRes{Data: dtoUserSessions}, nil
}

func (l *SessionLogic) GetUserSession(uId, sId int64, claims baseDto.ThkClaims) (*dto.UserSession, error) {
	userSession, err := l.appCtx.UserSessionModel().GetUserSession(uId, sId)
	if err != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("GetUserSession, %v %v %v", uId, sId, err)
		return nil, err
	}
	dtoUserSession := l.convUserSession(userSession)
	return dtoUserSession, nil
}

func (l *SessionLogic) GetUserSessionByEntityId(req *dto.QueryUserSessionReq, claims baseDto.ThkClaims) (*dto.UserSession, error) {
	userSession, err := l.appCtx.UserSessionModel().FindUserSessionByEntityId(req.UId, req.EntityId, req.Type, false)
	if err != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("GetUserSessionByEntityId, %v %v", req, err)
		return nil, err
	}
	dtoUserSession := l.convUserSession(userSession)
	return dtoUserSession, nil
}

func (l *SessionLogic) convUserSession(userSession *model.UserSession) *dto.UserSession {
	return &dto.UserSession{
		SId:        userSession.SessionId,
		Type:       userSession.Type,
		Name:       userSession.Name,
		Remark:     userSession.Remark,
		Role:       userSession.Role,
		Mute:       userSession.Mute,
		Top:        userSession.Top,
		Status:     userSession.Status,
		EntityId:   userSession.EntityId,
		ExtData:    userSession.ExtData,
		NoteName:   userSession.NoteName,
		NoteAvatar: userSession.NoteAvatar,
		Deleted:    userSession.Deleted,
		CTime:      userSession.CreateTime,
		MTime:      userSession.UpdateTime,
	}
}
