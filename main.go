package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Traceableai/goagent"
	"github.com/Traceableai/goagent/config"
	"github.com/Traceableai/goagent/instrumentation/github.com/gorilla/traceablemux"
	"github.com/gorilla/mux"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["ip"] = r.RemoteAddr
	response["headers"] = r.Header
	response["query"] = r.URL.Query()
	response["form"] = r.Form
	response["url"] = r.URL
	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding body: ", err)
	} else {
		log.Println("Body: ", data)
	}
	response["data"] = data
	str_data := fmt.Sprintf("%v", data)
	reg, _ := regexp.Compile("(?i:(sleep|wait|receive))")
	if reg.FindString(str_data) != "" {
		log.Println("Sleeping")
		time.Sleep(time.Second * 5)
	} else {
		log.Println("Not sleeping")
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

var (
	port *int
)

func init() {
	port = flag.Int("port", 3000, "port number")
}

func main() {
	cfg := config.Load()
	cfg.Tracing.ServiceName = config.String("goservice")

	shutdown := goagent.Init(cfg)
	defer shutdown()

	router := mux.NewRouter()
	router.Use(traceablemux.NewMiddleware())

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	router.HandleFunc("/", echoHandler)
	router.HandleFunc("/*", echoHandler)
	porti := os.Getenv("PORT")
	flag.Parse()
	if porti == "" {
		porti = strconv.Itoa(*port)
	}
	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:" + porti,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server on port: " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
