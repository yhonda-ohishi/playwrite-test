# --- ビルドステージ (builder) ---
# Playwright の公式イメージ (Ubuntu Nobleベース) を使用します。
# このイメージには Node.js、Playwright ライブラリ、ブラウザバイナリがプリインストールされています。
# ただし、Go コンパイラは含まれていません。
# FROM mcr.microsoft.com/playwright:v1.53.0-noble AS builder
FROM golang:1.24.4-alpine AS go_builder
# Goのバージョンを1.22に設定
# AlpineベースのGoイメージを使用します。

# 作業ディレクトリを設定
WORKDIR /app

# Go コンパイラをインストールする
# Ubuntuベースなので apt-get を使用します。
# Go 1.24.4 に対応するバージョンをインストールします。
# リポジトリによっては最新版が提供されないこともあるので、go.mod のバージョンと合うか確認してください。
# もし特定のバージョンが必要なら、goenv のようなツールを使う手もありますが、ここではシンプルに。
# RUN apt-get update && apt-get install -y --no-install-recommends \
#     golang-go \
#     git \
#     ca-certificates \
#     && rm -rf /var/lib/apt/lists/*

RUN apk add --no-cache git build-base
# Goのビルドツールをインストール
# Alpineベースなので apk を使用します。


# Goモジュールのバージョンが go.mod と一致することを確認する（推奨）
# この行はデバッグ用であり、通常は不要かもしれません。
RUN go version

# プライベートモジュールに依存している場合のみ残す
# もし go.mod にプライベートモジュールへの require がないなら、これらも削除してください。
# ENV GOPRIVATE github.com/yhonda-ohishi/playwrite-test/* # 必要なら修正
# ARG GITHUB_TOKEN
# RUN if [ -n "$GITHUB_TOKEN" ]; then \
#     git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com".insteadOf "https://github.com"; \
#     fi

# Goモジュールファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
# CGO_ENABLED=0 は静的リンクされたバイナリを作成するために推奨されます。
# GOOS=linux はこのイメージもLinuxベースなので問題ありません。
# installsuffix nocgo は CGO を無効にするためのオプションです。
# これにより、Goのビルドが軽量になり、依存関係が少なくなります。
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o server .

# --- 最終イメージステージ ---
# アプリケーションが実行時にもPlaywrightを使うなら、同じイメージをベースにすべきです。
# Playwright は Node.js プロセスやブラウザバイナリに依存するため、
# 実行環境にもそれらが必要です。
# FROM mcr.microsoft.com/playwright:v1.53.0-noble

# --- Playwrightを動かすためのNode.js環境 ---
# FROM node:20-alpine AS playwright_runner
FROM alpine:latest


# 必要なパッケージのインストール
# Playwrightはいくつかのブラウザに依存するため、それらもインストールします
# Chromiumの依存関係は特に重要です。
# https://playwright.dev/docs/intro#requirements を参照
RUN apk add --no-cache \
    nss \
    freetype \
    harfbuzz \
    brotli \
    libstdc++ \
    udev \
    ttf-freefont \
    # その他の一般的なユーティリティ（ログ出力、ファイル操作など）
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*



# 作業ディレクトリを設定します。
WORKDIR /app


ENV TZ=Asia/Tokyo

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドステージからコンパイル済みバイナリをコピー
# COPY --from=builder /app/server .
COPY --from=go_builder /usr/local/bin/server .

# アプリケーションがリッスンするポートを公開 (情報提供のみ)
EXPOSE 8080

# アプリケーションの実行コマンド
CMD ["./server"]