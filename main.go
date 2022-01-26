package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

var (
	bitbucketUrl     = os.Getenv("BITBUCKET_URL")
	LincenseEndpoint = "/rest/api/1.0/admin/license"
	//IdleConnTimeout  = os.Getenv("IDLE_CONNECTION_TIMEOUT")
)

func getLicenseCount(c *http.Client) ([]byte, error) {
	if *&bitbucketUrl == "" {
		return nil, errors.New("BITBUCKET_URL env not set")
	}
	url := *&bitbucketUrl + *&LincenseEndpoint
	resp, err := c.Get(url)
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
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		MaxConnsPerHost:    1,
	}
	client := &http.Client{Transport: tr}
	resp, err := getLicenseCount(client)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	body := string(resp)
	log.Printf(body)

}
