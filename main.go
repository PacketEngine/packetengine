package packetengine

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/go-resty/resty/v2"
)

type PacketEngineClient struct {
	apiKey string
	client *resty.Client
}

var apiURL = "https://api.packetengine.co.uk"

func NewPacketEngineClient(apiKey string) (*PacketEngineClient, error) {
	if len(apiKey) == 0 {
		return nil, errors.New("API key is required.")
	}

	apiURLEnv := os.Getenv("PACKETENGINE_API_URL")

	if apiURLEnv != "" {
		apiURL = apiURLEnv
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(string(apiKey)).
		Get(apiURL + "/v1/verify")

	if resp.StatusCode() == 401 {
		return nil, errors.New("The provided API key is invalid.")
	}

	if err != nil {
		return nil, err
	}

	return &PacketEngineClient{apiKey: apiKey, client: client}, nil
}

func (c *PacketEngineClient) GetSubdomains(domain string, withoutTags string, allSubdomains bool) ([]string, error) {
	if len(domain) == 0 {
		return nil, errors.New("A domain is required.")
	}

	clientResp := c.client.R()

	if len(withoutTags) > 0 {
		clientResp.SetQueryParams(map[string]string{
			"withoutTags": withoutTags,
		})
	}

	if allSubdomains == true {
		clientResp.SetQueryParams(map[string]string{
			"all": "true",
		})
	}

	resp, err := clientResp.
		SetHeader("Accept", "application/json").
		SetAuthToken(string(c.apiKey)).
		Get(apiURL + "/v1/domains/" + domain + "/subdomains")

	if resp.StatusCode() == 401 {
		return nil, errors.New("The API key is invalid.")
	}

	if err != nil {
		return nil, err
	}

	// Parse the response body into a map[string]interface{}
	var result map[string]interface{}

	err = json.Unmarshal(resp.Body(), &result)

	// Check for any errors
	if err != nil {
		return nil, err
	}

	if errorMessage, hasError := result["error"]; hasError {
		return nil, errors.New(errorMessage.(string))
	}

	var subdomains []string

	if subdomainsInterface, ok := result["subdomains"]; ok {
		for _, subdomain := range subdomainsInterface.([]interface{}) {
			subdomains = append(subdomains, subdomain.(string))
		}
	}

	return subdomains, nil
}
