package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func RunServer() {
	server := gin.Default()

	server.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, POST", // Only GET and POST methods are allowed for this API.
		RequestHeaders: "Origin,Authorization,Content-Type,Access-Control-Allow-Origin",
		MaxAge:         50 * time.Second,
	}))

	registerAppRoutes(server)

	server.Run(
		fmt.Sprintf(":8080"),
	)
}
