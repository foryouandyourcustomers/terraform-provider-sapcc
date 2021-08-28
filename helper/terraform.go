package helper

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

func ResourceTest(t *testing.T, tfscript string, resourceName string) *tfjson.StateResource {
	tmpDir, err := ioutil.TempDir("", "tfinstall")

	if err != nil {
		t.Fatalf("error creating temp dir: %s", err)
	}

	defer os.RemoveAll(tmpDir)

	tfexecPath := os.Getenv("TF_EXEC_PATH")

	tf, err := tfexec.NewTerraform(tmpDir, tfexecPath)
	if err != nil {
		t.Fatalf("error running NewTerraform: %s", err)
	}

	err = writeScript(path.Join(tmpDir, "terraform.tf"), tfscript)
	if err != nil {
		log.Fatal(err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		t.Fatalf("error running Init: %s", err)
	}

	err = tf.Apply(context.Background())
	if err != nil {
		t.Fatalf("error running Apply: %s", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		t.Fatalf("error running state: %s", err)
	}

	return findResource(state.Values.RootModule.Resources, resourceName)
}

func writeScript(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func findResource(resources []*tfjson.StateResource, resourceName string) *tfjson.StateResource {
	for _, r := range resources {
		if strings.EqualFold(r.Name, resourceName) {
			return r
		}
	}

	return nil
}
