package packetengine

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/go-resty/resty/v2"
)

type PacketEngineClient struct {
	apiToken string
	client   *resty.Client
}

var apiURL = "https://api.packetengine.co.uk"

func NewPacketEngineClient(apiToken string) (*PacketEngineClient, error) {
	if len(apiToken) == 0 {
		return nil, errors.New("API token is required.")
	}

	apiURLEnv := os.Getenv("PACKETENGINE_API_URL")

	if apiURLEnv != "" {
		apiURL = apiURLEnv
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(string(apiToken)).
		Get(apiURL + "/v1/verify")

	if resp.StatusCode() == 401 {
		return nil, errors.New("The provided API token is invalid.")
	}

	if err != nil {
		return nil, err
	}

	return &PacketEngineClient{apiToken: apiToken, client: client}, nil
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

	if allSubdomains {
		clientResp.SetQueryParams(map[string]string{
			"all": "true",
		})
	}

	resp, err := clientResp.
		SetHeader("Accept", "application/json").
		SetAuthToken(string(c.apiToken)).
		Get(apiURL + "/v1/domains/" + domain + "/subdomains")

	if resp.StatusCode() == 401 {
		return nil, errors.New("The API token is invalid.")
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

func (c *PacketEngineClient) GetIPs(domain string, withoutTags string) ([]string, error) {
	if len(domain) == 0 {
		return nil, errors.New("A domain is required.")
	}

	clientResp := c.client.R()

	if len(withoutTags) > 0 {
		clientResp.SetQueryParams(map[string]string{
			"withoutTags": withoutTags,
		})
	}

	resp, err := clientResp.
		SetHeader("Accept", "application/json").
		SetAuthToken(string(c.apiToken)).
		Get(apiURL + "/v1/domains/" + domain + "/ips")

	if resp.StatusCode() == 401 {
		return nil, errors.New("The API token is invalid.")
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

	var ips []string

	if ipsInterface, ok := result["ips"]; ok {
		for _, ip := range ipsInterface.([]interface{}) {
			ips = append(ips, ip.(string))
		}
	}

	return ips, nil
}
