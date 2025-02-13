package sdk

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"net/http"
	"time"
)

const (
	jsonContentType = "application/json"
)

type (
	SystemApi interface {
		CreateSession(req *dto.CreateSessionReq, claims baseDto.ThkClaims) (*dto.CreateSessionRes, error)
		UpdateSessionType(req *dto.UpdateSessionTypeReq, claims baseDto.ThkClaims) error
		SysDelSessionUser(sessionId int64, req *dto.SessionDelUserReq, claims baseDto.ThkClaims) error
		SysAddSessionUser(sessionId int64, req *dto.SessionAddUserReq, claims baseDto.ThkClaims) error
		SysQueryLatestSessionUsers(sessionId int64, req *dto.QuerySessionUsersReq, claims baseDto.ThkClaims) (*dto.QuerySessionUsersRes, error)
		PushMessage(req *dto.PushMessageReq, claims baseDto.ThkClaims) (*dto.PushMessageRes, error)
		SendSysMessage(req *dto.SendSysMessageReq, claims baseDto.ThkClaims) (*dto.SendSysMessageRes, error)
		SendSessionMessage(req *dto.SendMessageReq, claims baseDto.ThkClaims) (*dto.SendMessageRes, error)
		KickOffUser(req *dto.KickUserReq, claims baseDto.ThkClaims) error
		QueryUsersOnlineStatus(req *dto.QueryUsersOnlineStatusReq, claims baseDto.ThkClaims) (*dto.QueryUsersOnlineStatusRes, error)
		PostUserOnlineStatus(req *dto.PostUserOnlineReq, claims baseDto.ThkClaims) error
	}

	SessionApi interface {
		QuerySessionUserCount(sessionId int64, claims baseDto.ThkClaims) (*dto.SessionUserCountRes, error)
		DelSessionUser(sessionId int64, req *dto.SessionDelUserReq, claims baseDto.ThkClaims) error
		AddSessionUser(sessionId int64, req *dto.SessionAddUserReq, claims baseDto.ThkClaims) error
		DelSession(sessionId int64, req *dto.DelSessionReq, claims baseDto.ThkClaims) error
		UpdateSession(sessionId int64, req *dto.UpdateSessionReq, claims baseDto.ThkClaims) error
		QueryLatestSessionUsers(sessionId int64, req *dto.QuerySessionUsersReq, claims baseDto.ThkClaims) (*dto.QuerySessionUsersRes, error)
		QuerySessionUser(sessionId, userId int64, claims baseDto.ThkClaims) (*dto.SessionUser, error)
		UpdateSessionUser(sessionId int64, req *dto.SessionUserUpdateReq, claims baseDto.ThkClaims) error
	}

	UserSessionApi interface {
		DeleteUserSession(userId, sessionId int64, claims baseDto.ThkClaims) error
		UpdateUserSession(req *dto.UpdateUserSessionReq, claims baseDto.ThkClaims) error
		QueryUserSession(req *dto.QueryUserSessionReq, claims baseDto.ThkClaims) (*dto.UserSession, error)
		SearchUserSession(req *dto.SearchUserSessionReq, claims baseDto.ThkClaims) (*dto.SearchUserSessionRes, error)
		QueryLatestUserSession(req *dto.QueryLatestUserSessionReq, claims baseDto.ThkClaims) (*dto.QueryLatestUserSessionsRes, error)
	}

	UserMessageApi interface {
		ReadUserMessage(req *dto.ReadUserMessageReq, claims baseDto.ThkClaims) error
	}

	MsgApi interface {
		SystemApi
		SessionApi
		UserSessionApi
		UserMessageApi
	}

	defaultMsgApi struct {
		endpoint string
		logger   *logrus.Entry
		client   *resty.Client
	}
)

func NewMsgApi(sdk conf.Sdk, logger *logrus.Entry) MsgApi {
	return defaultMsgApi{
		endpoint: sdk.Endpoint,
		logger:   logger.WithField("rpc", sdk.Name),
		client: resty.New().
			SetTransport(&http.Transport{
				MaxIdleConns:    10,
				MaxConnsPerHost: 10,
				IdleConnTimeout: 30 * time.Second,
			}).
			SetTimeout(5 * time.Second).
			SetRetryCount(3).
			SetRetryWaitTime(15 * time.Second).
			SetRetryMaxWaitTime(5 * time.Second),
	}
}
