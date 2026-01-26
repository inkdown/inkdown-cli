package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Release struct {
	ID        int    `json:"id"`
	TagName   string `json:"tag_name"`
	Name      string `json:"name"`
	UploadURL string `json:"upload_url"` // "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets{?name,label}"
}

// GetReleases fetches all releases for a repository
func GetReleases(token, owner, repo string) ([]Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get releases: %s", resp.Status)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	return releases, nil
}

// GetReleaseByTag fetches a specific release by tag
func GetReleaseByTag(token, owner, repo, tag string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil // Not found
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get release by tag: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// CreateRelease creates a new release
func CreateRelease(token, owner, repo, tag, name, body string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	payload := map[string]interface{}{
		"tag_name": tag,
		"name":     name,
		"body":     body,
		"draft":    false,
	}
	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create release: %s", string(b))
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// DeleteRelease deletes a release by ID
func DeleteRelease(token, owner, repo string, id int) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%d", owner, repo, id)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("failed to delete release: %s", resp.Status)
	}

	return nil
}

// DeleteTag deletes a tag reference
func DeleteTag(token, owner, repo, tag string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/tags/%s", owner, repo, tag)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 && resp.StatusCode != 422 { // 422 usually means tag doesn't exist as ref but might have been deleted with release
		return fmt.Errorf("failed to delete tag: %s", resp.Status)
	}

	return nil
}

// UploadReleaseAsset uploads a file to a release
func UploadReleaseAsset(token, uploadUrlTemplate, filePath, contentType string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	fileName := filepath.Base(filePath)
	
	// Upload URL comes as "https://.../assets{?name,label}", we need to remove the template part
	// and add ?name=filename
	cleanUrl := uploadUrlTemplate
	if idx := len(cleanUrl) - 13; idx > 0 && cleanUrl[idx:] == "{?name,label}" {
         cleanUrl = cleanUrl[:idx]
    }
	
	url := fmt.Sprintf("%s?name=%s", cleanUrl, fileName)

	// GitHub API for uploads requires raw binary body, but with correct content-length and type
	stat, err := f.Stat()
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", url, f)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", contentType)
	req.ContentLength = stat.Size() // Set the length explicitly on the request struct
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload asset %s: %s (%s)", fileName, resp.Status, string(b))
	}

	return nil
}
