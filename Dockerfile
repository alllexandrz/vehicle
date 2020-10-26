FROM golang:1.15.2
WORKDIR /go/src/app/
COPY *.go .
RUN go get -d -v -t
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=0 /go/src/app/app .
CMD ["./app"]