package main

import (
	"context"
	"path"

	"go.uber.org/multierr"
)

// DumpCmd is commands to export data from firestore to local json
type DumpCmd struct {
	FirestorePath string `arg required name:"firestore_path" help:"Target firestore path. The path containts colleaction's path and document id. Document id allowed wildcard character (*). (e.g. collectionName/docId, collecntionName/*)"`
	Credentials   string `optional name:"cred" help:"Set target firestore's credentail file path."`
	Path          string `optional default:"./" help:"Export local path."`
	ShowGoStruct  bool   `optional default:"false" help:"Show go struct mode without json file exportation." optional`
	// cp IsDelete     bool `help:"delete source document data after dump." optional`
}

// Run is main function
func (d *DumpCmd) Run(opt *Option) error {
	Debugf("dump from %v to %v \n", d.FirestorePath, d.Path)
	ctx := context.Background()
	fs, err := NewFirebase(ctx, opt, OptWithCred(d.Credentials))
	if err != nil {
		return err
	}
	if d.ShowGoStruct {
		return d.showGoStruct(ctx, opt, fs)
	}

	readerList, err := fs.Scan(ctx, d.FirestorePath)
	if err != nil {
		return err
	}
	for k, reader := range readerList {
		Debugf("write outStream at: %v\n", k)

		fn := path.Join(d.Path, k+".json")
		err = multierr.Append(err, writeFile(fn, reader))
	}
	if err != nil {
		return err
	}
	PrintInfof(opt.Stdout, "Dump complete! \n\n")
	return nil
}

func (d *DumpCmd) showGoStruct(ctx context.Context, opt *Option, fs *Firestore) error {
	return fs.ToStruct(ctx, d.FirestorePath, opt.Stdout)
}
