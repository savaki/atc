package pipes

import (
	"io"
	"net/http"

	"github.com/concourse/atc"
)

func (s *Server) ReadPipe(w http.ResponseWriter, r *http.Request) {
	pipeID := r.FormValue(":pipe_id")

	dbPipe, err := s.db.GetPipe(pipeID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	closed := w.(http.CloseNotifier).CloseNotify()

	if dbPipe.URL == s.url {
		s.pipesL.RLock()
		pipe, found := s.pipes[pipeID]
		s.pipesL.RUnlock()

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)

		w.(http.Flusher).Flush()

		copied := make(chan struct{})
		go func() {
			io.Copy(w, pipe.read)
			close(copied)
		}()

	dance:
		for {
			select {
			case <-copied:
				break dance
			case <-closed:
				// connection died; terminate the pipe
				pipe.write.Close()
			}
		}

		s.pipesL.Lock()
		delete(s.pipes, pipeID)
		s.pipesL.Unlock()
	} else {
		response, err := s.forwardRequest(w, r, dbPipe.URL, atc.ReadPipe, dbPipe.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.(http.Flusher).Flush()

		copied := make(chan struct{})
		go func() {
			io.Copy(w, response.Body)
			close(copied)
		}()

	danceMore:
		for {
			select {
			case <-copied:
				break danceMore
			case <-closed:
				// connection died; terminate the pipe
				w.WriteHeader(http.StatusGatewayTimeout)
			}
		}
	}
}
