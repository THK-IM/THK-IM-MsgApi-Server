package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-msgapi-server/pkg/app"
	userSdk "github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func RegisterMsgApiHandlers(appCtx *app.Context) {
	httpEngine := appCtx.HttpEngine()
	ipAuth := middleware.WhiteIpAuth(appCtx.Context)
	userApi := appCtx.UserApi()
	userTokenAuth := userSdk.UserTokenAuth(userApi, appCtx.Logger())

	var authMiddleware gin.HandlerFunc
	if appCtx.Config().DeployMode == conf.DeployExposed {
		authMiddleware = userTokenAuth
	} else if appCtx.Config().DeployMode == conf.DeployBackend {
		authMiddleware = ipAuth
	} else {
		panic(errors.New("check deployMode conf"))
	}

	sessionRoute := httpEngine.Group("/session")
	sessionRoute.Use(authMiddleware)
	{
		sessionRoute.POST("", createSession(appCtx))                      // 创建/获取session
		sessionRoute.PUT("/:id", updateSession(appCtx))                   // 修改session相关信息
		sessionRoute.GET("/:id/user", getSessionUser(appCtx))             // 会话成员查询
		sessionRoute.POST("/:id/user", addSessionUser(appCtx))            // 会话增员
		sessionRoute.DELETE("/:id/user", deleteSessionUser(appCtx))       // 会话减员
		sessionRoute.PUT("/:id/user", updateSessionUser(appCtx))          // 会话成员修改
		sessionRoute.GET("/:id/message", getSessionMessages(appCtx))      // 获取session下的消息列表
		sessionRoute.DELETE("/:id/message", deleteSessionMessage(appCtx)) // 删除session下的消息列表

		// 如果提供内置对象存储服务，则开放接口
		if appCtx.ObjectStorage() != nil {
			sessionRoute.GET("/object/upload_params", getObjectUploadParams(appCtx)) // 获取对象上传参数
			sessionRoute.GET("/object/download_url", getObjectDownloadUrl(appCtx))   // 获取对象,鉴权后重定向到签名后的minio地址
		}
	}

	userSessionRoute := httpEngine.Group("/user_session")
	userSessionRoute.Use(authMiddleware)
	{
		userSessionRoute.GET("/latest", getUserSessions(appCtx))         // 用户获取自己最近的session列表
		userSessionRoute.GET("/:uid/:sid", getUserSession(appCtx))       // 用户获取自己的session
		userSessionRoute.PUT("", updateUserSession(appCtx))              // 用户修改自己的session
		userSessionRoute.DELETE("/:uid/:sid", deleteUserSession(appCtx)) // 用户删除自己的session
	}

	messageRoute := httpEngine.Group("/message")
	messageRoute.Use(authMiddleware)
	{
		messageRoute.GET("/latest", getUserLatestMessages(appCtx)) // 获取最近消息
		messageRoute.POST("", sendMessage(appCtx))                 // 发送消息
		messageRoute.DELETE("", deleteUserMessage(appCtx))         // 删除消息
		messageRoute.POST("/ack", ackUserMessages(appCtx))         // 用户消息设置ack(已接收) 不支持超级群
		messageRoute.POST("/read", readUserMessage(appCtx))        // 用户消息设置已读 不支持超级群
		messageRoute.POST("/revoke", revokeUserMessage(appCtx))    // 用户消息撤回
		messageRoute.POST("/reedit", reeditUserMessage(appCtx))    // 更新用户消息
		messageRoute.POST("/forward", forwardUserMessage(appCtx))  // 转发用户消息
	}

	systemRoute := httpEngine.Group("/system")
	systemRoute.Use(ipAuth)
	{
		systemRoute.POST("/user/online", updateUserOnlineStatus(appCtx)) // 更新用户在线状态
		systemRoute.GET("/user/online", getUsersOnlineStatus(appCtx))    // 获取用户上线状态
		systemRoute.POST("/user/kickoff", kickOffUser(appCtx))           // 踢下线用户
		systemRoute.POST("/message/send", sendSystemMessage(appCtx))     // 发送会话中的系统消息
		systemRoute.POST("/message/push", pushExtendedMessage(appCtx))   // 推送消息(用户消息/好友消息/群组消息/自定义消息)
	}
}