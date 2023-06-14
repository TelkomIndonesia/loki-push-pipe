package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/grafana/loki/pkg/loghttp/push"
	"github.com/grafana/loki/pkg/logql/syntax"
	_ "github.com/grafana/loki/pkg/push" // to be able to pin the version inside go.mod https://github.com/grafana/loki/blob/v2.8.1/go.mod#L336
)

type nooplogger struct{}

func (l nooplogger) Log(keyvals ...interface{}) error {
	return nil
}

type data struct {
	Labels    map[string]string `json:"labels"`
	Timestamp time.Time         `json:"timestamp"`
	Line      string            `json:"line"`
	TenantID  string            `json:"tenant_id"`
}

func main() {
	var noopl nooplogger
	var stdout = log.New(os.Stdout, "", 0)
	http.HandleFunc("/loki/api/v1/push", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("http panic", r)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			return
		}

		req, err := push.ParseRequest(noopl, "", r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("parse request failed:", err)
			return
		}

		d := data{TenantID: r.Header.Get("X-Scope-OrgID")}
		for _, s := range req.Streams {
			lbs, err := syntax.ParseLabels(s.Labels)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("parse labels failed,", err)
				return
			}

			d.Labels = lbs.Map()
			for _, e := range s.Entries {
				d.Line = e.Line
				d.Timestamp = e.Timestamp
				b, err := json.Marshal(d)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					log.Println("json encoding failed:", err)
					return
				}
				stdout.Println(string(b))
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}))

	log.Println("listening on :3100")
	if err := http.ListenAndServe(":3100", nil); err != nil {
		log.Fatal("http server exited:", err)
	}
}
