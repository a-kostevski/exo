package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	isDay  bool
	noOpen bool
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new zet note",
	Long:  `Creates a new zet note with the given name in ZETDIR.`,
	Args:  cobra.RangeArgs(0, 1),
	Run:   runNew,
}

func runNew(cmd *cobra.Command, args []string) {
	if isDay {
		createDaily()
		return
	}

	if len(args) == 0 {
		fmt.Println("Error: no note name provided")
		os.Exit(1)
	}
	createNote(args[0])
}

func getEnvVar(key string) string {
	value := os.Getenv(key)

	if value == "" {
		fmt.Printf("Error: %s environment variable not set\n", key)
		os.Exit(1)
	}
	return value
}

func formatFileName(name string) string {
	filename := strings.ToLower(name)
	filename = strings.ReplaceAll(filename, " ", "-")

	return filename
}

func templateExists(templatepath string) bool {
	_, err := os.Stat(templatepath)
	return err == nil
}

func appendToDailyNote(dailyPath, filename string) {
	file, err := os.OpenFile(dailyPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening daily note: %v\n", err)
		return
	}
	defer file.Close()

	link := fmt.Sprintf("\n- [[%s]]\n", filename)
	if _, err := file.WriteString(link); err != nil {
		fmt.Printf("Error appending to daily note: %v\n", err)
	}
}

func getDateStrings() (string, string, string) {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	return today.Format("2006-01-02"), yesterday.Format("2006-01-02"), tomorrow.Format("2006-01-02")
}

func createNote(name string) {
	LIFE := getEnvVar("LIFE")
	ZETDIR := getEnvVar("ZETDIR")

	today := time.Now()
	todayStr := today.Format("2006-01-02")
	dailyPath := fmt.Sprintf("%s/periodic/daily/%s.md", LIFE, todayStr)

	// Create daily note if it doesn't exist
	if _, err := os.Stat(dailyPath); os.IsNotExist(err) {
		createDaily()
	}

	filename := formatFileName(name)

	path := fmt.Sprintf("%s/%s.md", ZETDIR, filename)

	template := fmt.Sprintf("# %s\n\nContent\n\n## Links\n", name)
	templatePath := filepath.Join(ZETDIR, "templates", "note.md")

	if _, err := os.Stat(templatePath); err == nil {
		createFileFromTemplate(path, templatePath, map[string]string{
			"title": name,
		})
	} else {
		createFile(path, template)
	}

	appendToDailyNote(dailyPath, filename)

	if !noOpen {
		if err := openInEditor(path); err != nil {
			fmt.Printf("Error opening file in editor: %v\n", err)
		}
	}
}

func createDaily() {
	LIFE := getEnvVar("LIFE")

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	todayStr := today.Format("2006-01-02")
	yesterdayStr := yesterday.Format("2006-01-02")
	tomorrowStr := tomorrow.Format("2006-01-02")

	path := fmt.Sprintf("%s/periodic/daily/%s.md", LIFE, todayStr)

	template := fmt.Sprintf("# %s\n\n## [[%s]] - [[%s]]\n\n....\n\n## Notes\n",
		todayStr, yesterdayStr, tomorrowStr)

	templatePath := filepath.Join(LIFE, "templates", "daily.md")
	if _, err := os.Stat(templatePath); err == nil {
		createFileFromTemplate(path, templatePath, map[string]string{
			"date.today":     todayStr,
			"date.yesterday": yesterdayStr,
			"date.tomorrow":  tomorrowStr,
		})
	} else {
		createFile(path, template)
	}

	if !noOpen {
		if err := openInEditor(path); err != nil {
			fmt.Printf("Error opening file in editor: %v\n", err)
		}
	}
}

func createFileFromTemplate(destPath, templatePath string, replacements map[string]string) {
	if fileExists(destPath) {
		fmt.Printf("Error: file already exists: %s\n", destPath)
		os.Exit(1)
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		os.Exit(1)
	}

	templateContent := string(content)
	for key, value := range replacements {
		templateContent = strings.ReplaceAll(templateContent, "{%"+key+"}", value)
	}

	if err := os.WriteFile(destPath, []byte(templateContent), 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created new note from template: %s\n", destPath)
}

func createFile(filepath string, content string) {
	if fileExists(filepath) {
		fmt.Printf("Error: file already exists: %s\n", filepath)
		os.Exit(1)
	}

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created new note: %s\n", filepath)
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func openInEditor(filepath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("EDITOR environment variable not set")
	}

	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	newCmd.Flags().BoolVar(&isDay, "day", false, "Create a daily note")
	newCmd.Flags().BoolVarP(&noOpen, "no-open", "n", false, "Don't open the note after creation")
	rootCmd.AddCommand(newCmd)
}
