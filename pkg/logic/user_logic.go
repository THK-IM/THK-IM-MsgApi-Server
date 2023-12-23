package logic

import (
	"context"
	"encoding/json"
	"fmt"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/event"
	"github.com/thk-im/thk-im-base-server/utils"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"time"
)

type UserLogic struct {
	appCtx *app.Context
}

func NewUserLogic(appCtx *app.Context) UserLogic {
	return UserLogic{
		appCtx: appCtx,
	}
}

func (l *UserLogic) UpdateUserOnlineStatus(req *dto.PostUserOnlineReq, claims baseDto.ThkClaims) error {
	key := fmt.Sprintf(userOnlineKey, l.appCtx.Config().Name, req.UId)
	var err error
	if req.Online {
		timeout := time.Duration(l.appCtx.Config().IM.OnlineTimeout)
		dtoUserOnlineStatus := &dto.UserOnlineStatus{
			UId:       req.UId,
			Platform:  req.Platform,
			ConnId:    req.ConnId,
			Timestamp: req.Timestamp,
			NodeId:    req.NodeId,
		}
		jsonBytes, errJson := json.Marshal(dtoUserOnlineStatus)
		if errJson != nil {
			return errJson
		}
		err = l.appCtx.RedisCache().Set(context.Background(), key, string(jsonBytes), timeout*time.Second).Err()
	} else {
		js := l.appCtx.RedisCache().Get(context.Background(), key).String()
		if js != "" {
			dtoUserOnlineStatus := &dto.UserOnlineStatus{}
			err = json.Unmarshal([]byte(js), dtoUserOnlineStatus)
			if err == nil {
				if dtoUserOnlineStatus.ConnId == req.ConnId {
					l.appCtx.RedisCache().Del(context.Background(), key)
				}
			}
		}
	}
	return err
}

func (l *UserLogic) GetUsersOnlineStatus(uIds []int64, claims baseDto.ThkClaims) (*dto.QueryUsersOnlineStatusRes, error) {
	uidOnlineKeys := make([]string, 0)
	for _, uid := range uIds {
		uidOnlineKey := fmt.Sprintf(userOnlineKey, l.appCtx.Config().Name, uid)
		uidOnlineKeys = append(uidOnlineKeys, uidOnlineKey)
	}
	onlineUsers, err := utils.BatchGet(l.appCtx.RedisCache(), uidOnlineKeys)
	if err != nil {
		return nil, err
	}
	dtoUsersOnlineStatus := make([]*dto.UserOnlineStatus, 0)
	for _, onlineUser := range onlineUsers {
		if onlineUser == nil {
			continue
		}
		if jsonString, ok := onlineUser.(string); ok {
			userOnlineStatus := &dto.UserOnlineStatus{}
			errJson := json.Unmarshal([]byte(jsonString), userOnlineStatus)
			if errJson == nil {
				dtoUsersOnlineStatus = append(dtoUsersOnlineStatus, userOnlineStatus)
			}
		}
	}
	return &dto.QueryUsersOnlineStatusRes{UsersOnlineStatus: dtoUsersOnlineStatus}, nil
}

func (l *UserLogic) KickUser(req *dto.KickUserReq, claims baseDto.ThkClaims) error {
	ids := []int64{req.UId}
	if idsStr, err := json.Marshal(ids); err != nil {
		return err
	} else {
		msg := make(map[string]interface{})
		msg[event.PushEventTypeKey] = event.SignalKickOffUser
		msg[event.PushEventReceiversKey] = string(idsStr)
		msg[event.PushEventBodyKey] = "kickOff"
		return l.appCtx.MsgPusherPublisher().Pub("", msg)
	}
}
