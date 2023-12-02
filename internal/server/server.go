package server

import (
	"encoding/json"
	"github.com/chrisfregly/tictactoe"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"sync"
)

type Server struct {
	game      tictactoe.TicTacToe
	gameMutex sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		game: tictactoe.NewTicTacToe(),
	}
}

func (s *Server) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/game", s.readHandler)
	r.Post("/game/move", s.moveHandler)
	r.Delete("/game", s.resetHandler)
	return r
}

type gameStateResponse struct {
	Board  [3][3]*tictactoe.Player `json:"board"`
	Turn   tictactoe.Player        `json:"turn"`
	Winner *tictactoe.Player       `json:"winner,omitempty"`
}

func (s *Server) read() gameStateResponse {
	s.gameMutex.RLock()
	defer s.gameMutex.RUnlock()
	return gameStateResponse{
		Board:  s.game.GetBoard(),
		Turn:   s.game.GetTurn(),
		Winner: s.game.GetWinner(),
	}
}

func (s *Server) readHandler(w http.ResponseWriter, _ *http.Request) {
	respond(w, http.StatusOK, s.read())
}

type moveRequest struct {
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Player string `json:"player"`
}

func (s *Server) move(request moveRequest) error {
	s.gameMutex.Lock()
	defer s.gameMutex.Unlock()
	return s.game.Move(tictactoe.Player(request.Player), request.Row, request.Column)
}

func (s *Server) moveHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	var request moveRequest
	if err = json.Unmarshal(b, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	if err = s.move(request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respond(w, http.StatusCreated, s.read())
}

func (s *Server) reset() {
	s.gameMutex.Lock()
	defer s.gameMutex.Unlock()
	s.game = tictactoe.NewTicTacToe()
}

func (s *Server) resetHandler(w http.ResponseWriter, _ *http.Request) {
	s.reset()
	respond(w, http.StatusNoContent, nil)
}

func respond(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}

type errorResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, statusCode int, err error) {
	respond(w, statusCode, errorResponse{err.Error()})
}
