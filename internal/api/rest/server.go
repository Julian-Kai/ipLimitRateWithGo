package rest

import (
	"net/http"

	"ip.limit.rate/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var ipRateLimiter = util.NewIPRateLimiter()

type Router struct {
	router *gin.Engine
}

func NewRouter() *Router {
	router := gin.New()

	router.Use(LimitMiddleware)
	router.GET("/ping", Ping)

	return &Router{
		router: router,
	}
}

func (r *Router) GetRouter() *gin.Engine {
	return r.router
}

func (r *Router) Run() error {
	return r.router.Run()
}

func LimitMiddleware(ctx *gin.Context) {
	bkt := ipRateLimiter.GetLimiter(ctx.Request.RemoteAddr)
	if count := bkt.TakeAvailable(1); count == 0 {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"message": "Error",
		})
		ctx.Abort()
		return
	}

	ctx.Next()
	return
}

func Ping(ctx *gin.Context) {
	key := ctx.Request.RemoteAddr
	bkt := ipRateLimiter.GetLimiter(key)

	currentAvailableToken := bkt.AvailableTokens()
	currentUsedToken := bkt.Capacity() - currentAvailableToken

	logrus.Infof("ip:%s, currentAvailableToken:%d, currentUsedToken:%d", key, currentAvailableToken, currentUsedToken)
	ctx.JSON(200, gin.H{
		"message": currentUsedToken,
	})
}
