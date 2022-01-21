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

// Go doesn't have Set in stdlib so use a map.
var services = map[string]struct{}{
	"riff-raff": {},
	"teamcity":  {},
	"amigo":     {},
}

func main() {
	m := flag.String("m", "Describe your problem here.", "Optional description of task.")
	flag.Parse()
	service := flag.Arg(0)

	if _, ok := services[service]; !ok {
		fmt.Printf("Unrecognised or missing service: '%s'\n", service)
		fmt.Println("Supported services are: riff-raff, teamcity, amigo.")
		os.Exit(1)
	}

	out, err := exec.Command("git", "config", "user.email").CombinedOutput()
	check(err, fmt.Sprintf("Unable to get user's git email: %s", string(out)))

	data :=
		fmt.Sprintf(`
responsible: %s
service: %s
----
%s
`, strings.TrimSpace(string(out)), service, *m)

	data = strings.TrimSpace(data)

	println(data)
	os.Exit(0)

	home, _ := os.UserHomeDir()
	toilHome := filepath.Join(home, "toil")

	// if ~/toil doesn't exist clone into it
	_, err = os.Stat(toilHome)
	if os.IsNotExist(err) {
		out, err := execGit(home, "clone", "git@github.com:guardian/toil-records.git", "toil")
		check(err, fmt.Sprintf("Unable to clone from git@github.com:guardian/toil-records.git: %s\n%v", string(out), err))
	}

	fileName := time.Now().Format(time.RFC3339)
	filePath := filepath.Join(toilHome, fileName)

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
