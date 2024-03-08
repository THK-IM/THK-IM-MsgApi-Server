package logic

const (
	sessionCreateLockKey     = "%s:se:c:%d:%d"
	sessionUpdateLockKey     = "%s:se:m:%d"
	userSessionUpdateLockKey = "%s:u:se:m:%d:%d"

	userOnlineKey = "%s:olu:%s:%d"

	PlatformAndroid = "Android"
	PlatformIOS     = "IOS"
	PlatformWeb     = "Web"
)
