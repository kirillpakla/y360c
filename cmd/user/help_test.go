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
		{Cmd, []string{"Examples:", "y360c user add --help", "y360c user set --help", "Y360_TOKEN"}},
		{addCmd, []string{
			"Examples:", "Required:", "job title", "YYYY-MM-DD",
			"--nickname", "--department-id", "--first-name", "--last-name", "--middle-name Middle",
			"--position Engineer", "--display-name", "--about", "--birthday 1990-01-31",
			"--gender unspecified", "--language ru", "--timezone Europe/Moscow", "--external-id ext-100",
			"--password TempPass-1", "--output creds.csv",
		}},
		{setCmd, []string{
			"Examples:", "not sent", "job title", "YYYY-MM-DD", "yes or no",
			"--email person@example.com", "--position Engineer", "--department-id 42",
			"--display-name", "--about", "--birthday 1990-01-31", "--external-id ext-100",
			"--first-name Person", "--last-name Example", "--middle-name Middle", "--nickname person",
			"--gender unspecified", "--language ru", "--timezone Europe/Moscow",
			"--id user-100", "--is-enabled no", "--is-admin no",
		}},
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
		if strings.Contains(text, "21-school.ru") || strings.Contains(text, "test-tool") {
			t.Errorf("%s help contains a private example\n%s", tc.cmd.Name(), text)
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
