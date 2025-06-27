# --- ステージ1: Goアプリケーションのビルド ---
FROM golang:1.24.4-alpine AS go_builder

# Goアプリケーションのビルドに必要なツールをインストールします。
# git は go mod download がプライベートリポジトリから取得する場合などに必要です。
# build-base は Go の CGO_ENABLED=0 ビルドに必要です。
RUN apk add --no-cache git build-base

# 作業ディレクトリを設定します。
WORKDIR /app

# Goモジュールファイル (go.mod, go.sum) をコピーし、依存関係をダウンロードします。
# これにより、ソースコード変更時に依存関係のダウンロードがキャッシュされる可能性が高まります。
COPY go.mod go.sum ./
RUN go mod download

# Goアプリケーションのソースコードをコピーします。
COPY . .

# Goアプリケーションをビルドします。
# CGO_ENABLED=0: 静的リンクされたバイナリを作成し、実行環境のCライブラリへの依存をなくします。
# GOOS=linux: Linux環境向けのバイナリを生成します。
# -a: 完全に静的なリンクを強制します。
# -installsuffix nocgo: CGO を無効にするためのオプションです。
# -o server: ビルドされたバイナリの名前を 'server' に指定します。
# . : 現在のディレクトリにあるGoモジュールをビルドします。
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o server .

# --- ステージ2: 最終的な実行イメージ (Playwright-Goのランタイム依存関係を含む) ---
FROM alpine:latest

# Playwrightがブラウザ (Chromiumを想定) を実行するために必要なシステム依存関係をインストールします。
# さらに、Playwrightドライバを実行するために必要な Node.js と npm も追加します。
RUN apk add --no-cache \
    nodejs \
    npm \
    nss \
    freetype \
    harfbuzz \
    brotli \
    libstdc++ \
    udev \
    ttf-freefont \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

# アプリケーションの作業ディレクトリを設定します。
WORKDIR /app

# タイムゾーンを設定します。
ENV TZ=Asia/Tokyo

# Goアプリケーションのビルド済みバイナリをコピーします。
COPY --from=go_builder /app/server .

# アプリケーションがリッスンするポートを公開 (情報提供のみ)
EXPOSE 8080

# コンテナ起動時に実行されるコマンドを設定します。
# ビルドステージで指定したバイナリ名 (`server`) を使用します。
CMD ["./server"]
