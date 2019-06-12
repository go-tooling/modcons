// +build cli_tests

package cli_tests

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

func handler(w http.ResponseWriter, r *http.Request) {
	rules, _ := ioutil.ReadFile("../test_data/rules.modcons")

	w.Write(rules)
}

func TestMain(m *testing.M) {

	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", handler) // set router
	go func() {
		srv.ListenAndServe()
	}()
	defer srv.Shutdown(nil)

	os.Exit(m.Run())
}

func Test_Local_Paths_Ok(t *testing.T) {
	cmd := exec.Command("../artefacts/modcons", "--rulepath=../test_data/rules.modcons", "--modpath=../go.mod")
	cmd.Stdout = os.Stdout
	cmd.Run()

	if cmd.ProcessState.ExitCode() != 0 {
		t.Fatal()
	}
}

func Test_Local_ParseOnly_Ok(t *testing.T) {
	cmd := exec.Command("../artefacts/modcons", "--rulepath=../test_data/rules.modcons", "--modpath=../go.mod", "--parseOnly=true")
	cmd.Stdout = os.Stdout

	cmd.Run()

	if cmd.ProcessState.ExitCode() != 0 {
		t.Fatal()
	}
}

func Test_Local_ParseOnly_NotOk(t *testing.T) {
	cmd := exec.Command("../artefacts/modcons", "--rulepath=../test_data/rules_bad.modcons", "--modpath=../go.mod", "--parseOnly=true")
	cmd.Stdout = os.Stdout
	cmd.Run()

	if cmd.ProcessState.ExitCode() != 1 {
		t.Fatal()
	}
}

func Test_Local_Paths_NotOk(t *testing.T) {
	cmd := exec.Command("../artefacts/modcons", "--rulepath=../test_data/rules2.modcons", "--modpath=../go.mod")
	cmd.Stdout = os.Stdout
	cmd.Run()

	if cmd.ProcessState.ExitCode() != 1 {
		t.Fatal()
	}
}

func Test_Http_Rules_Ok(t *testing.T) {
	cmd := exec.Command("../artefacts/modcons", "--rulepath=http://localhost:8080", "--modpath=../go.mod")
	cmd.Stdout = os.Stdout
	cmd.Run()

	if cmd.ProcessState.ExitCode() != 0 {
		t.Fatal()
	}
}
