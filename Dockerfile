FROM golang:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
ADD /tmp/bin/liche /go/bin/liche
# ENV PATH "/go/bin:$PATH"
RUN liche --help
CMD liche
