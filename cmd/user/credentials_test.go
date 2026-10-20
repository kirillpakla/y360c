package user

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendCredentialsBuildsTableForSeveralAccounts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "creds.csv")
	first := credential{Nickname: "person", Email: "person@example.com", UserID: "user-100", Password: "synthetic-password"}
	if err := appendCredentials(path, []credential{first}); err != nil {
		t.Fatal(err)
	}
	second := credential{Nickname: "other", Email: "other@example.com", UserID: "user-200", Password: "other-password"}
	if err := appendCredentials(path, []credential{second}); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %o", info.Mode().Perm())
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(body)
	for _, part := range []string{
		"nickname,email,user_id,password\n",
		"person,person@example.com,user-100,synthetic-password\n",
		"other,other@example.com,user-200,other-password\n",
	} {
		if !strings.Contains(got, part) {
			t.Fatalf("file missing %q\n%s", part, got)
		}
	}
	if strings.Count(got, "nickname,email,user_id,password") != 1 {
		t.Fatalf("header repeated:\n%s", got)
	}
}
