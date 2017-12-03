package serve

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	logger "github.com/Sirupsen/logrus"
)

func ServeHttp(addr string) {
	router := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	router.POST("/peers/:name", search)
	router.GET("/peers/:name", peersGET)
	logger.Infof("http server listening at %s", addr)
	router.Run(addr)
}

func search(c *gin.Context) {
	name := c.Param("name")
	port, _ := strconv.Atoi(c.DefaultQuery("port", "0"))

	decService := service.GetService(name)
	if decService != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Service already exists."})
		return
	}

	service.AddService(name, uint32(port))
	defer service.StopService(name)

	c.Stream(func(w io.Writer) bool {
		//c.SSEvent("message", <-listener)
		return true
	})
}

func peersGET(c *gin.Context) {
	name := c.Param("name")
	decService := service.GetService(name)
	if decService == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service doesn't exist."})
		return
	}
	c.JSON(200, decService.GetPeers())
}