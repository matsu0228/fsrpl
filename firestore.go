package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/ChimeraCoder/gojson"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

// Firestore :データストア
type Firestore struct {
	firebase        *firebase.App
	firestoreClient *firestore.Client
}

// NewFirebase connect firebase
func NewFirebase(ctx context.Context, crtFile string) (*Firestore, error) {

	_, err := os.Stat(crtFile)
	if err != nil {
		return nil, errors.Wrapf(err, "not found secret file:%v", crtFile)
	}
	opt := option.WithCredentialsFile(crtFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &Firestore{
		firebase:        app,
		firestoreClient: client,
	}, nil
}

// parsePath is parser of path to collection path and documentID.
// path should containts even number.
func (f *Firestore) parsePath(path string) (string, string, error) {
	sep := "/"
	paths := strings.Split(path, sep)
	if len(paths)%2 != 0 {
		return "", "", fmt.Errorf("path should containts even namber of IDs:%v", path)
	}
	return strings.Join(paths[:len(paths)-1], sep), paths[len(paths)-1], nil
}

// SaveData :
func (f *Firestore) SaveData(ctx context.Context, path string, data map[string]interface{}) error {
	colID, docID, err := f.parsePath(path)
	if err != nil {
		return err
	}
	log.Printf("[INFO] SaveDoc() path:%v w/ %#v", path, data)
	doc := f.firestoreClient.Collection(colID).Doc(docID)
	_, err = doc.Set(ctx, data)
	return err
}

// Scan is function scan stream data from firestore path
func (f *Firestore) Scan(ctx context.Context, path string) (io.Reader, error) {

	colID, docID, err := f.parsePath(path)
	if err != nil {
		return nil, err
	}
	snap, err := f.firestoreClient.Collection(colID).Doc(docID).Get(ctx)
	if err != nil {
		return nil, err
	}

	data := snap.Data()
	s, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(s)
	return reader, nil
}

// ReaderToStruct :
func ReaderToStruct(reader io.Reader, outStream io.Writer) error {
	parser := gojson.ParseJson
	name := "JsonStruct"
	pkg := "main"
	tagList := []string{"json"}
	subStruct := false
	convertFloats := true
	output, err := gojson.Generate(reader, parser, name, pkg, tagList, subStruct, convertFloats)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(outStream, string(output))
	return err
}

// ToStruct is converter from firestore data to GoStruct
func (f *Firestore) ToStruct(ctx context.Context, path string, outStream io.Writer) error {
	reader, err := f.Scan(ctx, path)
	if err != nil {
		return err
	}
	return ReaderToStruct(reader, outStream)
}
