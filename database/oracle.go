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
