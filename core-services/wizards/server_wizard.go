package wizards

import "github.com/gin-gonic/gin"

func RegisterServer(router *gin.Engine) {
	// Setup proxy for the post-service
	// The NewReverseProxy function is defined in wizards/proxy.go
	postProxy := NewReverseProxy(Config.Server.PostServiceURL)
	interactionProxy := NewReverseProxy(Config.Server.InteractionServiceURL)

	apiGroup := router.Group("/api")
	{
		// Auth routes handled by core-service
		apiGroup.POST("/register", UserHttp.Register)
		apiGroup.POST("/login", UserHttp.Login)

		// All routes below this point require authentication
		apiGroup.Use(AuthMiddleware.UseAuthMiddleware())
		// User and Follow routes handled by core-service
		userRoutes := apiGroup.Group("/users")
		{
			userRoutes.GET("/me", UserHttp.GetMe)
			userRoutes.GET("/:username", UserHttp.GetProfile)
			userRoutes.PUT("", UserHttp.UpdateUser) // PUT /api/users (updates the authenticated user)
			userRoutes.DELETE("", UserHttp.DeleteUser) // DELETE /api/users (deletes the authenticated user)
			userRoutes.POST("/follow/:username", FollowHttp.FollowUser)
			userRoutes.POST("/unfollow/:username", FollowHttp.UnfollowUser)
		}
		// --- PROXY ROUTES ---
		// The authentication middleware is already applied to /api.
		// We add another middleware specifically for proxied routes to add required headers.

		// Any request to /api/posts/* is proxied to the post-service.
		postRoutes := apiGroup.Group("/posts")
		{
			// 1. The ProxyHeadersMiddleware runs first, adding the X-User-ID header.
			// 2. The request is then passed to the proxy handler.
			postRoutes.Use(ProxyMiddleware)
			postRoutes.Any("/*proxyPath", func(c *gin.Context) {
				postProxy.ServeHTTP(c.Writer, c.Request)
			})
		}

		interactionRoutes := apiGroup.Group("/interactions")
		{
			interactionRoutes.Use(ProxyMiddleware)
			interactionRoutes.Any("/*proxyPath", func(c *gin.Context) {
				interactionProxy.ServeHTTP(c.Writer, c.Request)
			})
		}
	}
}
