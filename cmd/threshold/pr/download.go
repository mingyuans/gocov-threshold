package pr

import (
	"fmt"
	"github.com/mingyuans/gocov-threshold/cmd/threshold/log"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"time"
)

func (srv Service) DownloadDiff() ([]byte, error) {
	diffURL := fmt.Sprintf("%s/repos/%s/%s/pulls/%d",
		srv.env.APIServerURL,
		srv.env.RepositoryOwner,
		srv.env.RepositoryName,
		srv.info.Number)

	log.Get().Debug("Downloading diff from URL", zap.String("diffURL", diffURL))

	// Create request
	req, err := http.NewRequest("GET", diffURL, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication headers
	req.Header.Set("Authorization", "Bearer "+srv.inputArg.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.diff")
	req.Header.Set("User-Agent", "GitHub-Action-PR-Diff-Downloader")

	// Send request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to download diff: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("diff download failed with status %d", resp.StatusCode)
	}

	// Read response content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read diff content: %w", err)
	}

	log.Get().Debug("Downloading diff from URL", zap.String("diffContent", string(content)))

	return content, nil
}

func saveDiffToFile(content []byte, filename string) error {
	return os.WriteFile(filename, content, 0644)
}

func (srv Service) DownloadAndSaveDiff(filename string) error {
	diffContent, err := srv.DownloadDiff()
	if err != nil {
		return fmt.Errorf("failed to download diff: %w", err)
	}

	if err := saveDiffToFile(diffContent, filename); err != nil {
		return fmt.Errorf("failed to save diff to file: %w", err)
	}

	log.Get().Info("Diff downloaded and saved successfully", zap.String("filename", filename))
	return nil
}
