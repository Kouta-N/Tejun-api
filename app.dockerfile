FROM golang:1.18

# appディレクトリを基準にDockerfileに書いた操作を実行
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

CMD ["go", "run", "main.go"]