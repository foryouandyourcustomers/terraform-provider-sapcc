package helper

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	tfjson "github.com/hashicorp/terraform-json"
)

func ResourceTest(t *testing.T, tfscript string, resourceName string) (*tfjson.StateResource, []string, []string) {
	cwd, _ := os.Getwd()
	tmpScript, _ := ioutil.TempDir(cwd, "tfscript")
	tmpInstall, err := ioutil.TempDir("", "tfinstall")
	var stdErrBuf strings.Builder
	var errorList, warnList []string

	fmt.Printf("%s", tmpScript)

	defer os.RemoveAll(tmpInstall)
	defer os.RemoveAll(tmpScript)

	tfexecPath := os.Getenv("TF_ACC_TERRAFORM_EXEC_PATH")
	tfVersion := os.Getenv("TF_ACC_TERRAFORM_VERSION")

	if tfVersion != "" {
		tfexecPath, err = tfinstall.Find(context.Background(), tfinstall.ExactVersion(tfVersion, tmpInstall))
		if err != nil {
			log.Fatalf("error locating Terraform binary: %s", err)
		}
	} else {
		if tfexecPath == "" {
			tfexecPath, err = tfinstall.Find(context.Background(), tfinstall.LatestVersion(tmpInstall, false))
			if err != nil {
				t.Fatalf("error locating Terraform binary: %s", err)
			}
		}
	}

	tf, err := tfexec.NewTerraform(tmpScript, tfexecPath)
	if err != nil {
		t.Fatalf("error running NewTerraform: %s", err.Error())
	}

	tf.SetStderr(&stdErrBuf)
	ctx := context.Background()
	defer tf.Destroy(context.Background())

	err = writeScript(path.Join(tmpScript, "terraform.tf"), tfscript)
	if err != nil {
		log.Fatal(err)
	}

	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		t.Fatalf("error running Init: %+v", err)
	}

	err = tf.Apply(ctx)
	if err != nil {

		for _, str := range strings.Split(stdErrBuf.String(), "\n") {

			if str != "" && strings.Contains(str, "Error") {
				errorList = append(errorList, str)
			}

			if str != "" && strings.Contains(str, "Warn") {
				warnList = append(warnList, str)
			}
		}
	}

	state, err := tf.Show(ctx)

	if err != nil {
		t.Fatalf("error showing state: %+v", err)
	}

	return findResource(state.Values.RootModule.Resources, resourceName), errorList, warnList
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
