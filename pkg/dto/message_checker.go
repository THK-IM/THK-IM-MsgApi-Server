package dto

type CheckMessageReq struct {
	SessionType    int    `json:"session_type"`
	SessionId      int64  `json:"session_id"`
	FunctionFlag   int64  `json:"function_flag"`
	FromUId        int64  `json:"from_u_id"`
	EntityId       int64  `json:"entity_id"`
	MessageType    int    `json:"message_type"`
	MessageContent string `json:"message_content"`
}
