package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	http.Handle("/",
		&MethodMux{
			map[string]func(http.ResponseWriter, *http.Request){
				"GET":  GetRequestHandler,
				"POST": PostRequestHandler,
			},
			NotImplementedHandler})

	fmt.Println("listening on " + os.Getenv("PORT") + "...")

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		panic(err)
	}
}

func GetRequestHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Asked for ", req.URL.Path[1:])

}

func PostRequestHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var entry Entry

	if err := decoder.Decode(&entry); err != nil {
		http.Error(res, err.Error(), 500)
	}

	fmt.Fprintln(res, entry.Uri)
}

func NotImplementedHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Method Not Implemented", 501)
}

type Entry struct {
	Uri string `json:"uri"`
}
