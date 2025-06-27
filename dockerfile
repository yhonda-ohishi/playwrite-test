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

# --- ステージ2: 最終的な実行イメージ (軽量なDebianベース) ---
# Playwrightの公式イメージの代わりに、より軽量なDebianベースのイメージを使用します。
# これにより、必要なランタイム依存関係を最小限に抑えつつ、glibc互換性を維持します。
FROM debian:stable-slim

# パッケージリストを更新し、Node.jsとPlaywrightに必要なランタイム依存をインストールします。
# 必要なライブラリは、Playwrightの公式ドキュメントやChromiumの依存関係リストを参照し、
# 最小限に絞り込んでいます。
RUN apt-get update && apt-get install -y --no-install-recommends \
    nodejs \
    npm \
    libnss3 \
    libfontconfig1 \
    libgbm1 \
    libglib2.0-0 \
    libgdk-pixbuf2.0-0 \
    libgtk-3-0 \
    libasound2 \
    libgconf-2-4 \
    libnotify4 \
    libxss1 \
    libatk-bridge2.0-0 \
    libcups2 \
    libdrm2 \
    libxkbcommon0 \
    libxi6 \
    libxcomposite1 \
    libxdamage1 \
    libxfixes3 \
    libxrandr2 \
    libxcursor1 \
    libxrender1 \
    libxtst6 \
    libegl1 \
    libgstreamer-plugins-base1.0-0 \
    libgstreamer1.0-0 \
    libharfbuzz-icu0 \
    libjpeg-turbo8 \
    libpng16-16 \
    libwebp6 \
    libxshmfence6 \
    libvulkan1 \
    # フォントも重要 (ブラウザがUIを正しくレンダリングするため)
    fonts-liberation \
    fonts-noto-color-emoji \
    # その他のユーティリティ
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# タイムゾーンを設定します。
ENV TZ=Asia/Tokyo

# アプリケーションの作業ディレクトリを設定します。
WORKDIR /app

# Goアプリケーションのビルド済みバイナリをコピーします。
COPY --from=go_builder /app/server .

# アプリケーションがリッスンするポートを公開 (情報提供のみ)
EXPOSE 8080

# コンテナ起動時に実行されるコマンドを設定します。
CMD ["./server"]
