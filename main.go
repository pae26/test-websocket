package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Sticky struct {
	Id       int
	Locate_x int
	Locate_y int
	Text     string
}

var Db *sql.DB

func init() {
	log.Println("db start")
	var err error
	err = godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		panic(err)
	}

	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	PROTOCOL_NAME := "tcp(127.0.0.1:3306)"
	DB_STR := DB_USER + ":" + DB_PASS + "@" + PROTOCOL_NAME + "/" + DB_NAME
	Db, err = sql.Open("mysql", DB_STR)
	if err != nil {
		panic(err)
	}
	log.Println(DB_STR)
	log.Println("db connected")
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	rows, e := Db.Query("select * from sticky")
	if e != nil {
		log.Println(e.Error())
	}

	var result []Sticky

	for rows.Next() {
		s := Sticky{}
		if er := rows.Scan(&s.Id, &s.Locate_x, &s.Locate_y, &s.Text); er != nil {
			log.Println(er)
		}
		result = append(result, s)
	}

	log.Println("Locate_x:", result[0].Locate_x)

	defer rows.Close()

	data := map[string]string{
		"text":       result[0].Text,
		"Location_x": strconv.Itoa(result[0].Locate_x),
		"Location_y": strconv.Itoa(result[0].Locate_y),
	}

	t, err := template.ParseFiles(
		"static/home.html",
	)
	if err != nil {
		log.Fatalln("テンプレートファイルを読み込めません:", err.Error())
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatalln("エラー!:", err.Error())
	}
}

func main() {
	log.Println("Web server started")

	server := http.Server{
		Addr: ":9001",
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", templateHandler)
	server.ListenAndServe()
}
