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

const (
	sessionUrl = "/session"
)

func (d defaultMsgApi) DelSession(sessionId int64, req *dto.DelSessionReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("DelSession: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/%d", d.endpoint, sessionUrl, sessionId)
	request := d.client.R()
	for k, v := range claims {
		vs := v.(string)
		request.SetHeader(k, vs)
	}
	res, errRequest := request.
		SetHeader("Content-Type", jsonContentType).
		SetBody(dataBytes).
		Delete(url)
	if errRequest != nil {
		return errRequest
	}
	if res.StatusCode() != http.StatusOK {
		e := errors.New(string(res.Body()))
		d.logger.Errorf("DelSession: %v %v", req, e)
		return e
	} else {
		d.logger.Errorf("DelSession: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) UpdateSession(sessionId int64, req *dto.UpdateSessionReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("UpdateSession: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/%d", d.endpoint, sessionUrl, sessionId)
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
		d.logger.Errorf("UpdateSession: %v %v", req, e)
		return e
	} else {
		d.logger.Errorf("UpdateSession: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) QuerySessionUsers(sessionId int64, req *dto.QuerySessionUsersReq, claims baseDto.ThkClaims) (*dto.QuerySessionUsersRes, error) {
	url := fmt.Sprintf("%s%s/%d/user?count=%d&m_time=%d", d.endpoint, sessionUrl, sessionId, req.Count, req.MTime)
	if req.Role != nil {
		url += fmt.Sprintf("&role=%d", *req.Role)
	}
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
			d.logger.Errorf("QuerySessionUsers: %v %v", req, e)
			return nil, e
		} else {
			return nil, errRes
		}
	} else {
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.Info("QuerySessionUsers: %v %s", req, "Body is nil")
			return nil, nil
		}
		resp := &dto.QuerySessionUsersRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("QuerySessionUsers: %v %v", req, e)
			return nil, e
		} else {
			d.logger.Infof("QuerySessionUsers: %v %v", req, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) QuerySessionUser(sessionId, userId int64, claims baseDto.ThkClaims) (*dto.SessionUser, error) {
	url := fmt.Sprintf("%s%s/%d/user/%d", d.endpoint, sessionUrl, sessionId, userId)

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
			d.logger.Errorf("QuerySessionUsers: %d %d %v", sessionId, userId, e)
			return nil, e
		} else {
			return nil, errRes
		}
	} else {
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.Info("QuerySessionUsers: %d %d %s", sessionId, userId, "Body is nil")
			return nil, nil
		}
		resp := &dto.SessionUser{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("QuerySessionUsers: %d %d %v", sessionId, userId, e)
			return nil, e
		} else {
			d.logger.Infof("QuerySessionUsers: %d %d %v", sessionId, userId, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) DelSessionUser(sessionId int64, req *dto.SessionDelUserReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("DelSessionUser: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/%d/user", d.endpoint, sessionUrl, sessionId)
	request := d.client.R()
	for k, v := range claims {
		vs := v.(string)
		request.SetHeader(k, vs)
	}
	res, errRequest := request.
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

func (d defaultMsgApi) UpdateSessionUser(sessionId int64, req *dto.SessionUserUpdateReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("UpdateSessionUser: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/%d/user", d.endpoint, sessionUrl, sessionId)
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
		d.logger.Errorf("UpdateSessionUser: %v %v", req, e)
		return e
	} else {
		d.logger.Errorf("UpdateSessionUser: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) AddSessionUser(sessionId int64, req *dto.SessionAddUserReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("AddSessionUser: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/%d/user", d.endpoint, sessionUrl, sessionId)
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
		e := errors.New(string(res.Body()))
		d.logger.Errorf("AddSessionUser: %v %v", req, e)
		return e
	} else {
		d.logger.Errorf("AddSessionUser: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) CreateSession(req *dto.CreateSessionReq, claims baseDto.ThkClaims) (*dto.CreateSessionRes, error) {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("CreateSession: %v %v", req, err)
		return nil, err
	}
	url := fmt.Sprintf("%s%s", d.endpoint, sessionUrl)
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
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.Infof("CreateSession: %v %s", req, "Body is nil")
			return nil, nil
		}
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
