package dept

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestDepartmentHelpDescribesFlagsAndShowsExamples(t *testing.T) {
	cases := []struct {
		cmd  *cobra.Command
		want []string
	}{
		{DeptCmd, []string{"Examples:", "y360c dept ls"}},
		{lsCmd, []string{"Examples:", "--parent-id", "--order-by", "parent", "y360c dept ls --parent-id 10"}},
	}
	for _, tc := range cases {
		var buf bytes.Buffer
		tc.cmd.SetOut(&buf)
		tc.cmd.SetErr(&buf)
		if err := tc.cmd.Help(); err != nil {
			t.Fatal(err)
		}
		text := buf.String()
		for _, part := range tc.want {
			if !strings.Contains(text, part) {
				t.Errorf("%s help missing %q\n%s", tc.cmd.Name(), part, text)
			}
		}
	}
}
