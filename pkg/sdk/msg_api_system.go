package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"net/http"
)

const (
	systemUrl = "/system"
)

func (d defaultMsgApi) PostUserOnlineStatus(req *dto.PostUserOnlineReq) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("PostUserOnlineStatus: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s%s", d.endpoint, systemUrl, "/user/online")
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
