package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func returnPlaylist(num int) http.HandlerFunc {
	var playlists = [3]string{`SELECT * FROM One;`, `SELECT * FROM Two;`, `SELECT * FROM Three;`}

	playlist := playlists[num]

	return func(w http.ResponseWriter, r *http.Request) {
		Ret[num] = nil
		var songData = Playlist{}

		rows, e := db.Query(playlist)

		ErrorCheck(e)

		defer rows.Close()

		for rows.Next() {
			e = rows.Scan(&songData.Id, &songData.Artist, &songData.Song, &songData.Path)
			ErrorCheck(e)
			fmt.Println(songData)
			Ret[num] = append(Ret[num], songData)
		}

		json.NewEncoder(w).Encode(Ret[num])
	}
}
