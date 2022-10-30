package crawler

import (
	"awesomeProject/database"
	data "awesomeProject/type"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"strconv"
)

func Crawler(date string) bool {
	url := "https://sports.news.naver.com/wfootball/news/index?isphoto=N"
	if len(date) > 0 {
		url += "&date="
	}

	// 옵션 브라우저 창 띄우지 않고 진행
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", true))

	contextVar, cancelFunc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelFunc()

	contextVar, cancelFunc = chromedp.NewContext(contextVar)
	defer cancelFunc()

	//
	var sportsNews []data.SportsNews
	// 뉴스 링크
	var links []*cdp.Node

	// 메인 링크 경로를 아이디로 사용
	var mainLinks []string

	var index = 1
	var err error
Loop:
	for {
		var tmpUrl = url
		if index != 1 {
			tmpUrl += "&page=" + strconv.Itoa(index)
		}

		fmt.Println(tmpUrl + "파싱중....")
		err = chromedp.Run(contextVar, chromedp.Navigate(tmpUrl),
			chromedp.WaitVisible(`div.news_list`),
			chromedp.Nodes("div.news_list > ul > li > a", &links),
		)

		// 검색되는 정보가 없는 경우 반복문 종료(동작하지 않음 네이버는 페이지 넘버가 넘어가도 마지막 페이지를 요청함)
		if len(links) == 0 {
			fmt.Println("모든 페이지 파싱 수행 완료")
			break Loop
		}

		//fmt.Println(strVar)
		for _, info := range links {
			link, _ := info.Attribute("href")

			// 메인 링크가 이미 들어있을 경우 종료
			if findMainLink(mainLinks, link) {
				fmt.Println("모든 페이지 파싱 수행 완료")
				break Loop
			}

			// 중단점을 찾기위한 부분
			mainLinks = append(mainLinks, link)

			for _, child := range info.Children {
				img, _ := child.Attribute("src")
				title, _ := child.Attribute("alt")

				sportsNews = append(sportsNews, data.SportsNews{Link: link, Thumnail: img, Title: title})
				break
			}
		}
		index++
		fmt.Println("파싱완료")
	}
	database.InsertDB(sportsNews)
	//fmt.Println("정보 출력 진행")
	//
	//time.Sleep(5000)
	//for _, news := range sportsNews {
	//	fmt.Println(news.Link, news.Title, news.Thumnail)
	//}

	if err != nil {
		return false
	}

	return true
}

func findMainLink(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
