FROM golang
COPY . /go/src/github.com/raviqqe/liche
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux go get /go/src/github.com/raviqqe/liche

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/liche /liche
ENTRYPOINT ["/liche"]
