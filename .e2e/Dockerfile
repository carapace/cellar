FROM golang:1.11 AS builder


# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/carapace/cellar/.e2e
COPY . ./
RUN GO111MODULE="on" CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app main.go

FROM alpine
RUN apk update && apk upgrade
# Neither scratch or alpine comes with certs. Scratch is a lot harder to configure (no bash), so use alpine
COPY --from=builder /app ./
ENTRYPOINT ["./app"]