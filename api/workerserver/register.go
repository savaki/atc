package workerserver

import (
	"encoding/json"
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

type IntMetric int

func (i IntMetric) String() string {
	return strconv.Itoa(int(i))
}

var workerContainers = expvar.NewMap("WorkerContainers")

func (s *Server) RegisterWorker(w http.ResponseWriter, r *http.Request) {
	var registration atc.Worker
	err := json.NewDecoder(r.Body).Decode(&registration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(registration.Addr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "missing address")
		return
	}

	var ttl time.Duration

	ttlStr := r.URL.Query().Get("ttl")
	if len(ttlStr) > 0 {
		ttl, err = time.ParseDuration(ttlStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "malformed ttl")
			return
		}
	}

	workerContainers.Set(registration.Addr, IntMetric(registration.ActiveContainers))

	err = s.db.SaveWorker(db.WorkerInfo{
		Addr:             registration.Addr,
		ActiveContainers: registration.ActiveContainers,
		ResourceTypes:    registration.ResourceTypes,
		Platform:         registration.Platform,
		Tags:             registration.Tags,
	}, ttl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
