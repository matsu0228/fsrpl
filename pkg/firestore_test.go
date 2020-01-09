// +build integration

package fsrpl

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/kylelemons/godebug/pretty"
)

var (
	fsEmulator   *Firestore
	testDataList = map[string]map[string]interface{}{
		"one": map[string]interface{}{
			"string": "value",
			"int":    100,
			"time":   time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC),
			"bool":   true,
			"array":  []interface{}{"hoge", "fuga"},
			"map":    map[string]interface{}{"one": 1, "two": 2},
		},
		"two": map[string]interface{}{
			"complex": map[string]interface{}{"one": []interface{}{"1", "11", "111"}, "two": 2},
		},
	}
)

func TestMain(m *testing.M) {

	var err error
	println("before all...")

	fsEmulator, err = connectFirebaseEmulator()
	if err != nil {
		log.Fatalf("[ERROR] cant connect firebase emulator")
	}

	err = fsEmulator.importTestdata()
	if err != nil {
		log.Fatalf("[ERROR] cant import testdata err:%v", err)
	}

	code := m.Run()

	println("after all...")

	os.Exit(code)
}

func connectFirebaseEmulator() (*Firestore, error) {
	var app *firebase.App
	var client *firestore.Client
	var err error

	ctx := context.Background()
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

func (f *Firestore) importTestdata() error {
	ctx := context.Background()
	pathNode := "testData/"
	for k, testData := range testDataList {
		path := pathNode + k
		if err := f.SaveData(ctx, path, testData); err != nil {
			return err
		}
	}
	return nil
}

func TestScan(t *testing.T) {

	ctx := context.Background()
	path := "testData/*"
	readers, err := fsEmulator.Scan(ctx, path)
	if err != nil {
		t.Errorf("cant scan from %v err:%v", path, err)
	}

	for k, reader := range readers {
		var org map[string]interface{}
		if err := json.NewDecoder(reader).Decode(&org); err != nil {
			t.Errorf("cant decode documents  key:%v err:%v", k, err)
		} else {
			if _, ok := testDataList[k]; !ok {
				log.Printf("[TEST] undefined data key:%v", k)
				continue
			}
			got := InterpretationEachValueForTime(org)
			if diff := pretty.Compare(got, testDataList[k]); diff != "" {
				t.Errorf("invalid decoded json: %s", pretty.Compare(got, testDataList[k]))
			}
		}
	}
}
