package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// ImportDataFromJSONFiles import from json
func ImportDataFromJSONFiles(ctx context.Context, opt *Option, fs *Firestore, importPath, exportPath string) error {

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
		documentData := fs.InterpretationEachValueForTime(org)

		Debugf("import:%v to %v data: %#v", basefn, exportPath, documentData)
		return fs.SaveDataWithSubdocumentID(ctx, opt, exportPath, basefn, documentData)
	})
	return err
}

func writeFile(fn string, contents io.Reader) error {
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0666)
	defer func() {
		if cErr := file.Close(); cErr != nil {
			Debugf("cant close:%v", cErr)
		}
	}()
	if err != nil {
		return err
	}
	_, err = io.Copy(file, contents)
	return err
}
