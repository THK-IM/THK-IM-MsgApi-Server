package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"net/http"
)

const userSessionUrl = "user_session"

func (d defaultMsgApi) UpdateUserSession(req *dto.UpdateUserSessionReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("UpdateUserSession: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s", d.endpoint, userSessionUrl)
	request := d.client.R()
	for k, v := range claims {
		vs := v.(string)
		request.SetHeader(k, vs)
	}
	res, errRequest := request.
		SetHeader("Content-Type", jsonContentType).
		SetBody(dataBytes).
		Put(url)
	if errRequest != nil {
		return errRequest
	}
	if res.StatusCode() != http.StatusOK {
		e := errors.New(string(res.Body()))
		d.logger.Errorf("UpdateUserSession: %v %v", req, e)
		return e
	} else {
		d.logger.Errorf("UpdateUserSession: %v %s", req, "success")
		return nil
	}
}
