package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/rpc"
	"github.com/thk-im/thk-im-msg-api-server/pkg/app"
	"net"
	"strings"
)

const (
	tokenKey = "Token"
	uidKey   = "Uid"
)

func userTokenAuth(appCtx *app.Context) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get(tokenKey)
		if token == "" {
			appCtx.Logger().Warn("token nil error")
			dto.ResponseUnauthorized(context)
			context.Abort()
			return
		}
		req := rpc.GetUserIdByTokenReq{Token: token}
		res, err := appCtx.RpcUserApi().GetUserIdByToken(req)
		if err != nil {
			appCtx.Logger().Warn("token error")
			dto.ResponseUnauthorized(context)
			context.Abort()
		} else {
			context.Set(uidKey, res.UserId)
			context.Next()
		}
	}
}

func whiteIpAuth(appCtx *app.Context) gin.HandlerFunc {
	ipWhiteList := appCtx.Config().IpWhiteList
	ips := strings.Split(ipWhiteList, ",")
	return func(context *gin.Context) {
		ip := context.ClientIP()
		appCtx.Logger().Infof("RemoteAddr: %s", ip)
		if isIpValid(ip, ips) {
			dto.ResponseForbidden(context)
			context.Abort()
		} else {
			context.Next()
		}
	}
}

func isIpValid(clientIp string, whiteIpList []string) bool {
	ip := net.ParseIP(clientIp)
	for _, whiteIp := range whiteIpList {
		_, ipNet, err := net.ParseCIDR(whiteIp)
		if err != nil {
			return false
		}
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}
