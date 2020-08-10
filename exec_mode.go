package main

import (
	"os"

	"github.com/pkg/errors"
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
