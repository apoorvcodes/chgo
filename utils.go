package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

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

// playLesson plays the selected lesson using `mpv`
func playLesson(url string) {
  player, err := exec.LookPath("mpv")

  // TODO: notify the user about missing dependencies
  if err != nil {
    return
  }

  cmd := exec.Command(player, url)
  cmd.Start()
}
