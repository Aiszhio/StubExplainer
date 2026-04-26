FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o /bin/explainer ./cmd/explainer

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /bin/explainer /bin/explainer

CMD ["/bin/explainer"]
