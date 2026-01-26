package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"inkdown-cli/utils"
)

type GitHubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

type PullRequestResponse struct {
	HTMLURL string `json:"html_url"`
	Number  int    `json:"number"`
}

type UpdateFilePayload struct {
	Message string `json:"message"`
	Content string `json:"content"`
	SHA     string `json:"sha"`
	Branch  string `json:"branch"`
}

type PullRequestPayload struct {
	Title string `json:"title"`
	Head  string `json:"head"`
	Base  string `json:"base"`
	Body  string `json:"body"`
}

const (
	ORG_NAME  = "inkdown"
	REPO_NAME = "inkdown-community"
	BRANCH    = "main"
)

func ForkRepo(token string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/forks", ORG_NAME, REPO_NAME)

	body := map[string]string{}
	payload, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "community-cli")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro ao criar fork: %s", string(body))
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	fullName, ok := res["full_name"].(string)
	if !ok {
		return "", errors.New("não conseguiu pegar full_name do fork")
	}

	return fullName, nil // ex: "usuario/community-plugins"
}

func GetBranchSHA(token string, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/git/refs/heads/%s", repo, BRANCH)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "community-cli")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// If status is not OK, return body to help debugging
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("erro ao pegar SHA: status %d, body: %s", resp.StatusCode, string(body))
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta ao pegar SHA: %v, body: %s", err, string(body))
	}

	object, ok := res["object"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("erro ao pegar SHA: objeto 'object' não encontrado na resposta: %s", string(body))
	}
	sha, ok := object["sha"].(string)
	if !ok {
		return "", fmt.Errorf("SHA não encontrado na resposta: %s", string(body))
	}
	return sha, nil
}

func CreateBranch(token string, repo string, newBranch string, baseSHA string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/git/refs", repo)
	utils.Info("Using the url: %s", url)

	body := map[string]string{
		"ref": "refs/heads/" + newBranch,
		"sha": baseSHA,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "community-cli")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro criando branch: %s", string(b))
	}

	return nil
}

func GetFileContent(token string, owner string, branch string, path string) (string, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=%s", owner, path, branch)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "community-cli")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	contentBase64, ok := res["content"].(string)
	if !ok {
		return "", "", errors.New("conteudo não encontrado")
	}
	sha, ok := res["sha"].(string)
	if !ok {
		return "", "", errors.New("SHA não encontrado")
	}

	decoded, _ := utils.DecodeBase64(contentBase64) // função helper para Base64
	return decoded, sha, nil
}

func UpdateFile(token string, owner string, branch string, path string, newContent, sha string, message string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", owner, path)
	contentB64 := utils.EncodeBase64(newContent)

	body := map[string]string{
		"message": message,
		"content": contentB64,
		"branch":  branch,
		"sha":     sha,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "community-cli")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro atualizando arquivo: %s", string(b))
	}
	return nil
}

func CreatePR(token string, headBranch string, title string, body string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", ORG_NAME+"/"+REPO_NAME)
	payload := map[string]string{
		"title": title,
		"head":  headBranch,
		"base":  BRANCH,
		"body":  body,
	}
	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "community-cli")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro criando PR: %s", string(b))
	}

	var prResp PullRequestResponse

	if err := json.NewDecoder(resp.Body).Decode(&prResp); err != nil {
		return "", err
	}

	return prResp.HTMLURL, nil
}

func GetGitHubUsername(token string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "community-cli")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("erro ao buscar usuário: %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}

	return user.Login, nil
}

func AppendPlugin(existing string, pluginJSON string) string {
	existing = strings.TrimSpace(existing)
	existing = strings.TrimSuffix(existing, "]")

	if !strings.HasSuffix(existing, "[") {
		existing += ",\n"
	}

	existing += pluginJSON + "\n]"
	return existing
}
