package user

import (
	"bytes"
	"strings"
	"testing"

	"github.com/foxsoft2005/y360c/cmd/user/alias"
	"github.com/spf13/cobra"
)

func TestHelpDescribesFlagsAndShowsExamples(t *testing.T) {
	cases := []struct {
		cmd  *cobra.Command
		want []string
	}{
		{Cmd, []string{"Examples:", "y360c user add --help", "Y360_TOKEN"}},
		{addCmd, []string{"Examples:", "--nickname", "--department-id", "department id", "y360c user add --nickname person --department-id 42", "Y360_PASSWORD"}},
		{setCmd, []string{"Examples:", "--department-id", "--is-enabled", "y360c user set --email person@example.com --department-id 42"}},
		{listCmd, []string{"Examples:", "--email", "--dept-id", "y360c user ls --dept-id 42"}},
		{infoCmd, []string{"Examples:", "--email", "y360c user info --email person@example.com"}},
		{alias.Cmd, []string{"Examples:", "--alias", "y360c user alias ls --email person@example.com", "y360c user alias add --email person@example.com --alias alias@example.com"}},
	}
	for _, tc := range cases {
		text := helpText(t, tc.cmd)
		for _, part := range tc.want {
			if !strings.Contains(text, part) {
				t.Errorf("%s help missing %q\n%s", tc.cmd.Name(), part, text)
			}
		}
	}
}

func helpText(t *testing.T, cmd *cobra.Command) string {
	t.Helper()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	if err := cmd.Help(); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}
