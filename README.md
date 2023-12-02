# tictacgoserver

## How to run

```bash
$ go run main.go
```

Alternatively, you can build and run with Docker:

```bash
$ docker build -t tictacgoserver .
$ docker run -p 8080:8080 tictacgoserver
```

## How to play

`GET /game` to get the current game state
```bash
$ curl http://localhost:8080/game
{"board":[[null,null,null],[null,null,null],[null,null,null]],"turn":"X"}
```
`POST /game/move` to make a move, returning the new game state
```bash
$ curl -X POST http://localhost:8080/game/move -d '{"player": "X", "row": 0, "column": 0}'
{"board":[["X",null,null],[null,null,null],[null,null,null]],"turn":"O"}
```

`DELETE /game` to reset the game
```bash
$ curl -X DELETE http://localhost:8080/game
$ curl http://localhost:8080/game
{"board":[[null,null,null],[null,null,null],[null,null,null]],"turn":"X"}
```

## How to test

```bash
$ go test ./...
```