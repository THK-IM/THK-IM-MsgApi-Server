package dto

const (
	FuncTextFlag    = 1  // 文本消息包括emoji表情
	FuncAudioFlag   = 2  // 录音消息
	FuncEmojiFlag   = 4  // 自定义表情
	FuncImageFlag   = 8  // 图片消息
	FuncVideoFlag   = 16 // 视频消息
	FuncForwardFlag = 32 // 转发
	FuncReadFlag    = 64 // 已读
	FuncAll         = FuncTextFlag | FuncAudioFlag | FuncEmojiFlag | FuncImageFlag | FuncVideoFlag | FuncForwardFlag | FuncReadFlag
)

type CreateSessionReq struct {
	UId            int64    `json:"u_id" binding:"required"`
	UserNoteName   string   `json:"user_note_name"`
	UserNoteAvatar string   `json:"user_note_avatar"`
	Type           int      `json:"type" binding:"required"`
	EntityId       int64    `json:"entity_id"`      // 单聊type为对方用户id,group或supergroup为群id
	Members        []int64  `json:"members"`        // 单聊时不用提交该字段
	MemberAvatars  []string `json:"member_avatars"` // 单聊时不用提交该字段
	MemberNames    []string `json:"member_names"`   // 单聊时不用提交该字段
	ExtData        *string  `json:"ext_data"`       // 业务方扩展字段
	Name           string   `json:"name"`           // session名
	Remark         string   `json:"remark"`
	FunctionFlag   int64    `json:"function_flag"`
}

type CreateSessionRes struct {
	SId          int64   `json:"s_id"`
	ParentId     int64   `json:"parent_id"`
	EntityId     int64   `json:"entity_id"`
	Type         int     `json:"type"`
	Name         string  `json:"name"`
	Remark       string  `json:"remark"`
	FunctionFlag int64   `json:"function_flag"`
	ExtData      *string `json:"ext_data"`
	NoteName     string  `json:"note_name"`
	NoteAvatar   string  `json:"note_avatar"`
	Mute         int     `json:"mute"`
	Role         int     `json:"role"`
	CTime        int64   `json:"c_time"`
	MTime        int64   `json:"m_time"`
	Top          int64   `json:"top"`
	Status       int     `json:"status"`
	IsNew        bool    `json:"is_new"` // 如果之前已经创建，false
}

type UpdateSessionReq struct {
	Id           int64   `json:"id"`
	Mute         *int    `json:"mute"`
	Name         *string `json:"name"`
	Remark       *string `json:"remark"`
	FunctionFlag *int64  `json:"function_flag"`
	ExtData      *string `json:"ext_data"`
}

type UpdateSessionTypeReq struct {
	Id   int64 `json:"id"`
	Type int   `json:"type"`
}

type DelSessionReq struct {
	Id int64 `json:"id"`
}

type UpdateUserSessionReq struct {
	UId        int64   `json:"u_id"`
	SId        int64   `json:"s_id"`
	NoteName   *string `json:"note_name"`
	NoteAvatar *string `json:"note_avatar"`
	Top        *int64  `json:"top"`
	Status     *int    `json:"status"`
	ParentId   *int64  `json:"parent_id"`
}

type SearchUserSessionReq struct {
	UId      int64   `json:"u_id" form:"u_id"`
	Offset   int     `json:"offset" form:"offset"`
	Count    int     `json:"count" form:"count"`
	Keywords *string `json:"m_time" form:"m_time"`
	Types    []int   `json:"types" form:"types"`
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
	SId          int64   `json:"s_id"`
	Name         string  `json:"name"`
	Remark       string  `json:"remark"`
	FunctionFlag int64   `json:"function_flag"`
	Type         int     `json:"type"`
	Status       int     `json:"status"`
	Role         int     `json:"role"`
	Mute         int     `json:"mute"`
	Top          int64   `json:"top"`
	NoteName     string  `json:"note_name"`
	NoteAvatar   string  `json:"note_avatar"`
	Deleted      int8    `json:"deleted"`
	EntityId     int64   `json:"entity_id"`
	ExtData      *string `json:"ext_data,omitempty"`
	CTime        int64   `json:"c_time"`
	MTime        int64   `json:"m_time"`
}

type SessionUser struct {
	SId        int64  `json:"s_id"`
	UId        int64  `json:"u_id"`
	Type       int    `json:"type"`
	Mute       int    `json:"mute"`
	Role       int    `json:"role"`
	Status     int    `json:"status"`
	NoteAvatar string `json:"note_avatar"`
	NoteName   string `json:"note_name"`
	Deleted    int8   `json:"deleted"`
	CTime      int64  `json:"c_time"`
	MTime      int64  `json:"m_time"`
}

type SearchUserSessionRes struct {
	Data  []*UserSession `json:"data"`
	Total int            `json:"total"`
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
	EntityId    int64    `json:"entity_id" binding:"required"`
	UIds        []int64  `json:"u_ids" binding:"required"`
	NoteNames   []string `json:"note_names"`
	NoteAvatars []string `json:"note_avatars"`
	Role        int      `json:"role" binding:"required"`
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

type SessionUserCountRes struct {
	Count int `json:"count"`
}
