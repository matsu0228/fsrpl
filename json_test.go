package main

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"google.golang.org/genproto/googleapis/type/latlng"
)

func getImportTestData() map[string]map[string]interface{} {
	// for import tests
	// NOTE: plaese match testData/*.json
	return map[string]map[string]interface{}{
		"user": map[string]interface{}{
			"_created_at": time.Date(2019, time.August, 26, 5, 0, 0, 0, time.UTC),
			"coin":        0,
			"favorites":   []interface{}{"1", "2"},
			"isDeleted":   true,
			"mapData":     map[string]interface{}{"b": true, "name": "mName"},
			"name":        "user",
		},
		"typelist": map[string]interface{}{
			"geo":    &latlng.LatLng{Latitude: 80, Longitude: 80},
			"nil":    nil,
			"ref":    fsEmulator.Client.Doc("dest/dog"),
			"refSub": fsEmulator.Client.Doc("original/dog/subCol/subDoc"),
			"str":    "string",
			"time":   time.Date(2020, time.August, 31, 15, 0, 0, 0, time.UTC),
		},
	}
}

func TestImportDataFromJSONFiles(t *testing.T) {
	ctx := context.Background()
	importPath := "./testData/"
	firestorePath := "import/"
	importDataList := getImportTestData()

	if err := ImportDataFromJSONFiles(ctx, testOpt, fsEmulator, importPath, fmt.Sprintf("%s*", firestorePath)); err != nil {
		t.Errorf("failed import err:%v", err)
	}

	for docID := range importDataList {
		gotSnap, err := fsEmulator.Client.Doc(fmt.Sprintf("%s%s", firestorePath, docID)).Get(ctx)
		if err != nil {
			t.Errorf("failed get %s err:%v", docID, err)
		}
		got := gotSnap.Data()
		if docID == "typelist" { // NOTE: typelistでcrashするためDeepEqualで確認
			if !reflect.DeepEqual(got, importDataList[docID]) {
				t.Errorf("invalid imported data:%s \ngot=%#v, \nwant=%#v", docID, got, importDataList[docID])
			}
			continue
		}
		if diff := pretty.Compare(got, importDataList[docID]); diff != "" {
			t.Errorf("invalid decoded json: %s", pretty.Compare(got, importDataList[docID]))
		}
	}
}
