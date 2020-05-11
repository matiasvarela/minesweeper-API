# minesweeper-API

## Docker

### Build image
````
$ docker build -t minesweeper:1.0.0 .
````

### Run in local
````
$ docker run -e ENV=local -p 8080:8080 -d minesweeper:1.0.0
````

### Run in production
````
$ docker run -e ENV=production -e AWS_ACCESS_KEY_ID=${access_key} -e AWS_SECRET_ACCESS_KEY=${secret_access_key} -p 8080:8080 -d minesweeper:1.0.0
````