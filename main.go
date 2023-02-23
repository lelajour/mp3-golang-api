package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type Playlist struct {
	Id     int    `json:"id"`
	Artist string `json:"Artist"`
	Song   string `json:"Song"`
	Path   string `json:"Path"`
}

type PostRequest struct {
	FileName string `json:"FileName"`
	Artist   string `json:"Artist"`
	Song     string `json:"Song"`
	Data     []byte `json:"Data"`
}

type DeleteRequest struct {
	Id   int    `json:"Id"`
	Path string `json:"Path"`
}

var db *sql.DB

var Ret = make([][]Playlist, 3)

var db_user = os.Getenv("PLAYLIST_DB_USER")
var db_psw = os.Getenv("PLAYLIST_DB_PSW")
var token = os.Getenv("PLAYLIST_BEARER_TOKEN")
var playlists_path = os.Getenv("PLAYLISTS_PATH")

func ErrorCheck(e error) {
	if e != nil {
		log.Println(e.Error())
	}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")

		var Auth = r.Header.Get("Authorization")

		if Auth != "Bearer "+token {
			json.NewEncoder(w).Encode("Wrong auth token")
			fmt.Println("Wrong auth token from " + r.RemoteAddr)

			return
		}
		next.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: home")
	json.NewEncoder(w).Encode("Endpoint hit : home")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(Middleware)
	myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/one", addToPlaylist(0)).Methods("POST")
	myRouter.HandleFunc("/one", deleteFromPlaylist(0)).Methods("DELETE")
	myRouter.HandleFunc("/one", returnPlaylist(0))
	myRouter.HandleFunc("/two", addToPlaylist(1)).Methods("POST")
	myRouter.HandleFunc("/two", deleteFromPlaylist(1)).Methods("DELETE")
	myRouter.HandleFunc("/two", returnPlaylist(1))
	myRouter.HandleFunc("/three", addToPlaylist(2)).Methods("POST")
	myRouter.HandleFunc("/three", deleteFromPlaylist(2)).Methods("DELETE")
	myRouter.HandleFunc("/three", returnPlaylist(2))
	myRouter.PathPrefix("/playlist1/").Handler(http.StripPrefix("/playlist1/", http.FileServer(http.Dir(playlists_path+"Playlist1"))))
	myRouter.PathPrefix("/playlist2/").Handler(http.StripPrefix("/playlist2/", http.FileServer(http.Dir(playlists_path+"Playlist2"))))
	myRouter.PathPrefix("/playlist3/").Handler(http.StripPrefix("/playlist3/", http.FileServer(http.Dir(playlists_path+"Playlist3"))))

	log.Fatal(http.ListenAndServe(":10000", myRouter))

}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	var err error

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	login := db_user + ":" + db_psw
	db, err = sql.Open("mysql", login+"@tcp(localhost:3306)/MP3-Playlists")

	if err != nil {
		log.Println(err.Error())
	} else {
		fmt.Println("Connected to database")
		handleRequests()
	}

	defer db.Close()
}
