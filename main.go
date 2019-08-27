package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var crtFile, firestorePath, jsonFilePath, outputPath string

func init() {
	flag.StringVar(&crtFile, "secret", "", "set secrets json for firestore")
	flag.StringVar(&firestorePath, "p", "sports/rugby_topleague/games/4ayoQUGbmNAC1l9lITjB", "set firetstore path(containts collection's path and documentID")
	flag.StringVar(&outputPath, "outputPath", "", "set output firetstore path(containts collection's path and documentID")
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

	log.Printf("[INFO] connect firestore with: %v", crtFile)

	ctx := context.Background()
	repo, err := NewFirebase(ctx, crtFile)
	errorCheck(err)
	log.Printf("[DEBUG] repo %#v", repo)

	outStream := os.Stdout
	// err = repo.ToStruct(ctx, firestorePath, outStream)
	// errorCheck(err)

	reader, err := repo.Scan(ctx, firestorePath)
	errorCheck(err)
	// _, err = io.Copy(outStream, reader)
	// errorCheck(err)

	var m map[string]interface{}
	err = json.NewDecoder(reader).Decode(&m)
	errorCheck(err)

	om := InterpretationEachValueForTime(m)

	if outputPath != "" {
		err = repo.SaveData(ctx, outputPath, m)
		errorCheck(err)
		return
	}
	// _, err = fmt.Fprint(outStream, fmt.Sprintf("%#v", m))
	_, err = fmt.Fprint(outStream, fmt.Sprintf("%#v", om))
	errorCheck(err)

}
