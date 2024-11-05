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

const userSessionUrl = "/user_session"

func (d defaultMsgApi) DeleteUserSession(userId, sessionId int64, claims baseDto.ThkClaims) error {
	url := fmt.Sprintf("%s%s/%d/%d", d.endpoint, userSessionUrl, userId, sessionId)
	request := d.client.R()
	for k, v := range claims {
		vs := v.(string)
		request.SetHeader(k, vs)
	}
	res, errRequest := request.
		SetHeader("Content-Type", jsonContentType).
		Delete(url)
	if errRequest != nil {
		return errRequest
	}
	if res.StatusCode() != http.StatusOK {
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("DeleteUserSession: %v", e)
		return e
	} else {
		d.logger.WithFields(logrus.Fields(claims)).Infof("DeleteUserSession: %v", "success")
		return nil
	}
}

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
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("UpdateUserSession: %v %v", req, e)
		return e
	} else {
		d.logger.WithFields(logrus.Fields(claims)).Infof("UpdateUserSession: %v %s", req, "success")
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
		e := errorx.NewErrorXFromResp(res)
		d.logger.WithFields(logrus.Fields(claims)).Errorf("QueryUserSession: %v %v", req, e)
		return nil, e
	} else {
		resp := &dto.UserSession{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.WithFields(logrus.Fields(claims)).Errorf("QueryUserSession: %v %s", req, e)
			return nil, e
		} else {
			d.logger.WithFields(logrus.Fields(claims)).Infof("QueryUserSession: %v %v", req, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) SearchUserSession(req *dto.SearchUserSessionReq, claims baseDto.ThkClaims) (*dto.SearchUserSessionRes, error) {
	url := fmt.Sprintf("%s%s/search?u_id=%d&offset=%d&count=%d", d.endpoint, userSessionUrl, req.UId, req.Offset, req.Count)
	for _, sessionTypes := range req.Types {
		url += fmt.Sprintf("&types=%d", sessionTypes)
	}
	if req.Keywords != nil {
		url += fmt.Sprintf("&keywords=%s", *req.Keywords)
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
		d.logger.WithFields(logrus.Fields(claims)).Errorf("SearchUserSession: %v %v", req, e)
		return nil, e
	} else {
		resp := &dto.SearchUserSessionRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.WithFields(logrus.Fields(claims)).Errorf("SearchUserSession: %v %s", req, e)
			return nil, e
		} else {
			d.logger.WithFields(logrus.Fields(claims)).Infof("SearchUserSession: %v %v", req, resp)
			return resp, nil
		}
	}
}

func (d defaultMsgApi) QueryLatestUserSession(req *dto.QueryLatestUserSessionReq, claims baseDto.ThkClaims) (*dto.QueryLatestUserSessionsRes, error) {
	url := fmt.Sprintf("%s%s/latest?u_id=%d&m_time=%d&offset=%d&count=%d", d.endpoint, userSessionUrl, req.UId, req.MTime, req.Offset, req.Count)
	for _, sessionTypes := range req.Types {
		url += fmt.Sprintf("&types=%d", sessionTypes)
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
		d.logger.WithFields(logrus.Fields(claims)).Errorf("QueryLatestUserSession: %v %v", req, e)
		return nil, e
	} else {
		resp := &dto.QueryLatestUserSessionsRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.WithFields(logrus.Fields(claims)).Errorf("QueryLatestUserSession: %v %s", req, e)
			return nil, e
		} else {
			d.logger.WithFields(logrus.Fields(claims)).Infof("QueryLatestUserSession: %v %v", req, resp)
			return resp, nil
		}
	}
}
