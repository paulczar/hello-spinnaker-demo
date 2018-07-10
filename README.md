# go-hello-world

A simple hello world API service for demonstrating various kubernetes concepts.

The following API endpoints are available via GET method:

* `/` - Will attempt to serve the file `index.html` if not present will serve "hello world"
* `/health` - Will return health status in JSON
* `/???` - will serve a page that says "hello ???" where ??? is whatever string you used.


## Build Docker image

```
$ docker build -t hello-world .
...
...
$ docker run -p
```

## Usage

### Run hello-world

```
$ go run main.go
starting hello world app
```

or

```
$ docker run --rm -p 8080:8080 hello-world
starting hello world app
```

### Interact with hello-world

```
$ go run main.go &
curl localhost:8080      
<html><head><title>hello world</title></head><body>hello world!</body></html>

$ curl localhost:8080/health/
{"status":"UP"}

$ curl localhost:8080/apple
<html><head><title>hello apple</title></head><body>hello apple!</body></html>

$ echo "TEST" > index.html
$ curl localhost:8080        
TEST
```
