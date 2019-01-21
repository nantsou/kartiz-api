# first stage
FROM golang:alpine AS build-env

# install required tools
RUN apk update && apk add git && go get gopkg.in/natefinch/lumberjack.v2

# add the source codes
ADD src /go/src

# download dependencies
ENV GO111MODULE=on
RUN cd /go/src && go mod download

# build the executable file
RUN cd /go/src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kartiz

# final stage
FROM centurylink/ca-certs
COPY --from=build-env /go/src/kartiz /
ENTRYPOINT ["/kartiz"]
EXPOSE 8080
