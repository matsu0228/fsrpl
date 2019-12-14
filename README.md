# fsrpl

English | [日本語](https://github.com/matsu0228/fsrpl/blob/master/README_ja.md)

fsrpl is CloudFirestore replication tool.

Features:

- Replicate document data from some node to another node. With Wildcar option, Replicate all document data from some collection node to another collenction node.
- Replicate document data from some projectId's Firestore to another projectId's Firestore.
- Backup document data from some node to local JSON file, and Restore  document data from local JSON file.

This tool is Beta version, so if use it BE CAREFUL to [Limitation](#limitation)

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
    - [replicate some documents](#replicate-some-documents)
    - [export data from some documents](#export-data-from-some-documents)
    - [import data from some JSON files](#import-data-from-some-json-files)
    - [generate Go struct from some document](#generate-go-struct-from-some-document)
  - [Limitation](#limitation)
    - [Supporting data-types](#supporting-data-types)
    <!-- /MarkdownTOC -->

## DEMO

replicate `specific one document` and `each documents with wildcard option`

![fsrpl_demo_190829_02](https://user-images.githubusercontent.com/5501329/63935971-a6dfc280-ca99-11e9-8d8c-1e4e93516602.gif)

## SETUP

### homebrew

you can use `homebrew` for macOS

```
# add informal fomura
brew tap matsu0228/homebrew-fsrpl

brew install fsrpl
```

### go get


```
go get github.com/matsu0228/fsrpl
```

### Download


download here to get binary.
https://github.com/matsu0228/fsrpl/releases


### Firestore private key

- you should set firestore's private key(JSON file).
  - you can get private key from console. see [official document](https://firebase.google.com/docs/admin/setup?authuser=0)
- You have two options.
  - set enveronment variable: `FIRESTORE_SECRET`
  - use `--secret` option

## USAGE

### replicate some documents

- replicate some documents with `-d` option.

```
fsrpl [input document path] [output document path]

e.g.

fsrpl "inputData/user" -d "new/user"
fsrpl "inputData/*" -d "outputData/*"
```

### export data from some documents

- export data from some documents with `-f` option.

```
fsrpl [input document path] -f [json file directory path]

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


### import data from some JSON files

- import data from JSON files with `-i` option.


```
fsrpl [import document path] -i [inport JSON file directory path]

e.g.

fsrpl "importData/*" -i "./" --verbose

[INFO] save to importData / dog. document of map[string]interface {}{"_created_at":time.Time{wall:0x0, ext:63702392400, loc:(*time.Location)(nil)}, "coin":0, "favorites":[]interface {}{"1", "2"}, "isDeleted":true, "mapData":map[string]interface {}{"b":true, "name":"mName"}, "name":"dog"}
...
[INFO] import:user of map[string]interface {}{"favorites":[]interface {}{"1", "2"}, "isDeleted":true, "mapData":map[string]interface {}{"b":true, "name":"mName"}, "name":"user", "_created_at":time.Time{wall:0x0, ext:63702392400, loc:(*time.Location)(nil)}, "coin":0}
```

### generate Go struct from some document

- generate Go struct from some document with `-s` option

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

## Limitation

### Supporting data-types


Firestore's data-types list is [here (official document)](https://firebase.google.com/docs/firestore/manage-data/data-types) .
This tool ONLY support below data-types.

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

At this time, the following data types are not supported, so please be careful when using them. It will be supported in the future.
```
- Bytes
- Geographical point
- Reference
```