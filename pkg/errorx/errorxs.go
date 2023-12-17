package errorx

import "github.com/thk-im/thk-im-base-server/errorx"

var (
	ErrSessionInvalid        = errorx.NewErrorX(4004001, "Invalid session")
	ErrSessionAlreadyDeleted = errorx.NewErrorX(4004002, "group has been deleted")
	ErrSessionType           = errorx.NewErrorX(4004003, "Session type error")
	ErrSessionMessageInvalid = errorx.NewErrorX(4004004, "Invalid session message")
	ErrMessageTypeNotSupport = errorx.NewErrorX(4004005, "Message type not support")
	ErrSessionMuted          = errorx.NewErrorX(4004101, "Session muted")
	ErrUserMuted             = errorx.NewErrorX(4004102, "User muted")
	ErrUserReject            = errorx.NewErrorX(4004103, "user reject your message")
	ErrMessageDeliveryFailed = errorx.NewErrorX(5004001, "Message delivery failed")
)
