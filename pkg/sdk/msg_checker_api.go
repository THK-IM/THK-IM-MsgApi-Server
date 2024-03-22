package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"net/http"
	"time"
)

const (
	messageCheckerUrl = "/system/message/check"
)

type (
	MsgCheckerApi interface {
		CheckMessage(req *dto.CheckMessageReq, claims baseDto.ThkClaims) error
	}

	defaultMsgCheckerApi struct {
		endpoint string
		logger   *logrus.Entry
		client   *resty.Client
	}
)

func (d defaultMsgCheckerApi) CheckMessage(req *dto.CheckMessageReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("CheckMessage: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s", d.endpoint, messageCheckerUrl)
	request := d.client.R()
	for k, v := range claims {
		vs := v.(string)
		request.SetHeader(k, vs)
	}
	res, errRequest := request.
		SetHeader("Content-Type", jsonContentType).
		SetBody(dataBytes).
		Post(url)
	if errRequest != nil {
		return errRequest
	}
	if res.StatusCode() != http.StatusOK {
		e := errorx.NewErrorXFromResp(res)
		d.logger.Errorf("CheckMessage: %v %v", req, e)
		return e
	} else {
		d.logger.Infof("CheckMessage: %v %s", req, "success")
		return nil
	}
}

func NewMsgCheckerApi(sdk conf.Sdk, logger *logrus.Entry) MsgCheckerApi {
	return defaultMsgCheckerApi{
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
