FROM golang:alpine AS builder

# Git is required for fetching the dependencies
RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/alfcope/checkout-test/

COPY . .

# Using go mod.
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

RUN cd cmd/checkoutserver && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/checkout-service

# Create serviceuser
RUN adduser -D -g '' serviceuser



FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /bin/checkout-service /bin/checkout/checkout-service
COPY ./config/*.yml ./config/*.json /bin/checkout/config/

# Use the unprivileged user serviceuser
USER serviceuser

WORKDIR /bin/checkout/
ENTRYPOINT ["./checkout-service"]
