package database

import (
	data "awesomeProject/type"
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	_ "github.com/sijms/go-ora/v2"
)

var autonomousDB = map[string]string{
	"service":  "g49ca3436b646ab_dovb_high.adb.oraclecloud.com",
	"username": "admin",
	"server":   "adb.us-ashburn-1.oraclecloud.com",
	"port":     "1522",
	"password": "!Qkswodms0626",
}

/** DB 초기화 */
func initDB() *sql.DB {
	connectionString := "oracle://" + autonomousDB["username"] + ":"
	connectionString += autonomousDB["password"] + "@" + autonomousDB["server"] + ":"
	connectionString += autonomousDB["port"] + "/" + autonomousDB["service"]
	path, _ := filepath.Abs("./database/Wallet_DOVB")
	repPath := strings.ReplaceAll(path, "\\", "/")
	fmt.Println(repPath)
	if repPath != "" {
		connectionString += "?TRACE FILE=trace.log&SSL=enable&SSL Verify=false&WALLET=" + url.QueryEscape(repPath)
	}

	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		panic("Sql Connection Fail")
	}

	return db
}

// InsertDB /** 데이터 삽입 */
func InsertDB(sportsNews []data.SportsNews) {
	db := initDB()
	var sequence int

	fmt.Println("오늘날짜 데이터 삭제진행..")
	// 오늘날짜에 대한 정보 제거
	db.Exec("DELETE FROM TB_SPORT_NEWS where TO_CHAR(REGIST_DATE, 'YYYYMMDD') = TO_CHAR(SYSDATE, 'YYYYMMDD') ")
	fmt.Println("오늘날짜 데이터 삭제완료")

	fmt.Println("데이터 삽입중")
	for _, news := range sportsNews {
		// 시퀀스 값 지정
		err := db.QueryRow("SELECT news_seq.NEXTVAL from dual").Scan(&sequence)
		if err != nil {
			panic(err)
		}

		// 정보 Insert
		_, err = db.Exec("INSERT INTO TB_SPORT_NEWS VALUES (:1, :2, :3, :4, CURRENT_TIMESTAMP)", sequence, news.Title, news.Thumnail, news.Link)

		if err != nil {
			panic(err)
		}

		fmt.Println("Data Insert Success!")
	}
	fmt.Println("삽입 완료")
}

func GetNews(date string) []data.ResponseMsg {
	db := initDB()

	var sportsNews []data.ResponseMsg
	query := "SELECT title, image, link FROM TB_SPORT_NEWS"

	if len(date) > 0 {
		query += " WHERE TO_CHAR(new_dtm, 'YYYYMMDD') = " + date
	}
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var title string
		var image string
		var link string

		err = rows.Scan(&title, &image, &link)

		sportsNews = append(sportsNews, data.ResponseMsg{Link: link, Image: image, Title: title})
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return sportsNews
}
