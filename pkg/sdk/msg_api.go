package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-msg-api-server/pkg/dto"
	"net/http"
	"time"
)

const (
	createSessionUrl              = "/session"
	msgApiPostUserOnlineStatusUrl = "/system/user/online"
	jsonContentType               = "application/json"
)

type (
	MsgApi interface {
		CreateSession(req *dto.CreateSessionReq) (*dto.CreateSessionRes, error)
		PostUserOnlineStatus(req *dto.PostUserOnlineReq) error
	}

	defaultMsgApi struct {
		endpoint string
		logger   *logrus.Entry
		client   *resty.Client
	}
)

func (d defaultMsgApi) CreateSession(req *dto.CreateSessionReq) (*dto.CreateSessionRes, error) {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("CreateSession: %v %v", req, err)
		return nil, err
	}
	url := fmt.Sprintf("%s%s", d.endpoint, createSessionUrl)
	res, errRequest := d.client.R().
		SetHeader("Content-Type", jsonContentType).
		SetBody(dataBytes).
		Post(url)
	if errRequest != nil {
		d.logger.Errorf("CreateSession: %v %v", req, errRequest)
		return nil, errRequest
	}
	if res.StatusCode() != http.StatusOK {
		errRes := &errorx.ErrorX{}
		e := json.Unmarshal(res.Body(), errRes)
		if e != nil {
			d.logger.Errorf("CreateSession: %v %v", req, e)
			return nil, e
		} else {
			return nil, errRes
		}
	} else {
		resp := &dto.CreateSessionRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("CreateSession: %v %v", req, e)
			return nil, e
		} else {
			d.logger.Infof("CreateSession: %v %v", req, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) PostUserOnlineStatus(req *dto.PostUserOnlineReq) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("PostUserOnlineStatus: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s", d.endpoint, msgApiPostUserOnlineStatusUrl)
	res, errRequest := d.client.R().
		SetHeader("Content-Type", jsonContentType).
		SetBody(dataBytes).
		Post(url)
	if errRequest != nil {
		return errRequest
	}
	if res.StatusCode() != http.StatusOK {
		e := errors.New(string(res.Body()))
		d.logger.Errorf("PostUserOnlineStatus: %v %v", req, e)
		return e
	} else {
		d.logger.Errorf("PostUserOnlineStatus: %v %s", req, "success")
		return nil
	}
}

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
