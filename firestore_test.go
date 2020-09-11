package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestScan(t *testing.T) {

	ctx := context.Background()
	path := fmt.Sprintf("%s*", scanFirestorePath)
	readers, err := fsEmulator.Scan(ctx, path)
	if err != nil {
		t.Errorf("cant scan from %v err:%v", path, err)
	}

	for k, reader := range readers {
		var org map[string]interface{}
		if err := json.NewDecoder(reader).Decode(&org); err != nil {
			t.Errorf("cant decode documents  key:%v err:%v", k, err)
		} else {
			if _, ok := scanDataList[k]; !ok {
				log.Printf("[TEST] undefined data key:%v", k)
				continue
			}
			got := fsEmulator.InterpretationEachValueForTime(org)
			if diff := pretty.Compare(got, scanDataList[k]); diff != "" {
				t.Errorf("invalid decoded json: %s", pretty.Compare(got, scanDataList[k]))
			}
			log.Printf("[TEST] got data: %v = %v", k, got)
		}
	}
}
