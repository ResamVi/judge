FROM golang:1.24
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

ENTRYPOINT ["/server"]
