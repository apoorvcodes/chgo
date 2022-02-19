package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zshbunni/chgo/query"
	"github.com/zshbunni/chgo/types"
	"gopkg.in/yaml.v2"
)

var (
	config   types.Config
	helpText = "Failed to login. Try logging in."
	baseURL  = "https://coursehunter.net"
)

// setConfig sets the config global object using the config file if present
func setConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	file := filepath.Join(homeDir, ".config", "chgo", "config.yaml")
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config)
	return err
}

// isTokenExpires verifies whether the token in the config is expired or not
func isTokenExpired() error {
	tokens := strings.Split(config.AccessToken, ".")

	if len(tokens) != 3 {
		return fmt.Errorf("received malformed access token")
	}

	var token struct {
		Exp int64 `json:"exp"`
	}

	jwt, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		return err
	}

	err = json.Unmarshal(jwt, &token)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	if now > token.Exp {
		return fmt.Errorf("token expired. try logging in")
	}

	return nil
}

// createConfig creats the config in `.config/chgo/config.yaml` location
// if not present or overrides the previous config
func createConfig(config types.Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	configPath := filepath.Join(homeDir, ".config", "chgo")
	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return nil
	}

	file, err := os.Create(filepath.Join(configPath, "config.yaml"))
	if err != nil {
		return nil
	}
	defer file.Close()

	_, err = file.Write(data)

	return err
}

// login logins using the email and password and updates the config file
func login(email, password string) (types.Config, error) {
	var config types.Config

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

func searchCourses(title string) ([]types.Course, error) {
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

	courses, err := query.ExtractCourses(resp.Body)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func main() {
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	email := loginCmd.String("u", "", "email")
	password := loginCmd.String("p", "", "password")

	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	title := searchCmd.String("title", "", "course title")

	if len(os.Args) < 2 {
		fmt.Println(`expected 'login' or 'search' subcommand`)
		return
	}

	switch os.Args[1] {
	case "login":
		loginCmd.Parse(os.Args[2:])

		if len(*email) == 0 || len(*password) == 0 {
			fmt.Println("Missing credentials")
			return
		}

		c, err := login(*email, *password)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = createConfig(c)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Logged in")
	case "search":
		searchCmd.Parse(os.Args[2:])

		if err := setConfig(); err != nil {
			fmt.Println(err)
			return
		}

		if err := isTokenExpired(); err != nil {
			fmt.Println(err)
			return
		}

		_, err := searchCourses(*title)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println(`expected 'login' or 'search' subcommand`)
		return
	}
}
