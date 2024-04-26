package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/event"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	userDto "github.com/thk-im/thk-im-user-server/pkg/dto"
	"strings"
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
	var err error = nil
	var cacheOnlineStatus *dto.UserOnlineStatus = nil
	key := fmt.Sprintf(userOnlineKey, l.appCtx.Config().Name, req.Platform, req.UId)
	cacheStatus, errCache := l.appCtx.RedisCache().Get(context.Background(), key).Result()
	if errCache != nil && !errors.Is(errCache, redis.Nil) {
		return errCache
	}
	if !strings.EqualFold("", cacheStatus) {
		_ = json.Unmarshal([]byte(cacheStatus), &cacheOnlineStatus)
	}
	if req.Online {
		timeout := time.Duration(l.appCtx.Config().IM.OnlineTimeout)
		newCacheOnlineStatus := &dto.UserOnlineStatus{
			UId:         req.UId,
			Platform:    req.Platform,
			ConnId:      req.ConnId,
			TimestampMs: req.TimestampMs,
			NodeId:      req.NodeId,
		}
		jsonBytes, errJson := json.Marshal(newCacheOnlineStatus)
		if errJson != nil {
			return errJson
		}
		err = l.appCtx.RedisCache().Set(context.Background(), key, string(jsonBytes), timeout*time.Second).Err()
		if err == nil {
			minute := (req.TimestampMs - req.LoginTime) / (1000 * 3600)
			if cacheOnlineStatus == nil || ((minute)%10) == 5 { // 缓存未空或，每隔5分钟通知一次
				l.notify(req, claims)
			}
		}
	} else {
		l.appCtx.Logger().Infof("cacheOnlineStatus: %v", cacheOnlineStatus)
		if cacheOnlineStatus != nil && cacheOnlineStatus.ConnId == req.ConnId {
			// 缓存有数据 并且连接id一致，则删除缓存
			err = l.appCtx.RedisCache().Del(context.Background(), key).Err()
			if err == nil {
				l.notify(req, claims)
			}
		} else {
			l.notify(req, claims)
		}
	}
	return err
}

func (l *UserLogic) notify(req *dto.PostUserOnlineReq, claims baseDto.ThkClaims) {
	userOnlineStatusReq := &userDto.UserOnlineStatusReq{
		UserId:      req.UId,
		IsOnline:    req.Online,
		TimestampMs: req.TimestampMs,
		ConnId:      req.ConnId,
		Platform:    req.Platform,
	}
	userApi := l.appCtx.UserApi()
	if userApi != nil {
		_ = l.appCtx.UserApi().PostUserOnlineStatus(userOnlineStatusReq, claims)
	}
}

func (l *UserLogic) GetUsersOnlineStatus(uIds []int64, claims baseDto.ThkClaims) (*dto.QueryUsersOnlineStatusRes, error) {
	uidOnlineKeys := make([]string, 0)
	for _, uid := range uIds {
		uidOnlineKeys = append(uidOnlineKeys, fmt.Sprintf(userOnlineKey, l.appCtx.Config().Name, PlatformAndroid, uid))
		uidOnlineKeys = append(uidOnlineKeys, fmt.Sprintf(userOnlineKey, l.appCtx.Config().Name, PlatformIOS, uid))
		uidOnlineKeys = append(uidOnlineKeys, fmt.Sprintf(userOnlineKey, l.appCtx.Config().Name, PlatformWeb, uid))
	}
	onlineUsers, err := l.appCtx.RedisCache().MGet(context.Background(), uidOnlineKeys...).Result()
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
	if idsStr, err := json.Marshal(req.UIds); err != nil {
		return err
	} else {
		msg := make(map[string]interface{})
		msg[event.PushEventTypeKey] = event.SignalKickOffUser
		msg[event.PushEventReceiversKey] = string(idsStr)
		msg[event.PushEventBodyKey] = "kickOff"
		return l.appCtx.MsgPusherPublisher().Pub("", msg)
	}
}
