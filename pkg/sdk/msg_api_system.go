package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"net/http"
	"strings"
)

const (
	systemUrl = "/system"
)

func (d defaultMsgApi) PushExtendedMessage(req *dto.PushMessageReq) (*dto.PushMessageRes, error) {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("PushExtendedMessage: %v %v", req, err)
		return nil, err
	}
	url := fmt.Sprintf("%s%s%s", d.endpoint, systemUrl, "/message/push")
	res, errRequest := d.client.R().
		SetHeader("Content-Type", jsonContentType).
		SetBody(dataBytes).
		Post(url)
	if errRequest != nil {
		d.logger.Errorf("PushExtendedMessage: %v %v", req, errRequest)
		return nil, errRequest
	}
	if res.StatusCode() != http.StatusOK {
		errRes := &errorx.ErrorX{}
		e := json.Unmarshal(res.Body(), errRes)
		if e != nil {
			d.logger.Errorf("PushExtendedMessage: %v %v", req, e)
			return nil, e
		} else {
			return nil, errRes
		}
	} else {
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.Infof("PushExtendedMessage: %v %s", req, "Body is nil")
			return nil, nil
		}
		resp := &dto.PushMessageRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("PushExtendedMessage: %v %v", req, e)
			return nil, e
		} else {
			d.logger.Infof("PushExtendedMessage: %v %v", req, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) SendSystemMessage(req *dto.SendMessageReq) (*dto.SendMessageRes, error) {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("SendSystemMessage: %v %v", req, err)
		return nil, err
	}
	url := fmt.Sprintf("%s%s%s", d.endpoint, systemUrl, "/message/send")
	res, errRequest := d.client.R().
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
		resp := &dto.SendMessageRes{}
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

func (d defaultMsgApi) KickOffUser(req *dto.KickUserReq) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("KickOffUser: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s%s", d.endpoint, systemUrl, "/user/kickoff")
	res, errRequest := d.client.R().
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
		d.logger.Errorf("KickOffUser: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) QueryUsersOnlineStatus(req *dto.QueryUsersOnlineStatusReq) (*dto.QueryUsersOnlineStatusRes, error) {
	uIds := make([]string, 0)
	for _, id := range req.UIds {
		uIds = append(uIds, fmt.Sprintf("%d", id))
	}
	query := "u_ids=" + strings.Join(uIds, ",")
	url := fmt.Sprintf("%s%s%s?%s", d.endpoint, systemUrl, "/user/online", query)
	res, errRequest := d.client.R().
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
