
FROM golang:1.25.5-alpine AS builder


RUN apk add --no-cache make gcc musl-dev


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY Makefile ./


RUN make build


FROM alpine:3.18


WORKDIR /app


COPY --from=builder /app/bin/* ./main


RUN chmod +x ./main


CMD ["./main"]