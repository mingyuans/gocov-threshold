package pr

import (
	"encoding/json"
	"fmt"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/arg"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func getPREnvironment() Environment {
	repository := os.Getenv("GITHUB_REPOSITORY")

	var repositoryName = ""
	parts := strings.Split(repository, "/")
	if len(parts) == 2 {
		repositoryName = parts[1]
	}

	return Environment{
		Repository:      repository,
		RepositoryName:  repositoryName,
		EventName:       os.Getenv("GITHUB_EVENT_NAME"),
		EventPath:       os.Getenv("GITHUB_EVENT_PATH"),
		RefName:         os.Getenv("GITHUB_REF_NAME"),
		RepositoryOwner: os.Getenv("GITHUB_REPOSITORY_OWNER"),
		SHA:             os.Getenv("GITHUB_SHA"),
		Actor:           os.Getenv("GITHUB_ACTOR"),
		ServerURL:       getEnvWithDefault("GITHUB_SERVER_URL", "https://github.com"),
		APIServerURL:    getEnvWithDefault("GITHUB_API_URL", "https://api.github.com"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (srv Service) GetEnvironment() Environment {
	return srv.env
}

type Environment struct {
	Repository      string // GITHUB_REPOSITORY
	RepositoryName  string
	EventName       string // GITHUB_EVENT_NAME
	EventPath       string // GITHUB_EVENT_PATH
	RepositoryOwner string // GITHUB_REPOSITORY_OWNER
	RefName         string // GITHUB_REF_NAME
	SHA             string // GITHUB_SHA
	Actor           string // GITHUB_ACTOR
	ServerURL       string // GITHUB_SERVER_URL
	APIServerURL    string // GITHUB_API_URL
}

type GitHubPRInfo struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	HtmlURL  string `json:"html_url"`
	DiffURL  string `json:"diff_url"`
	PatchURL string `json:"patch_url"`
	Head     struct {
		SHA string `json:"sha"`
		Ref string `json:"ref"`
	} `json:"head"`
	Base struct {
		SHA string `json:"sha"`
		Ref string `json:"ref"`
	} `json:"base"`
	State string `json:"state"`
}

type PullRequestEvent struct {
	Number      int          `json:"number"`
	PullRequest GitHubPRInfo `json:"pull_request"`
	Repository  struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Owner    struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
	Action string `json:"action"`
}

func gettPRInfo(env Environment, arg arg.Arg) (GitHubPRInfo, error) {
	if env.EventPath != "" {
		if prInfo, err := getPRInfoFromEventFile(env.EventPath); err == nil {
			return prInfo, nil
		}
		log.Get().Debug("⚠️  Failed to read from event file, trying API...")
	}

	return getPRInfoFromAPI(env, arg)
}

func (srv Service) GetPRInfo() GitHubPRInfo {
	return srv.info
}

func getPRInfoFromEventFile(eventPath string) (GitHubPRInfo, error) {
	prInfo := GitHubPRInfo{}
	data, err := os.ReadFile(eventPath)
	if err != nil {
		return prInfo, fmt.Errorf("failed to read event file: %w", err)
	}

	var event PullRequestEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return prInfo, fmt.Errorf("failed to parse event JSON: %w", err)
	}

	return event.PullRequest, nil
}

func getPRInfoFromAPI(env Environment, arg arg.Arg) (GitHubPRInfo, error) {
	prInfo := GitHubPRInfo{}
	//The for: refs/pull/{pr_number}/merge
	prNumber, err := extractPRNumber(env.RefName, env.SHA)
	if err != nil {
		return prInfo, fmt.Errorf("failed to extract PR number: %w", err)
	}

	apiURL := fmt.Sprintf("%s/repos/%s/pulls/%d", env.APIServerURL, env.Repository, prNumber)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return prInfo, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+arg.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "mingyuans/gocov-threshold")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return prInfo, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return prInfo, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// 解析响应
	if decodeErr := json.NewDecoder(resp.Body).Decode(&prInfo); decodeErr != nil {
		return prInfo, fmt.Errorf("failed to decode response: %w", decodeErr)
	}

	return prInfo, nil
}

func extractPRNumber(refName, sha string) (int, error) {
	if strings.HasPrefix(refName, "refs/pull/") {
		parts := strings.Split(refName, "/")
		if len(parts) >= 3 {
			return strconv.Atoi(parts[2])
		}
	}

	if prNum, err := strconv.Atoi(refName); err == nil {
		return prNum, nil
	}

	return 0, fmt.Errorf("unable to extract PR number from ref: %s", refName)
}
