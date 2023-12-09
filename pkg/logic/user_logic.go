package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/event"
	"github.com/thk-im/thk-im-base-server/rpc"
	"github.com/thk-im/thk-im-base-server/utils"
	"github.com/thk-im/thk-im-msg-api-server/pkg/app"
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

func (l *UserLogic) UpdateUserOnlineStatus(req *dto.PostUserOnlineReq) error {
	key := fmt.Sprintf(userOnlineKey, l.appCtx.Config().Name, req.UId)
	var err error
	if req.Online {
		timeout := time.Duration(l.appCtx.Config().IM.OnlineTimeout)
		err = l.appCtx.RedisCache().Set(context.Background(), key, req.ConnId, timeout*time.Second).Err()
	} else {
		_, err = utils.DelKeyByValue(l.appCtx.RedisCache(), key, req.ConnId)
	}
	if req.IsLogin { // 登录写库登录记录
		err = l.appCtx.UserOnlineStatusModel().UpdateUserOnlineStatus(req.UId, req.Timestamp, req.ConnId, req.Platform)
		go func() {
			onlineReq := rpc.PostUserOnlineReq{
				UserId:    req.UId,
				IsOnline:  req.Online,
				Timestamp: req.Timestamp,
				ConnId:    req.ConnId,
				Platform:  req.Platform,
			}
			if l.appCtx.RpcUserApi() != nil {
				if e := l.appCtx.RpcUserApi().PostUserOnlineStatus(onlineReq); e != nil {
					l.appCtx.Logger().Errorf("UpdateUserOnlineStatus, RpcUserApi, call err: %s", e.Error())
				}
			}
		}()
	}
	return err
}

func (l *UserLogic) GetUsersOnlineStatus(uIds []int64) (*dto.GetUsersOnlineStatusRes, error) {
	usersOnlineStatus, err := l.appCtx.UserOnlineStatusModel().GetUsersOnlineStatus(uIds)
	if err != nil {
		return nil, err
	} else {
		dtoUsersOnlineStatus := make([]*dto.UserOnlineStatus, 0)
		for _, user := range usersOnlineStatus {
			dtoUserOnlineStatus := &dto.UserOnlineStatus{
				UId:            user.UserId,
				Platform:       user.Platform,
				LastOnlineTime: user.OnlineTime,
			}
			dtoUsersOnlineStatus = append(dtoUsersOnlineStatus, dtoUserOnlineStatus)
		}
		return &dto.GetUsersOnlineStatusRes{UsersOnlineStatus: dtoUsersOnlineStatus}, nil
	}
}

func (l *UserLogic) user() {

}

func (l *UserLogic) KickUser(req *dto.KickUserReq) error {
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
