FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
ADD bin/liche /go/bin/liche
# ENV PATH "/go/bin:$PATH"
RUN liche --help
CMD liche
