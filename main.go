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
	defer resp.Body.Close()
	if err != nil {
		panic("request error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	m := map[string]interface{}{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		panic("json unmarshal error")
	}

	basic := m["basic"].(map[string]interface{})
	if len(os.Args) == 3 {
		fmt.Println(basic["us-phonetic"])
	}

	ex := basic["explains"].([]interface{})
	for _, v := range ex {
		fmt.Println(v)
	}
}
