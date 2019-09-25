FROM golang:alpine AS builder
RUN apk add --no-cache git && mkdir -p $GOPATH/src/github.com/DTherHtun/pet-with-peer
WORKDIR $GOPATH/src/github.com/DTherHtun/pet-with-peer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/pet-with-peer .
FROM scratch
COPY --from=builder /go/bin/pet-with-peer /go/bin/pet-with-peer
ENTRYPOINT ["/go/bin/pet-with-peer"]
EXPOSE 8080
