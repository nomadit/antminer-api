package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nomadit/antminer-api/api/ctrls"
	"time"
)

func RouteHandler() *gin.Engine {
	r := gin.Default()
	r.Use(getCors())

	proxyCtrl := ctrls.NewProxyCtrl()
	pcCtrl := ctrls.NewPcCtrl()
	commandCtrl := ctrls.NewCommandCtrl()

	beatApi := r.Group("/api/beat")
	beatApi.GET("/pc/all/by_net_id/:proxyID", pcCtrl.PcListByNetIDIncludeDeleted)
	beatApi.PUT("/pc/by_net_id/:proxyID", pcCtrl.UpsertPC)
	beatApi.PUT("/pc/invalid/by_mac/:mac", pcCtrl.UpdateInvalidByMac)
	beatApi.PUT("/pc/status/:id", pcCtrl.UpdatePcStatus)
	beatApi.PUT("/pc/name/:id", pcCtrl.UpdatePcName)
	beatApi.PUT("/pc/config/frequency/:id", pcCtrl.UpsertPcFrequencyConfig)
	beatApi.GET("/proxy/by_mac/:mac", proxyCtrl.ProxyByMac)
	beatApi.PUT("/proxy", proxyCtrl.UpdateProxy)
	beatApi.POST("/proxy", proxyCtrl.CreateProxy)
	beatApi.GET("/command/list/by_pc_ids", commandCtrl.CommandListByPcIDs)
	beatApi.PUT("/command/status/:id", commandCtrl.UpdateCommandStatus)

	return r
}

func getCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost",
		},
		AllowMethods:  []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:  []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	})
}

// remove the auth
//func authUserID() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if _, exists := c.Get("JWT_PAYLOAD"); !exists {
//			log.Println("JWT_PAYLOAD not exist")
//			return
//		}
//		jwtClaims, _ := c.Get("JWT_PAYLOAD")
//		c.Set("authUserID", jwtClaims.(jwt.MapClaims)["id"])
//	}
//}
