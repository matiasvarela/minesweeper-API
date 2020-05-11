# minesweeper-API
minesweeper-API is a rest api for the popular game Minesweeper. It provides the following endpoints:

1. Create a new game
2. Get a game by id
3. Reveal a cell
4. Mark a cell with a flag 

## Notes about this project
- It uses GinGonic to manage routing
- It uses an error handling library written by my own
- It follows a Hexagonal Architecture
- It uses a DynamoDB to persist the games
- It provides a Dockerfile to build and run the application
- It provides a demo application
- It provides a client lib written in python 
- Testing with ~80% of code coverage. It follows the Table Driven Test pattern (https://github.com/golang/go/wiki/TableDrivenTests)

## Decisions
- All the endpoints that return the game representation hides the cells with bombs replacing them with empty cells, so it is impossible for clients of this API to know the positions of those bombs.

## Demo

```
$ curl http://ec2-3-14-1-190.us-east-2.compute.amazonaws.com:8080/ping
```

## Getting started

### Docker

#### Build image
````
$ docker build -t minesweeper:1.0.0 .
````

#### Run in local environment
To run this application locally is necessary to run the local version of dynamodb in port 8000. See https://hub.docker.com/r/amazon/dynamodb-local/

````
$ docker run -e ENV=local -p 8080:8080 -d minesweeper:1.0.0
````

#### Run in production environment
````
$ docker run -e ENV=production -e AWS_ACCESS_KEY_ID=${access_key} -e AWS_SECRET_ACCESS_KEY=${secret_access_key} -p 8080:8080 -d minesweeper:1.0.0
````

### Local
To run this application locally is necessary to run the local version of dynamodb in port 8000. See https://hub.docker.com/r/amazon/dynamodb-local/

```
$ go run cmd/restserver/main.go
```

## API Documentation

### Create a new game

```http
POST /games
```

Body
```json
{
    "rows": 4,
    "columns": 4,
    "bombs_number": 5
}
```

Response

The following json correspond with a `game` and from now on we will call it `game_json` 

```json
{
  "id": "7ecbe4ee-4f1d-426a-bf8a-0d2382d61805",
  "board": [
    ["e","e","e","e"],
    ["e","e","e","e"],
    ["e","e","e","e"],
    ["e","e","e","e"]
  ],
  "settings": {
    "rows": 10,
    "columns": 10,
    "bombs_number": 5
  },
  "state": "new",
  "started_at": "0001-01-01T00:00:00Z",
  "ended_at": "0001-01-01T00:00:00Z"
}
```

The `id` attribute is the unique id for the created game.

The `board` attribute is a matrix of cells that represents the board of the game.

| Cell | Description |
| :--- | :--- |
| e | covered cell |
| E | empty revealed cell |
| X | marked cell with a flag |
| B | cell with a bomb | 

The `settings` attribute contains the settings used to create the game.

The `state` attribute indicates the current state of the game 

| State | Description |
| :--- | :--- |
| new | the game has not began  |
| ongoing | the game has began and has not finished |
| lost | the game is over and resulted lost because a bomb has been revealed  |
| won | the game is over and resulted won because all the empty cells has been revealed | 

The `started_at` attribute indicates the time when the first cell has been revealed.

The `ended_at` attribute indicates the time when the game ended.

### Get a game by id
Get a previously created game given its unique id. 

```http
GET /games/:id
```

Response

1. `game_json` if the game exists
2. Not found
```json
{
  "status": 404,
  "code": "not_found",
  "message": "the game has not been found"
}
``` 

### Mark a cell with a flag
Mark a cell with a flag. A cell marked by flag means that that particular cell cannot be revealed unless it is unmarked. 

```http
PUT /games/:id/mark
```
Body

```json
{
    "row": 2,
    "column": 2
}
```

The attributes `row` and `column` refers to a particular position within the board.

Response

1. `game_json` is game has been marked/unmarked successfully
2. Not found
```json
{
  "status": 404,
  "code": "not_found",
  "message": "the game has not been found"
}
``` 
3. Invalid row and column
```json
 {
   "status": 400,
   "code": "invalid_input",
   "message": "invalid row and column parameters"
 }
 ``` 

### Reveal a cell
Reveals a particular cell. If there is no adjacent bombs then all the adjacent (except those marked with a flag) will be revealed repeating this process until no other cell can be revealed. 

```http
PUT /games/:id/reveal
```
Body

```json
{
    "row": 2,
    "column": 2
}
```

The attributes `row` and `column` refers to a particular position within the board.

Response

1. `game_json` is game has been revealed successfully
2. Not found
```json
{
  "status": 404,
  "code": "not_found",
  "message": "the game has not been found"
}
``` 
3. Invalid row and column
```json
 {
   "status": 400,
   "code": "invalid_input",
   "message": "invalid row and column parameters"
 }
 ``` 



