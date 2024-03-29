package model

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/snowflake"
	"gorm.io/gorm"
	"time"
)

const (
	// MsgTypeRevoke 撤回消息
	MsgTypeRevoke = 100
	// MsgTypeReceived 已接收消息
	MsgTypeReceived = -1
	// MsgTypeRead 已读消息
	MsgTypeRead = -2
	// MsgTypeReedit 重编辑消息
	MsgTypeReedit = -3
)

type (
	SessionMessage struct {
		Id         int64   `gorm:"id" json:"id"`
		MsgId      int64   `gorm:"msg_id" json:"msg_id"`
		ClientId   int64   `gorm:"client_id" json:"client_id"`
		SessionId  int64   `gorm:"session_id" json:"session_id"`
		FromUserId int64   `gorm:"from_user_id" json:"from_user_id"`
		MsgType    int     `gorm:"msg_type" json:"msg_type"`
		MsgContent string  `gorm:"msg_content" json:"msg_content"`
		AtUsers    *string `gorm:"at_users" json:"at_users"`
		ReplyMsgId *int64  `gorm:"reply_msg_id" json:"reply_msg_id"`
		ExtData    *string `gorm:"ext_data" json:"ext_data"`
		CreateTime int64   `gorm:"create_time" json:"create_time"`
		UpdateTime int64   `gorm:"update_time" json:"update_time"`
		Deleted    int8    `gorm:"deleted" json:"deleted"`
	}

	SessionMessageModel interface {
		NewMsgId() int64
		UpdateSessionMessageContent(sessionId, msgId, fUid int64, content string) (int64, error)
		DeleteSessionMessage(sessionId, msgId int64, fUid int64) (int64, error)
		DelMessages(sessionId int64, messageIds []int64, from, to int64) error
		InsertMessage(clientId int64, fromUserId int64, sessionId int64, msgId int64, msgContent string, extData *string,
			msgType int, atUserIds *string, replayMsgId *int64, creatTime int64) (*SessionMessage, error)
		FindMessageByClientId(sessionId, clientId, fromUId int64) (*SessionMessage, error)
		FindSessionMessage(sessionId, msgId, fUid int64) (*SessionMessage, error)
		GetSessionMessages(sessionId, ctime int64, offset, count int, msgIds []int64, asc int8) ([]*SessionMessage, error)
	}

	defaultSessionMessageModel struct {
		shards        int64
		db            *gorm.DB
		logger        *logrus.Entry
		snowflakeNode *snowflake.Node
	}
)

func (d defaultSessionMessageModel) NewMsgId() int64 {
	return d.snowflakeNode.Generate().Int64()
}

func (d defaultSessionMessageModel) UpdateSessionMessageContent(sessionId, msgId, fUid int64, content string) (int64, error) {
	sqlStr := fmt.Sprintf("update %s set msg_content = ?, status , update_time = ?  where session_id = ? and msg_id = ? and from_user_id = ? ", d.genSessionMessageTableName(sessionId))
	tx := d.db.Exec(sqlStr, content, time.Now().UnixMilli(), sessionId, msgId, fUid)
	return tx.RowsAffected, tx.Error
}

func (d defaultSessionMessageModel) DeleteSessionMessage(sessionId, msgId int64, fUid int64) (int64, error) {
	sqlStr := fmt.Sprintf("update %s set deleted = 1 where session_id = ? and msg_id = ? and from_user_id = ? and deleted = 0", d.genSessionMessageTableName(sessionId))
	tx := d.db.Exec(sqlStr, sessionId, msgId, fUid)
	return tx.RowsAffected, tx.Error
}

func (d defaultSessionMessageModel) DelMessages(sessionId int64, messageIds []int64, from, to int64) error {
	if len(messageIds) > 0 {
		sqlStr := fmt.Sprintf("update %s set deleted = 1 where session_id = ? and msg_id in ? and create_time >= ? and create_time <= ? ", d.genSessionMessageTableName(sessionId))
		err := d.db.Exec(sqlStr, sessionId, messageIds, from, to).Error
		return err
	} else {
		sqlStr := fmt.Sprintf("update %s set deleted = 1 where session_id = ? and create_time >= ? and create_time <= ?", d.genSessionMessageTableName(sessionId))
		err := d.db.Exec(sqlStr, sessionId, from, to).Error
		return err
	}
}

func (d defaultSessionMessageModel) InsertMessage(clientId int64, fromUserId int64, sessionId int64, msgId int64,
	msgContent string, extData *string, msgType int, atUserIds *string, replayMsgId *int64, creatTime int64) (*SessionMessage, error) {
	currTime := time.Now().UnixMilli()
	sessionMessage := &SessionMessage{
		MsgId:      msgId,
		ClientId:   clientId,
		SessionId:  sessionId,
		FromUserId: fromUserId,
		AtUsers:    atUserIds,
		MsgType:    msgType,
		ExtData:    extData,
		MsgContent: msgContent,
		ReplyMsgId: replayMsgId,
		CreateTime: creatTime,
		UpdateTime: currTime,
		Deleted:    0,
	}
	tx := d.db.Table(d.genSessionMessageTableName(sessionId)).Create(sessionMessage)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return sessionMessage, nil
}

func (d defaultSessionMessageModel) FindSessionMessage(sessionId, msgId, fUid int64) (*SessionMessage, error) {
	result := &SessionMessage{}
	strSql := "select * from " + d.genSessionMessageTableName(sessionId) + " where session_id = ? and msg_id = ? and from_user_id = ?"
	err := d.db.Raw(strSql, sessionId, msgId, fUid).Scan(result).Error
	return result, err
}

func (d defaultSessionMessageModel) FindMessageByClientId(sessionId, clientId, fromUId int64) (*SessionMessage, error) {
	result := &SessionMessage{}
	strSql := "select * from " + d.genSessionMessageTableName(sessionId) + " where session_id = ? and client_id = ? and from_user_id = ?"
	err := d.db.Raw(strSql, sessionId, clientId, fromUId).Scan(result).Error
	return result, err
}

func (d defaultSessionMessageModel) GetSessionMessages(sessionId, ctime int64, offset, count int, msgIds []int64, asc int8) ([]*SessionMessage, error) {

	sqlBuffer := bytes.NewBufferString("select * from " + d.genSessionMessageTableName(sessionId) + " where session_id = ? ")
	if len(msgIds) > 0 {
		sqlBuffer.WriteString("and msg_id in ? ")
	}
	if asc == 0 { // 降序
		sqlBuffer.WriteString("and create_time <= ? order by create_time desc ")
	} else {
		sqlBuffer.WriteString("and create_time >= ? order by create_time ")
	}
	sqlBuffer.WriteString("limit ?, ? ")
	result := make([]*SessionMessage, 0)
	if len(msgIds) > 0 {
		err := d.db.Raw(sqlBuffer.String(), sessionId, msgIds, ctime, offset, count).Scan(&result).Error
		return result, err
	} else {
		err := d.db.Raw(sqlBuffer.String(), sessionId, ctime, offset, count).Scan(&result).Error
		return result, err
	}
}

func (d defaultSessionMessageModel) genSessionMessageTableName(sessionId int64) string {
	return fmt.Sprintf("session_message_%d", sessionId%(d.shards))
}

func NewSessionMessageModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) SessionMessageModel {
	return defaultSessionMessageModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
