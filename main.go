package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

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
	log.Println("db connected")
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"static/home.html",
	)
	if err != nil {
		log.Fatalln("テンプレートファイルを読み込めません:", err.Error())
	}
	if err := t.Execute(w, nil); err != nil {
		log.Fatalln("エラー!:", err.Error())
	}
}

func setCookies(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "hoge",
		Value: "bar",
	}
	http.SetCookie(w, cookie)

	fmt.Fprintf(w, "Cookieの設定ができたよ")
}

func showCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("hoge")

	if err != nil {
		log.Fatal("Cookie: ", err)
	}

	log.Println(cookie)
	// tmpl := template.Must(template.ParseFiles("./cookie.html"))
	// tmpl.Execute(w, cookie)

}

func getStickiesInfo(w http.ResponseWriter, r *http.Request) {
	rows, e := Db.Query("select * from sticky")
	if e != nil {
		log.Println(e.Error())
	}

	var stickies []Sticky

	for rows.Next() {
		sticky := Sticky{}
		if er := rows.Scan(
			&sticky.Id,
			&sticky.Locate_x,
			&sticky.Locate_y,
			&sticky.Text,
		); er != nil {
			log.Println(er)
		}
		stickies = append(stickies, sticky)
	}

	defer rows.Close()

	result, err := json.Marshal(stickies)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func updateSticky(w http.ResponseWriter, r *http.Request) {
	var sticky Sticky
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	if err := json.Unmarshal(body[:len], &sticky); err != nil {
		log.Fatalln("エラー!:", err)
	}

	sql, err := Db.Prepare("update sticky set locate_x=?, locate_y=? where id=?")
	if err != nil {
		log.Fatalln(err)
	}
	sql.Exec(sticky.Locate_x, sticky.Locate_y, sticky.Id)

	res, err := json.Marshal("{200, \"ok\"}")
	if err != nil {
		log.Fatalln(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	log.Println("Web server started")
	r := newRoom()
	http.Handle("/room", r)
	go r.run()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", templateHandler)
	http.HandleFunc("/stickies", getStickiesInfo)
	http.HandleFunc(("/update-sticky"), updateSticky)
	http.HandleFunc("/set_cookie", setCookies)
	http.HandleFunc("/show_cookie", showCookie)

	server := http.Server{
		Addr: ":9001",
	}
	server.ListenAndServe()
}
