package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

// CopyCmd is commands to copy data from some firestore path to another firestore path
type CopyCmd struct {
	FirestorePath            string `arg required name:"firestore_path" help:"Source firestore path. The path containts colleaction's path and document id. Document id allowed wildcard character (*). (e.g. collectionName/docId, collecntionName/*)"`
	Credentials              string `optional name:"cred" help:"Set source firestore's credentail file path."`
	DestinationFirestorePath string `required name:"dest" help:"Destination firestore path. The path containts colleaction's path and document id. Document id allowed wildcard character (*). (e.g. collectionName/docId, collecntionName/*)"`
	DestinationCredentials   string `optional name:"dest_cred" help:"Set destination firestore's credentail file path."`
	IsDelete                 bool   `optional default:"false" help:"delete source document data after dump." optional`
}

// Run is main function
func (c *CopyCmd) Run(opt *Option) error {
	Debugf("copy from %v to %v isAnotherFirestore:%v", c.FirestorePath, c.DestinationFirestorePath, len(c.DestinationCredentials) != 0)
	ctx := context.Background()
	srcFs, err := NewFirebase(ctx, opt, c.Credentials)
	if err != nil {
		return err
	}

	if len(c.DestinationCredentials) != 0 {
		PrintInfof(opt.Stdout, "use destination firestore")
		dstFs, err := NewFirebase(ctx, opt, c.DestinationCredentials)
		if err != nil {
			return err
		}
		return c.Replicate(ctx, opt, srcFs, dstFs)
	}

	return c.Replicate(ctx, opt, srcFs, srcFs)
}

// Replicate from some firestore path to another firestore path
func (c *CopyCmd) Replicate(ctx context.Context, opt *Option, srcFs, dstFs *Firestore) error {
	var err error
	if c.IsDelete {
		PrintInfof(opt.Stdout, "delete original document? (y/n) \n")
		yes := askForConfirmation(opt)
		if !yes {
			return errors.New("exit")
		}
	}

	readerList, err := srcFs.Scan(ctx, c.FirestorePath)
	if err != nil {
		return err
	}
	for k, reader := range readerList {
		dstPath := strings.Replace(c.DestinationFirestorePath, "*", k, -1)
		srcPath := strings.Replace(c.FirestorePath, "*", k, -1)
		Debugf("save with : %v from %v \n", srcPath, srcPath)

		var m map[string]interface{}
		err = json.NewDecoder(reader).Decode(&m)
		if err != nil {
			return err
		}
		om := InterpretationEachValueForTime(m)

		err = dstFs.SaveData(ctx, opt, dstPath, om)
		if err != nil {
			return err
		}

		if c.IsDelete {
			err = dstFs.DeleteData(ctx, opt, srcPath)
			if err != nil {
				return err
			}
		}
	}
	PrintInfof(opt.Stdout, "Copy complete! \n\n")
	return nil
}
