package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
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

var entries = make(map[string]interface{})

func main() {
	port := os.Getenv("PORT")
	rand.Seed(time.Now().UTC().UnixNano())

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
	shorty := req.URL.Path[1:]

	if entry := entries[shorty]; entry == nil {
		http.Error(res, "Entity not found", 404)
	} else {
		decoded := 0
		multi := 1
		num := entry.(string)
		for len(num) > 0 {
			digit := num[len(num)-1]
			decoded = decoded + multi*strings.Index(alphabet, string(digit))
			multi = multi * len(alphabet)
			num = num[:len(num)-1]
		}

		entity := Entry{
			entries[num].(string)}

		encoder := json.NewEncoder(res)
		if err := encoder.Encode(entity); err != nil {
			http.Error(res, err.Error(), 500)
		}
	}
}

func PostRequestHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var entry Entry

	if err := decoder.Decode(&entry); err != nil {
		http.Error(res, err.Error(), 415)
	}

	uniqueId := rand.Int31()
	baseCount := int32(len(alphabet))
	encoded := string("")

	log.Println("Unique request number: ", uniqueId)

	baseId := uniqueId
	for baseId >= baseCount {
		mod := baseId % baseCount
		div := baseId / baseCount

		encoded = encoded + string(alphabet[mod])
		baseId = div
	}

	entries[fmt.Sprintf("%v", uniqueId)] = entry.Uri

	fmt.Fprintln(res, encoded)
}

func NotImplementedHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Method Not Allowed", 405)
}

type Entry struct {
	Uri string `json:"uri"`
}

var alphabet = string("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")
