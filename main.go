package main

import (
	"awesomeProject/crawler"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	"net/http"
	"time"
)

func main() {
	// 고루틴 동시 실행
	// 스케쥴러 실행
	startScheduler()
	// 서버 실행
	//startServer()
}

func startScheduler() {
	fmt.Println("스케줄 시작합니다.")
	gocron.Every(1).Day().At("18:00").Do(func() {
		fmt.Println(time.Now())
		fmt.Println("파싱작업 스케줄 진행")
		crawler.Crawler("")
	})

	<-gocron.Start()
}

func startServer() {
	router := gin.Default()
	router.GET("/getNews", getNewsData)
	router.GET("cralwer", naverCralwer)
	router.GET("/", home)
	router.Run("localhost:80")
}

func home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Welcome My Server!"})
}

func getNewsData(c *gin.Context) {
	// 날짜 파라미터
	date := c.Param("date")

	news := crawler.GetNewData(date)

	c.JSON(http.StatusOK, news)

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
