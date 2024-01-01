package dto

type CreateSessionReq struct {
	UId      int64   `json:"u_id" binding:"required"`
	Type     int     `json:"type" binding:"required"`
	EntityId int64   `json:"entity_id"` // 单聊type为对方用户id,group或supergroup为群id
	Members  []int64 `json:"members"`   // 单聊时不用提交该字段
	ExtData  *string `json:"ext_data"`  // 业务方扩展字段
	Name     string  `json:"name"`      // session名
	Remark   string  `json:"remark"`
}

type CreateSessionRes struct {
	SId      int64   `json:"s_id"`
	ParentId int64   `json:"parent_id"`
	EntityId int64   `json:"entity_id"`
	Type     int     `json:"type"`
	Name     string  `json:"name"`
	Remark   string  `json:"remark"`
	ExtData  *string `json:"ext_data"`
	Mute     int     `json:"mute"`
	Role     int     `json:"role"`
	CTime    int64   `json:"c_time"`
	MTime    int64   `json:"m_time"`
	Top      int64   `json:"top"`
	Status   int     `json:"status"`
	IsNew    bool    `json:"is_new"` // 如果之前已经创建，false
}

type UpdateSessionReq struct {
	Id      int64   `json:"id"`
	Mute    *int    `json:"mute"`
	Name    *string `json:"name"`
	Remark  *string `json:"remark"`
	ExtData *string `json:"ext_data"`
}

type UpdateSessionTypeReq struct {
	Id   int64 `json:"id"`
	Type int   `json:"type"`
}

type DelSessionReq struct {
	Id int64 `json:"id"`
}

type UpdateUserSessionReq struct {
	UId      int64   `json:"u_id"`
	SId      int64   `json:"s_id"`
	NoteName *string `json:"note_name"`
	Top      *int64  `json:"top"`
	Status   *int    `json:"status"`
	ParentId *int64  `json:"parent_id"`
}

type QueryLatestUserSessionReq struct {
	UId    int64 `json:"u_id" form:"u_id"`
	Offset int   `json:"offset" form:"offset"`
	Count  int   `json:"count" form:"count"`
	MTime  int64 `json:"m_time" form:"m_time"`
	Types  []int `json:"types" form:"types"`
}

type QueryUserSessionReq struct {
	UId      int64 `json:"u_id" form:"u_id"`
	EntityId int64 `json:"entity_id" form:"entity_id"`
	Type     int   `json:"type" form:"type"`
}

type UserSession struct {
	SId      int64   `json:"s_id"`
	Name     string  `json:"name"`
	Remark   string  `json:"remark"`
	Type     int     `json:"type"`
	Status   int     `json:"status"`
	Role     int     `json:"role"`
	Mute     int     `json:"mute"`
	Top      int64   `json:"top"`
	NoteName string  `json:"note_name"`
	Deleted  int8    `json:"deleted"`
	EntityId int64   `json:"entity_id"`
	ExtData  *string `json:"ext_data,omitempty"`
	CTime    int64   `json:"c_time"`
	MTime    int64   `json:"m_time"`
}

type SessionUser struct {
	SId      int64  `json:"s_id"`
	UId      int64  `json:"u_id"`
	Type     int    `json:"type"`
	Mute     int    `json:"mute"`
	Role     int    `json:"role"`
	Status   int    `json:"status"`
	NoteName string `json:"note_name"`
	Deleted  int8   `json:"deleted"`
	CTime    int64  `json:"c_time"`
	MTime    int64  `json:"m_time"`
}

type QueryLatestUserSessionsRes struct {
	Data []*UserSession `json:"data"`
}

type QuerySessionUsersReq struct {
	SId   int64 `json:"s_id" form:"s_id"`
	Role  *int  `json:"role" form:"role"`
	MTime int64 `json:"m_time" form:"m_time"`
	Count int   `json:"count" form:"count"`
}

type QuerySessionUsersRes struct {
	Data []*SessionUser `json:"data"`
}

type SessionAddUserReq struct {
	EntityId int64   `json:"entity_id" binding:"required"`
	UIds     []int64 `json:"u_ids" binding:"required"`
	Role     int     `json:"role" binding:"required"`
}

type SessionDelUserReq struct {
	UIds []int64 `json:"u_ids" binding:"required"`
}

type SessionUserUpdateReq struct {
	SId  int64   `json:"s_id" binding:"required"`
	UIds []int64 `json:"u_ids" binding:"required"`
	Role *int    `json:"role"`
	Mute *int    `json:"mute"`
}
