FROM golang:1.18

# appディレクトリを基準にDockerfileに書いた操作を実行
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
# go mod downloadはgo.modファイルで定義された依存関係をダウンロードし、モジュールキャッシュに保存
# go mod verify は各モジュールの.zipファイルと展開されたディレクトリをハッシュし、それらのハッシュをモジュールが最初にダウンロードされたときに記録されたハッシュと比較(モジュールキャッシュが最新であることを保証)

CMD ["go", "run", "main.go"]