package loader

import (
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-msgapi-server/pkg/sdk"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func LoadSdks(sdkConfigs []conf.Sdk, logger *logrus.Entry) map[string]interface{} {
	sdkMap := make(map[string]interface{})
	for _, sdkConfig := range sdkConfigs {
		if sdkConfig.Name == "login_api" {
			loginApi := userSdk.NewLoginApi(sdkConfig, logger)
			sdkMap[sdkConfig.Name] = loginApi
		} else if sdkConfig.Name == "user_api" {
			userApi := userSdk.NewUserApi(sdkConfig, logger)
			sdkMap[sdkConfig.Name] = userApi
		} else if sdkConfig.Name == "msg_check_api" {
			checkApi := sdk.NewMsgCheckerApi(sdkConfig, logger)
			sdkMap[sdkConfig.Name] = checkApi
		}
	}
	return sdkMap
}
