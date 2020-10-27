package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/ChimeraCoder/gojson"
	"github.com/antonholmquist/jason"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

// Firestore is datastore object
type Firestore struct {
	Client    *firestore.Client
	ProjectID string
}

// FirestoreConnectionOption  specifies a connection method such as a secret key or emulator.
type FirestoreConnectionOption struct {
	CredentialFilePath string
	EmulatorProjectID  string
}

// OptWithEmulatorProjectID generate firestore connection option
func OptWithEmulatorProjectID(projectID string) FirestoreConnectionOption {
	return FirestoreConnectionOption{
		EmulatorProjectID: projectID,
	}
}

// OptWithCred generate firestore connection option
func OptWithCred(cred string) FirestoreConnectionOption {
	return FirestoreConnectionOption{
		CredentialFilePath: cred,
	}
}

func loadProjectIDFromCredFile(credFilePath string) (string, error) {
	file, err := os.Open(credFilePath)
	if err != nil {
		return "", err
	}
	v, err := jason.NewObjectFromReader(file)
	if err != nil {
		return "", err
	}
	return v.GetString("project_id")
}

// NewFirebase is constoractor. connect firebase and firestore
func NewFirebase(ctx context.Context, cliOpt *Option, conOpt FirestoreConnectionOption) (*Firestore, error) {

	var err error
	var fs *firestore.Client
	var projectID string
	cred := conOpt.CredentialFilePath

	envCred := os.Getenv(EnvCredentials)
	envEmu := os.Getenv(EnvEmulatorHost)
	Debugf("connect firestore with: %v env-cred:%s, env-emulator%s", cred, envCred, envEmu)

	if len(envCred) != 0 {
		cred = envCred
		Debugf("set cred: %s", cred)
	}

	if len(envEmu) != 0 {
		projectID = conOpt.EmulatorProjectID
		if len(projectID) == 0 {
			projectID = "emulator"
		}
		fs, err = firestore.NewClient(ctx, projectID)
		if err != nil {
			return nil, err
		}
		PrintInfof(cliOpt.Stdout, "\nconnected emulator (projectID: %v) \n\n", highlight(projectID))
		return &Firestore{
			Client:    fs,
			ProjectID: projectID,
		}, nil
	}

	_, err = os.Stat(cred)
	if err != nil {
		return nil, errors.Wrapf(err, "not found secret file:%v", cred)
	}

	if pjID, loadErr := loadProjectIDFromCredFile(cred); loadErr == nil {
		projectID = pjID
	}
	opt := option.WithCredentialsFile(cred)
	fs, err = firestore.NewClient(ctx, projectID, opt)
	if err != nil {
		return nil, err
	}

	PrintInfof(cliOpt.Stdout, "\nconnected firestore (projectID: %v) \n\n", highlight(projectID))
	return &Firestore{
		Client:    fs,
		ProjectID: projectID,
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
	doc := f.Client.Collection(colID).Doc(docID)
	return doc, nil
}

// DeleteData :
func (f *Firestore) DeleteData(ctx context.Context, opt *Option, path string) error {
	PrintInfof(opt.Stdout, "delete document from %v \n", path)
	doc, err := f.getDocumentRefWithPath(path)
	if err != nil {
		return err
	}
	_, err = doc.Delete(ctx)
	return err
}

// SaveDataWithSubdocumentID save with collection+document path and subDocumentID
func (f *Firestore) SaveDataWithSubdocumentID(ctx context.Context, opt *Option, path, subDocID string, data map[string]interface{}) error {
	colID, docID, err := f.parsePath(path)
	if err != nil {
		return err
	}
	if docID == "*" {
		docID = subDocID
	}
	return f.ImportData(ctx, opt, colID, docID, data)
}

// SaveData save with collection+document path
func (f *Firestore) SaveData(ctx context.Context, opt *Option, path string, data map[string]interface{}) error {
	colID, docID, err := f.parsePath(path)
	if err != nil {
		return err
	}
	return f.ImportData(ctx, opt, colID, docID, data)
}

// ImportData setdata
func (f *Firestore) ImportData(ctx context.Context, opt *Option, colID, docID string, data map[string]interface{}) error {
	PrintInfof(opt.Stdout, "save to %s/ %s (doc:%v) \n\n", colID, docID, data)
	doc := f.Client.Collection(colID).Doc(docID)
	_, err := doc.Set(ctx, data)
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
	Debugf("Scan() col:%v, doc:%v", colID, docID)

	if docID == "*" {
		return f.ScanAll(ctx, colID, docID)
	}

	snap, err := f.Client.Collection(colID).Doc(docID).Get(ctx)
	if err != nil {
		Debugf("Get err :%v", err)
		return nil, err
	}

	data := snap.Data()
	r, err := f.dataToStream(data)
	rs[snap.Ref.ID] = r
	return rs, err
}

// ScanAll scan stream data with all document
func (f *Firestore) ScanAll(ctx context.Context, colID, docID string) (map[string]io.Reader, error) {
	rs := map[string]io.Reader{}

	dRefs, err := f.Client.Collection(colID).DocumentRefs(ctx).GetAll()
	Debugf("getall col:%v, doc:%v, len:%v, err:%v", colID, docID, len(dRefs), err)
	if err != nil {
		return rs, err
	}
	for _, d := range dRefs {
		var snap *firestore.DocumentSnapshot
		snap, err = d.Get(ctx)
		if err != nil {
			return rs, err
		}
		Debugf("id:%v, doc:%#v", d.ID, snap.Data())

		r, rErr := f.dataToStream(snap.Data())
		if rErr != nil {
			return rs, rErr
		}
		rs[d.ID] = r
	}
	return rs, err
}

// ReaderToStruct :
func ReaderToStruct(reader io.Reader, outStream io.Writer) error {
	var parser gojson.Parser = func(input io.Reader) (interface{}, error) {
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
	readerList, err := f.Scan(ctx, path)
	if err != nil {
		return err
	}
	for k, reader := range readerList {
		err = ReaderToStruct(reader, outStream)
		if err != nil {
			PrintAlertf(outStream, "cant convert err: %v at %v \n", err, k)
		}
	}
	return err
}
