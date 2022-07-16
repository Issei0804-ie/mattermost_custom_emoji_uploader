# mattermost_custom_emoji_uploader

## 概要

- mattermost で面倒な custom emoji の登録を簡単に実行できるバイナリ
- personal access token が必要

## 使い方

### インストール方法

`go install` か release からバイナリをダウンロードしてください。

```go
go install github.com/Issei0804-ie/mattermost_custom_emoji_uploader@latest
```

### 実行方法

upload したい絵文字を手元に準備してください。

```shell
mattermost_custom_emoji_uploader http://your-mattermost.example [personal access token] [path of dir or image file]
```

下記のように画像ファイルを指定しても構いませんし、

```shell
mattermost_custom_emoji_uploader https://sample.com TOKEN sample.png
```

ディレクトリごと指定しても大丈夫です。

```shell
mattermost_custom_emoji_uploader https://sample.com TOKEN images
```

拡張子は、png, jpeg, jpg, gif に対応しています。それ以外だとエラーがでます。


## issueの報告方法について

- 日本語優先でお願いします。英語でも対応しますが返信が遅れる可能性があります。ご了承ください。
