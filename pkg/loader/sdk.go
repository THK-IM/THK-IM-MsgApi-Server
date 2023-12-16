package loader

import (
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
)

func LoadSdks(sdkConfigs []conf.Sdk, logger *logrus.Entry) map[string]interface{} {
	sdkMap := make(map[string]interface{})
	//for _, c := range sdkConfigs {
	//	if c.Name == "user-api" {
	//		userApi := rpc.NewUserApi(c, logger)
	//		sdkMap[c.Name] = userApi
	//	} else if c.Name == "msg-api" {
	//		userApi := rpc.NewMsgApi(c, logger)
	//		sdkMap[c.Name] = userApi
	//	}
	//}
	return sdkMap
}
