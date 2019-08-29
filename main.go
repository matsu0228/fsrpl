package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	fsrpl "github.com/matsu0228/fsrpl/pkg"
)

var crtFile, firestorePath, jsonFilePath, outputPath string
var isOutputFile, isOutputGoStruct bool

func init() {
	flag.StringVar(&crtFile, "secret", "", "set secrets json for firestore")
	flag.StringVar(&firestorePath, "p", "", "set firetstore path(containts collection's path and documentID")
	// output
	flag.StringVar(&outputPath, "d", "", "destination firestore path(containts collection's path and documentID")
	flag.BoolVar(&isOutputFile, "f", false, "output json data to file")
	flag.BoolVar(&isOutputGoStruct, "s", false, "output go struct to stdout")
}

func errorCheck(err error) {
	if err != nil {
		errorExit(err)
	}
}

func errorExit(err error) {
	log.Fatalf("[ERROR] %v", err)
}

func main() {

	flag.Parse()

	// validate
	if crtFile == "" {
		crtFile = os.Getenv("FIRESTORE_SECRET")
	}
	if crtFile == "" {
		errorCheck(fmt.Errorf("set secret file path: --secret ** %v", ""))
	}

	// validate of options

	log.Printf("[INFO] connect firestore with: %v", crtFile)

	ctx := context.Background()
	repo, err := fsrpl.NewFirebase(ctx, crtFile)
	errorCheck(err)
	log.Printf("[DEBUG] repo %#v", repo)

	outStream := os.Stdout

	if isOutputFile {

	}

	if isOutputGoStruct {
		err = repo.ToStruct(ctx, firestorePath, outStream)
		errorCheck(err)
	}

	readerList, err := repo.Scan(ctx, firestorePath)
	errorCheck(err)

	for k, reader := range readerList {

		if isOutputFile {
			log.Printf("[INFO] write with : %v ----------------\n", k)
			_, err = io.Copy(outStream, reader)
			errorCheck(err)
			continue
		}

		if outputPath != "" {
			path := strings.Replace(outputPath, "*", k, -1)
			log.Printf("[INFO] save with : %v ---------------- \n", path)

			var m map[string]interface{}
			err = json.NewDecoder(reader).Decode(&m)
			errorCheck(err)
			om := fsrpl.InterpretationEachValueForTime(m)

			err = repo.SaveData(ctx, path, om)
			errorCheck(err)
		}
		// _, err = fmt.Fprint(outStream, fmt.Sprintf("%#v", m))
		// errorCheck(err)
	}
}
