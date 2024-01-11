package xrcm_pf

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// SignIn - Get a new token for user
func (c *Client) SignIn() (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}
	rb, err := json.Marshal(c.Auth)
	if err != nil {
		return nil, err
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/signin", c.HostURL), strings.NewReader(string(rb)))
	// fmt.Println("Connecting server: ", fmt.Sprintf("%s/oauth/token?client_id=test&audience=test", c.HostURL))

	// TODO - Need to data drive config with self signed TLS vs valid certificates

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // Self signed certificate flag, TODO to remove or data drive it latter
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/oauth/token?client_id=test&audience=test", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	str1 := fmt.Sprintf("%s", body)
	// fmt.Println("body auth", str1)
	ar := AuthResponse{}
	var respmap map[string]string
	err = json.Unmarshal([]byte(str1), &respmap)
	ar.Token = "Bearer " + respmap["access_token"]
	// fmt.Println("token", ar.Token)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

// SignOut - Revoke the token for a user
func (c *Client) SignOut() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signout", c.HostURL), strings.NewReader(string("")))
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if string(body) != "Signed out user" {
		return errors.New(string(body))
	}

	return nil
}
