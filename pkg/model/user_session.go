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
	// MutedAllBitInUserSessionStatus 全员被禁言标志位
	MutedAllBitInUserSessionStatus = 1 << 0
	// MutedSingleBitInUserSessionStatus 用户被禁言标志位
	MutedSingleBitInUserSessionStatus = 1 << 1
	// RejectBitInUserSessionStatus 拒收标志位
	RejectBitInUserSessionStatus = 1 << 0
	// SilenceBitInUserSessionStatus 静音标志位
	SilenceBitInUserSessionStatus = 1 << 1
)

type (
	UserSession struct {
		Id         int64   `gorm:"id" json:"id"`
		SessionId  int64   `gorm:"session_id" json:"session_id"`
		UserId     int64   `gorm:"user_id" json:"user_id"`
		ParentId   int64   `gorm:"parent_id" json:"parent_id"`
		Type       int     `gorm:"type" json:"type"`
		EntityId   int64   `gorm:"entity_id" json:"entity_id"`
		Name       string  `gorm:"name" json:"name"`
		Remark     string  `gorm:"remark" json:"remark"`
		ExtData    *string `json:"ext_data" json:"ext_data"`
		Top        int64   `gorm:"top" json:"top"`
		Role       int     `gorm:"role" json:"role"`
		Mute       int     `gorm:"mute" json:"mute"`
		Status     int     `gorm:"status" json:"status"`
		NoteName   string  `gorm:"note_name" json:"note_name"`
		CreateTime int64   `gorm:"create_time" json:"create_time"`
		UpdateTime int64   `gorm:"update_time" json:"update_time"`
		Deleted    int8    `gorm:"deleted" json:"deleted"`
	}

	UserSessionModel interface {
		FindUserSessionByEntityId(userId, entityId int64, sessionType int, containDeleted bool) (*UserSession, error)
		UpdateUserSessionType(userIds []int64, sessionId int64, sessionType int) error
		UpdateUserSession(userIds []int64, sessionId int64, sessionName, sessionRemark, mute, extData, noteName *string, top *int64, status, role *int, parentId *int64) error
		FindEntityIdsInUserSession(userId, sessionId int64) []int64
		QueryLatestUserSessions(userId, mTime int64, offset, count int, types []int) ([]*UserSession, error)
		GetUserSession(userId, sessionId int64) (*UserSession, error)
		GenUserSessionTableName(userId int64) string
	}

	defaultUserSessionModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultUserSessionModel) FindUserSessionByEntityId(userId, entityId int64, sessionType int, containDeleted bool) (*UserSession, error) {
	userSession := &UserSession{}
	sqlStr := "select * from " + d.GenUserSessionTableName(userId) + " where user_id = ? and entity_id = ? and type = ?"
	if !containDeleted {
		sqlStr += " and deleted = 0"
	}
	err := d.db.Raw(sqlStr, userId, entityId, sessionType).Scan(&userSession).Error
	return userSession, err
}

func (d defaultUserSessionModel) UpdateUserSessionType(userIds []int64, sessionId int64, sessionType int) (err error) {
	// 分表uid数组
	sharedUIds := make(map[int64][]int64)
	for _, uId := range userIds {
		share := uId % d.shards
		if sharedUIds[share] == nil {
			sharedUIds[share] = make([]int64, 0)
		}
		sharedUIds[share] = append(sharedUIds[share], uId)
	}

	tx := d.db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	now := time.Now().UnixMilli()
	for k, v := range sharedUIds {
		sql := fmt.Sprintf("update %s set type = ?, update_time = ? where session_id = ? and user_id in ?  ", d.GenUserSessionTableName(k))
		err = tx.Exec(sql, sessionType, now, sessionId, v).Error
		if err != nil {
			return err
		}
	}
	return
}

func (d defaultUserSessionModel) UpdateUserSession(userIds []int64, sessionId int64, sessionName, sessionRemark, mute, extData, noteName *string, top *int64, status, role *int, parentId *int64) (err error) {
	if sessionName == nil && sessionRemark == nil && top == nil && status == nil && mute == nil && role == nil {
		return
	}
	// 分表uid数组
	sharedUIds := make(map[int64][]int64)
	for _, uId := range userIds {
		share := uId % d.shards
		if sharedUIds[share] == nil {
			sharedUIds[share] = make([]int64, 0)
		}
		sharedUIds[share] = append(sharedUIds[share], uId)
	}

	tx := d.db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	for k, v := range sharedUIds {
		sqlBuffer := bytes.Buffer{}
		sqlBuffer.WriteString(fmt.Sprintf("update %s set ", d.GenUserSessionTableName(k)))
		if sessionName != nil {
			sqlBuffer.WriteString(fmt.Sprintf("name = '%s', ", *sessionName))
		}
		if sessionRemark != nil {
			sqlBuffer.WriteString(fmt.Sprintf("remark = '%s', ", *sessionRemark))
		}
		if top != nil {
			sqlBuffer.WriteString(fmt.Sprintf("top = %d, ", *top))
		}
		if status != nil {
			sqlBuffer.WriteString(fmt.Sprintf("status = %d, ", *status))
		}
		if mute != nil {
			sqlBuffer.WriteString(fmt.Sprintf("mute = %s, ", *mute))
		}
		if extData != nil {
			sqlBuffer.WriteString(fmt.Sprintf("ext_data = %s, ", *extData))
		}
		if noteName != nil {
			sqlBuffer.WriteString(fmt.Sprintf("note_name = %s, ", *noteName))
		}
		if role != nil {
			sqlBuffer.WriteString(fmt.Sprintf("role = %d, ", *role))
		}
		if parentId != nil {
			sqlBuffer.WriteString(fmt.Sprintf("parent_id = %d, ", *parentId))
		}
		sqlBuffer.WriteString(fmt.Sprintf("update_time = %d ", time.Now().UnixMilli()))
		sqlBuffer.WriteString("where session_id = ? and user_id in ? ")
		err = tx.Exec(sqlBuffer.String(), sessionId, v).Error
		if err != nil {
			return err
		}
	}
	return
}

func (d defaultUserSessionModel) FindEntityIdsInUserSession(userId, sessionId int64) []int64 {
	entityIds := make([]int64, 0)
	sqlStr := fmt.Sprintf("select entity_id from %s where user_id = ? and session_id = ? and deleted = 0", d.GenUserSessionTableName(userId))
	_ = d.db.Raw(sqlStr, userId, sessionId).Scan(&entityIds).Error
	return entityIds
}

func (d defaultUserSessionModel) QueryLatestUserSessions(userId, mTime int64, offset, count int, types []int) ([]*UserSession, error) {
	var (
		err          error
		userSessions = make([]*UserSession, 0)
	)
	if len(types) > 0 {
		sqlStr := "select * from " + d.GenUserSessionTableName(userId) + " where user_id = ? and type in ? and update_time > ? order by update_time asc limit ? offset ?"
		err = d.db.Raw(sqlStr, userId, types, mTime, count, offset).Scan(&userSessions).Error
	} else {
		sqlStr := "select * from " + d.GenUserSessionTableName(userId) + " where user_id = ? and update_time > ? order by update_time asc limit ? offset ?"
		err = d.db.Raw(sqlStr, userId, mTime, count, offset).Scan(&userSessions).Error
	}
	if err != nil {
		return nil, err
	}
	return userSessions, nil
}

func (d defaultUserSessionModel) GetUserSession(userId, sessionId int64) (*UserSession, error) {
	userSession := &UserSession{}
	sqlStr := "select * from " + d.GenUserSessionTableName(userId) + " where user_id = ? and session_id = ?"
	err := d.db.Raw(sqlStr, userId, sessionId).Scan(userSession).Error
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (d defaultUserSessionModel) GenUserSessionTableName(userId int64) string {
	return fmt.Sprintf("user_session_%d", userId%(d.shards))
}

func NewUserSessionModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserSessionModel {
	return defaultUserSessionModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
