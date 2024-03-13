package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"net/http"
)

const userMessageUrl = "/message"

func (d defaultMsgApi) ReadUserMessage(req *dto.ReadUserMessageReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.WithFields(logrus.Fields(claims)).Errorf("ReadUserMessage: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/read", d.endpoint, userMessageUrl)
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
		d.logger.WithFields(logrus.Fields(claims)).Errorf("ReadUserMessage: %v", e)
		return e
	} else {
		d.logger.WithFields(logrus.Fields(claims)).Infof("ReadUserMessage: %v", "success")
		return nil
	}
}
