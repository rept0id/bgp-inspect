package main

import (
	"fmt"
	"os/exec"
	"strings"
	"log"
)

func getASName(as uint32, logging bool) (name string) {
	WHOIS_AS_NAME_KEYS := []string{
		"as-name:", "AS-NAME:", "AS-Name:",
		"ASName:",
		"AS이름:", "b. [AS Name]                    ", "AS Name            :",
	}

	asPrefixed := fmt.Sprintf("AS%d", as)

	cmd := exec.Command("whois", asPrefixed)
	output, err := cmd.CombinedOutput()
	if err != nil {
		name = "ERR"

		if logging {
			log.Print("WHOIS " + asPrefixed + ": " + name)
		}
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		for _, p :=range WHOIS_AS_NAME_KEYS {
			if strings.HasPrefix(line, p) {
				name = strings.TrimSpace(strings.SplitN(line, p, 2)[1])

				break
			}
		}
	}

	if len(name) == 0 { name = "-" }

	if logging {
		log.Print("WHOIS " + asPrefixed + ": " + name)
	}

	return
}
