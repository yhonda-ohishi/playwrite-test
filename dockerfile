# --- ビルドステージ (builder) ---
# Playwright の公式イメージ (Ubuntu Nobleベース) を使用します。
# このイメージには Node.js、Playwright ライブラリ、ブラウザバイナリがプリインストールされています。
# ただし、Go コンパイラは含まれていません。
FROM mcr.microsoft.com/playwright:v1.53.0-noble AS builder

# 作業ディレクトリを設定
WORKDIR /app

# Go コンパイラをインストールする
# Ubuntuベースなので apt-get を使用します。
# Go 1.24.4 に対応するバージョンをインストールします。
# リポジトリによっては最新版が提供されないこともあるので、go.mod のバージョンと合うか確認してください。
# もし特定のバージョンが必要なら、goenv のようなツールを使う手もありますが、ここではシンプルに。
RUN apt-get update && apt-get install -y --no-install-recommends \
    golang-go \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

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
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o server .

# --- 最終イメージステージ ---
# アプリケーションが実行時にもPlaywrightを使うなら、同じイメージをベースにすべきです。
# Playwright は Node.js プロセスやブラウザバイナリに依存するため、
# 実行環境にもそれらが必要です。
FROM mcr.microsoft.com/playwright:v1.53.0-noble

ENV TZ Asia/Tokyo

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドステージからコンパイル済みバイナリをコピー
COPY --from=builder /app/server .

# アプリケーションがリッスンするポートを公開 (情報提供のみ)
EXPOSE 8080

# アプリケーションの実行コマンド
CMD ["./server"]