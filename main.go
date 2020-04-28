package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

// Daily is the object
type Daily struct {
	gorm.Model
	Date     string `json:"date"`
	Today    string `json:"today"`
	Tomorrow string `json:"tomorrow"`
	Point    string `json:"point"`
}

//DB初期化
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("DBの初期化に失敗しました（dbInit）")
	}
	db.AutoMigrate(&Daily{})
	defer db.Close()
}

//DB追加
func dbInsert(daily Daily) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("DBを開けませんでした（dbInsert)")
	}
	db.Create(&daily)
	defer db.Close()
}

//DB更新
func dbUpdate(daily Daily) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("DBを開けませんでした（dbUpdate)")
	}
	db.Save(&daily)
	db.Close()
}

//DB削除
func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("DBを開けませんでした（dbDelete)")
	}
	var daily Daily
	db.First(&daily, id)
	db.Delete(&daily)
	db.Close()
}

//DB全取得
func dbGetAll() []Daily {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！(dbGetAll())")
	}
	var dailies []Daily
	db.Order("created_at desc").Find(&dailies)
	db.Close()
	return dailies
}

//DB一つ取得
func dbGetOne(date string) Daily {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！(dbGetOne())")
	}
	var daily Daily
	// db.First(&daily, date)
	db.Where("date = ?", date).First(&daily)
	db.Close()
	return daily
}

// メイン関数
func main() {
	r := gin.Default()

	dbInit() // DB初期化syokika

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/dailies", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		var dailies []Daily
		dailies = dbGetAll()

		c.JSON(200, gin.H{
			"dailies": dailies,
		})
	})

	r.POST("/daily/create", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		var daily Daily
		daily.Date = c.PostForm("date")
		daily.Today = c.PostForm("today")
		daily.Tomorrow = c.PostForm("tomorrow")
		daily.Point = c.PostForm("point")

		newDaily := dbGetOne(daily.Date)

		if newDaily.Date != "" && daily.Today != "" && daily.Tomorrow != "" && daily.Point != "" {
			newDaily.Today = daily.Today
			newDaily.Tomorrow = daily.Tomorrow
			newDaily.Point = daily.Point
			dbUpdate(newDaily)
			c.JSON(200, gin.H{
				"message": "dbInsert done!",
				"daily":   newDaily,
				"status":  "OK",
			})
		} else {
			if daily.Today != "" && daily.Tomorrow != "" && daily.Point != "" && daily.Date != "" {
				dbInsert(daily)
				c.JSON(200, gin.H{
					"message": "dbInsert done!",
					"daily":   daily,
					"status":  "OK",
				})
				log.Printf(daily.Today, daily.Tomorrow, daily.Point)
				log.Printf("OK!")
			} else {
				c.JSON(200, gin.H{
					"message": "dbInsert failed...",
					"status":  "NG",
				})
				// log.Printf(daily.Today, daily.Tomorrow, daily.Point)
			}
		}

	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
