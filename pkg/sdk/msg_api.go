package sdk

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"net/http"
	"time"
)

const (
	jsonContentType = "application/json"
)

type (
	SystemApi interface {
		PushMessage(req *dto.PushMessageReq) (*dto.PushMessageRes, error)
		SendSysMessage(req *dto.SendSysMessageReq) (*dto.SendSysMessageRes, error)
		SendSessionMessage(req *dto.SendMessageReq) (*dto.SendMessageRes, error)
		KickOffUser(req *dto.KickUserReq) error
		QueryUsersOnlineStatus(req *dto.QueryUsersOnlineStatusReq) (*dto.QueryUsersOnlineStatusRes, error)
		PostUserOnlineStatus(req *dto.PostUserOnlineReq) error
	}

	SessionApi interface {
		DelSession(sessionId int64, req *dto.DelSessionReq) error
		UpdateSession(sessionId int64, req *dto.UpdateSessionReq) error
		QuerySessionUsers(sessionId int64, req *dto.QuerySessionUsersReq) (*dto.QuerySessionUsersRes, error)
		QuerySessionUser(sessionId, userId int64) (*dto.SessionUser, error)
		UpdateSessionUser(sessionId int64, req *dto.SessionUserUpdateReq) error
		DelSessionUser(sessionId int64, req *dto.SessionDelUserReq) error
		AddSessionUser(sessionId int64, req *dto.SessionAddUserReq) error
		CreateSession(req *dto.CreateSessionReq) (*dto.CreateSessionRes, error)
	}

	MsgApi interface {
		SystemApi
		SessionApi
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
