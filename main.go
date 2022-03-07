package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var home, _ = os.UserHomeDir()
var toilHome = filepath.Join(home, "toil")

func main() {

	m := flag.String("m", "", "Description of task.")
	dryRun := flag.Bool("dryRun", false, "Prints out toil file and exits without writing to toil home and remote.")
	h := flag.Bool("h", false, "Help info.")

	flag.Parse()
	service := flag.Arg(0)

	if *m == "" {
		fmt.Println("Required flag '-m' missing (note, must be *before* the service ID).")
		os.Exit(1)
	}

	if *h {
		fmt.Println("$ toil [-flag ...] [service]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	ensureToilHome()
	validateService(service)

	out, err := exec.Command("git", "config", "user.email").CombinedOutput()
	check(err, fmt.Sprintf("Unable to get user's git email: %s", string(out)))

	data :=
		fmt.Sprintf(`
responsible: %s
service: %s
----
%s
`, strings.TrimSpace(string(out)), service, *m)

	data = strings.TrimSpace(data) + "\n"

	if *dryRun {
		println(data)
		os.Exit(0)
	}

	fileName := time.Now().Format(time.RFC3339)
	filePath := filepath.Join(toilHome, fileName)

	_, err = execGit(toilHome, "pull")
	check(err, fmt.Sprintf("Unable to pull from main: '%v'", err))

	err = os.WriteFile(filePath, []byte(data), 0777)
	check(err, fmt.Sprintf("Unable to write toil file: '%v'", err))

	out, err = execGit(toilHome, "add", fileName)
	check(err, fmt.Sprintf("Unable to git add new toil file: %s\n%v", string(out), err))

	out, err = execGit(toilHome, "commit", "-m", fmt.Sprintf("Add %s", fileName))
	check(err, fmt.Sprintf("Unable to git commit new toil file: %s\n%v", string(out), err))

	out, err = execGit(toilHome, "push", "-u", "origin", "HEAD")
	check(err, fmt.Sprintf("Unable to git push new toil file: %s\n%v", out, err))
}

func execGit(workingDir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = workingDir
	out, err := cmd.CombinedOutput()
	return out, err
}

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func ensureToilHome() {
	_, err := os.Stat(toilHome)
	if os.IsNotExist(err) {
		out, err := execGit(home, "clone", "git@github.com:guardian/toil-records.git", "toil")
		check(err, fmt.Sprintf("Unable to clone from git@github.com:guardian/toil-records.git: %s\n%v", string(out), err))
	}
}

func validateService(service string) {
	services := readServicesConfig()

	for _, s := range services {
		if s == service {
			return
		}
	}

	fmt.Printf("Unrecognised or missing service arg: '%s'. Must exist in ~/toil/services.txt. Use -h for help.\n", service)
	os.Exit(1)
}

func readServicesConfig() []string {
	servicesData, err := os.ReadFile(filepath.Join(toilHome, "services.txt"))
	check(err, fmt.Sprintf("Unable to read ~/toil/services.txt file. Should contain a list of services separated by newline: %v", err))

	trimmed := strings.TrimSpace(string(servicesData))
	return strings.Split(trimmed, "\n")
}
