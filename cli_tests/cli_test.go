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
	rules, _ := ioutil.ReadFile("../test_data/rules.modcop")

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
	cmd := exec.Command("../artefacts/modcop", "--rulepath=../test_data/rules.modcop", "--modpath=../go.mod")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}

func Test_Http_Rules_Ok(t *testing.T) {
	cmd := exec.Command("../artefacts/modcop", "--rulepath=http://localhost:8080", "--modpath=../go.mod")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}
