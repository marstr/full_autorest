FROM golang:1.10 AS builder

# Populate builder with current source
WORKDIR $GOPATH/src/github.com/Azure/full_autorest
ADD . .

# Package Management
RUN go get -u github.com/golang/dep/cmd/dep && \
    dep ensure && \
    go install

FROM marstr/autorest
COPY --from=builder /go/bin/full_autorest .
EXPOSE 80
ENTRYPOINT ["./full_autorest", "start"]