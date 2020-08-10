package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/pkg/errors"
)

// ReplicateData from some firestore path to another firestore path
func ReplicateData(ctx context.Context, fs *Firestore, readerList map[string]io.Reader) error {
	var err error
	if *isDelete {
		fmt.Printf("delete original document? \n")
		yes := askForConfirmation()
		if !yes {
			return errors.New("exit")
		}
	}

	var destFs *Firestore
	if *hasDestinationFirestore {
		destFs, err = NewFirebase(ctx, *destCrtFile)
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
		om := InterpretationEachValueForTime(m)

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
