package integrations

import (
	. "dashboard/types"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const NGINX_CONF_TEMPLATE = `events {}

http {
	server {
		listen %s;

		%s
	}
}`

const NGINX_LOCATION_TEMPLATE = `location %s {
	proxy_pass %s;
}`

func GenerateNginxConfig(port string, resources []Resource) string {
	locations := ""
	for _, resource := range resources {
		locations += fmt.Sprintf(NGINX_LOCATION_TEMPLATE, resource.Route, resource.Source)
		locations += "\n"
	}

	config := fmt.Sprintf(NGINX_CONF_TEMPLATE, port, locations)
	return FormatNginxConfig(config)
}

func FormatNginxConfig(nginxConfig string) string {
	lines := strings.Split(nginxConfig, "\n")

	formattedConfig := ""
	depth := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 && line[0] == '}' {
			depth--
		}

		formattedConfig += fmt.Sprintf("%s%s\n", strings.Repeat("  ", depth), line)

		if len(line) > 0 && line[len(line)-1] == '{' {
			depth++
		}
	}

	return formattedConfig
}

func ReloadNginxConfig(confPath string, newConf string) {
	// Update nginx.conf with new configuration
	log.Println("Updating nginx.conf")
	f, err := os.Create(confPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(newConf)
	if err != nil {
		panic(err)
	}
	f.Sync()

	// Reload nginx
	log.Println("Reloading nginx")
	cmd := exec.Command("nginx", "-s", "reload")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	log.Println(string(output))
}
