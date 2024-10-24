# file-packer

## 用途

ストックイラストサイトに素材をアップロードする際、サイトごとにアップロードできるファイルフォーマットが異なっている。
必要な画像形式は作家が作成した前提で、各サービスの仕様に合わせてファイルをzip化する部分を自動化する。

## 制限

メッセージボックスの表示部分でWindowsに依存したコードがあるため、Windowsでしか動作しない。動作確認もWindowsでしか行っていない。

## 利用方法

ビルド済みバイナリは配布していませんので、手元でビルドする必要があります。

### 【すでに入っている場合は不要】Golang および Git のインストール

公式サイトの手順に従い Go, Git をインストールしてください。

- [Golang](https://go.dev/doc/install)
- [Git for Windows](https://gitforwindows.org/)


### Clone と ビルド

以下のコマンドを実行すると`file-packer.exe`が作られます。

```cmd
git clone https://github.com/uoya/file-packer.git
cd file-packer
go build .
```

### 作業用フォルダの設定

作業用フォルダは既定では`file-packer.exe`と同階層にある`work`フォルダです。これは`config.json`の`workDir`プロパティで指定されています。別名のフォルダを指定する場合は`config.json`を書き換えてください。

### 実行

`file-packer.exe`と同階層に`work`フォルダと`config.json`が存在する状況で`file-packer.exe`を実行すると、即座に`work`フォルダ内の画像に対して処理を開始します。

### config.json のリファレンス

[config.goのConfig構造体](./config.go)を参照のこと