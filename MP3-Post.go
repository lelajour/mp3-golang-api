package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func createNewMp3File(playlist string, post PostRequest) bool {

	err := ioutil.WriteFile(filepath.Join(playlist, filepath.Base(post.FileName)),
		post.Data, 0644)

	if err != nil {
		return false
	}
	fmt.Println(post.FileName + " -> Created")
	return true
}

func addToPlaylist(num int) http.HandlerFunc {
	var playlists = [3]string{"/playlist1/", "/playlist2/", "/playlist3/"}

	playlist := playlists[num]
	return func(w http.ResponseWriter, r *http.Request) {
		var Post PostRequest

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&Post)

		if err != nil {
			log.Println(err.Error())
		}

		ret := createNewMp3File("."+playlist, Post)

		if ret != false {
			pl := [3]string{"One", "Two", "Three"}

			_, err2 := db.Exec("insert into "+pl[num]+"(artist, song, path) values(?, ?, ?)",
				Post.Artist, Post.Song, "http://192.168.1.148:10000"+playlist+Post.FileName)

			if err2 != nil {
				fmt.Println("Failed")
			}
		}

		json.NewEncoder(w).Encode("Succesfully added: " + Post.FileName)
	}
}
