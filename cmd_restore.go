package main

import (
	"context"
)

// RestoreCmd is commands to import data from local json to firestore
type RestoreCmd struct {
	FirestorePath     string `arg required name:"firestore_path" help:"The destination firestore path to restore to. The path containts colleaction's path and document id. Document id allowed wildcard character (*). (e.g. collectionName/docId, collecntionName/*)"`
	Path              string `required help:"The path to the local file containing the data to be restored."`
	Credentials       string `optional name:"cred" help:"Set target firestore's credentail file path."`
	EmulatorProjectID string `optional name:"emulators-project-id" help:"Set projectID of firestore emulator."`
}

// Run is main function
func (r *RestoreCmd) Run(opt *Option) error {
	Debugf("restore %v, %v, %#v", r.FirestorePath, r.Path, opt)
	ctx := context.Background()

	conOpt := OptWithCred(r.Credentials)
	if len(r.EmulatorProjectID) != 0 {
		conOpt = OptWithEmulatorProjectID(r.EmulatorProjectID)
	}
	fs, err := NewFirebase(ctx, opt, conOpt)
	if err != nil {
		return err
	}
	err = ImportDataFromJSONFiles(ctx, opt, fs, r.Path, r.FirestorePath)
	if err != nil {
		return err
	}
	PrintInfof(opt.Stdout, "Restore complete! \n\n")
	return nil
}
