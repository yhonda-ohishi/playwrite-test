# Playwright Test Docker Image

デジタコのデータをplaywrightでダウンロードするためのDockerイメージです。

## 使い方
### Dockerイメージのビルド
```
docker login ghcr.io -u YOUR_GITHUB_USERNAME
```
上記を実行し、パスワードにはGitHubで生成したトークンを入力してください。

```
pasword: githubで生成したトークン
```

```
docker run -d --net posgres-net --name dtako_server ghcr.io/yhonda-ohishi/playwrite-test:latest
```