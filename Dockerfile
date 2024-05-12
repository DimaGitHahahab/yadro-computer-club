FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

ENV GOOS linux

RUN go build -o comp-club ./cmd/main.go

FROM alpine AS runner

WORKDIR /root/

COPY --from=builder /app/comp-club .

ENTRYPOINT ["./comp-club"]