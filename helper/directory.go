package helper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/foxsoft2005/y360c/model"
)

// DirectoryClient calls the employee methods against a caller-supplied host.
// Tests pass a local server. An empty base keeps the documented host.
type DirectoryClient struct {
	DirectoryBase string
	CloudBase     string
	Token         string
	OrgID         int
	HTTP          *http.Client
}

// CreateEmployeeInput is the employee create contract.
// Password is a value already read from the process environment, never a command flag.
type CreateEmployeeInput struct {
	Nickname     string
	DepartmentID int
	First        string
	Last         string
	Middle       string
	Password     string
	About        string
	Birthday     string
	DisplayName  string
	ExternalID   string
	Gender       string
	Position     string
	Language     string
	Timezone     string
}

// PatchInput is the official user PATCH body except administrator and password changes.
// DepartmentID is sent only when the caller set it.
type PatchInput struct {
	About        string
	Birthday     string
	DisplayName  string
	ExternalID   string
	Position     string
	Nickname     string
	Gender       string
	Language     string
	Timezone     string
	First        string
	Last         string
	Middle       string
	DepartmentID *int
	IsAdmin      EnumYesNo
	IsEnabled    EnumYesNo
}

// AliasResult is the part of an alias response the CLI prints.
type AliasResult struct {
	ID      string
	Aliases []string
	Name    model.UserName
}

func NewDirectoryClient(orgID int, token string) DirectoryClient {
	return DirectoryClient{
		DirectoryBase: DirectoryBase(),
		CloudBase:     CloudBase(),
		Token:         token,
		OrgID:         orgID,
		HTTP:          &http.Client{Timeout: 60 * time.Second},
	}
}

func (c DirectoryClient) directoryBase() string {
	if c.DirectoryBase != "" {
		return strings.TrimRight(c.DirectoryBase, "/")
	}
	return BaseUrl
}

func (c DirectoryClient) cloudBase() string {
	if c.CloudBase != "" {
		return strings.TrimRight(c.CloudBase, "/")
	}
	return BaseUrlV2
}

// CreateEmployee sends POST /directory/v1/org/{orgId}/users on api360.yandex.net.
func (c DirectoryClient) CreateEmployee(in CreateEmployeeInput) (*model.User, error) {
	if in.Password == "" {
		return nil, errors.New("password must be specified")
	}
	if in.Nickname == "" || in.DepartmentID == 0 || in.First == "" || in.Last == "" {
		return nil, errors.New("nickname, department id, first name and last name must be specified")
	}
	body := map[string]any{
		"nickname":     in.Nickname,
		"departmentId": in.DepartmentID,
		"name":         nameBody(in.First, in.Last, in.Middle),
		"password":     in.Password,
	}
	putString(body, "about", in.About)
	putString(body, "birthday", in.Birthday)
	putString(body, "displayName", in.DisplayName)
	putString(body, "externalId", in.ExternalID)
	putString(body, "gender", in.Gender)
	putString(body, "position", in.Position)
	putString(body, "language", in.Language)
	putString(body, "timezone", in.Timezone)

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/directory/v1/org/%d/users", c.directoryBase(), c.OrgID)
	resp, err := c.do(http.MethodPost, endpoint, payload)
	if err != nil {
		return nil, err
	}
	user, err := decodeDirectoryUser(resp.Body)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// PatchEmployee sends PATCH /directory/v1/org/{orgId}/users/{userId} on api360.yandex.net.
func (c DirectoryClient) PatchEmployee(userID string, in PatchInput) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("user id must be specified")
	}
	body := map[string]any{}
	putString(body, "about", in.About)
	putString(body, "birthday", in.Birthday)
	putString(body, "displayName", in.DisplayName)
	putString(body, "externalId", in.ExternalID)
	putString(body, "position", in.Position)
	putString(body, "nickname", in.Nickname)
	putString(body, "gender", in.Gender)
	putString(body, "language", in.Language)
	putString(body, "timezone", in.Timezone)
	if in.First != "" || in.Last != "" || in.Middle != "" {
		body["name"] = nameBody(in.First, in.Last, in.Middle)
	}
	if in.DepartmentID != nil {
		body["departmentId"] = *in.DepartmentID
	}
	if value := EnumYesNoToBool(in.IsAdmin); value != nil {
		body["isAdmin"] = *value
	}
	if value := EnumYesNoToBool(in.IsEnabled); value != nil {
		body["isEnabled"] = *value
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/directory/v1/org/%d/users/%s", c.directoryBase(), c.OrgID, pathSegment(userID))
	resp, err := c.do(http.MethodPatch, endpoint, payload)
	if err != nil {
		return nil, err
	}
	user, err := decodeDirectoryUser(resp.Body)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// AddAlias sends POST to the cloud-api alias method.
func (c DirectoryClient) AddAlias(userID, alias string) (*AliasResult, error) {
	payload, err := json.Marshal(struct {
		Alias string `json:"alias"`
	}{Alias: alias})
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/v1/api360/directory/org/%d/users/%s/aliases", c.cloudBase(), c.OrgID, pathSegment(userID))
	resp, err := c.do(http.MethodPost, endpoint, payload)
	if err != nil {
		return nil, err
	}
	user, err := decodeDirectoryUser(resp.Body)
	if err != nil {
		return nil, err
	}
	return &AliasResult{ID: user.Id, Aliases: user.Aliases, Name: user.Name}, nil
}

// ListAliases reads aliases from the directory user resource on api360.yandex.net.
func (c DirectoryClient) ListAliases(userID string) (*model.User, error) {
	endpoint := fmt.Sprintf("%s/directory/v1/org/%d/users/%s", c.directoryBase(), c.OrgID, pathSegment(userID))
	resp, err := c.do(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	user, err := decodeDirectoryUser(resp.Body)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// RemoveAlias sends DELETE to the cloud-api alias method.
func (c DirectoryClient) RemoveAlias(userID, alias string) (*model.RmAliasResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/api360/directory/org/%d/users/%s/aliases/%s", c.cloudBase(), c.OrgID, pathSegment(userID), pathSegment(alias))
	resp, err := c.do(http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}
	var data model.RmAliasResponse
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (c DirectoryClient) do(method, endpoint string, payload []byte) (*ApiResponse, error) {
	if c.Token == "" {
		return nil, ErrTokenMissing
	}
	if c.OrgID == 0 {
		return nil, errors.New("organization id must be specified")
	}
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "OAuth "+c.Token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}
	client := c.HTTP
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	apiResp := &ApiResponse{HttpCode: resp.StatusCode, Body: data}
	if err := GetErrorText(apiResp); err != nil {
		return nil, err
	}
	return apiResp, nil
}

func pathSegment(value string) string {
	return strings.ReplaceAll(url.PathEscape(value), "@", "%40")
}

func putString(body map[string]any, key, value string) {
	if value = strings.TrimSpace(value); value != "" {
		body[key] = value
	}
}

func nameBody(first, last, middle string) map[string]string {
	name := map[string]string{}
	if first != "" {
		name["first"] = first
	}
	if last != "" {
		name["last"] = last
	}
	if middle != "" {
		name["middle"] = middle
	}
	return name
}

func decodeDirectoryUser(body []byte) (model.User, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		return model.User{}, err
	}
	if raw, ok := fields["id"]; ok && len(bytes.TrimSpace(raw)) > 0 && raw[0] != '"' {
		encoded, err := json.Marshal(rawID(raw))
		if err != nil {
			return model.User{}, err
		}
		fields["id"] = encoded
		var err2 error
		body, err2 = json.Marshal(fields)
		if err2 != nil {
			return model.User{}, err2
		}
	}
	var user model.User
	if err := json.Unmarshal(body, &user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func rawID(raw json.RawMessage) string {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	if raw[0] == '"' {
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return ""
		}
		return value
	}
	return string(raw)
}
