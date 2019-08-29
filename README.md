# fsrpl

fsrpl is CloudFirestore replication tool.

Features:

- Replicate CloudFirestore's some node data to another node.
- Make Go struct from CloudFirestore's some node data.

<!-- MarkdownTOC -->

- [fsrpl](#fsrpl)
  - [DEMO](#demo)
  - [SETUP](#setup)
  - [USAGE](#usage)
    - [replicate some documents](#replicate-some-documents)
    - [export data from some documents](#export-data-from-some-documents)
    - [generate Go struct from some document](#generate-go-struct-from-some-document)
    <!-- /MarkdownTOC -->

## DEMO

replicate `specific one document` and `each documents with wildcard option`

![fsrpl_demo_190829_02](https://user-images.githubusercontent.com/5501329/63935971-a6dfc280-ca99-11e9-8d8c-1e4e93516602.gif)

## SETUP

write later..

- you should set firestore's secret(JSON file). You have two options.
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
fsrpl [input document path] [output document path]

e.g.

fsrpl "inputData/user" -f | jq
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



fsrpl "inputData/*" -f | jq


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
{
  "name": "dog"
  ...
}
...

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
