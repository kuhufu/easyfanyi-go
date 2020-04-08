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
	err := Call(ctx, &m)
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

	if len(os.Args) == 3 {
		fmt.Println(m.Basic.USPhonetic)
		for _, v := range m.Basic.Explains {
			fmt.Println(v)
		}
	}
}

func Call(ctx context.Context, v interface{}) error {
	done := make(chan error)

	go func() {
		err := call(v)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return errors.New("request time out")
	case err := <-done:
		return err
	}
}

func call(v interface{}) error {
	domain := "http://fanyi.youdao.com"

	keyfrom := "easyfanyi"
	key := "1929637537"
	q := url.QueryEscape(os.Args[1])
	myurl := domain + "/openapi.do?keyfrom=" + keyfrom + "&key=" + key + "&type=data&doctype=json&version=1.1&q=" + q

	resp, err := http.Get(myurl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}
