<p align="center"><img width="50%" src="assets/logo.png" /></p>

# fsrpl

![test](https://github.com/matsu0228/fsrpl/workflows/test/badge.svg)

[English](https://github.com/matsu0228/fsrpl/blob/master/README.md) | 日本語


`fsrpl` (CloudFirestore replication tool) は CloudFirestore のデータをコピー・バックアップ・復元することができるCLIツールです.

## Features

- `copy` 特定のドキュメントを、別のCollection配下にコピーできる。また、ワイルドカードを利用することでコレクション配下のすべてのドキュメントをコピーできる。
  - さらに、特定のProjectのFirebaseから他のFirebaseへ、特定のドキュメントのデータをコピーできる
- `dump` 特定のドキュメントをローカルのJSONファイルとしてバックアップができる。
- `restore` ローカルのJSONファイルからドキュメントを復元できる。`firestore emulator` へもデータの復元ができるため、テストデータ作成にも利用できる


Table Of Contents:
<!-- MarkdownTOC -->
- [fsrpl](#fsrpl)
  - [Features](#features)
  - [DEMO](#demo)
  - [SETUP](#setup)
    - [homebrew](#homebrew)
    - [go get](#go-get)
    - [Download](#download)
    - [Firestore private key](#firestore-private-key)
  - [USAGE](#usage)
    - [copy:特定のドキュメントをコピーする](#copy特定のドキュメントをコピーする)
    - [dump:ドキュメントからローカルJSONファイルとしてデータをバックアップする](#dumpドキュメントからローカルjsonファイルとしてデータをバックアップする)
    - [restore:ローカルのJSONファイルからデータを復元する](#restoreローカルのjsonファイルからデータを復元する)
    - [restore:firestore emulatorへデータを復元する](#restorefirestore-emulatorへデータを復元する)
    - [(Go開発者向け機能)ドキュメントのデータからGoの構造体を生成する](#go開発者向け機能ドキュメントのデータからgoの構造体を生成する)

## DEMO

| copy                                                 |
| ---------------------------------------------------- |
| <image width="600" src="assets/copy.gif" alt="copy"> |

| restore                                                    | dump                                                 |
| ---------------------------------------------------------- | ---------------------------------------------------- |
| <image width="400" src="assets/restore.gif" alt="restore"> | <image width="400" src="assets/dump.gif" alt="dump"> |


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
  - 環境変数 `FSRPL_CREDENTIALS` に、秘密鍵のPATHを指定する方法
  - CLIのオプション `--cred` で、秘密鍵のPATHを指定する方法

## USAGE

### copy:特定のドキュメントをコピーする

- 特定のドキュメントをコピーするために、 `copy` コマンドを利用します
  - ワイルドカード`*`の指定ですべてのドキュメントをコピーできます

```
fsrpl copy [コピー元のドキュメントを指定] --dest [コピー先のドキュメントを指定] 

e.g.

fsrpl copy "inputData/user" --dest "new/user"
fsrpl copy "inputData/*" --dest "outputData/*"
```


### dump:ドキュメントからローカルJSONファイルとしてデータをバックアップする

- 特定のドキュメントからデータをバックアップするために、 `dump` コマンドを利用します

```
fsrpl dump [コピー元のドキュメントを指定] --path [バックアップファイルを保存するディレクトリを指定]

e.g.

fsrpl dump "inputData/user" --path ./
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



fsrpl dump "inputData/*" --path ./

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


### restore:ローカルのJSONファイルからデータを復元する

- ローカルのJSONファイルからデータを復元するために、 `restore` コマンドを利用します


```
fsrpl restore [コピー元のドキュメントを指定] --path [復元対象のJSONファイルを指定]

e.g.

fsrpl restore "importData/*" --path "./"

save to importData/ dog. data: map[string]interface {}{"_created_at":time.Time{wall:0x0, ext:63702392400, loc:(*time.Location)(nil)}, "coin":0, "favorites":[]interface {}{"1", "2"}, "isDeleted":true, "mapData":map[string]interface {}{"b":true, "name":"mName"}, "name":"pig"}
...
```

### restore:firestore emulatorへデータを復元する

- firestore emulatorを起動し、`FIRESTORE_EMULATOR_HOST` の環境変数を設定した状態で、 `restore` コマンドを利用すると、emulatorへの復元ができる
  - `--emulators-project-id` オプションにてprojectId指定してrestoreできる。emulatorではprojectIdにてデータを別管理できるため、並列で複数のtestを実行する場合には、testごとに固有のprojectIdを指定することでデータの競合を避けられる。
  - テストコードでの利用例はこちらの [examples](/examples) を参照のこと

```
FIRESTORE_EMULATOR_HOST=**your_firestore_emulator**  fsrpl restore [コピー元のドキュメントを指定] --path [復元対象のJSONファイルを指定] --emulators-project-id [testごとに固有のId]

e.g.

FIRESTORE_EMULATOR_HOST=localhost:8080 fsrpl restore "importData/*" --path "./" --emulators-project-id emulator-integration-test
```


### (Go開発者向け機能)ドキュメントのデータからGoの構造体を生成する

- 特定のドキュメントのデータ形式に対応するGoの構造体を生成するには、`--show-go-struct` オプションを利用します

```
e.g.

fsrpl dump "inputData/user" --show-go-struct

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

