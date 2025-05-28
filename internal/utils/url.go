package utils

import (
	"fmt"
	"net/url"
	"strings"
)

type GitURLParts struct {
	Host  string
	Owner string
	Repo  string
}

func ParseGitURL(rawurl string) (GitURLParts, error) {
	var u GitURLParts
	urlStr := rawurl
	if strings.HasSuffix(urlStr, ".git") {
		urlStr = strings.TrimSuffix(urlStr, ".git")
	}
	if strings.HasPrefix(urlStr, "git@") {
		parts := strings.SplitN(urlStr, ":", 2)
		if len(parts) != 2 {
			return u, fmt.Errorf("invalid ssh url: %s", rawurl)
		}
		host := strings.TrimPrefix(parts[0], "git@")
		pathParts := strings.Split(parts[1], "/")
		if len(pathParts) < 2 {
			return u, fmt.Errorf("invalid ssh url path: %s", rawurl)
		}
		u.Host = host
		u.Owner = pathParts[0]
		u.Repo = pathParts[1]
		return u, nil
	}
	if strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://") {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return u, err
		}
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(parts) < 2 {
			return u, fmt.Errorf("invalid url path: %s", rawurl)
		}
		u.Host = parsed.Host
		u.Owner = parts[0]
		u.Repo = parts[1]
		return u, nil
	}
	return u, fmt.Errorf("unsupported git url format: %s", rawurl)
}
