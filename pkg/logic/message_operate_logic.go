package logic

import (
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"time"
)

func (l *MessageLogic) AckUserMessages(req dto.AckUserMessagesReq, claims baseDto.ThkClaims) error {
	err := l.appCtx.UserMessageModel().AckUserMessages(req.UId, req.SId, req.MsgIds)
	if err != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("AckUserMessages err:%v, %v", req, err)
	}
	return err
}

func (l *MessageLogic) ReadUserMessages(req dto.ReadUserMessageReq, claims baseDto.ThkClaims) error {
	session, errSession := l.appCtx.SessionModel().FindSession(req.SId)
	if errSession != nil {
		return errorx.ErrSessionInvalid
	}
	if session.Type == model.SuperGroupSessionType {
		// 更新session的update_time，客户端同步消息时超级群消息时间小于update_time时认为已读
		return l.appCtx.UserSessionModel().UpdateUserSession([]int64{req.UId}, session.Id, nil,
			nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	} else {
		errUpdate := l.appCtx.UserMessageModel().UpdateUserMessage(req.UId, req.SId, req.MsgIds, model.MsgStatusRead, nil)
		if errUpdate != nil {
			return errUpdate
		}
		if session.FunctionFlag&dto.ReadFlag > 0 {
			if userMessages, err := l.appCtx.UserMessageModel().FindUserMessages(req.UId, req.SId, req.MsgIds); err == nil {
				for _, userMessage := range userMessages {
					if userMessage.MsgId == 0 {
						return errorx.ErrSessionMessageInvalid
					}
					if userMessage.MsgType < 0 || userMessage.Status&model.MsgStatusRead == 1 { // 小于0的类型消息为状态操作消息或者已经是已读了，不需要发送已读
						continue
					}
					sendMessageReq := dto.SendMessageReq{
						CId:       l.genClientId(),
						SId:       req.SId,
						Type:      model.MsgTypeRead,
						FUid:      req.UId,
						CTime:     time.Now().UnixMilli(),
						RMsgId:    &userMessage.MsgId,
						Receivers: []int64{userMessage.FromUserId, req.UId}, // 发送给对方和自己
					}
					// 对消息发件人发送已读消息
					if _, err = l.SendUserMessage(session, sendMessageReq, claims); err != nil {
						l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("ReadUserMessages err:%v, %v", req, err)
					}
				}
			} else {
				l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("ReadUserMessages err:%v, %v", req, err)
			}
		}
	}
	return nil
}

func (l *MessageLogic) RevokeUserMessage(req dto.RevokeUserMessageReq, claims baseDto.ThkClaims) error {
	if userMessage, err := l.appCtx.UserMessageModel().FindUserMessage(req.UId, req.SId, req.MsgId); err == nil {
		if userMessage.SessionId == 0 || userMessage.Deleted == 1 {
			return errorx.ErrSessionMessageInvalid
		}
		if userMessage.MsgType < 0 { // 小于0的类型消息为状态操作消息，不能重新编辑
			return errorx.ErrMessageTypeNotSupport
		}
		if userMessage.Deleted == 1 { // 被删除了则不做处理
			return nil
		}
		sendMessageReq := dto.SendMessageReq{
			CId:    l.genClientId(),
			SId:    req.SId,
			Type:   model.MsgTypeRevoke,
			FUid:   req.UId,
			CTime:  time.Now().UnixMilli(),
			RMsgId: &req.MsgId,
		} // 发送给session下的所有人
		if _, err = l.SendMessage(sendMessageReq, claims); err != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("RevokeUserMessage err:%v, %v", req, err)
			return err
		}
		return nil
	} else {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("RevokeUserMessage err:%v, %v", req, err)
		return err
	}
}

func (l *MessageLogic) ReeditUserMessage(req dto.ReeditUserMessageReq, claims baseDto.ThkClaims) error {
	session, errSession := l.appCtx.SessionModel().FindSession(req.SId)
	if errSession != nil {
		return errorx.ErrSessionInvalid
	}
	if session.Type == model.SuperGroupSessionType {
		sessionMessage, err := l.appCtx.SessionMessageModel().FindSessionMessage(req.SId, req.MsgId, req.UId)
		if err != nil {
			return err
		}
		if sessionMessage == nil {
			return errorx.ErrSessionMessageInvalid
		}
	} else {
		userMessage, err := l.appCtx.UserMessageModel().FindUserMessage(req.UId, req.SId, req.MsgId)
		if err == nil {
			if userMessage.SessionId == 0 || userMessage.Deleted == 1 {
				return errorx.ErrSessionMessageInvalid
			}
			if userMessage.MsgType < 0 { // 小于0的类型消息为状态操作消息，不能重新编辑
				return errorx.ErrMessageTypeNotSupport
			}
		} else {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("ReeditUserMessage err:%v, %v", req, err)
		}
	}

	sendMessageReq := dto.SendMessageReq{
		CId:   l.genClientId(),
		SId:   req.SId,
		Type:  model.MsgTypeReedit,
		FUid:  req.UId,
		CTime: time.Now().UnixMilli(),
		Body:  req.Content,
	}
	// 发送给session下的所有人
	if _, err := l.SendMessage(sendMessageReq, claims); err != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("ReeditUserMessage err:%v %v", req, err)
		return err
	}
	return nil
}

func (l *MessageLogic) ForwardUserMessages(req dto.ForwardUserMessageReq, claims baseDto.ThkClaims) (*dto.SendMessageRes, error) {
	if len(req.ForwardClientIds) > 0 && len(req.ForwardClientIds) > 0 {
		ids, err := l.appCtx.SessionObjectModel().AddSessionObjects(req.ForwardSId, req.ForwardFromUIds, req.ForwardClientIds, req.FUid, req.CId, req.SId)
		if err != nil {
			return nil, err
		}
		if ids != nil {
			if err = l.appCtx.ObjectModel().AddSessions(ids, req.SId); err != nil {
				return nil, err
			}
		}
	}
	return l.SendMessage(req.SendMessageReq, claims)
}

func (l *MessageLogic) genClientId() int64 {
	return l.appCtx.SnowflakeNode().Generate().Int64()
}
