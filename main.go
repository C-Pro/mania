package main

import (
	"log"
	"net/http"

	"github.com/golang/protobuf/jsonpb"

	"mania/intents"

	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

const timeFormat = "2006-01-02 15:04:05"

func writeError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("error"))
}

func logRequest(r *http.Request) {
	log.Printf("%s %s: %s %s\n", r.RemoteAddr, r.Referer(), r.Method, r.RequestURI)
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	req := dialogflow.WebhookRequest{}

	if err := jsonpb.Unmarshal(r.Body, &req); err != nil {
		log.Println(err)
		writeError(w)
		return
	}
	defer r.Body.Close()

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	m := jsonpb.Marshaler{Indent: "  "}
	s, _ := m.MarshalToString(&req)
	log.Printf("REQ:\n%s\n", s)

	h, err := intents.GetHandler(req.QueryResult.Intent.DisplayName)
	if err != nil {
		log.Println(err)
		writeError(w)
		return
	}

	resp, err := h(req)
	if err != nil {
		log.Println(err)
		writeError(w)
		return
	}

	s, _ = m.MarshalToString(&resp)
	log.Printf("RESP:\n%s\n", s)

	m.Marshal(w, &resp)
}

func main() {
	http.HandleFunc("/", WebhookHandler)
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		panic(err)
	}
}
