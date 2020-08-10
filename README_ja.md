# fsrpl ![test](https://github.com/matsu0228/fsrpl/workflows/test/badge.svg)

[English](https://github.com/matsu0228/fsrpl/blob/master/README.md) | 日本語

`fsrpl` (CloudFirestore replication tool) は CloudFirestore のデータをコピー・バックアップ・復元ができるCLIツールです.

Features:

- 特定のドキュメントを、別のCollection配下にコピーできる。また、ワイルドカードを利用することでコレクション配下のすべてのドキュメントをコピーすることもできる。
- 特定のProjectのFirebaseから他のFirebaseへ、特定のドキュメントのデータをコピーできる
- 特定のドキュメントをローカルのJSONファイルとしてバックアップができる。また、ローカルのJSONファイルからドキュメントを復元できる

本ツールはβ版です。ご利用時には、[制限](#%e5%88%b6%e9%99%90)について注意してください。

Table Of Contents:
<!-- MarkdownTOC -->

- [fsrpl](#fsrpl)
  - [DEMO](#demo)
  - [SETUP](#setup)
    - [homebrew](#homebrew)
    - [go get](#go-get)
    - [Download](#download)
    - [Firestore private key](#firestore-private-key)
  - [USAGE](#usage)
    - [特定のドキュメントをコピーする](#%e7%89%b9%e5%ae%9a%e3%81%ae%e3%83%89%e3%82%ad%e3%83%a5%e3%83%a1%e3%83%b3%e3%83%88%e3%82%92%e3%82%b3%e3%83%94%e3%83%bc%e3%81%99%e3%82%8b)
    - [ドキュメントからローカルJSONファイルとしてデータをバックアップする](#%e3%83%89%e3%82%ad%e3%83%a5%e3%83%a1%e3%83%b3%e3%83%88%e3%81%8b%e3%82%89%e3%83%ad%e3%83%bc%e3%82%ab%e3%83%abjson%e3%83%95%e3%82%a1%e3%82%a4%e3%83%ab%e3%81%a8%e3%81%97%e3%81%a6%e3%83%87%e3%83%bc%e3%82%bf%e3%82%92%e3%83%90%e3%83%83%e3%82%af%e3%82%a2%e3%83%83%e3%83%97%e3%81%99%e3%82%8b)
    - [ローカルのJSONファイルかでデータをインポートする](#%e3%83%ad%e3%83%bc%e3%82%ab%e3%83%ab%e3%81%aejson%e3%83%95%e3%82%a1%e3%82%a4%e3%83%ab%e3%81%8b%e3%81%a7%e3%83%87%e3%83%bc%e3%82%bf%e3%82%92%e3%82%a4%e3%83%b3%e3%83%9d%e3%83%bc%e3%83%88%e3%81%99%e3%82%8b)
    - [(Go開発者向け機能)ドキュメントのデータからGoの構造体を生成する](#go%e9%96%8b%e7%99%ba%e8%80%85%e5%90%91%e3%81%91%e6%a9%9f%e8%83%bd%e3%83%89%e3%82%ad%e3%83%a5%e3%83%a1%e3%83%b3%e3%83%88%e3%81%ae%e3%83%87%e3%83%bc%e3%82%bf%e3%81%8b%e3%82%89go%e3%81%ae%e6%a7%8b%e9%80%a0%e4%bd%93%e3%82%92%e7%94%9f%e6%88%90%e3%81%99%e3%82%8b)
  - [制限](#%e5%88%b6%e9%99%90)
    - [対応しているFirestoreの型](#%e5%af%be%e5%bf%9c%e3%81%97%e3%81%a6%e3%81%84%e3%82%8bfirestore%e3%81%ae%e5%9e%8b)
    <!-- /MarkdownTOC -->

## DEMO

特定のドキュメントのデータをコピーするデモです。「ワイルドカード(*)」を使うことで複数のドキュメントを一括コピーすることができます。

![fsrpl_demo_190829_02](https://user-images.githubusercontent.com/5501329/63935971-a6dfc280-ca99-11e9-8d8c-1e4e93516602.gif)

## SETUP

### homebrew

`homebrew` を使ってインストールできます

```
# add informal fomura
brew tap matsu0228/homebrew-fsrpl

brew install fsrpl
```

### go get

go が利用できる環境であれば、`go get`でのインストールができます

```
go get github.com/matsu0228/fsrpl
```

### Download


バイナリファイルから利用する場合は、こちらからダウンロードできます。Linux,Windowsでも利用できます。
https://github.com/matsu0228/fsrpl/releases


### Firestore private key

- 本ツールを利用するためには、Firestoreの秘密鍵(JSON file)が必要です
  - 秘密鍵は、Firebase Consoleから取得できます。詳しくは、 [official document](https://firebase.google.com/docs/admin/setup?authuser=0)から。
- 下記いずれかの方法で、秘密鍵を指定できます
  - 環境変数 `FIRESTORE_SECRET` に、秘密鍵のPATHを指定する方法
  - CLIのオプション `--secret` で、秘密鍵のPATHを指定する方法

## USAGE

### 特定のドキュメントをコピーする

- 特定のドキュメントをコピーするために、 `-d` オプションを利用します
  - ワイルドカード`*`の指定ですべてのドキュメントをコピーできます

```
fsrpl [コピー元のドキュメントを指定] -d [コピー先のドキュメントを指定]

e.g.

fsrpl "inputData/user" -d "new/user"
fsrpl "inputData/*" -d "outputData/*"
```

### ドキュメントからローカルJSONファイルとしてデータをバックアップする

- 特定のドキュメントからデータをバックアップするために、`-f` オプションを利用します

```
fsrpl [コピー元のドキュメントを指定] -f [バックアップファイルを保存するディレクトリを指定]

e.g.

fsrpl "inputData/user" -f ./
cat user.json
{
  "_created_at": "2019-08-26T05:00:00Z",
  "coin": 0,
  "favorites": [
    "1",
    "2"
  ],
  "isDeleted": true,
  "mapData": {
    "isMan": true,
    "name": "subName"
  },
  "name": "user"
}



fsrpl "inputData/*" -f ./

cat cat.json | jq
{
  "_created_at": "2019-08-26T05:00:00Z",
  "coin": 0,
  "favorites": [
    "1",
    "2"
  ],
  "isDeleted": true,
  "mapData": {
    "isBrownhair": true,
    "name": "calico"
  },
  "name": "cat"
}

cat dog.json | jq
{
  "name": "dog"
  ...
}
...

```


### ローカルのJSONファイルかでデータをインポートする

- ローカルのJSONファイルからデータをインポートするために、 `-i` オプションを利用します


```
fsrpl [コピー元のドキュメントを指定] -i [インポート対象のJSONファイルを指定]

e.g.

fsrpl "importData/*" -i "./" --verbose

[INFO] save to importData / dog. document of map[string]interface {}{"_created_at":time.Time{wall:0x0, ext:63702392400, loc:(*time.Location)(nil)}, "coin":0, "favorites":[]interface {}{"1", "2"}, "isDeleted":true, "mapData":map[string]interface {}{"b":true, "name":"mName"}, "name":"dog"}
...
[INFO] import:user of map[string]interface {}{"favorites":[]interface {}{"1", "2"}, "isDeleted":true, "mapData":map[string]interface {}{"b":true, "name":"mName"}, "name":"user", "_created_at":time.Time{wall:0x0, ext:63702392400, loc:(*time.Location)(nil)}, "coin":0}
```

### (Go開発者向け機能)ドキュメントのデータからGoの構造体を生成する

- 特定のドキュメントのデータ形式に対応するGoの構造体を生成するには、`-s` オプションを利用します

```
e.g.

fsrpl -p "inputData/user" -s

package main

type JsonStruct struct {
	CreatedAt string   `json:"_created_at"`
	Coin      int64    `json:"coin"`
	Favorites []string `json:"favorites"`
	IsDeleted bool     `json:"isDeleted"`
	MapData   struct {
		B    bool   `json:"b"`
		Name string `json:"name"`
	} `json:"mapData"`
	Name string `json:"name"`
}

```

## 制限

### 対応しているFirestoreの型


[Firestoreの型](https://firebase.google.com/docs/firestore/manage-data/data-types) のうち、本ツールで対応している型は下記です


```
- Boolean
- Text string
- Integer
- Array
- Map
- Date and time
- Null
- Floating-point number
```

現時点で、下記のデータ型には対応できていないので、ご利用時には注意ください。今後対応予定です。
```
- Bytes
- Geographical point
- Reference
```