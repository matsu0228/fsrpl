package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

// ImportDataFromJSONFiles import from json
func ImportDataFromJSONFiles(ctx context.Context, fs *Firestore, importPath, exportPath string) error {

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
		documentData := InterpretationEachValueForTime(org)

		log.Printf("[INFO] import:%v of %#v", basefn, documentData)
		return fs.SaveDataWithSubdocumentID(ctx, exportPath, basefn, documentData)
	})
	return err
}

func writeFile(fn string, contents io.Reader) error {
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0666)
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("[WARN] cant close:%v", err)
		}
	}()
	if err != nil {
		return err
	}
	_, err = io.Copy(file, contents)
	return err
}
