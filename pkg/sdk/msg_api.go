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
	MsgApi interface {
		UpdateSession(sessionId int64, req *dto.UpdateSessionReq) error
		QuerySessionUsers(sessionId int64, req *dto.QuerySessionUsersReq) (*dto.QuerySessionUsersRes, error)
		QuerySessionUser(sessionId, userId int64) (*dto.SessionUser, error)
		UpdateSessionUser(sessionId int64, req *dto.SessionUserUpdateReq) error
		DelSessionUser(sessionId int64, req *dto.SessionDelUserReq) error
		AddSessionUser(sessionId int64, req *dto.SessionAddUserReq) error
		CreateSession(req *dto.CreateSessionReq) (*dto.CreateSessionRes, error)
		PostUserOnlineStatus(req *dto.PostUserOnlineReq) error
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
