# apps-gcp-gae-app-cloudbuild-cicd-slackbot-migrate-traffic
このリポジトリは「[GAE アプリを Cloud Bulid で CI/CD して Slackbot でトラフィック移行する！](https://apps-gcp.com/gae-app-cloudbuild-cicd-slackbot-migrate-traffic)」の記事内で利用したソースコードを管理しています。

## .env.yaml について
今回 Cloud Functions で利用した環境変数については以下の構成となっております。
```YAML
# Bot User OAuth Access Token
BOT_USER_OAUTH_ACCESS_TOKEN: xoxb-foo
# Verification Token
VERIFICATION_TOKEN: bar
# Channel ID
CHANNEL_ID: hoge
```
