package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
		w.Write([]byte(err.Error()))
	}
	log.Print(response)
	json.NewEncoder(w).Encode(response)
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
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	router.HandleFunc("/", echoHandler)
	router.HandleFunc("/*", echoHandler)
	port := os.Getenv("ECHOPORT")
	if port == "" {
		port = "8081"
	}
	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server on port: " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
