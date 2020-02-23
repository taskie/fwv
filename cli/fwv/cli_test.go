package fwv_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/taskie/fwv/cli/fwv"
	"github.com/taskie/ose"
	"github.com/taskie/ose/coli"
)

func setUp(args ...string) (w *ose.FakeWorld, cl *coli.Coli, cmd *cobra.Command) {
	w = ose.NewFakeWorld()
	ose.SetWorld(w)
	cl = coli.NewColiInThisWorld()
	cmd = fwv.NewCommand(cl)
	cmd.SetArgs(args)
	return
}

func TestCli(t *testing.T) {
	w, _, cmd := setUp("--no-color", "-t", "csv", "--debug")
	w.FakeIO.InBuf.WriteString("ab   c    d\ne   fg   h\n")
	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	actual := w.FakeIO.OutBuf.String()
	if "ab,c,d\ne,fg,h\n" != actual {
		t.Fatalf("invalid output: %s", actual)
	}
}
