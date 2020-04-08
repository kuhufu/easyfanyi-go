package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dic <world>")
		os.Exit(0)
	}

	domain := "http://fanyi.youdao.com"

	keyfrom := "easyfanyi"
	key := "1929637537"
	q := url.QueryEscape(os.Args[1])
	myurl := domain + "/openapi.do?keyfrom=" + keyfrom + "&key=" + key + "&type=data&doctype=json&version=1.1&q=" + q

	resp, err := http.Get(myurl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	m := struct {
		Translation []string `json:"translation"`
		Basic       struct {
			USPhonetic string   `json:"us-phonetic"`
			Explains   []string `json:"explains"`
		} `json:"basic"`
	}{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		panic("json unmarshal error")
	}

	if len(os.Args) == 2 {
		for _, v := range m.Translation {
			fmt.Println(v)
		}
		return
	}

	if len(os.Args) == 3 {
		fmt.Println(m.Basic.USPhonetic)
		for _, v := range m.Basic.Explains {
			fmt.Println(v)
		}
	}
}
