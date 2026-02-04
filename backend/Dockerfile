# 開発用 Dockerfile
FROM golang:1.25-alpine

WORKDIR /app

# 必要なツールをインストール
RUN apk add --no-cache git gcc musl-dev

# Air（ホットリロード）をインストール
RUN go install github.com/air-verse/air@latest

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードはボリュームマウントで共有
COPY . .

EXPOSE 8080

# デフォルトはairでホットリロード
CMD ["air", "-c", ".air.toml"]