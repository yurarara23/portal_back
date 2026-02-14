# ステージ1: ビルド用
FROM architecture-dependent-golang-image AS builder 
# ↑ 実際は golang:1.22-alpine など
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 静的バイナリとしてビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ステージ2: 実行用（OSすら入っていない極小イメージ）
FROM alpine:latest
WORKDIR /app
# ビルドしたバイナリだけを持ってくる
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]