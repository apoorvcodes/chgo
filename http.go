package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// login logins using the email and password and updates the config file
func login(email, password string) (Config, error) {
	var config Config

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	payload := url.Values{}
	payload.Set("e_mail", email)
	payload.Set("password", password)

	req, err := http.NewRequest("POST", baseURL+"/sign-in", strings.NewReader(payload.Encode()))

	if err != nil {
		return config, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return config, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 302 {
		return config, fmt.Errorf("Failed to login. Make sure you entered valid credentials.")
	}

	config.Locale = "en"
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "user_ident" {
			config.UserIdent = cookie.Value
		}

		if cookie.Name == "accessToken" {
			config.AccessToken = cookie.Value
		}
	}

	return config, nil
}

func prepareRequest(endpoint string) (*http.Request, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{
		Name:  "accessToken",
		Value: config.AccessToken,
	})

	req.AddCookie(&http.Cookie{
		Name:  "user_ident",
		Value: config.UserIdent,
	})

	req.AddCookie(&http.Cookie{
		Name:  "locale",
		Value: config.Locale,
	})

	return req, nil
}

func searchCourses(title string) ([]Course, error) {
	var client http.Client

	endpoint := baseURL + "/search?q=" + title
	req, err := prepareRequest(endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	courses, err := ExtractCourses(resp.Body)
	if err != nil {
		return nil, err
	}

	return courses, nil
}
