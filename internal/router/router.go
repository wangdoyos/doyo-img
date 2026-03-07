package router

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wangdoyos/doyo-img/internal/config"
	"github.com/wangdoyos/doyo-img/internal/handler"
	"github.com/wangdoyos/doyo-img/internal/middleware"
	"github.com/wangdoyos/doyo-img/internal/service"
)

// Setup 注册所有路由：API 接口 + 图片直链 + 前端 SPA fallback
func Setup(r *gin.Engine, svc *service.ImageService, cfg *config.Config, frontendFS fs.FS) {
	// CORS 中间件全局生效
	r.Use(middleware.CORS(&cfg.CORS))

	// 初始化处理器
	uploadHandler := handler.NewUploadHandler(svc, cfg)
	imageHandler := handler.NewImageHandler(svc, cfg)
	configHandler := handler.NewConfigHandler(cfg)

	// API 路由组 —— 限流仅作用于 API，不影响静态资源和图片直链
	api := r.Group("/api")
	if cfg.RateLimit.Enabled {
		rl := middleware.NewRateLimiter(&cfg.RateLimit)
		api.Use(rl.Middleware())
	}
	{
		api.POST("/upload", uploadHandler.Upload)
		api.GET("/image/:id", imageHandler.GetImageInfo)
		api.DELETE("/image/:id", imageHandler.DeleteImage)
		api.GET("/recent", imageHandler.ListRecent)
		api.GET("/config", configHandler.GetPublicConfig)
	}

	// 图片直链路由
	r.GET("/i/*id", imageHandler.ServeImage)

	// 前端 SPA —— 从嵌入的文件系统提供静态资源
	if frontendFS != nil {
		fileServer := http.FileServer(http.FS(frontendFS))
		r.NoRoute(func(c *gin.Context) {
			// 先尝试匹配静态文件
			path := c.Request.URL.Path
			f, err := frontendFS.Open(path[1:]) // 去掉前导斜杠
			if err == nil {
				f.Close()
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
			// 匹配不到则回退到 index.html（SPA 路由支持）
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
		})
	}
}
