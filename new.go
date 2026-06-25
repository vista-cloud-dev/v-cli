package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/vista-cloud-dev/v-pkg/clikit"
)

// NewCmd is `v new <domain>` (CQ4): scaffold a new v domain tool. It enforces
// the plain-language rule on the domain name and emits a minimal,
// convention-conforming skeleton built into `v new`. A domain
// is born with: a Go module github.com/vista-cloud-dev/v-<domain>, an importable
// <domain>cli command package the umbrella mounts, and a Makefile.
type NewCmd struct {
	Domain string `arg:"" help:"New domain name — a plain-language noun (db, config, rpc, job…), never a VistA product name."`
	Dir    string `help:"Parent directory in which to create v-<domain>/." default:"."`
}

func (c *NewCmd) Run(cc *clikit.Context) error {
	if bad := plainLangViolations([]string{c.Domain}); len(bad) > 0 {
		return clikit.Fail(clikit.ExitUsage, "VISTA_ESE",
			fmt.Sprintf("domain %q uses insider VistA term(s): %v", c.Domain, bad[c.Domain]),
			"use a plain-language noun (db not fileman, pkg not kids, config not xpar)")
	}
	root := filepath.Join(c.Dir, "v-"+c.Domain)
	if _, err := os.Stat(root); err == nil {
		return clikit.Fail(clikit.ExitUsage, "EXISTS", root+" already exists", "remove it or choose another name")
	}
	files, err := scaffoldDomain(c.Domain)
	if err != nil {
		return clikit.Fail(clikit.ExitRuntime, "SCAFFOLD_FAILED", err.Error(), "")
	}
	written := make([]string, 0, len(files))
	for rel, content := range files {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return clikit.Fail(clikit.ExitRuntime, "SCAFFOLD_FAILED", err.Error(), "")
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			return clikit.Fail(clikit.ExitRuntime, "SCAFFOLD_FAILED", err.Error(), "")
		}
		written = append(written, rel)
	}
	sort.Strings(written)
	return cc.Result(map[string]any{"domain": c.Domain, "dir": root, "files": written}, func() {
		fmt.Fprintln(cc.Stdout, cc.Success(fmt.Sprintf("scaffolded v-%s at %s", c.Domain, root)))
		for _, f := range written {
			fmt.Fprintln(cc.Stdout, "  "+cc.Faint(f))
		}
		fmt.Fprintln(cc.Stdout, cc.Faint("next: pin it in the v umbrella's go.mod and mount "+c.Domain+"cli.Commands"))
	})
}

// scaffoldDomain returns the relative-path → content map for a new domain
// skeleton. Kept minimal and convention-conforming; the skeleton is built into
// `v new` (there is no standalone template repo).
func scaffoldDomain(domain string) (map[string]string, error) {
	mod := "github.com/vista-cloud-dev/v-" + domain
	pkg := domain + "cli"
	files := map[string]string{
		"go.mod": fmt.Sprintf("module %s\n\ngo 1.26.3\n", mod),
		"README.md": fmt.Sprintf("# v-%s\n\nThe `v %s` domain of the v CLI. Scaffolded by `v new`.\n\n"+
			"Offline verbs run standalone; the `v` umbrella mounts `%s.Commands` as `v %s <verb>`.\n",
			domain, domain, pkg, domain),
		"Makefile": "BIN ?= v-" + domain + "\n\nbuild:\n\tgo build -o dist/$(BIN) .\n\n" +
			"test:\n\tgo test -race -cover ./...\n\nlint:\n\tgolangci-lint run ./...\n\n" +
			"check: lint test build\n",
		pkg + "/commands.go": fmt.Sprintf("// Package %s is the importable command surface of the v %s domain,\n"+
			"// mounted by the v umbrella as `v %s <verb>`.\npackage %s\n\n"+
			"// Commands is the %s verb set. Add domain verbs as fields.\n"+
			"type Commands struct {\n}\n", pkg, domain, domain, pkg, domain),
	}
	return files, nil
}
