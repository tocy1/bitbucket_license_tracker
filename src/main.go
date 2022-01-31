package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	Auth       *auth
	apiBaseURL *url.URL
	HttpClient *http.Client
}

type auth struct {
	user, password string
	bearerToken    string
}

var (
	bitbucketUrl     = os.Getenv("BITBUCKET_URL")
	LicenseEndpoint  = "/rest/api/1.0/admin/license"
	OAuthbearerToken = os.Getenv("BEARER_TOKEN")
	UserName         = os.Getenv("USER_NAME")
	Password         = os.Getenv("PASSWORD")
	IdleConnTimeout  = getEnv("IDLE_CONNECTION_TIMEOUT", "2s")
	MaxConnsPerHost  = getEnvInt("MAX_CONNECTION_PER_HOST", "2")
	MaxIdleConns     = getEnvInt("MAX_IDLE_CONNECTIONS", "10")
)

func (c *Client) getAvailableLicenseCount(req *http.Request) ([]byte, error) {
	u, err := url.Parse(LicenseEndpoint)
	if err != nil {
		return nil, errors.New("Unable to parse Lincense URI endpoint")
	}
	req.Method = "GET"

	url := c.apiBaseURL.ResolveReference(u)
	req.URL = url
	req.Host = ""
	resp, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func main() {
	c, err := newClient()
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		log.Fatalln(err)
	}
	c.authenticateRequest(req)
	if err != nil {
		log.Fatalln(err)
	}
	// ...
	resp, err := c.getAvailableLicenseCount(req)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	body := string(resp)
	log.Printf(body)

}

func newClient() (*Client, error) {
	d, err := time.ParseDuration(IdleConnTimeout)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		MaxIdleConns:       MaxIdleConns,
		IdleConnTimeout:    d * time.Second,
		DisableCompression: true,
		MaxConnsPerHost:    MaxConnsPerHost,
	}
	client := &http.Client{Transport: tr}
	if bitbucketUrl == "" {
		return nil, errors.New("BITBUCKET_URL env not set")
	}
	base, err := url.Parse(bitbucketUrl)
	if err != nil {
		return nil, err
	}
	auth, err := newAuth()

	if err != nil {
		return nil, err
	}
	return &Client{
		apiBaseURL: base,
		Auth:       auth,
		HttpClient: client,
	}, nil
}

func newAuth() (*auth, error) {
	if OAuthbearerToken != "" {
		return NewOAuthbearerToken(OAuthbearerToken), nil
	}
	if UserName != "" && Password != "" {
		return NewBasicAuth(UserName, Password), nil
	}
	return nil, errors.New("Unable to unathenticate, Please make sure  either environmental variable BEARER_TOKEN is set or USER_NAME and  PASSWORD is set")

}

func NewOAuthbearerToken(token string) *auth {
	return &auth{bearerToken: token}
}
func NewBasicAuth(username, password string) *auth {
	return &auth{user: username, password: password}
}

func (c *Client) authenticateRequest(req *http.Request) {
	if c.Auth.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.Auth.bearerToken)

	} else if c.Auth.user != "" && c.Auth.password != "" {
		req.SetBasicAuth(c.Auth.user, c.Auth.password)
	}
	return
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
		log.Output(1, key+" env variable  does not exist, using the default value "+value)
	}
	return value
}

func getEnvInt(key, fallback string) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
		log.Output(1, key+" env variable does not exist, using the default value "+value)
	}
	val, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	return val
}
