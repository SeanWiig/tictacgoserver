package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func requestRead(t *testing.T, s *Server) *httptest.ResponseRecorder {
	t.Helper()
	req, err := http.NewRequest("GET", "/game", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	s.readHandler(rr, req)
	return rr
}

func TestServer_readHandler(t *testing.T) {
	type move struct {
		player string
		row    int
		column int
	}
	tests := []struct {
		name         string
		setupMoves   []move
		expectedBody string
	}{
		{
			name:         "empty board",
			expectedBody: `{"board":[[null,null,null],[null,null,null],[null,null,null]],"turn":"X"}`,
		},
		{
			name: "O's turn",
			setupMoves: []move{
				{"X", 0, 0},
			},
			expectedBody: `{"board":[["X",null,null],[null,null,null],[null,null,null]],"turn":"O"}`,
		},
		{
			name: "win state",
			setupMoves: []move{
				{"X", 0, 0},
				{"O", 1, 0},
				{"X", 0, 1},
				{"O", 1, 1},
				{"X", 0, 2},
			},
			expectedBody: `{"board":[["X","X","X"],["O","O",null],[null,null,null]],"turn":"X","winner":"X"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer()
			for _, m := range tt.setupMoves {
				requestMove(t, s, m.player, m.row, m.column)
			}
			rr := requestRead(t, s)
			assert.Equal(t, http.StatusOK, rr.Code)
			require.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
	s := NewServer()
	rr := requestRead(t, s)
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `{"board":[[null,null,null],[null,null,null],[null,null,null]],"turn":"X"}`
	require.JSONEq(t, expected, rr.Body.String())
}

func requestMove(t *testing.T, s *Server, player string, row, column int) *httptest.ResponseRecorder {
	t.Helper()
	body := fmt.Sprintf(`{"row":%d,"column":%d,"player":"%s"}`, row, column, player)
	req, err := http.NewRequest("POST", "/game/move", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	s.moveHandler(rr, req)
	return rr
}

func TestServer_moveHandler(t *testing.T) {
	type move struct {
		player string
		row    int
		column int
	}
	tests := []struct {
		name         string
		setupMoves   []move
		move         move
		expectedCode int
	}{
		{
			name: "valid move",
			setupMoves: []move{
				{"X", 0, 0},
			},
			move:         move{"O", 0, 1},
			expectedCode: http.StatusCreated,
		},
		{
			name: "illegal move",
			setupMoves: []move{
				{"X", 0, 0},
			},
			move:         move{"O", 0, 0},
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer()
			for _, m := range tt.setupMoves {
				requestMove(t, s, m.player, m.row, m.column)
			}
			rr := requestMove(t, s, tt.move.player, tt.move.row, tt.move.column)
			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func TestServer_resetHandler(t *testing.T) {
	s := NewServer()
	requestMove(t, s, "X", 0, 0)
	req, err := http.NewRequest("DELETE", "/game", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	s.resetHandler(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	expected := `{"board":[[null,null,null],[null,null,null],[null,null,null]],"turn":"X"}`
	require.JSONEq(t, expected, requestRead(t, s).Body.String())
}
