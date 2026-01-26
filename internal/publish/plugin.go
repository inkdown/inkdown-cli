package publish

import (
	"bufio"
	"encoding/json"
	"fmt"
	"inkdown-cli/config"
	"inkdown-cli/internal/github"
	"inkdown-cli/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type PackageJSON interface {
	GetName() string
	GetVersion() string
	GetDescription() string
}

type Package struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Main        string `json:"main"`
}

func (p Package) GetName() string {
	return p.Name
}

func (p Package) GetVersion() string {
	return p.Version
}

func (p Package) GetDescription() string {
	return p.Description
}

func PublishPlugin(dir *string) (string, error) {
	utils.Info("Started the publish process...")

	manifestPath := filepath.Join(*dir, "manifest.json")
	rawManifest, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", fmt.Errorf("could not read manifest.json: %v (make sure you are in the plugin root)", err)
	}

	var manifest Package
	if err := json.Unmarshal(rawManifest, &manifest); err != nil {
		return "", fmt.Errorf("invalid manifest.json: %v", err)
	}

	if manifest.Version == "" {
		return "", fmt.Errorf("manifest.json missing 'version'")
	}
	if manifest.Name == "" {
		return "", fmt.Errorf("manifest.json missing 'name'")
	}

	utils.Info("Detected Plugin: %s v%s", manifest.Name, manifest.Version)

	// 1.5 Build Plugin
	pkgJsonPath := filepath.Join(*dir, "package.json")
	if _, err := os.Stat(pkgJsonPath); err == nil {
		utils.Info("Building plugin...")
		
		// npm install
		utils.Info("Running 'npm install'...")
		installCmd := exec.Command("bun", "install")
		installCmd.Dir = *dir
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return "", fmt.Errorf("failed to run 'npm install': %v", err)
		}

		// npm run build
		utils.Info("Running 'npm run build'...")
		buildCmd := exec.Command("bun", "run", "build")
		buildCmd.Dir = *dir
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		if err := buildCmd.Run(); err != nil {
			return "", fmt.Errorf("failed to run 'npm run build': %v", err)
		}

		// Verify main.js
		mainJsPath := filepath.Join(*dir, "main.js")
		if _, err := os.Stat(mainJsPath); os.IsNotExist(err) {
			return "", fmt.Errorf("build completed but 'main.js' was not found")
		}
		utils.Success("Build successful!")
	} else {
		utils.Warn("No package.json found. Skipping build step (expecting pre-built assets).")
	}

	env := config.LoadEnv()
	token, err := github.LoadToken()
	if err == nil {
		if err := github.ValidateToken(token); err == nil {
			fmt.Println("Using saved GitHub token")
		} else {
			token = ""
		}
	}

	if token == "" {
		code, err := github.RequestDeviceCode(env.ClientID)
		if err != nil {
			return "", err
		}

		utils.Info("To authorize this application, open: %s", code.VerificationURI)
		utils.Info("And enter the code: %s", code.UserCode)

		go func() {
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "darwin":
				cmd = exec.Command("open", code.VerificationURI)
			case "windows":
				cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", code.VerificationURI)
			default:
				cmd = exec.Command("xdg-open", code.VerificationURI)
			}
			_ = cmd.Start()
		}()

		token, err = github.PollForToken(env.ClientID, code.DeviceCode, code.Interval)
		if err != nil {
			return "", err
		}
		_ = github.SaveToken(token)
	}

	if err := github.ValidateToken(token); err != nil {
		return "", err
	}

	// 3. User & Repo Info
	// Get Repo Info from Git Config
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = *dir
	out, err := cmd.Output()
	var repoOwner, repoName string

	if err == nil {
		remoteURL := strings.TrimSpace(string(out))
		// Handle SSH: git@github.com:owner/repo.git
		// Handle HTTPS: https://github.com/owner/repo.git
		
		remoteURL = strings.TrimSuffix(remoteURL, ".git")
		
		parts := strings.Split(remoteURL, "/")
		if len(parts) >= 2 {
			repoName = parts[len(parts)-1]
			repoOwner = parts[len(parts)-2]
			
			// Handle potential SSH prefix in owner (git@github.com:owner)
			if strings.Contains(repoOwner, ":") {
				ownerParts := strings.Split(repoOwner, ":")
				repoOwner = ownerParts[len(ownerParts)-1]
			}
		}
	}

	// Fallback or validation
	if repoName == "" || repoOwner == "" {
		// Try to use authenticated user and guessed name as fallback, but warn
		var err error
		repoOwner, err = github.GetGitHubUsername(token)
		if err != nil {
			return "", err
		}
		repoName = strings.ReplaceAll(manifest.Name, " ", "-")
		utils.Warn("Could not detect git remote. defaulting to %s/%s", repoOwner, repoName)
	}

	username := repoOwner
	userRepoName := repoName

	tagName := "v" + manifest.Version
	if strings.HasPrefix(manifest.Version, "v") {
		tagName = manifest.Version
	}

	utils.Info("Checking for existing release %s in %s/%s...", tagName, username, userRepoName)
	existingRelease, err := github.GetReleaseByTag(token, username, userRepoName, tagName)

	if existingRelease != nil {
		utils.Warn("Release %s already exists!", tagName)
		reader := bufio.NewReader(os.Stdin)
		utils.Prompt("Do you want to overwrite it? ALL ASSETS WILL BE REPLACED. (y/N): ")

		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer != "y" && answer != "yes" {
			utils.Info("Aborting.")
			return "", nil
		}

		utils.Info("Deleting old release...")
		if err := github.DeleteRelease(token, username, userRepoName, existingRelease.ID); err != nil {
			return "", fmt.Errorf("failed to delete release: %v", err)
		}
		_ = github.DeleteTag(token, username, userRepoName, tagName)
	}

	utils.Info("Creating release %s...", tagName)
	newRelease, err := github.CreateRelease(
		token,
		username,
		userRepoName,
		tagName,
		fmt.Sprintf("%s %s", manifest.Name, manifest.Version),
		manifest.Description,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create release: %v", err)
	}

	assets := []string{"main.js", "manifest.json"}

	if _, err := os.Stat(filepath.Join(*dir, "styles.css")); err == nil {
		assets = append(assets, "styles.css")
	}

	for _, asset := range assets {
		assetPath := filepath.Join(*dir, asset)
		utils.Info("Uploading %s...", asset)

		contentType := "application/javascript"
		if strings.HasSuffix(asset, ".json") {
			contentType = "application/json"
		} else if strings.HasSuffix(asset, ".css") {
			contentType = "text/css"
		}

		if err := github.UploadReleaseAsset(token, newRelease.UploadURL, assetPath, contentType); err != nil {
			return "", fmt.Errorf("failed to upload %s: %v", asset, err)
		}
	}

	utils.Success("Release published successfully!")

	utils.Info("Proceeding to update Community Registry...")

	if _, err := github.ForkRepo(token); err != nil {
		return "", err
	}

	branch := fmt.Sprintf("add-plugin/%s", userRepoName)
	communityRepo := fmt.Sprintf("%s/inkdown-community", username)

	sha, err := github.GetBranchSHA(token, communityRepo)
	if err != nil {
		return "", err
	}

	_ = github.CreateBranch(token, communityRepo, branch, sha)

	content, contentSha, err := github.GetFileContent(token, communityRepo, branch, "plugins.json")
	if err != nil {
		return "", err
	}

	pluginEntry := fmt.Sprintf(`{
  "id": "%s",
  "name": "%s",
  "author": "%s",
  "version": "%s",
  "description": "%s",
  "repo": "%s/%s"
}`,
		userRepoName,
		manifest.Name,
		username,
		manifest.Version,
		manifest.Description,
		username,
		userRepoName,
	)

	utils.Info("Updating plugins.json...")
	updated := github.AppendPlugin(content, pluginEntry)

	if err := github.UpdateFile(
		token,
		communityRepo,
		branch,
		"plugins.json",
		updated,
		contentSha,
		fmt.Sprintf("feat: add plugin %s v%s", manifest.Name, manifest.Version),
	); err != nil {
		return "", err
	}

	prBody := fmt.Sprintf(
		"\n### New Plugin (v%s)\n\n"+
			"- **Name:** %s\n"+
			"- **Description:** %s\n\n"+
			"Published via Inkdown CLI.",
		manifest.Version,
		manifest.Name,
		manifest.Description,
	)

	utils.Note("Creating the following PR:\n%s", prBody)
	reader := bufio.NewReader(os.Stdin)
	utils.Prompt("Please provide a title for your PR: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	prURL, err := github.CreatePR(token, username+":"+branch, title, prBody)
	if err != nil {
		return "", err
	}

	return prURL, nil
}
