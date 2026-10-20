package helper

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDirectoryClientIgnoresEnvironmentProxy(t *testing.T) {
	t.Setenv("HTTP_PROXY", "http://proxy.example.com:8080")
	t.Setenv("HTTPS_PROXY", "http://proxy.example.com:8080")
	t.Setenv("ALL_PROXY", "http://proxy.example.com:8080")

	client := NewDirectoryClient(1001, "synthetic-token")
	transport, ok := client.HTTP.Transport.(*http.Transport)
	if !ok || transport.Proxy != nil {
		t.Fatal("directory client must dial Yandex directly")
	}
}

func TestFindUserByEmailUsesCurrentDirectoryFilter(t *testing.T) {
	var gotPath, gotQuery, gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		gotQuery = r.URL.RawQuery
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"limit":10,"offset":0,"total":1,"items":[{"id":100,"email":"person@example.com","aliases":["alias@example.com"],"name":{"first":"Person","last":"Example"}}]}`)
	}))
	defer server.Close()

	user, err := (DirectoryClient{
		CloudBase: server.URL,
		Token:     "synthetic-token",
		OrgID:     1001,
		HTTP:      server.Client(),
	}).FindUserByEmail("person@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/v1/directory/organizations/1001/users" || gotQuery != "email=person%40example.com&limit=10" {
		t.Fatalf("request = %s?%s", gotPath, gotQuery)
	}
	if gotAuth != "OAuth synthetic-token" {
		t.Fatal("authorization header mismatch")
	}
	if user == nil || user.Id != "100" || user.Email != "person@example.com" || len(user.Aliases) != 1 || user.Aliases[0] != "alias@example.com" {
		t.Fatalf("user = %+v", user)
	}
}
