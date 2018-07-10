FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app .

FROM alpine
COPY --from=0 /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
