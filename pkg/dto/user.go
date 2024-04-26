package dto

type PostUserOnlineReq struct {
	NodeId      int64  `json:"node_id" binding:"required"`
	ConnId      int64  `json:"conn_id" binding:"required"`
	Online      bool   `json:"online"`
	LoginTime   int64  `json:"login_time"`
	IsLogin     bool   `json:"is_login"`
	UId         int64  `json:"u_id" binding:"required"`
	Platform    string `json:"platform" binding:"required"`
	TimestampMs int64  `json:"timestamp_ms" binding:"required"`
}

type QueryUsersOnlineStatusReq struct {
	UIds []int64 `json:"u_ids" form:"u_ids"`
}

type UserOnlineStatus struct {
	UId         int64  `json:"u_id"`
	ConnId      int64  `json:"conn_id"`
	Platform    string `json:"platform"`
	NodeId      int64  `json:"node_id"`
	TimestampMs int64  `json:"timestamp_ms"`
}

type QueryUsersOnlineStatusRes struct {
	UsersOnlineStatus []*UserOnlineStatus `json:"data"`
}
