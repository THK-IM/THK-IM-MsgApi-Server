package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"net/http"
)

const (
	sessionUrl = "/session"
)

func (d defaultMsgApi) DelSessionUser(sessionId int64, req *dto.SessionDelUserReq) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("DelSessionUser: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/%d/user", d.endpoint, sessionUrl, sessionId)
	res, errRequest := d.client.R().
		SetHeader("Content-Type", jsonContentType).
		SetBody(dataBytes).
		Delete(url)
	if errRequest != nil {
		return errRequest
	}
	if res.StatusCode() != http.StatusOK {
		e := errors.New(string(res.Body()))
		d.logger.Errorf("DelSessionUser: %v %v", req, e)
		return e
	} else {
		d.logger.Errorf("DelSessionUser: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) AddSessionUser(sessionId int64, req *dto.SessionAddUserReq) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("AddSessionUser: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/%d/user", d.endpoint, sessionUrl, sessionId)
	res, errRequest := d.client.R().
		SetHeader("Content-Type", jsonContentType).
		SetBody(dataBytes).
		Post(url)
	if errRequest != nil {
		return errRequest
	}
	if res.StatusCode() != http.StatusOK {
		e := errors.New(string(res.Body()))
		d.logger.Errorf("AddSessionUser: %v %v", req, e)
		return e
	} else {
		d.logger.Errorf("AddSessionUser: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) CreateSession(req *dto.CreateSessionReq) (*dto.CreateSessionRes, error) {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("CreateSession: %v %v", req, err)
		return nil, err
	}
	url := fmt.Sprintf("%s%s", d.endpoint, sessionUrl)
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
