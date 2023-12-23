package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"net/http"
	"strings"
)

const (
	systemUrl = "/system"
)

func (d defaultMsgApi) CreateSession(req *dto.CreateSessionReq, claims baseDto.ThkClaims) (*dto.CreateSessionRes, error) {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("CreateSession: %v %v", req, err)
		return nil, err
	}
	url := fmt.Sprintf("%s%s/session", d.endpoint, systemUrl)
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

func (d defaultMsgApi) UpdateSessionType(req *dto.UpdateSessionTypeReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("UpdateSessionType: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/session", d.endpoint, systemUrl)
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
		d.logger.Errorf("UpdateSessionType: %v %v", req, e)
		return e
	} else {
		d.logger.Infof("UpdateSessionType: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) SysDelSessionUser(sessionId int64, req *dto.SessionDelUserReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("SysDelSessionUser: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/%s/user", d.endpoint, systemUrl, fmt.Sprintf("session/%d", sessionId))
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
		d.logger.Errorf("SysDelSessionUser: %v %v", req, e)
		return e
	} else {
		d.logger.Infof("SysDelSessionUser: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) SysAddSessionUser(sessionId int64, req *dto.SessionAddUserReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("SysAddSessionUser: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s/%s/user", d.endpoint, systemUrl, fmt.Sprintf("session/%d", sessionId))
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
		d.logger.Errorf("SysAddSessionUser: %v %v", req, e)
		return e
	} else {
		d.logger.Infof("SysAddSessionUser: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) PushMessage(req *dto.PushMessageReq, claims baseDto.ThkClaims) (*dto.PushMessageRes, error) {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("PushMessage: %v %v", req, err)
		return nil, err
	}
	url := fmt.Sprintf("%s%s%s", d.endpoint, systemUrl, "/push_message")
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
		d.logger.Errorf("PushMessage: %v %v", req, errRequest)
		return nil, errRequest
	}
	if res.StatusCode() != http.StatusOK {
		errRes := &errorx.ErrorX{}
		e := json.Unmarshal(res.Body(), errRes)
		if e != nil {
			d.logger.Errorf("PushMessage: %v %v", req, e)
			return nil, e
		} else {
			return nil, errRes
		}
	} else {
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.Infof("PushMessage: %v %s", req, "Body is nil")
			return nil, nil
		}
		resp := &dto.PushMessageRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("PushMessage: %v %v", req, e)
			return nil, e
		} else {
			d.logger.Infof("PushMessage: %v %v", req, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) SendSysMessage(req *dto.SendSysMessageReq, claims baseDto.ThkClaims) (*dto.SendSysMessageRes, error) {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("SendSystemMessage: %v %v", req, err)
		return nil, err
	}
	url := fmt.Sprintf("%s%s%s", d.endpoint, systemUrl, "/sys_message")
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
		d.logger.Errorf("SendSystemMessage: %v %v", req, errRequest)
		return nil, errRequest
	}
	if res.StatusCode() != http.StatusOK {
		errRes := &errorx.ErrorX{}
		e := json.Unmarshal(res.Body(), errRes)
		if e != nil {
			d.logger.Errorf("SendSystemMessage: %v %v", req, e)
			return nil, e
		} else {
			return nil, errRes
		}
	} else {
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.Infof("SendSystemMessage: %v %s", req, "Body is nil")
			return nil, nil
		}
		resp := &dto.SendSysMessageRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("SendSystemMessage: %v %v", req, e)
			return nil, e
		} else {
			d.logger.Infof("SendSystemMessage: %v %v", req, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) SendSessionMessage(req *dto.SendMessageReq, claims baseDto.ThkClaims) (*dto.SendMessageRes, error) {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("SendSessionMessage: %v %v", req, err)
		return nil, err
	}
	url := fmt.Sprintf("%s%s%s", d.endpoint, systemUrl, "/session_message")
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
		d.logger.Errorf("SendSessionMessage: %v %v", req, errRequest)
		return nil, errRequest
	}
	if res.StatusCode() != http.StatusOK {
		errRes := &errorx.ErrorX{}
		e := json.Unmarshal(res.Body(), errRes)
		if e != nil {
			d.logger.Errorf("SendSessionMessage: %v %v", req, e)
			return nil, e
		} else {
			return nil, errRes
		}
	} else {
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.Infof("SendSessionMessage: %v %s", req, "Body is nil")
			return nil, nil
		}
		resp := &dto.SendMessageRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("SendSessionMessage: %v %v", req, e)
			return nil, e
		} else {
			d.logger.Infof("SendSessionMessage: %v %v", req, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) KickOffUser(req *dto.KickUserReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("KickOffUser: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s%s", d.endpoint, systemUrl, "/user/kickoff")
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
		d.logger.Errorf("KickOffUser: %v %v", req, e)
		return e
	} else {
		d.logger.Infof("KickOffUser: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) QueryUsersOnlineStatus(req *dto.QueryUsersOnlineStatusReq, claims baseDto.ThkClaims) (*dto.QueryUsersOnlineStatusRes, error) {
	uIds := make([]string, 0)
	for _, id := range req.UIds {
		uIds = append(uIds, fmt.Sprintf("%d", id))
	}
	query := "u_ids=" + strings.Join(uIds, ",")
	url := fmt.Sprintf("%s%s%s?%s", d.endpoint, systemUrl, "/user/online", query)
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
			d.logger.Errorf("QueryUsersOnlineStatus: %v %v", req, e)
			return nil, e
		} else {
			return nil, errRes
		}
	} else {
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.Info("QueryUsersOnlineStatus: %v %s", req, "Body is nil")
			return nil, nil
		}
		resp := &dto.QueryUsersOnlineStatusRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("QueryUsersOnlineStatus: %v %v", req, e)
			return nil, e
		} else {
			d.logger.Infof("QueryUsersOnlineStatus: %v %v", req, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) PostUserOnlineStatus(req *dto.PostUserOnlineReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("PostUserOnlineStatus: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s%s", d.endpoint, systemUrl, "/user/online")
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
		d.logger.Errorf("PostUserOnlineStatus: %v %v", req, e)
		return e
	} else {
		d.logger.Infof("PostUserOnlineStatus: %v %s", req, "success")
		return nil
	}
}
