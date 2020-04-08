package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dic <world>")
		os.Exit(0)
	}

	m := struct {
		Translation []string `json:"translation"`
		Basic       struct {
			USPhonetic string   `json:"us-phonetic"`
			Explains   []string `json:"explains"`
		} `json:"basic"`
	}{}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	word := url.QueryEscape(os.Args[1])
	data, err := Call(ctx, word)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) == 3 && os.Args[2] == "-d" {
		fmt.Println(string(data))
		return
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) == 2 {
		for _, v := range m.Translation {
			fmt.Println(v)
		}
		return
	}

	if len(os.Args) == 3 && os.Args[2] == "-d" {
		fmt.Println(m.Basic.USPhonetic)
		for _, v := range m.Basic.Explains {
			fmt.Println(v)
		}
	}
}

func Call(ctx context.Context, word string) (data []byte, err error) {
	done := make(chan error)

	go func() {
		data, err = query(word)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, errors.New("request time out")
	case <-done:
		return
	}
}

func query(word string) ([]byte, error) {
	domain := "http://fanyi.youdao.com"

	keyfrom := "easyfanyi"
	key := "1929637537"
	q := word
	myurl := domain + "/openapi.do?keyfrom=" + keyfrom + "&key=" + key + "&type=data&doctype=json&version=1.1&q=" + q

	resp, err := http.Get(myurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
