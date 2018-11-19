FROM golang:latest

RUN CGO_ENABLED=0 GOOS=linux go get -u github.com/raviqqe/liche

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/bin/liche /go/bin/liche
ENV PATH "/go/bin:$PATH"
CMD liche
