package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

func deleteFromPlaylist(num int) http.HandlerFunc {
	playlist := [3]string{"One", "Two", "Three"}

	return func(w http.ResponseWriter, r *http.Request) {
		var Del DeleteRequest

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&Del)

		if err != nil {
			log.Println(err.Error())
		}

		path := strings.TrimPrefix(Del.Path, "http://192.168.1.148:10000")

		if _, err := os.Stat(path); err == nil {

			err = os.RemoveAll(Del.Path)

			if err == nil {
				_, err2 := db.Exec("DELETE FROM "+playlist[num]+" WHERE id = ?", Del.Id)
				if err2 != nil {
					log.Println("Failed to delete file")
				}
			}
		}

		json.NewEncoder(w).Encode("Entry succesfully deleted")
	}
}
