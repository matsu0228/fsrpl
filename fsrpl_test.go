package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

var (
	fsEmulator *Firestore
	testOpt    = &Option{Debug: true, Stdout: os.Stdout, Stderr: os.Stderr}

	// for scan tests
	scanFirestorePath = "scan/"
	scanDataList      = map[string]map[string]interface{}{
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

	ctx := context.Background()
	println("before all...", os.Getenv("FIRESTORE_EMULATOR_HOST"))
	if err := setupEmulator(ctx); err != nil {
		log.Fatalf("failed setupEmulator: %v", err)
	}

	// defer teardown()
	os.Exit(m.Run())
}

func setupEmulator(ctx context.Context) error {
	var err error
	conOpt := FirestoreConnectionOption{}
	fsEmulator, err = NewFirebase(ctx, testOpt, conOpt)
	if err != nil {
		return err
	}

	return fsEmulator.importTestdata()
}

func (f *Firestore) importTestdata() error {
	ctx := context.Background()
	for k, testData := range scanDataList {
		path := scanFirestorePath + k
		if err := f.SaveData(ctx, testOpt, path, testData); err != nil {
			return err
		}
	}
	return nil
}
