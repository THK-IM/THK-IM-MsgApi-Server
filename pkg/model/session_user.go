package model

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/snowflake"
	"gorm.io/gorm"
	"time"
)

const (
	SessionMember     = 1 // 普通成员 可以查询session信息，session会话历史消息
	SessionAdmin      = 2 // 管理员 模式下可以修改session基本信息，禁言单个用户, 添加/删除普通成员
	SessionSuperAdmin = 3 // 超级管理员 可以全员禁言, 添加/删除管理员
	SessionOwner      = 4 // 拥有者 可以添加超级管理员, 删除管理员，删除session
)

var ErrMemberCnt = errors.New("member count error")

type (
	SessionUser struct {
		Id         int64  `gorm:"id" json:"id"`
		SessionId  int64  `gorm:"session_id" json:"session_id"`
		UserId     int64  `gorm:"user_id" json:"user_id"`
		Type       int    `gorm:"type" json:"type"`
		Role       int    `gorm:"role" json:"role"`
		Mute       int    `gorm:"mute" json:"mute"`
		Status     int    `gorm:"status" json:"status"`
		NoteName   string `gorm:"note_name" json:"note_name"`
		NoteAvatar string `gorm:"note_name" json:"note_avatar"`
		CreateTime int64  `gorm:"create_time" json:"create_time"`
		UpdateTime int64  `gorm:"update_time" json:"update_time"`
		Deleted    int8   `gorm:"deleted" json:"deleted"`
	}

	SessionUserModel interface {
		FindSessionUsersByMTime(sessionId, mTime int64, role *int, count int) ([]*SessionUser, error)
		FindAllSessionUsers(sessionId int64) ([]*SessionUser, error)
		FindSessionUsers(sessionId int64, userIds []int64) ([]*SessionUser, error)
		FindSessionUser(sessionId, userId int64) (*SessionUser, error)
		FindSessionUserCount(sessionId int64) (int, error)
		FindUIdsInSessionWithoutStatus(sessionId int64, status int, uIds []int64) []*SessionUser
		FindUIdsInSessionContainStatus(sessionId int64, status int, uIds []int64) []*SessionUser
		AddUser(session *Session, entityIds []int64, userIds []int64, role []int, noteNames, noteAvatars []string, maxCount int) ([]*UserSession, error)
		DelUser(session *Session, userIds []int64) (err error)
		UpdateType(sessionId int64, sessionType int) (err error)
		UpdateUser(sessionId int64, userIds []int64, role, status *int, noteName, noteAvatar, mute *string) (err error)
		DelSession(sessionId int64) error
	}

	defaultSessionUserModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultSessionUserModel) FindSessionUsersByMTime(sessionId, mTime int64, role *int, count int) ([]*SessionUser, error) {
	sessionUser := make([]*SessionUser, 0)
	tableName := d.genSessionUserTableName(sessionId)
	var err error
	if role == nil {
		sqlStr := fmt.Sprintf("select * from %s where session_id = ? and update_time >= ? order by update_time asc limit 0, ?", tableName)
		err = d.db.Raw(sqlStr, sessionId, mTime, count).Scan(&sessionUser).Error
	} else {
		sqlStr := fmt.Sprintf("select * from %s where session_id = ? and role >= ? and update_time >= ? order by update_time asc limit 0, ?", tableName)
		err = d.db.Raw(sqlStr, sessionId, *role, mTime, count).Scan(&sessionUser).Error
	}
	return sessionUser, err
}

func (d defaultSessionUserModel) FindAllSessionUsers(sessionId int64) ([]*SessionUser, error) {
	sessionUser := make([]*SessionUser, 0)
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select * from %s where session_id = ?", tableName)
	err := d.db.Raw(sqlStr, sessionId).Scan(&sessionUser).Error
	return sessionUser, err
}

func (d defaultSessionUserModel) FindSessionUsers(sessionId int64, userIds []int64) ([]*SessionUser, error) {
	sessionUser := make([]*SessionUser, 0)
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select * from %s where session_id = ? and user_id in ? and deleted = 0", tableName)
	err := d.db.Raw(sqlStr, sessionId, userIds).Scan(&sessionUser).Error
	return sessionUser, err
}

func (d defaultSessionUserModel) FindSessionUser(sessionId, userId int64) (*SessionUser, error) {
	sessionUser := &SessionUser{}
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select * from %s where session_id = ?  and user_id = ? and deleted = 0", tableName)
	err := d.db.Raw(sqlStr, sessionId, userId).Scan(sessionUser).Error
	return sessionUser, err
}

func (d defaultSessionUserModel) FindSessionUserCount(sessionId int64) (int, error) {
	count := 0
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select count(0) from %s where session_id = ?  and deleted = 0", tableName)
	err := d.db.Raw(sqlStr, sessionId).Scan(&count).Error
	return count, err
}

func (d defaultSessionUserModel) FindUIdsInSessionWithoutStatus(sessionId int64, status int, userIds []int64) []*SessionUser {
	sessionUsers := make([]*SessionUser, 0)
	uIdsCondition := ""
	if len(userIds) > 0 {
		uIdsCondition = " and user_id in ? "
	}
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ? %s and status & ? = 0 and deleted = 0",
		d.genSessionUserTableName(sessionId), uIdsCondition)
	if len(userIds) > 0 {
		tx := d.db.Raw(sqlStr, sessionId, userIds, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	} else {
		tx := d.db.Raw(sqlStr, sessionId, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	}
	return sessionUsers
}

func (d defaultSessionUserModel) FindUIdsInSessionContainStatus(sessionId int64, status int, userIds []int64) []*SessionUser {
	sessionUsers := make([]*SessionUser, 0)
	uIdsCondition := ""
	if len(userIds) > 0 {
		uIdsCondition = " and user_id in ? "
	}
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ? %s and status & ? > 0 and deleted = 0",
		d.genSessionUserTableName(sessionId), uIdsCondition)
	if len(userIds) > 0 {
		tx := d.db.Raw(sqlStr, sessionId, userIds, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	} else {
		tx := d.db.Raw(sqlStr, sessionId, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	}
	return sessionUsers
}

func (d defaultSessionUserModel) AddUser(session *Session, entityIds []int64, userIds []int64, role []int, noteNames, noteAvatars []string, maxCount int) (userSessions []*UserSession, err error) {
	tx := d.db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	count := 0
	tableName := d.genSessionUserTableName(session.Id)
	sqlStr := fmt.Sprintf("select count(0) from %s where session_id = ? and user_id not in ? and deleted = 0", tableName)
	if err = tx.Raw(sqlStr, session.Id, userIds).Scan(&count).Error; err != nil {
		return nil, err
	}

	if count > maxCount-len(userIds) {
		return nil, ErrMemberCnt
	}

	t := time.Now().UnixMilli()
	sql1 := "insert into " + d.genSessionUserTableName(session.Id) + " " +
		"(session_id, user_id, role, note_name, note_avatar, type, create_time, update_time) " +
		"values (?, ?, ?, ?, ?, ?, ?, ?) " +
		"on duplicate key update role = ?, note_name = ?, note_avatar = ?, deleted = ?, update_time = ? "
	userMute := 0
	if session.Mute == 1 {
		userMute = 1
	}
	userSessions = make([]*UserSession, 0)
	for index, id := range userIds {
		if err = tx.Exec(sql1, session.Id, id, role[index], noteNames[index], noteAvatars[index], session.Type, t, t,
			role[index], noteNames[index], noteAvatars[index], 0, t,
		).Error; err != nil {
			return nil, err
		}

		sql2 := "insert into " + d.genUserSessionTableName(id) + " " +
			"(session_id, user_id, type, entity_id, role, name, remark, function_flag, mute, ext_data, parent_id, note_name, " +
			"note_avatar, create_time, update_time) " +
			"values (?, ?, ?, ?, ?, ?, ?, ?, " +
			"?, ?, ?, ?, ?, ?, ?) " +
			"on duplicate key update top = ?, role = ?, name = ?, remark = ?, function_flag = ?, mute = ?, " +
			"deleted = ?, ext_data = ?, parent_id = ?, note_name = ?, note_avatar = ?, update_time = ? "
		if err = tx.Exec(
			sql2, session.Id, id, session.Type, entityIds[index], role[index], session.Name, session.Remark, session.FunctionFlag,
			userMute, session.ExtData, 0, noteNames[index], noteAvatars[index], t, 0,
			0, role[index], session.Name, session.Remark, session.FunctionFlag, userMute,
			0, session.ExtData, 0, noteNames[index], noteAvatars[index], t,
		).Error; err != nil {
			return nil, err
		}
		userSession := &UserSession{
			SessionId:    session.Id,
			UserId:       userIds[index],
			Type:         session.Type,
			EntityId:     entityIds[index],
			Name:         session.Name,
			Remark:       session.Remark,
			FunctionFlag: session.FunctionFlag,
			Top:          0,
			Role:         role[index],
			Mute:         userMute,
			ExtData:      session.ExtData,
			NoteName:     noteNames[index],
			NoteAvatar:   noteAvatars[index],
			ParentId:     0,
			Status:       0,
			CreateTime:   t,
			UpdateTime:   t,
			Deleted:      0,
		}
		userSessions = append(userSessions, userSession)
	}
	return
}

func (d defaultSessionUserModel) DelUser(session *Session, userIds []int64) (err error) {
	tx := d.db.Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()
	t := time.Now().UnixMilli()
	sql1 := "update " + d.genSessionUserTableName(session.Id) +
		" set deleted = ?, update_time = ? where session_id = ? and user_id = ?"
	for _, id := range userIds {
		if err = tx.Exec(sql1, 1, t, session.Id, id).Error; err != nil {
			return err
		}

		sql2 := "update " + d.genUserSessionTableName(id) +
			" set deleted = ?, update_time = ? where session_id = ? and user_id = ?"
		if err = tx.Exec(sql2, 1, t, session.Id, id).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d defaultSessionUserModel) UpdateType(sessionId int64, sessionType int) (err error) {
	t := time.Now().UnixMilli()
	sql := fmt.Sprintf("update %s set type = ?, update_time = ? where session_id = ? ", d.genSessionUserTableName(sessionId))
	return d.db.Exec(sql, sessionType, t, sessionId).Error
}

func (d defaultSessionUserModel) UpdateUser(sessionId int64, userIds []int64, role, status *int, noteName, noteAvatar, mute *string) (err error) {
	if role == nil && status == nil && mute == nil && noteName == nil && noteAvatar == nil {
		return nil
	}
	t := time.Now().UnixMilli()
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString(fmt.Sprintf("update %s set ", d.genSessionUserTableName(sessionId)))
	if role != nil {
		sqlBuffer.WriteString(fmt.Sprintf(" role = %d, ", *role))
	}
	if status != nil {
		sqlBuffer.WriteString(fmt.Sprintf(" status = %d, ", *status))
	}
	if noteName != nil {
		sqlBuffer.WriteString(fmt.Sprintf(" note_name = '%s', ", *noteName))
	}
	if noteAvatar != nil {
		sqlBuffer.WriteString(fmt.Sprintf(" note_avatar = '%s', ", *noteAvatar))
	}
	if mute != nil {
		sqlBuffer.WriteString(fmt.Sprintf(" mute = %s, ", *mute))
	}
	sqlBuffer.WriteString(fmt.Sprintf(" update_time = %d where session_id = ? and user_id in ? ", t))
	return d.db.Exec(sqlBuffer.String(), sessionId, userIds).Error
}

func (d defaultSessionUserModel) DelSession(sessionId int64) (err error) {
	tx := d.db.Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	sessionUser := make([]*SessionUser, 0)
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select * from %s where session_id = ?", tableName)
	err = d.db.Raw(sqlStr, sessionId).Scan(&sessionUser).Error
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	for _, su := range sessionUser {
		tableName = d.genUserSessionTableName(su.UserId)
		if su.Deleted == 0 {
			sqlStr = fmt.Sprintf("update %s set deleted = 1, update_time = ? where session_id = ? and user_id = ? ", tableName)
			err = tx.Exec(sqlStr, now, su.SessionId, su.UserId).Error
			if err != nil {
				return
			}
		}
	}

	tableName = d.genUserSessionTableName(sessionId)
	sqlStr = fmt.Sprintf("update %s set deleted = 1, update_time = ? where session_id = ?", tableName)
	err = tx.Exec(sqlStr, now, sessionId).Error
	if err != nil {
		return
	}

	tableName = d.genSessionTableName(sessionId)
	sqlStr = fmt.Sprintf("update %s set deleted = 1, update_time = ? where id = ?", tableName)
	err = tx.Exec(sqlStr, now, sessionId).Error
	return
}

func (d defaultSessionUserModel) genUserSessionTableName(userId int64) string {
	return fmt.Sprintf("user_session_%d", userId%(d.shards))
}

func (d defaultSessionUserModel) genSessionUserTableName(sessionId int64) string {
	return fmt.Sprintf("session_user_%d", sessionId%(d.shards))
}

func (d defaultSessionUserModel) genSessionTableName(sessionId int64) string {
	return fmt.Sprintf("session_%d", sessionId%(d.shards))
}

func NewSessionUserModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) SessionUserModel {
	return defaultSessionUserModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
