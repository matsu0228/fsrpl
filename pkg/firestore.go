package fsrpl

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

// Firestore is datastore object
type Firestore struct {
	firebase        *firebase.App
	FirestoreClient *firestore.Client
}

// NewFirebase is constoractor. connect firebase and firestore
func NewFirebase(ctx context.Context, crtFile string) (*Firestore, error) {

	var app *firebase.App
	var client *firestore.Client
	var err error
	log.Printf("[INFO] connect firestore with: %v", crtFile)

	if crtFile == "" { // local emurator などへの接続時
		app, err = firebase.NewApp(ctx, nil)
		if err != nil {
			return nil, err
		}
		client, err = app.Firestore(ctx)
		if err != nil {
			return nil, err
		}
		return &Firestore{
			firebase:        app,
			FirestoreClient: client,
		}, nil
	}

	_, err = os.Stat(crtFile)
	if err != nil {
		return nil, errors.Wrapf(err, "not found secret file:%v", crtFile)
	}

	opt := option.WithCredentialsFile(crtFile)
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	client, err = app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &Firestore{
		firebase:        app,
		FirestoreClient: client,
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

// getDocumentRefWithPath get documentRef with collection+document path
func (f *Firestore) getDocumentRefWithPath(path string) (*firestore.DocumentRef, error) {
	colID, docID, err := f.parsePath(path)
	if err != nil {
		return nil, err
	}
	doc := f.FirestoreClient.Collection(colID).Doc(docID)
	return doc, nil
}

// DeleteData :
func (f *Firestore) DeleteData(ctx context.Context, path string) error {
	log.Printf("[INFO] DeleteDoc() path:%v", path)
	doc, err := f.getDocumentRefWithPath(path)
	if err != nil {
		return err
	}
	_, err = doc.Delete(ctx)
	return err
}

// SaveData save with collection+document path
func (f *Firestore) SaveData(ctx context.Context, path string, data map[string]interface{}) error {
	log.Printf("[INFO] SaveDoc() path:%v w/ %#v", path, data)
	doc, err := f.getDocumentRefWithPath(path)
	if err != nil {
		return err
	}
	_, err = doc.Set(ctx, data)
	// log.Printf("[INFO] SaveDoc() result %#v.  path:%v ", wr, path)
	return err
}

func (f *Firestore) dataToStream(d map[string]interface{}) (io.Reader, error) {
	s, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(s)
	return reader, nil
}

// Scan is function scan stream data from firestore path
func (f *Firestore) Scan(ctx context.Context, path string) (map[string]io.Reader, error) {

	rs := map[string]io.Reader{}

	colID, docID, err := f.parsePath(path)
	if err != nil {
		return rs, err
	}

	if docID == "*" { // wildcard
		dRefs, err := f.FirestoreClient.Collection(colID).DocumentRefs(ctx).GetAll()
		if err != nil {
			return rs, err
		}
		for _, d := range dRefs {
			snap, err := d.Get(ctx)
			if err != nil {
				return rs, err
			}
			log.Printf("[DEBUG] id:%v, doc:%#v", d.ID, snap.Data())
			r, err := f.dataToStream(snap.Data())
			if err != nil {
				return rs, err
			}
			rs[d.ID] = r
		}

		return rs, err
	}

	snap, err := f.FirestoreClient.Collection(colID).Doc(docID).Get(ctx)
	if err != nil {
		return nil, err
	}

	data := snap.Data()
	r, err := f.dataToStream(data)
	rs[snap.Ref.ID] = r
	return rs, err
}

// ReaderToStruct :
func ReaderToStruct(reader io.Reader, outStream io.Writer) error {
	var parser gojson.Parser
	parser = func(input io.Reader) (interface{}, error) {
		var result interface{}
		if err := json.NewDecoder(input).Decode(&result); err != nil {
			return nil, err
		}
		return result, nil
	}

	name := "JsonStruct"
	pkg := "main"
	tagList := []string{"json"}
	subStruct := false
	// convertFloats := true

	output, err := gojson.Generate(reader, parser, name, pkg, tagList, subStruct) //, convertFloats)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(outStream, string(output))
	return err
}

// ToStruct is converter from firestore data to GoStruct
func (f *Firestore) ToStruct(ctx context.Context, path string, outStream io.Writer) error {
	readerList, err := f.Scan(ctx, path)
	if err != nil {
		return err
	}
	for k, reader := range readerList {
		err = ReaderToStruct(reader, outStream)
		if err != nil {
			log.Printf("[ERROR] %v w/%v", err, k)
		}
	}
	return err
}
