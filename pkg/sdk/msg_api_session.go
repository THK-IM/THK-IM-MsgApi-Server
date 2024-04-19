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

const (
	sessionUrl = "/session"
)

func (d defaultMsgApi) QuerySessionUserCount(sessionId int64, claims baseDto.ThkClaims) (*dto.SessionUserCountRes, error) {
	url := fmt.Sprintf("%s%s/%d/user/count", d.endpoint, sessionUrl, sessionId)
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
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("QuerySessionUserCount: %v %v", sessionId, e)
		return nil, e
	} else {
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.WithFields(logrus.Fields(claims)).Info("QuerySessionUserCount: %v %s", sessionId, "Body is nil")
			return nil, nil
		}
		resp := &dto.SessionUserCountRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.WithFields(logrus.Fields(claims)).Errorf("QuerySessionUserCount: %v %v", sessionId, e)
			return nil, e
		} else {
			d.logger.WithFields(logrus.Fields(claims)).Infof("QuerySessionUserCount: %v %v", sessionId, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) DelSessionUser(sessionId int64, req *dto.SessionDelUserReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.WithFields(logrus.Fields(claims)).Errorf("DelSessionUser: %v %v", req, err)
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
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("DelSessionUser: %v %v", req, e)
		return e
	} else {
		d.logger.WithFields(logrus.Fields(claims)).Infof("DelSessionUser: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) AddSessionUser(sessionId int64, req *dto.SessionAddUserReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.WithFields(logrus.Fields(claims)).Errorf("AddSessionUser: %v %v", req, err)
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
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("AddSessionUser: %v %v", req, e)
		return e
	} else {
		d.logger.WithFields(logrus.Fields(claims)).Infof("AddSessionUser: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) DelSession(sessionId int64, req *dto.DelSessionReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.WithFields(logrus.Fields(claims)).Errorf("DelSession: %v %v", req, err)
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
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("DelSession: %v %v", req, e)
		return e
	} else {
		d.logger.WithFields(logrus.Fields(claims)).Infof("DelSession: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) UpdateSession(sessionId int64, req *dto.UpdateSessionReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.WithFields(logrus.Fields(claims)).Errorf("UpdateSession: %v %v", req, err)
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
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("UpdateSession: %v %v", req, e)
		return e
	} else {
		d.logger.WithFields(logrus.Fields(claims)).Infof("UpdateSession: %v %s", req, "success")
		return nil
	}
}

func (d defaultMsgApi) QueryLatestSessionUsers(sessionId int64, req *dto.QuerySessionUsersReq, claims baseDto.ThkClaims) (*dto.QuerySessionUsersRes, error) {
	url := fmt.Sprintf("%s%s/%d/user/latest?count=%d&m_time=%d", d.endpoint, sessionUrl, sessionId, req.Count, req.MTime)
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
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("QueryLatestSessionUsers: %v %v", req, e)
		return nil, e
	} else {
		if res.Body() == nil || len(res.Body()) == 0 {
			d.logger.WithFields(logrus.Fields(claims)).Info("QueryLatestSessionUsers: %v %s", req, "Body is nil")
			return nil, nil
		}
		resp := &dto.QuerySessionUsersRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.WithFields(logrus.Fields(claims)).Errorf("QueryLatestSessionUsers: %v %v", req, e)
			return nil, e
		} else {
			d.logger.WithFields(logrus.Fields(claims)).Infof("QueryLatestSessionUsers: %v %v", req, resp)
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
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("QuerySessionUser: %d %d %v", sessionId, userId, e)
		return nil, e
	} else {
		resp := &dto.SessionUser{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			return nil, nil
		} else {
			d.logger.WithFields(logrus.Fields(claims)).Infof("QuerySessionUser: %d %d %v", sessionId, userId, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) UpdateSessionUser(sessionId int64, req *dto.SessionUserUpdateReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.WithFields(logrus.Fields(claims)).Errorf("UpdateSessionUser: %v %v", req, err)
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
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("UpdateSessionUser: %v %v", req, e)
		return e
	} else {
		d.logger.WithFields(logrus.Fields(claims)).Infof("UpdateSessionUser: %v %s", req, "success")
		return nil
	}
}
