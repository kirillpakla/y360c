package helper

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSettingsReadsOrgAndBaseFromFileAndTokenFromEnv(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	body := "" +
		"# local non-secrets\n" +
		"Y360_ORG_ID=1001\n" +
		"Y360_API_BASE=http://127.0.0.1:9\n" +
		"Y360_TOKEN=file-secret\n"
	if err := os.WriteFile(envFile, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	settings, err := LoadSettings(envFile, func(key string) (string, bool) {
		if key == "Y360_TOKEN" {
			return "env-secret", true
		}
		return "", false
	})
	if err != nil {
		t.Fatal(err)
	}
	if settings.OrgID != 1001 {
		t.Fatalf("org id = %d", settings.OrgID)
	}
	if settings.APIBase != "http://127.0.0.1:9" {
		t.Fatalf("api base = %q", settings.APIBase)
	}
	if settings.Token != "env-secret" {
		t.Fatal("token was not taken from the environment")
	}
}

func TestLoadSettingsIgnoresFileTokenWhenEnvIsEmpty(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	const fileToken = "file-secret"
	if err := os.WriteFile(envFile, []byte("Y360_ORG_ID=1001\nY360_TOKEN="+fileToken+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	settings, err := LoadSettings(envFile, func(string) (string, bool) { return "", false })
	if !errors.Is(err, ErrTokenMissing) {
		t.Fatalf("error = %v", err)
	}
	if settings.Token != "" {
		t.Fatal("file token was accepted")
	}
	if err != nil && strings.Contains(err.Error(), fileToken) {
		t.Fatal("missing token error echoed the file value")
	}
}

func TestLoadSettingsMissingTokenDoesNotEchoSecret(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("Y360_ORG_ID=1001\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadSettings(envFile, func(string) (string, bool) { return "", false })
	if !errors.Is(err, ErrTokenMissing) {
		t.Fatalf("error = %v", err)
	}
	if strings.Contains(err.Error(), "1001") {
		t.Fatal("missing token error included the organization id")
	}
}

func TestOrgIDStatusDoesNotRevealValue(t *testing.T) {
	if OrgIDStatus(0) != "unset" {
		t.Fatal("zero org id should be unset")
	}
	if OrgIDStatus(1001) != "set" {
		t.Fatal("present org id should be set")
	}
	if strings.Contains(OrgIDStatus(1001), "1001") {
		t.Fatal("status revealed the organization id")
	}
}
