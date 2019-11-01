FROM golang:1.12
COPY . /src
ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux
WORKDIR /src
RUN go build -o liche

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=0 /src/liche /liche
ENTRYPOINT ["/liche"]
