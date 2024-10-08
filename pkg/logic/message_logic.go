package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseErrorx "github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-base-server/event"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"time"
)

type MessageLogic struct {
	appCtx *app.Context
}

func NewMessageLogic(appCtx *app.Context) MessageLogic {
	return MessageLogic{
		appCtx: appCtx,
	}
}

func (l *MessageLogic) convSessionMessage2Message(sessionMsg *model.SessionMessage) *dto.Message {
	vo := dto.Message{
		CId:     sessionMsg.ClientId,
		FUid:    sessionMsg.FromUserId,
		SId:     sessionMsg.SessionId,
		MsgId:   sessionMsg.MsgId,
		CTime:   sessionMsg.CreateTime,
		Body:    sessionMsg.MsgContent,
		AtUsers: sessionMsg.AtUsers,
		ExtData: sessionMsg.ExtData,
		Type:    sessionMsg.MsgType,
		RMsgId:  sessionMsg.ReplyMsgId,
	}
	return &vo
}

func (l *MessageLogic) convUserMessage2Message(userMsg *model.UserMessage) *dto.Message {
	msg := dto.Message{
		CId:     userMsg.ClientId,
		SId:     userMsg.SessionId,
		Type:    userMsg.MsgType,
		MsgId:   userMsg.MsgId,
		FUid:    userMsg.FromUserId,
		CTime:   userMsg.CreateTime,
		RMsgId:  userMsg.ReplyMsgId,
		Body:    userMsg.MsgContent,
		ExtData: userMsg.ExtData,
		Status:  &userMsg.Status,
		AtUsers: userMsg.AtUsers,
	}
	return &msg
}

func (l *MessageLogic) GetUserMessages(req dto.GetMessageReq, claims baseDto.ThkClaims) (*dto.GetMessageRes, error) {
	userMessages, err := l.appCtx.UserMessageModel().GetUserMessages(req.UId, req.CTime, req.Offset, req.Count)
	if err != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("GetUserMessages %v, %v", req, err)
		return nil, err
	}
	messages := make([]*dto.Message, 0)
	for _, userMessage := range userMessages {
		message := l.convUserMessage2Message(userMessage)
		messages = append(messages, message)
	}
	return &dto.GetMessageRes{Data: messages}, nil
}

func (l *MessageLogic) GetSessionMessages(req dto.GetSessionMessageReq, claims baseDto.ThkClaims) (*dto.GetMessageRes, error) {
	msgIds := make([]int64, 0)
	if req.MsgIds != "" {
		strIds := strings.Split(req.MsgIds, ",")
		for _, str := range strIds {
			if id, err := strconv.ParseInt(str, 10, 64); err != nil {
				l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("GetSessionMessages ParseInt err: %s %v", strIds, err)
				return nil, err
			} else {
				msgIds = append(msgIds, id)
			}
		}
	}
	sessionMessages, err := l.appCtx.SessionMessageModel().GetSessionMessages(req.SId, req.CTime, req.Offset, req.Count, msgIds, req.Asc)
	if err != nil {
		return nil, err
	}
	messages := make([]*dto.Message, 0)
	for _, sessionMessage := range sessionMessages {
		message := l.convSessionMessage2Message(sessionMessage)
		messages = append(messages, message)
	}
	return &dto.GetMessageRes{Data: messages}, nil
}

func (l *MessageLogic) DelSessionMessage(req *dto.DelSessionMessageReq, claims baseDto.ThkClaims) error {
	err := l.appCtx.SessionMessageModel().DelMessages(req.SId, req.MsgIds, req.TimeFrom, req.TimeTo)
	if err != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("DelSessionMessage err: %v %v", req, err)
	}
	return err
}

func (l *MessageLogic) SendMessage(req dto.SendMessageReq, claims baseDto.ThkClaims) (*dto.SendMessageRes, error) {
	session, errSession := l.appCtx.SessionModel().FindSession(req.SId)
	if errSession != nil || session.Id <= 0 {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendMessage FindSession %v, %v", req, errSession)
		return nil, errorx.ErrSessionInvalid
	}
	// req.FUid为0是系统消息, 不需要校验是否能对session发送消息
	if req.FUid > 0 {
		userSession, errUserSession := l.appCtx.UserSessionModel().GetUserSession(req.FUid, req.SId)
		if errUserSession != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendMessage GetUserSession %v, %v", req, errUserSession)
			return nil, errorx.ErrSessionInvalid
		}
		if userSession.Deleted == 1 || userSession.UserId == 0 {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendMessage GetUserSession %v, %v", req, userSession)
			return nil, errorx.ErrSessionInvalid
		}
		msgCheckApi := l.appCtx.MessageCheckApi()
		if msgCheckApi != nil {
			// 检查消息是否可以发送[内容检测/建联逻辑等检查]
			checkReq := &dto.CheckMessageReq{
				SessionType:    session.Type,
				SessionId:      session.Id,
				FunctionFlag:   session.FunctionFlag,
				FromUId:        req.FUid,
				MessageType:    req.Type,
				MessageContent: req.Body,
				EntityId:       userSession.EntityId,
			}
			errCheck := msgCheckApi.CheckMessage(checkReq, claims)
			if errCheck != nil {
				l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendMessage CheckMessage %v, %v", checkReq, errCheck)
				return nil, errCheck
			}
		}
		if userSession.Mute&model.MutedSingleBitInUserSessionStatus > 0 {
			return nil, errorx.ErrUserMuted
		} else if userSession.Mute&model.MutedAllBitInUserSessionStatus > 0 && userSession.Role < model.SessionSuperAdmin {
			// 如果是超管或者是群主，全员被禁言情况下仍允许发言
			return nil, errorx.ErrSessionMuted
		}
	}

	// 如果是超级群 读扩散模型，写入session_message表
	if session.Type == model.SuperGroupSessionType {
		return l.SendSessionMessage(session, req, claims)
	}
	// 其他情况使用使用写扩散模型，写入user_message表
	return l.SendUserMessage(session, req, claims)
}

func (l *MessageLogic) SendSessionMessage(session *model.Session, req dto.SendMessageReq, claims baseDto.ThkClaims) (*dto.SendMessageRes, error) {
	receivers := l.appCtx.SessionUserModel().FindUIdsInSessionWithoutStatus(req.SId, model.RejectBitInUserSessionStatus, req.Receivers)
	if receivers == nil || len(receivers) == 0 {
		return nil, errorx.ErrUserReject
	}
	// 根据clientId和fromUserId查询是否已经发送过消息
	sessionMessage, errMessage := l.appCtx.SessionMessageModel().FindMessageByClientId(req.SId, req.CId, req.FUid)
	// 如果已经发送过，直接取数据库里的数据库, 没有发送过则插入数据库
	if sessionMessage == nil || sessionMessage.MsgId == 0 {
		// 插入数据库发送消息
		msgId := int64(l.appCtx.SnowflakeNode().Generate())
		sessionMessage, errMessage = l.appCtx.SessionMessageModel().InsertMessage(
			req.CId, req.FUid, req.SId, msgId, req.Body, req.ExtData, req.Type, req.AtUsers, req.RMsgId, req.CTime)
		if errMessage != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendMessage InsertMessage %v, %v", errMessage, req)
			return nil, errMessage
		}
	}

	dtoMsg := l.convSessionMessage2Message(sessionMessage)
	offlineReceiverIds := make([]int64, 0)
	receiverUIds := make([]int64, 0)
	for _, r := range receivers {
		receiverUIds = append(receiverUIds, r.UserId)
		if r.Status&model.SilenceBitInUserSessionStatus == 0 {
			offlineReceiverIds = append(offlineReceiverIds, r.UserId)
		}
	}
	if onlineUIds, offlineUIds, err := l.publishSendMessageEvents(dtoMsg, session.Type, receiverUIds, offlineReceiverIds, claims); err != nil {
		return nil, errorx.ErrMessageDeliveryFailed
	} else {
		return &dto.SendMessageRes{
			MsgId:      sessionMessage.MsgId,
			CreateTime: sessionMessage.CreateTime,
			OnlineIds:  onlineUIds,
			OfflineIds: offlineUIds,
		}, nil
	}
}

func (l *MessageLogic) SendUserMessage(session *model.Session, req dto.SendMessageReq, claims baseDto.ThkClaims) (*dto.SendMessageRes, error) {
	receivers := l.appCtx.SessionUserModel().FindUIdsInSessionWithoutStatus(req.SId, model.RejectBitInUserSessionStatus, req.Receivers)
	if receivers == nil || len(receivers) == 0 {
		return nil, errorx.ErrUserReject
	}
	// 根据clientId和fromUserId查询是否已经发送过消息
	userMessage, errMessage := l.appCtx.UserMessageModel().FindUserMessageByClientId(req.FUid, req.SId, req.CId)
	// 如果已经发送过，直接取数据库里的数据库, 没有发送过则插入数据库
	if userMessage == nil || userMessage.MsgId == 0 {
		// 插入数据库发送消息
		msgId := int64(l.appCtx.SnowflakeNode().Generate())
		now := time.Now().UnixMilli()
		userMessage = &model.UserMessage{
			MsgId:      msgId,
			ClientId:   req.CId,
			UserId:     req.FUid,
			SessionId:  req.SId,
			FromUserId: req.FUid,
			MsgType:    req.Type,
			MsgContent: req.Body,
			ReplyMsgId: req.RMsgId,
			AtUsers:    req.AtUsers,
			ExtData:    req.ExtData,
			Status:     model.MsgStatusAcked | model.MsgStatusRead,
			CreateTime: req.CTime,
			UpdateTime: now,
		}
		if userMessage.FromUserId > 0 { // fromUserId = 0为系统发给用户的消息，不用插入数据库
			// 插入发件人消息表
			errMessage = l.appCtx.UserMessageModel().InsertUserMessage(userMessage)
			if errMessage != nil {
				l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("SendMessage InsertMessage %v, %v", errMessage, req)
				return nil, errMessage
			}
		}
	}
	userMessage.Status = model.MsgStatusInit
	dtoMsg := l.convUserMessage2Message(userMessage)
	offlineReceiverIds := make([]int64, 0)
	receiverUIds := make([]int64, 0)
	for _, r := range receivers {
		receiverUIds = append(receiverUIds, r.UserId)
		if r.Status&model.SilenceBitInUserSessionStatus == 0 {
			offlineReceiverIds = append(offlineReceiverIds, r.UserId)
		}
	}
	if onlineUIds, offlineUIds, err := l.publishSendMessageEvents(dtoMsg, session.Type, receiverUIds, offlineReceiverIds, claims); err != nil {
		return nil, errorx.ErrMessageDeliveryFailed
	} else {
		return &dto.SendMessageRes{
			MsgId:      userMessage.MsgId,
			CreateTime: userMessage.CreateTime,
			OnlineIds:  onlineUIds,
			OfflineIds: offlineUIds,
		}, nil
	}
}

func (l *MessageLogic) SendSysMessage(req dto.SendSysMessageReq, claims baseDto.ThkClaims) (*dto.SendSysMessageRes, error) {
	if req.Receivers == nil || len(req.Receivers) == 0 {
		return nil, baseErrorx.ErrParamsError
	}
	msgId := l.appCtx.SessionMessageModel().NewMsgId()
	now := time.Now().UnixMilli()
	sessionType := 0
	sessionMessage := &model.SessionMessage{
		MsgId:      msgId,
		ClientId:   msgId,
		SessionId:  0,
		FromUserId: 0,
		MsgType:    req.Type,
		MsgContent: req.Body,
		AtUsers:    nil,
		ReplyMsgId: nil,
		ExtData:    req.ExtData,
		CreateTime: now,
		UpdateTime: now,
		Deleted:    0,
	}
	dtoMsg := l.convSessionMessage2Message(sessionMessage)
	if onlineUIds, offlineUIds, err := l.publishSendMessageEvents(dtoMsg, sessionType, req.Receivers, nil, claims); err != nil {
		return nil, errorx.ErrMessageDeliveryFailed
	} else {
		return &dto.SendSysMessageRes{
			MsgId:      msgId,
			CreateTime: now,
			OnlineIds:  onlineUIds,
			OfflineIds: offlineUIds,
		}, nil
	}

}

func (l *MessageLogic) publishSendMessageEvents(dtoMsg *dto.Message, sessionType int, receivers []int64, offlineReceivers []int64, claims baseDto.ThkClaims) ([]int64, []int64, error) {
	msgJson, err := json.Marshal(dtoMsg)
	if err != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("publishSendMessageEvents json err", err)
		return nil, nil, err
	}
	msgJsonStr := string(msgJson)
	deliverKey := fmt.Sprintf("session-%d", dtoMsg.SId)
	onlineUIds, offlineUIds, errPubPush := l.pubPushMessageEvent(event.SignalNewMessage, msgJsonStr, receivers, offlineReceivers, deliverKey, true, claims)
	if errPubPush != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("pubPushMessageEvent, publish err:", errPubPush)
		return nil, nil, errPubPush
	}
	if sessionType != model.SuperGroupSessionType {
		errPubSave := l.pubSaveMsgEvent(msgJsonStr, receivers, dtoMsg.SId)
		if errPubSave != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Error("pubSaveMsgEvent, err:", errPubSave)
			return nil, nil, errPubPush
		}
	}
	return onlineUIds, offlineUIds, nil
}

// PushMessage 业务消息推送
func (l *MessageLogic) PushMessage(req dto.PushMessageReq, claims baseDto.ThkClaims) (*dto.PushMessageRes, error) {
	deliverKey := "push"
	// 在线推送
	onlineUIds, offlineUIds, err := l.pubPushMessageEvent(req.Type, req.Body, req.UIds, req.UIds, deliverKey, req.OfflinePush, claims)
	if err == nil {
		rsp := &dto.PushMessageRes{}
		rsp.OnlineUIds = onlineUIds
		rsp.OfflineUIds = offlineUIds
		return rsp, err
	} else {
		return nil, err
	}
}

func (l *MessageLogic) pubSaveMsgEvent(msgBody string, receivers []int64, sessionId int64) error {
	if receiversStr, errJson := json.Marshal(receivers); errJson != nil {
		return errJson
	} else {
		m := make(map[string]interface{})
		m[event.SaveMsgEventKey] = msgBody
		m[event.SaveMsgUsersKey] = receiversStr
		return l.appCtx.MsgSaverPublisher().Pub(fmt.Sprintf("session-%d", sessionId), m)
	}
}

// 发布推送消息
func (l *MessageLogic) pubPushMessageEvent(t int, body string, uIds []int64, offlinePushUIds []int64, deliverKey string, offlinePushTag bool, claims baseDto.ThkClaims) ([]int64, []int64, error) {
	uidOnlineKeys := make([]string, 0)
	for _, uid := range uIds {
		uidOnlineKeys = append(uidOnlineKeys, fmt.Sprintf(userOnlineKey, l.appCtx.Config().Name, PlatformAndroid, uid))
		uidOnlineKeys = append(uidOnlineKeys, fmt.Sprintf(userOnlineKey, l.appCtx.Config().Name, PlatformIOS, uid))
		uidOnlineKeys = append(uidOnlineKeys, fmt.Sprintf(userOnlineKey, l.appCtx.Config().Name, PlatformWeb, uid))
	}
	cacheOnLineUIds := make([]int64, 0)
	redisOnlineUsers, err := l.appCtx.RedisCache().MGet(context.Background(), uidOnlineKeys...).Result()
	if err != nil {
		// 如果查询报错 默认全部用户为离线
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("pubPushMessageEvents error: %v", err)
	} else {
		for index, redisOnlineUser := range redisOnlineUsers {
			if redisOnlineUser != nil {
				cacheOnLineUIds = append(cacheOnLineUIds, uIds[index/3]) // 3代表三个平台android/ios/web
			}
		}
	}

	onlineUIdMap := make(map[int64]bool)
	for _, uid := range cacheOnLineUIds {
		onlineUIdMap[uid] = true
	}

	onlineUIds := make([]int64, 0)
	offlineUIds := make([]int64, 0)
	for _, uid := range uIds {
		online := onlineUIdMap[uid]
		if !online {
			if !slices.Contains(offlineUIds, uid) {
				offlineUIds = append(offlineUIds, uid)
			}
		} else {
			onlineUIds = append(onlineUIds, uid)
		}
	}
	receiverStr, errJson := json.Marshal(onlineUIds)
	if errJson != nil {
		return nil, nil, errJson
	}
	m := make(map[string]interface{})
	m[event.PushEventTypeKey] = t
	m[event.PushEventBodyKey] = body
	m[event.PushEventReceiversKey] = string(receiverStr)
	err = l.appCtx.MsgPusherPublisher().Pub(deliverKey, m)
	if err != nil {
		return nil, nil, err
	}

	if offlinePushTag && len(offlineUIds) > 0 {
		receiverStr, errJson = json.Marshal(offlineUIds)
		if errJson != nil {
			return nil, nil, errJson
		}
		offlineEvent := make(map[string]interface{})
		offlineEvent[event.PushEventTypeKey] = t
		offlineEvent[event.PushEventBodyKey] = body
		offlineEvent[event.PushEventReceiversKey] = string(receiverStr)
		err = l.appCtx.MsgOfflinePusherPublisher().Pub(deliverKey, offlineEvent)
	}

	return onlineUIds, offlineUIds, err
}

func (l *MessageLogic) DeleteUserMessage(req *dto.DeleteMessageReq, claims baseDto.ThkClaims) error {
	err := l.appCtx.UserMessageModel().DeleteMessages(req.UId, req.SId, req.MessageIds, req.TimeFrom, req.TimeTo)
	if err != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("DeleteUserMessage err: %v %v", req, err)
	}
	return err
}
