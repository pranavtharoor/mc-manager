FROM golang:1.16.6-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o bin/mcbot cmd/mcbot/mcbot.go

FROM mcr.microsoft.com/azure-cli
RUN apk add --no-cache ffmpeg
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl && \
    chmod a+rx /usr/local/bin/youtube-dl
WORKDIR /app
COPY --from=builder /app/bin/mcbot .
CMD ["./mcbot"]
