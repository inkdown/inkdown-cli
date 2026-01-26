package auth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"inkdown-cli/config"
)

const (
	apiBaseURL = "http://localhost:8080/api/v1"
)

type CLILoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type CLILoginResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Token       string   `json:"token"`
		TokenPrefix string   `json:"token_prefix"`
		Scopes      []string `json:"scopes"`
		CreatedAt   string   `json:"created_at"`
		Message     string   `json:"message"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

func Auth() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.IsAuthenticated() {
		fmt.Printf("✓ You are already authenticated as: %s\n", cfg.Email)
		fmt.Println("  Use 'ink logout' to sign out first.")
		return nil
	}

	fmt.Println("Inkdown CLI Authentication")
	fmt.Println("─────────────────────────────")

	reader := bufio.NewReader(os.Stdin)

	// Get email
	fmt.Print("Email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read email: %w", err)
	}

	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	fmt.Print("Password: ")
	password, err := readPassword()
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
	deviceName := fmt.Sprintf("%s-%s-%s", hostname, runtime.GOOS, runtime.GOARCH)

	fmt.Println("\n Authenticating...")

	token, err := login(email, password, deviceName)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	cfg.Token = token
	cfg.Email = email
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	fmt.Println("\n✓ Authentication successful!")
	fmt.Printf("  Logged in as: %s\n", email)
	fmt.Printf("  Config saved to: %s\n", config.ConfigPath())

	return nil
}

func Logout() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.IsAuthenticated() {
		fmt.Println("You are not currently logged in.")
		return nil
	}

	email := cfg.Email
	cfg.Token = ""
	cfg.Email = ""

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}

	fmt.Printf("✓ Logged out successfully from: %s\n", email)
	return nil
}

func Whoami() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.IsAuthenticated() {
		fmt.Println("Not authenticated. Use 'inkdown auth' to login.")
		return nil
	}

	fmt.Printf("Logged in as: %s\n", cfg.Email)
	return nil
}

func login(email, password, deviceName string) (string, error) {
	reqBody := CLILoginRequest{
		Email:    email,
		Password: password,
		Name:     deviceName,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiBaseURL+"/cli/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	var result CLILoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		errMsg := result.Error
		if errMsg == "" {
			errMsg = "unknown error"
		}
		return "", fmt.Errorf("%s", errMsg)
	}

	if result.Data.Token == "" {
		return "", fmt.Errorf("server did not return a token")
	}

	return result.Data.Token, nil
}

func readPassword() (string, error) {
	if runtime.GOOS == "windows" {
		reader := bufio.NewReader(os.Stdin)
		password, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(password), nil
	}

	cmd := exec.Command("stty", "-echo")
	cmd.Stdin = os.Stdin
	_ = cmd.Run()

	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')

	cmd = exec.Command("stty", "echo")
	cmd.Stdin = os.Stdin
	_ = cmd.Run()

	fmt.Println()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(password), nil
}
