# Stage: build
FROM golang:alpine AS build

RUN apk add --no-cache git upx

WORKDIR /go/src/app
COPY *.go ./
RUN go get -d -v ./...

RUN CGO_ENABLED=0 go build -o disable-check -v main.go
RUN upx --brute disable-check

# Stage: package
FROM scratch
LABEL maintainer="Ruben Vermeersch <ruben@rocketeer.be>"
COPY --from=build /go/src/app/disable-check /
ENTRYPOINT ["/disable-check"]
