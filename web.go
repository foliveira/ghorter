package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type MethodMux struct {
	methods  map[string]func(http.ResponseWriter, *http.Request)
	delegate func(http.ResponseWriter, *http.Request)
}

func (mm *MethodMux) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if f, ok := mm.methods[request.Method]; ok {
		f(response, request)
	} else {
		mm.delegate(response, request)
	}
}

var entries = make(map[string]string)

func main() {
	port := os.Getenv("PORT")

	http.Handle("/",
		&MethodMux{
			map[string]func(http.ResponseWriter, *http.Request){
				"GET":  GetRequestHandler,
				"POST": PostRequestHandler,
			},
			NotImplementedHandler})

	log.Println("listening on " + port + "...")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func GetRequestHandler(res http.ResponseWriter, req *http.Request) {
	uri := req.URL.Path[1:]
	http.Redirect(res, req, entries[uri], 302)
}

func PostRequestHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var entry Entry

	if err := decoder.Decode(&entry); err != nil {
		http.Error(res, err.Error(), 500)
	}

	encodedUri := base64.StdEncoding.EncodeToString([]byte(entry.Uri))
	entries[encodedUri] = entry.Uri

	fmt.Fprintln(res, encodedUri)
}

func NotImplementedHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Method Not Implemented", 501)
}

type Entry struct {
	Uri string `json:"uri"`
}
