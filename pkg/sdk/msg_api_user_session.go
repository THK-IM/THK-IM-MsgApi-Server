package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/errorx"
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

func (d defaultMsgApi) QueryUserSession(req *dto.QueryUserSessionReq, claims baseDto.ThkClaims) (*dto.UserSession, error) {
	url := fmt.Sprintf("%s%s?u_id=%d&type=%d&entity_id=%d", d.endpoint, userSessionUrl, req.UId, req.Type, req.EntityId)
	request := d.client.R()
	for k, v := range claims {
		vs := v.(string)
		request.SetHeader(k, vs)
	}
	res, errRequest := request.
		SetHeader("Content-Type", jsonContentType).
		Get(url)
	if errRequest != nil {
		return nil, errRequest
	}
	if res.StatusCode() != http.StatusOK {
		errRes := &errorx.ErrorX{}
		e := json.Unmarshal(res.Body(), errRes)
		if e != nil {
			d.logger.Errorf("QueryUserSession: %v %v", req, e)
			return nil, e
		} else {
			return nil, errRes
		}
	} else {
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.Info("QueryUserSession: %v %s", req, "Body is nil")
			return nil, nil
		}
		resp := &dto.UserSession{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("QueryUserSession: %v %s", req, e)
			return nil, e
		} else {
			d.logger.Infof("QueryUserSession: %v %v", req, resp)
			return resp, nil
		}
	}
}
