FROM golang
COPY . /go/src/github.com/raviqqe/liche
WORKDIR /go/src/github.com/raviqqe/liche
RUN CGO_ENABLED=0 GOOS=linux go get .

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/liche /liche
ENTRYPOINT ["/liche"]
