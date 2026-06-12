package main

import "strings"

// vistaEse is the insider-VistA vocabulary the `v` platform exists to hide
// (v-cli-platform.md §2, §7). No domain, command, or flag name may contain any
// of these — a developer should be able to guess the surface without knowing
// VistA. The plain-language lint (and `v new`) enforce it mechanically.
var vistaEse = []string{
	"fileman", "kids", "xpar", "mailman", "taskman", "duz",
	"dic", "die", "diq", "dik", "zis", "ztload", "xusec", "xushsh",
	"vistalink", "rpcbroker", "kernel", "hl7",
}

// plainLangViolations returns, for each given name, the vista-ese terms it
// contains (case-insensitive substring match). An empty result means clean.
// Used to vet a new domain name (`v new`) and to lint a domain's whole command
// surface (the conformance plain-language gate).
func plainLangViolations(names []string) map[string][]string {
	out := map[string][]string{}
	for _, name := range names {
		lower := strings.ToLower(name)
		var hits []string
		for _, term := range vistaEse {
			if strings.Contains(lower, term) {
				hits = append(hits, term)
			}
		}
		if len(hits) > 0 {
			out[name] = hits
		}
	}
	return out
}
