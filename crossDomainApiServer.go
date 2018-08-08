package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
				"bytes"
	"fmt"
)

type Header struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type Data struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type Message struct {
	Method  string  `json:"method"`
	Url string `json:"url"`
	ContentType string `json:"contentType"`
	Headers []Header `json:"headers"`
	Datas []Data `json:"datas"`
}


func main() {
	http.HandleFunc("/crossApi", handler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// enable CORS
	setupResponse(&w, r)

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	/*output, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}*/

	resData, err := requestNaverShopApi(&msg, w)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(resData)
}


func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}


func requestNaverShopApi(msg *Message, w http.ResponseWriter) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://openapi.naver.com/v1/search/shop.json?query=%EC%A3%BC%EC%8B%9D&display=10&start=1&sort=sim", bytes.NewReader(nil))
	req.Header.Add("X-Naver-Client-Id", "ROhO_F9xQh2k9QS0dToc")
	req.Header.Add("X-Naver-Client-Secret", "yTtowUx8Rj")

	if err != nil {
		http.Error(w, err.Error(), 500)
		return []byte(""), err
	}


	resp, err := client.Do(req)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	fmt.Printf("%s\n", string(data))
	return data , nil
}