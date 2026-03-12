package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type Command struct {
	Use         string   `json:"use"`
	Short       string   `json:"short"`
	Method      string   `json:"method"`
	Endpoint    string   `json:"endpoint"`
	Flags       []string `json:"flags"`
	ArrayFlags  []string `json:"array_flags"`
}

type Spec struct {
	Name       string    `json:"name"`
	Short      string    `json:"short"`
	AuthHeader string    `json:"auth_header"`
	AuthPrefix string    `json:"auth_prefix"`
	Commands   []Command `json:"commands"`
}

const cliTemplate = `package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Config struct {
	Token string ` + "`json:\"token\"`" + `
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "{{.Name}}", "config.json")
}

func loadToken() (string, error) {
	b, err := os.ReadFile(getConfigPath())
	if err != nil {
		return "", err
	}
	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return "", err
	}
	return cfg.Token, nil
}

func saveToken(token string) error {
	configPath := getConfigPath()
	os.MkdirAll(filepath.Dir(configPath), 0755)
	cfg := Config{Token: token}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(configPath, b, 0600)
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "{{.Name}}",
		Short: "{{.Short}}",
	}

	// AUTH COMMAND
	var authCmd = &cobra.Command{Use: "auth", Short: "Manage authentication"}
	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticate and save API token",
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := cmd.Flags().GetString("token")
			if token == "" {
				fmt.Println("{\"error\": \"--token is required\"}")
				os.Exit(1)
			}
			if err := saveToken(token); err != nil {
				fmt.Printf("{\"error\": \"Failed to save token: %v\"}\n", err)
				os.Exit(1)
			}
			fmt.Println("{\"success\": true, \"message\": \"Authenticated successfully\"}")
		},
	}
	loginCmd.Flags().String("token", "", "API Key")
	authCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(authCmd)

	{{range .Commands}}
	{{$cmdUse := .Use}}
	var {{.Use}}Cmd = &cobra.Command{
		Use:   "{{.Use}}",
		Short: "{{.Short}}",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := loadToken()
			if err != nil || token == "" {
				// Fallback to env var if config file is not found
				token = os.Getenv("API_TOKEN")
				if token == "" {
					fmt.Println("{\"error\": \"Not authenticated. Run '{{$.Name}} auth login --token <token>'\"}")
					os.Exit(1)
				}
			}

			reqURL := "{{.Endpoint}}"
			var req *http.Request
			var reqErr error

			{{if eq .Method "GET"}}
			// GET Request logic
			u, _ := url.Parse(reqURL)
			q := u.Query()
			{{range .Flags}}
			if val, _ := cmd.Flags().GetString("{{.}}"); val != "" {
				q.Set("{{.}}", val)
			}
			{{end}}
			u.RawQuery = q.Encode()
			req, reqErr = http.NewRequest("GET", u.String(), nil)
			{{else}}
			// POST/PUT Request logic (Multipart Form)
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			
			{{range .Flags}}
			if val, _ := cmd.Flags().GetString("{{.}}"); val != "" {
				writer.WriteField("{{.}}", val)
			}
			{{end}}

			{{range .ArrayFlags}}
			if arr, _ := cmd.Flags().GetStringSlice("{{.}}"); len(arr) > 0 {
				for _, v := range arr {
					writer.WriteField("{{.}}[]", v)
				}
			}
			{{end}}

			writer.Close()
			req, reqErr = http.NewRequest("{{.Method}}", reqURL, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			{{end}}

			if reqErr != nil {
				fmt.Printf("{\"error\": \"%v\"}\n", reqErr)
				os.Exit(1)
			}
			authHeader := "{{if $.AuthHeader}}{{$.AuthHeader}}{{else}}Authorization{{end}}"
			authPrefix := "{{if $.AuthPrefix}}{{$.AuthPrefix}}{{else}}Apikey {{end}}"
			req.Header.Set(authHeader, authPrefix+token)

			client := &http.Client{}
			resp, doErr := client.Do(req)
			if doErr != nil {
				fmt.Printf("{\"error\": \"%v\"}\n", doErr)
				os.Exit(1)
			}
			defer resp.Body.Close()

			respBody, _ := io.ReadAll(resp.Body)
			
			// Try to pretty-print JSON response
			var prettyJSON bytes.Buffer
			if parseErr := json.Indent(&prettyJSON, respBody, "", "  "); parseErr == nil {
				fmt.Println(string(prettyJSON.Bytes()))
			} else {
				fmt.Println(string(respBody))
			}
		},
	}
	{{range .Flags}}
	{{$cmdUse}}Cmd.Flags().String("{{.}}", "", "{{.}} parameter")
	{{end}}
	{{range .ArrayFlags}}
	{{$cmdUse}}Cmd.Flags().StringSlice("{{.}}", []string{}, "{{.}} parameter")
	{{end}}
	rootCmd.AddCommand({{.Use}}Cmd)
	{{end}}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
`

func main() {
	if len(os.Args) < 3 || os.Args[1] != "build" {
		fmt.Println("Usage: aif build <spec.json>")
		os.Exit(1)
	}

	specPath := os.Args[2]
	specData, err := os.ReadFile(specPath)
	if err != nil {
		panic(err)
	}

	var spec Spec
	if err := json.Unmarshal(specData, &spec); err != nil {
		panic(err)
	}

	// Create build directory
	buildDir := "build_tmp"
	os.RemoveAll(buildDir)
	os.MkdirAll(buildDir, 0755)
	
	mainFile := filepath.Join(buildDir, "main.go")
	f, err := os.Create(mainFile)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("cli").Parse(cliTemplate)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(f, spec); err != nil {
		panic(err)
	}
	f.Close()

	fmt.Println("⚙️ Generating Go code...")
	
	// Build the binary
	fmt.Println("📦 Compiling binary...")
	cmdInit := exec.Command("go", "mod", "init", spec.Name)
	cmdInit.Dir = buildDir
	cmdInit.Run()

	cmdGet := exec.Command("go", "get", "github.com/spf13/cobra")
	cmdGet.Dir = buildDir
	cmdGet.Run()

	cmdBuild := exec.Command("go", "build", "-o", "../"+spec.Name)
	cmdBuild.Dir = buildDir
	out, err := cmdBuild.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Build failed: %v\n%s\n", err, out)
		os.Exit(1)
	}

	os.RemoveAll(buildDir)
	fmt.Printf("✅ Success! Single binary CLI '%s' generated.\n", spec.Name)
}