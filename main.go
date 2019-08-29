package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	fsrpl "github.com/matsu0228/fsrpl/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// version
	version  = "0.0.1"
	revision = "0"

	app       = kingpin.New("fsrpl", "A firestore replication tool.")
	crtFile   = app.Flag("secret", "set secrets json for firestore").Default("").String()
	inputPath = app.Arg("targetPath", "target firestore path(containts collection's path and documentID)").Default("").String()

	// output
	outputPath       = app.Flag("destPath", "destination firestore path(containts collection's path and documentID)").Short('d').Default("").String()
	isExportFile     = app.Flag("isExportFile", "output json data to file").Short('f').Default("false").Bool()
	isOutputGoStruct = app.Flag("isGoStruct", "output go struct to stdout").Short('g').Default("false").Bool()
)

func errorCheck(err error) {
	if err != nil {
		errorExit(err)
	}
}

func errorExit(err error) {
	log.Fatalf("[ERROR] %v", err)
}

func validate() error {
	if *inputPath == "" {
		return errors.New("set first option")
	}

	if *crtFile == "" {
		*crtFile = os.Getenv("FIRESTORE_SECRET")
		if *crtFile == "" {
			return errors.New("set secret file path: --secret **")
		}
	}
	return nil
}

func run() error {
	ctx := context.Background()
	repo, err := fsrpl.NewFirebase(ctx, *crtFile)
	if err != nil {
		return err
	}

	outStream := os.Stdout
	if *isOutputGoStruct {
		err = repo.ToStruct(ctx, *inputPath, outStream)
		return err
	}

	readerList, err := repo.Scan(ctx, *inputPath)
	if err != nil {
		return err
	}

	for k, reader := range readerList {

		if *isExportFile {
			log.Printf("[INFO] write with : %v ----------------\n", k)
			_, err = io.Copy(outStream, reader)
			if err != nil {
				return err
			}
			continue
		}

		if *outputPath != "" {
			path := strings.Replace(*outputPath, "*", k, -1)
			log.Printf("[INFO] save with : %v ---------------- \n", path)

			var m map[string]interface{}
			err = json.NewDecoder(reader).Decode(&m)
			if err != nil {
				return err
			}
			om := fsrpl.InterpretationEachValueForTime(m)

			err = repo.SaveData(ctx, path, om)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {

	app.Version(fmt.Sprintf("%s\nRev:%s", version, revision))
	if _, err := app.Parse(os.Args[1:]); err != nil {
		app.FatalUsage(fmt.Sprintf("\n%s\n-------------\n", err.Error()))
	}

	// validate
	if err := validate(); err != nil {
		app.FatalUsage(fmt.Sprintf("\n%s\n-------------\n", err.Error()))
	}

	log.Printf("[INFO] connect firestore with: %v", crtFile)
	if err := run(); err != nil {
		errorExit(err)
	}

	log.Print("[INFO] success!")
}
