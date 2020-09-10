// +build integration

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kylelemons/godebug/pretty"
)

var (
	fsEmulator   *Firestore
	testOpt      = &Option{Debug: true, Stdout: os.Stdout, Stderr: os.Stderr}
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
	println("before all...", os.Getenv("FIRESTORE_EMULATOR_HOST"))

	fsEmulator, err = connectFirebaseEmulator()
	if err != nil {
		log.Fatalf("[ERROR] cant connect firebase emulator: %v", err)
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
	var client *firestore.Client
	var err error

	ctx := context.Background()
	client, err = firestore.NewClient(ctx, "fs-test")
	if err != nil {
		return nil, err
	}
	return &Firestore{
		Client: client,
	}, nil
}

func (f *Firestore) importTestdata() error {
	ctx := context.Background()
	pathNode := "testData/"
	for k, testData := range testDataList {
		path := pathNode + k
		if err := f.SaveData(ctx, testOpt, path, testData); err != nil {
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
				fmt.Printf("[TEST] undefined data key:%v", k)
				continue
			}
			got := fsEmulator.InterpretationEachValueForTime(org)
			if diff := pretty.Compare(got, testDataList[k]); diff != "" {
				t.Errorf("invalid decoded json: %s", pretty.Compare(got, testDataList[k]))
			}
		}
	}
}
