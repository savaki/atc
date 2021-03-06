package buildserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/present"
)

func (s *Server) CreateBuild(w http.ResponseWriter, r *http.Request) {
	var plan atc.Plan
	err := json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	build, err := s.db.CreateOneOffBuild()
	if err != nil {
		s.logger.Error("failed-to-create-one-off-build", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	engineBuild, err := s.engine.CreateBuild(build, plan)
	if err != nil {
		s.logger.Error("failed-to-start-build", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go engineBuild.Resume(s.logger)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(present.Build(build))
}
