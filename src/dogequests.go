package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var (
	Shibes ShibesData
)

type ShibesData struct {
	Shibes []string
	Total  int
	Cursor int
}

func init() {
}

func getShibes() string {
	if Shibes.Cursor >= Shibes.Total {
		Shibes.Cursor = 0
		Shibes.Total = 10
		resp, err := http.Get("http://shibe.online/api/shibes?count=10")
		if err == nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(body, &Shibes.Shibes)
		}
	}
	Shibes.Cursor++
	return Shibes.Shibes[Shibes.Cursor - 1]
}
