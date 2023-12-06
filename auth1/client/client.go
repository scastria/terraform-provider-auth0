package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	FormEncoded = "application/x-www-form-urlencoded"
	Bearer      = "Bearer"
	// KonnectDomain                = "api.konghq.com"
	// GlobalRegion                 = "global"
	// IdSeparator                  = ":"
	// FilterName                   = "filter[name]"
	// FilterNameContains           = "filter[name][contains]"
	// FilterFullName               = "filter[full_name]"
	// FilterFullNameContains       = "filter[full_name][contains]"
	// FilterEmail                  = "filter[email]"
	// FilterEmailContains          = "filter[email][contains]"
	// FilterActive                 = "filter[active]"
	// FilterRoleName               = "filter[role_name]"
	// FilterRoleNameContains       = "filter[role_name][contains]"
	// FilterEntityTypeName         = "filter[entity_type_name]"
	// FilterEntityTypeNameContains = "filter[entity_type_name][contains]"
)

type Client struct {
	domain       string
	clientId     string
	clientSecret string
	accessToken  string
	httpClient   *http.Client
}

func NewClient(ctx context.Context, domain string, clientId string, clientSecret string) (*Client, error) {
	c := &Client{
		domain:       domain,
		clientId:     clientId,
		clientSecret: clientSecret,
		httpClient:   &http.Client{},
	}
	tflog.Info(ctx, "Auth0 Management API: Obtaining access token...")
	requestURL := fmt.Sprintf("https://%s/%s", c.domain, OauthTokenPath)
	requestForm := url.Values{
		"grant_type":    []string{"client_credentials"},
		"audience":      []string{fmt.Sprintf("https://%s/api/v2/", c.domain)},
		"client_id":     []string{c.clientId},
		"client_secret": []string{c.clientSecret},
	}
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBufferString(requestForm.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set(headers.ContentType, FormEncoded)
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		tflog.Info(ctx, "Auth0 Management API:", map[string]interface{}{"error": err})
	} else {
		tflog.Info(ctx, "Auth0 Management API: ", map[string]interface{}{"request": string(requestDump)})
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, &RequestError{StatusCode: resp.StatusCode, Err: err}
		}
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: fmt.Errorf("%s", respBody.String())}
	}
	//Parse body to extract access_token
	token := &OauthToken{}
	err = json.NewDecoder(resp.Body).Decode(token)
	if err != nil {
		return nil, err
	}
	tflog.Info(ctx, "Auth0 Management API: Received access token: ", map[string]interface{}{"token": token.AccessToken})
	//Inject token as access_token for client for all future calls
	c.accessToken = token.AccessToken
	return c, nil
}

func (c *Client) HttpRequest(ctx context.Context, method string, path string, query url.Values, headerMap http.Header, body *bytes.Buffer) (*bytes.Buffer, error) {
	req, err := http.NewRequest(method, c.RequestPath(path), body)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	//Handle query values
	if query != nil {
		requestQuery := req.URL.Query()
		for key, values := range query {
			for _, value := range values {
				requestQuery.Add(key, value)
			}
		}
		req.URL.RawQuery = requestQuery.Encode()
	}
	//Handle header values
	if headerMap != nil {
		for key, values := range headerMap {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
	//Handle authentication
	req.Header.Set(headers.Authorization, Bearer+" "+c.accessToken)
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		tflog.Info(ctx, "Konnect API:", map[string]any{"error": err})
	} else {
		tflog.Info(ctx, "Konnect API: ", map[string]any{"request": string(requestDump)})
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	defer resp.Body.Close()
	respBody := new(bytes.Buffer)
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: err}
	}
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: fmt.Errorf("%s", respBody.String())}
	}
	return respBody, nil
}

func (c *Client) RequestPath(path string) string {
	return fmt.Sprintf("https://%s/api/v2/%s", c.domain, path)
}
