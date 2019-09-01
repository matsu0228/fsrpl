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

	app         = kingpin.New("fsrpl", "A firestore replication tool.")
	crtFile     = app.Flag("secret", "set secrets json for firestore").Default("").String()
	destCrtFile = app.Flag("dest-secret", "set secrets json for destination firestore").Default("").String()
	inputPath   = app.Arg("targetPath", "target firestore path(containts collection's path and documentID)").Default("").String()

	// output
	outputPath       = app.Flag("destPath", "destination firestore path(containts collection's path and documentID)").Short('d').Default("").String()
	isExportFile     = app.Flag("isExportFile", "output json data to file").Short('f').Default("false").Bool()
	isOutputGoStruct = app.Flag("isGoStruct", "output go struct to stdout").Short('g').Default("false").Bool()

	// option
	isDelete                = app.Flag("delete", "delete source document after replication").Default("false").Bool()
	hasDestinationFirestore = app.Flag("destfs", "replicate for another destination Firestore").Default("false").Bool()
)

// ExecMode is execute mode
type ExecMode int

const (
	// Unknown is Undefined mode
	Unknown ExecMode = iota
	// ReplicateMode is mode for replicate
	ReplicateMode
	// ExportMode is mode
	ExportMode
	// GenerateStructMode is mode
	GenerateStructMode
)

func errorCheck(err error) {
	if err != nil {
		errorExit(err)
	}
}

func errorExit(err error) {
	log.Fatalf("[ERROR] %v", err)
}

func validate() (ExecMode, error) {
	if *inputPath == "" {
		return Unknown, errors.New("set first option")
	}

	if *crtFile == "" {
		*crtFile = os.Getenv("FIRESTORE_SECRET")
		if *crtFile == "" {
			return Unknown, errors.New("set secret file path: --secret **")
		}
	}

	if *hasDestinationFirestore {
		if *destCrtFile == "" {
			*destCrtFile = os.Getenv("DESTINATION_FIRESTORE_SECRET")
			if *destCrtFile == "" {
				return Unknown, errors.New("set secret file path: --dest-secret **")
			}
		}
	}

	if *outputPath != "" {
		return ReplicateMode, nil
	}
	if *isExportFile {
		return ExportMode, nil
	}
	// 3. generate Go struct from some document
	if *isOutputGoStruct {
		return GenerateStructMode, nil
	}
	// 3. generate Go struct from some document
	if *isOutputGoStruct {
		return GenerateStructMode, nil
	}

	return Unknown, nil
}

func run(mode ExecMode) error {

	var readerList map[string]io.Reader
	ctx := context.Background()
	fs, err := fsrpl.NewFirebase(ctx, *crtFile)
	if err != nil {
		return err
	}

	outStream := os.Stdout

	if mode == ExportMode || mode == ReplicateMode {
		readerList, err = fs.Scan(ctx, *inputPath)
		if err != nil {
			return err
		}
	}

	switch mode {
	case GenerateStructMode:
		err = fs.ToStruct(ctx, *inputPath, outStream)
		return err

	case ExportMode:
		for k, reader := range readerList {
			log.Printf("[INFO] write with : %v ----------------\n", k)
			_, err = io.Copy(outStream, reader)
			if err != nil {
				return err
			}
		}

	case ReplicateMode:

		var destFs *fsrpl.Firestore
		if *hasDestinationFirestore {
			destFs, err = fsrpl.NewFirebase(ctx, *destCrtFile)
			if err != nil {
				return err
			}
		}

		for k, reader := range readerList {
			path := strings.Replace(*outputPath, "*", k, -1)
			srcPath := strings.Replace(*inputPath, "*", k, -1)
			log.Printf("[INFO] save with : %v from %v ---------------- \n", path, srcPath)

			var m map[string]interface{}
			err = json.NewDecoder(reader).Decode(&m)
			if err != nil {
				return err
			}
			om := fsrpl.InterpretationEachValueForTime(m)

			if *hasDestinationFirestore {
				err = destFs.SaveData(ctx, path, om)
				if err != nil {
					return err
				}
			} else {
				err = fs.SaveData(ctx, path, om)
				if err != nil {
					return err
				}
			}

			if *isDelete {
				err = fs.DeleteData(ctx, srcPath)
				if err != nil {
					return err
				}
			}
		}
	default:
		log.Printf("[INFO] dont set execute mode: %v from  \n", mode)
	}
	return nil
}

func main() {

	app.Version(fmt.Sprintf("%s\nRev:%s", version, revision))
	if _, err := app.Parse(os.Args[1:]); err != nil {
		app.FatalUsage(fmt.Sprintf("\n%s\n-------------\n", err.Error()))
	}

	// validate
	mode, err := validate()
	if err != nil {
		app.FatalUsage(fmt.Sprintf("\n%s\n-------------\n", err.Error()))
	}
	if mode == Unknown {
		app.FatalUsage(fmt.Sprintf("\n%s\n-------------\n", errors.New("please set execute mode (-d or -f or -g)")))
	}

	if err := run(mode); err != nil {
		errorExit(err)
	}

	log.Print("[INFO] success!")
}
