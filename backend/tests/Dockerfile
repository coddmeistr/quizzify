ARG ARCH=amd64
FROM golang:1.22-alpine AS builder

RUN apk add upx
# Add CA Certificates for those services communicating with outerworld
RUN apk add -U --no-cache ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOARCH=${ARCH} go build -ldflags="-s -w" -buildvcs=false ./cmd/tests-server
RUN upx ./tests-server

RUN adduser -SD quizzify

FROM scratch
WORKDIR /
COPY --from=builder  /app /tests-server
COPY --from=builder /etc/passwd /etc/passwd

USER quizzify
EXPOSE 80 8000

ENTRYPOINT ["/tests-server/tests-server"]
