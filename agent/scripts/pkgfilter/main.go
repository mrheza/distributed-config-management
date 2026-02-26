package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	sep := flag.String("sep", "space", "output separator: space or comma")
	flag.Parse()

	cmd := exec.Command("go", "list", "./...")
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	pkgs := make([]string, 0, len(lines))

	for _, pkg := range lines {
		pkg = strings.TrimSpace(pkg)
		if pkg == "" {
			continue
		}
		if strings.Contains(pkg, "/internal/mocks/") {
			continue
		}
		if strings.HasSuffix(pkg, "/docs") ||
			strings.Contains(pkg, "/docs/") ||
			strings.HasSuffix(pkg, "/scripts") ||
			strings.Contains(pkg, "/scripts/") {
			continue
		}
		pkgs = append(pkgs, pkg)
	}

	joiner := " "
	if *sep == "comma" {
		joiner = ","
	}

	fmt.Print(strings.Join(pkgs, joiner))
}
