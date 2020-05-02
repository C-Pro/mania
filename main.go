package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mania/intents"
	"mania/store"

	"mania/dialogflow"
)

const (
	maxTries = 10
)

func writeError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write([]byte("error")); err != nil {
		log.Printf("failed to write error response: %v", err)
	}
}

func logRequest(r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	log.Printf("%s %s: %s %s\nbody: '%s'\n",
		r.RemoteAddr,
		r.Referer(),
		r.Method,
		r.RequestURI,
		string(body))
}

// MakeWebhookHandler returns handler function with dispatcher
func MakeWebhookHandler(d *intents.Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		req := dialogflow.Request{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			writeError(w)
			return
		}
		defer r.Body.Close()

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		h, err := d.GetHandler(req.QueryResult.Intent.DisplayName) //???
		if err != nil {
			log.Printf("failed to get handler: %v", err)
			writeError(w)
			return
		}

		resp, err := h(req)
		if err != nil {
			log.Println(err)
			writeError(w)
			return
		}

		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			log.Printf("failed to write response to %s: %v", r.RemoteAddr, err)
		}
	}
}

func initCache(ctx context.Context) *store.Cache {
	var (
		s   *store.Cache
		err error
	)

	for i := 0; i < maxTries; i++ {
		s, err = store.NewCache(ctx)
		if err != nil {
			log.Printf("failed to create Cache instance: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		log.Fatalf("failed to initialize Cache after %d tries", maxTries)
	}

	return s
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	s := initCache(ctx)
	d := intents.NewDispatcher(ctx, s)
	handlerFunc := MakeWebhookHandler(d)

	http.HandleFunc("/", handlerFunc)

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		addr := fmt.Sprintf("0.0.0.0:%s", port)
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err)
		}
	}()

	<-sigs
	cancel()
	log.Printf("shutting down")
	time.Sleep(time.Second)
}
