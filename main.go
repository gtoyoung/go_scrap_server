package main

import (
	"awesomeProject/crawler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/crawler", naverCralwer)
	router.GET("/", home)
	router.Run("localhost:80")
}

func home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "welcome home!"})
}

func naverCralwer(c *gin.Context) {

	// 날짜 파라미터
	date := c.Param("date")

	// 네이버 해외 축구 뉴스기사 추출
	result := crawler.Crawler(date)

	if result {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "crawler Success!"})
	} else {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "crawler Failure"})
	}

}
