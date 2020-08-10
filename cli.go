package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	ExitCodeOK    int = 0
	ExitCodeError     = 10 + iota
	ExitCodeBadArgs
)

const EnvDebug = "FSRPL_DEBUG"

var (
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

func Debugf(format string, args ...interface{}) {
	if env := os.Getenv(EnvDebug); len(env) != 0 {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func PrintAlertf(w io.Writer, format string, args ...interface{}) {
	fmt.Fprint(w, fmt.Sprintf(format, args...))
}

type CLI struct {
	outStream, errStream io.Writer
}

func (cli *CLI) Run(args []string) int {

	app.Version(version)

	if _, err := app.Parse(args[1:]); err != nil {
		app.FatalUsage(fmt.Sprintf("\n%s\n-------------\n", err.Error()))
	}

	if *isDebug || *isVerbose {
		os.Setenv(EnvDebug, "1")
		Debugf("set DEBUG mode")
	}

	mode, err := validate()
	if err != nil {
		PrintAlertf(cli.errStream, "ERROR: %v\n", err)
		return ExitCodeBadArgs
	}
	if mode == Unknown {
		PrintAlertf(cli.errStream, "ERROR: %v\n", "please set execute mode (-d or -f or -g)")
		return ExitCodeBadArgs
	}

	if err := cli.RunWithMode(mode); err != nil {
		PrintAlertf(cli.errStream, "ERROR: %v\n", err)
		return ExitCodeError
	}

	return ExitCodeOK
}

func (cli *CLI) RunWithMode(mode ExecMode) error {

	var readerList map[string]io.Reader
	ctx := context.Background()
	fs, err := NewFirebase(ctx, *crtFile)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] projectID: %v", fs.ProjectID)

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
		log.Println("generate struct")
		return fs.ToStruct(ctx, *inputPath, cli.outStream)

	case ExportMode:
		errs := []error{}
		for k, reader := range readerList {
			log.Printf("[DEBUG] write outStream with : %v ----------------\n", k)

			fn := path.Join(*exportFilePath, k+".json")
			if err := writeFile(fn, reader); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) >= 1 {
			return fmt.Errorf("err: %#v", errs)
		}
		return nil

	case ReplicateMode:
		if err = ReplicateData(ctx, fs, readerList); err != nil {
			return err
		}
	default:
		log.Printf("[INFO] dont set execute mode: %v from  \n", mode)
	}
	return nil
}
