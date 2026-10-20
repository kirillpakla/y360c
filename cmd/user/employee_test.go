package user

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/foxsoft2005/y360c/helper"
)

func TestCreateEmployeePostsOfficialContractWithoutAdminOrPasswordFlag(t *testing.T) {
	if addCmd.Flags().Lookup("password") != nil {
		t.Fatal("create command accepts a password flag")
	}

	var gotMethod, gotPath, gotAuth string
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotAuth = r.Header.Get("Authorization")
		var err error
		gotBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"user-100","email":"person@example.com","nickname":"person","name":{"first":"Person","last":"Example"}}`)
	}))
	defer server.Close()

	_, err := (helper.DirectoryClient{
		DirectoryBase: server.URL,
		CloudBase:     server.URL,
		Token:         "synthetic-token",
		OrgID:         1001,
		HTTP:          server.Client(),
	}).CreateEmployee(helper.CreateEmployeeInput{
		Nickname:     "person",
		DepartmentID: 42,
		First:        "Person",
		Last:         "Example",
		Password:     "synthetic-password",
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPost || gotPath != "/directory/v1/org/1001/users" {
		t.Fatalf("request = %s %s", gotMethod, gotPath)
	}
	if gotAuth != "OAuth synthetic-token" {
		t.Fatal("authorization header mismatch")
	}
	body := string(gotBody)
	for _, part := range []string{
		`"nickname":"person"`,
		`"departmentId":42`,
		`"first":"Person"`,
		`"last":"Example"`,
		`"password":"synthetic-password"`,
	} {
		if !strings.Contains(body, part) {
			t.Fatalf("body missing %s", part)
		}
	}
	if strings.Contains(body, "isAdmin") || strings.Contains(body, "passwordChangeRequired") {
		t.Fatal("create body included an excluded field")
	}

	_, err = (helper.DirectoryClient{}).CreateEmployee(helper.CreateEmployeeInput{Nickname: "person", DepartmentID: 42, First: "Person", Last: "Example"})
	if err == nil || strings.Contains(err.Error(), "synthetic-password") {
		t.Fatalf("empty password error = %v", err)
	}
}

func TestPatchEmployeeSendsMissingFieldsAndEnabledState(t *testing.T) {
	for _, name := range []string{"nickname", "first-name", "last-name", "middle-name", "gender", "language", "timezone", "is-enabled"} {
		if setCmd.Flags().Lookup(name) == nil {
			t.Fatalf("flag %s is missing", name)
		}
	}

	var gotMethod, gotPath string
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		var err error
		gotBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		_, _ = io.WriteString(w, `{"id":"user-100","email":"person@example.com","isEnabled":false}`)
	}))
	defer server.Close()

	_, err := (helper.DirectoryClient{
		DirectoryBase: server.URL,
		Token:         "synthetic-token",
		OrgID:         1001,
		HTTP:          server.Client(),
	}).PatchEmployee("user-100", helper.PatchInput{
		Nickname:  "person",
		First:     "Person",
		Last:      "Example",
		Middle:    "Middle",
		Gender:    "unspecified",
		Language:  "ru",
		Timezone:  "Europe/Moscow",
		IsEnabled: "no",
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/directory/v1/org/1001/users/user-100" {
		t.Fatalf("request = %s %s", gotMethod, gotPath)
	}
	body := string(gotBody)
	for _, part := range []string{
		`"nickname":"person"`,
		`"first":"Person"`,
		`"last":"Example"`,
		`"middle":"Middle"`,
		`"gender":"unspecified"`,
		`"language":"ru"`,
		`"timezone":"Europe/Moscow"`,
		`"isEnabled":false`,
	} {
		if !strings.Contains(body, part) {
			t.Fatalf("body missing %s", part)
		}
	}
	if strings.Contains(body, "isAdmin") || strings.Contains(body, "password") || strings.Contains(body, "departmentId") {
		t.Fatal("patch body included an excluded field")
	}
}

func TestPatchEmployeeMovesUserToDepartment(t *testing.T) {
	if setCmd.Flags().Lookup("department-id") == nil {
		t.Fatal("set command has no department-id flag")
	}

	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		gotBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if r.Method != http.MethodPatch || r.URL.EscapedPath() != "/directory/v1/org/1001/users/user-100" {
			t.Errorf("request = %s %s", r.Method, r.URL.EscapedPath())
		}
		_, _ = io.WriteString(w, `{"id":"user-100","departmentId":42,"email":"person@example.com"}`)
	}))
	defer server.Close()

	departmentID := 42
	_, err := (helper.DirectoryClient{
		DirectoryBase: server.URL,
		Token:         "synthetic-token",
		OrgID:         1001,
		HTTP:          server.Client(),
	}).PatchEmployee("user-100", helper.PatchInput{DepartmentID: &departmentID})
	if err != nil {
		t.Fatal(err)
	}
	body := string(gotBody)
	if !strings.Contains(body, `"departmentId":42`) {
		t.Fatalf("body = %s", body)
	}
	if strings.Contains(body, "isAdmin") || strings.Contains(body, "password") {
		t.Fatal("department move included an excluded field")
	}
}

func TestAliasCommandsUseDocumentedHosts(t *testing.T) {
	type call struct {
		method string
		path   string
		body   string
	}
	var calls []call
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		calls = append(calls, call{method: r.Method, path: r.URL.EscapedPath(), body: string(body)})
		switch r.Method {
		case http.MethodDelete:
			_, _ = io.WriteString(w, `{"alias":"alias@example.com","removed":true}`)
		default:
			_, _ = io.WriteString(w, `{"id":100,"email":"person@example.com","aliases":["alias@example.com"],"name":{"first":"Person","last":"Example"}}`)
		}
	}))
	defer server.Close()

	client := helper.DirectoryClient{
		DirectoryBase: server.URL,
		CloudBase:     server.URL,
		Token:         "synthetic-token",
		OrgID:         1001,
		HTTP:          server.Client(),
	}
	added, err := client.AddAlias("user-100", "alias@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(added.Aliases) != 1 || added.Aliases[0] != "alias@example.com" || added.ID != "100" {
		t.Fatalf("add result id=%q aliases=%v", added.ID, added.Aliases)
	}
	if _, err := client.ListAliases("user-100"); err != nil {
		t.Fatal(err)
	}
	removed, err := client.RemoveAlias("user-100", "alias@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if !removed.Removed || removed.Alias != "alias@example.com" {
		t.Fatalf("remove result = %+v", removed)
	}
	want := []call{
		{http.MethodPost, "/v1/api360/directory/org/1001/users/user-100/aliases", `{"alias":"alias@example.com"}`},
		{http.MethodGet, "/directory/v1/org/1001/users/user-100", ""},
		{http.MethodDelete, "/v1/api360/directory/org/1001/users/user-100/aliases/alias%40example.com", ""},
	}
	if len(calls) != len(want) {
		t.Fatalf("calls = %d", len(calls))
	}
	for i := range want {
		if calls[i].method != want[i].method || calls[i].path != want[i].path || calls[i].body != want[i].body {
			t.Fatalf("call %d = %s %s %s", i, calls[i].method, calls[i].path, calls[i].body)
		}
	}
}
