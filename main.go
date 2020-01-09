package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/logutils"
	fsrpl "github.com/matsu0228/fsrpl/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// version
	version = "0.0.1"

	app         = kingpin.New("fsrpl", "A firestore replication tool.")
	isVerbose   = app.Flag("verbose", "show logs with verbose level").Default("false").Bool()
	isDebug     = app.Flag("debug", "show logs with debug level").Default("false").Bool()
	crtFile     = app.Flag("secret", "set secrets json for firestore").Default("").String()
	destCrtFile = app.Flag("dest-secret", "set secrets json for destination firestore").Default("").String()
	inputPath   = app.Arg("targetPath", "target firestore path(containts collection's path and documentID)").Default("").String()

	// import
	importFilePath = app.Flag("importFilePath", "input directory of json files").Short('i').Default("").String()
	// output
	outputPath       = app.Flag("destPath", "destination firestore path(containts collection's path and documentID)").Short('d').Default("").String()
	exportFilePath   = app.Flag("exportFilePath", "output directory as json file").Short('f').Default("").String()
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
	// ReplicateMode is mode for replicating from some firestore to another firestore
	ReplicateMode
	// ExportMode is mode of export from firestore to json file
	ExportMode
	// ImportMode is mode of importing from json files to firestore
	ImportMode
	// GenerateStructMode is mode generating Go Struct
	GenerateStructMode
)

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

	if *importFilePath != "" {
		return ImportMode, nil
	}
	if *outputPath != "" {
		return ReplicateMode, nil
	}
	if *exportFilePath != "" {
		return ExportMode, nil
	}
	if *isOutputGoStruct {
		return GenerateStructMode, nil
	}

	return Unknown, nil
}

// ImportDataFromJSONFiles import from json
func ImportDataFromJSONFiles(ctx context.Context, fs *fsrpl.Firestore, importPath, exportPath string) error {

	err := filepath.Walk(importPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fn := info.Name()
		basefn := filepath.Base(fn[:len(fn)-len(filepath.Ext(fn))])

		var org map[string]interface{}
		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		if err := json.NewDecoder(file).Decode(&org); err != nil {
			return err
		}
		documentData := fsrpl.InterpretationEachValueForTime(org)

		log.Printf("[INFO] import:%v of %#v", basefn, documentData)
		return fs.SaveDataWithSubdocumentID(ctx, exportPath, basefn, documentData)
	})
	return err
}

// ReplicateData from some firestore path to another firestore path
func ReplicateData(ctx context.Context, fs *fsrpl.Firestore, readerList map[string]io.Reader) error {
	var err error
	if *isDelete {
		fmt.Printf("delete original document? \n")
		yes := askForConfirmation()
		if !yes {
			return errors.New("exit")
		}
	}

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
		log.Printf("[DEBUG] save with : %v from %v ---------------- \n", path, srcPath)

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
	return nil
}

func run(mode ExecMode) error {

	var readerList map[string]io.Reader
	ctx := context.Background()
	fs, err := fsrpl.NewFirebase(ctx, *crtFile)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] projectID: %v", fs.ProjectID)

	outStream := os.Stdout

	if mode == ExportMode || mode == ReplicateMode {
		readerList, err = fs.Scan(ctx, *inputPath)
		if err != nil {
			return err
		}
	}

	switch mode {

	case ImportMode:
		if err := ImportDataFromJSONFiles(ctx, fs, *importFilePath, *inputPath); err != nil {
			return err
		}

	case GenerateStructMode:
		err = fs.ToStruct(ctx, *inputPath, outStream)
		return err

	case ExportMode:
		for k, reader := range readerList {
			log.Printf("[DEBUG] write outStream with : %v ----------------\n", k)

			fn := path.Join(*exportFilePath, k+".json")
			file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				return err
			}
			_, err = io.Copy(file, reader)
			if err != nil {
				return err
			}
		}

	case ReplicateMode:
		if err = ReplicateData(ctx, fs, readerList); err != nil {
			return err
		}
	default:
		log.Printf("[INFO] dont set execute mode: %v from  \n", mode)
	}
	return nil
}

func main() {

	app.Version(version)
	if _, err := app.Parse(os.Args[1:]); err != nil {
		app.FatalUsage(fmt.Sprintf("\n%s\n-------------\n", err.Error()))
	}

	// set logger
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("ERROR"),
		Writer:   os.Stderr,
	}
	log.SetFlags(log.Ltime)
	if *isVerbose {
		filter = &logutils.LevelFilter{
			Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
			MinLevel: logutils.LogLevel("INFO"),
			Writer:   os.Stderr,
		}
		log.SetFlags(0)
	}
	if *isDebug {
		filter = &logutils.LevelFilter{
			Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
			MinLevel: logutils.LogLevel("DEBUG"),
			Writer:   os.Stderr,
		}
		log.SetFlags(log.Ltime | log.Lshortfile)
	}
	log.SetOutput(filter)

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
