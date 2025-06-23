

# --- ビルドステージ (builder) ---
# このステージでGoアプリケーションをビルドします。
# 開発ツールやソースコードはここに存在しますが、最終イメージには含まれません。
# Goのコンパイラと開発環境を含むベースイメージ
FROM golang:1.22-alpine AS builder 

# 作業ディレクトリを設定
WORKDIR /app 
ENV GOPRIVATE  github.com/yhonda-ohishi/playwrite-test/*


# gitコマンドをインストールする
# apk add はAlpine Linuxでのパッケージインストールコマンドです。
RUN apk add --no-cache git


ARG GITHUB_TOKEN
RUN if [ -n "$GITHUB_TOKEN" ]; then \
    git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com".insteadOf "https://github.com"; \
    fi
# Goモジュールファイルをコピー
# go.mod と go.sum だけをコピーし、依存関係をダウンロードします。
# これにより、Goモジュールの変更がなければこのレイヤーがキャッシュされ、ビルドが速くなります。
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
# ここでアプリケーションのソースコードをコピーします。
# このソースコードは「builder」ステージ内でのみ利用されます。
COPY . .

# アプリケーションをビルド
# CGO_ENABLED=0: Cのコードに依存しない静的リンクされたバイナリを作成
# GOOS=linux: Linux環境向けにビルド
# -a: 完全に静的なバイナリを作成（システムのCライブラリに依存しない）
# -installsuffix nocgo: CGoを使用しない場合の標準的なサフィックス
# -o server: 実行ファイル名を 'server' に指定
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o server .

# --- 最終イメージステージ ---
# このステージはGoのコンパイラやソースコードを含まず、
# ビルド済みの実行ファイルと、アプリケーションの実行に必要な最小限の環境のみを含みます。
# 非常に軽量なベースイメージ（約5MB）
FROM alpine:latest 

# タイムゾーンデータをインストール (ログなどでタイムゾーンが正しく表示されるように、必要であれば)
RUN apk --no-cache add ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドステージからコンパイル済みバイナリをコピー
# ここが重要です。`COPY --from=builder /app/server .` は、
# 上の「builder」ステージで作成された `server` という実行ファイルだけを
# この最終ステージの `/root/` ディレクトリにコピーします。
COPY --from=builder /app/server .

# アプリケーションがリッスンするポートを公開 (情報提供のみ)
EXPOSE 8080

# アプリケーションの実行コマンド
CMD ["./server"]