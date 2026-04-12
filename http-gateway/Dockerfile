FROM golang:1.25.5-alpine AS builder   

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 

COPY . .
RUN go build -o httpgateway ./cmd/main.go

FROM alpine:latest

COPY --from=builder app/configs ./configs/
COPY --from=builder app/docs ./docs/
COPY --from=builder app/httpgateway ./

CMD ["./httpgateway"]
