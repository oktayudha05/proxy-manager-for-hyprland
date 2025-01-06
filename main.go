package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	defaultProxyIP   = "192.168.92.209"
	defaultProxyPort = "8183"
)

func getHomeDir() (string, error) {
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		usr, err := user.Lookup(sudoUser)
		if err != nil {
			return "", err
		}
		return usr.HomeDir, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

var (
	environmentFile = "/etc/environment"
)

func generateProxySettings(ip, port string) ([]string, []string) {
	proxySettings := []string{
		fmt.Sprintf("http_proxy=\"http://%s:%s\"", ip, port),
		fmt.Sprintf("https_proxy=\"http://%s:%s\"", ip, port),
		"no_proxy=\"localhost,127.0.0.1\"",
	}
	exportProxySettings := []string{
		fmt.Sprintf("export http_proxy=\"http://%s:%s\"", ip, port),
		fmt.Sprintf("export https_proxy=\"http://%s:%s\"", ip, port),
		"export no_proxy=\"localhost,127.0.0.1\"",
	}
	return proxySettings, exportProxySettings
}

func updateEnvironmentFile(activate bool, proxySettings []string) error {
	var content string
	if activate {
		content = `#
# This file is parsed by pam_env module
#
# Syntax: simple "KEY=VAL" pairs on separate lines
#

` + strings.Join(proxySettings, "\n")
	} else {
		content = `#
# This file is parsed by pam_env module
#
# Syntax: simple "KEY=VAL" pairs on separate lines
#
`
	}

	return os.WriteFile(environmentFile, []byte(content), 0644)
}

func updateZshrc(activate bool, exportProxySettings []string) error {
	homeDir, err := getHomeDir()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan home directory: %v", err)
	}

	zshrcFile := filepath.Join(homeDir, ".zshrc")
	content, err := os.ReadFile(zshrcFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	newLines := make([]string, 0)

	for _, line := range lines {
		if !strings.Contains(line, "export http_proxy") &&
			!strings.Contains(line, "export https_proxy") &&
			!strings.Contains(line, "export no_proxy") {
			newLines = append(newLines, line)
		}
	}

	if activate {
		newLines = append(newLines, "")
		newLines = append(newLines, exportProxySettings...)
	}

	for len(newLines) > 0 && newLines[len(newLines)-1] == "" {
		newLines = newLines[:len(newLines)-1]
	}
	newLines = append(newLines, "")

	return os.WriteFile(zshrcFile, []byte(strings.Join(newLines, "\n")), 0644)
}

func main() {
	activate := flag.Bool("aktif", false, "Aktifkan proxy")
	deactivate := flag.Bool("mati", false, "Matikan proxy")
	proxyIP := flag.String("ip", defaultProxyIP, "IP proxy yang akan digunakan")
	proxyPort := flag.String("port", defaultProxyPort, "Port proxy yang akan digunakan")
	flag.Parse()

	if !*activate && !*deactivate {
		fmt.Println("Gunakan flag -aktif untuk mengaktifkan proxy atau -mati untuk menonaktifkan proxy")
		return
	}

	if *activate && *deactivate {
		fmt.Println("Tidak bisa menggunakan kedua flag sekaligus")
		return
	}

	proxySettings, exportProxySettings := generateProxySettings(*proxyIP, *proxyPort)

	var err error
	if *activate {
		err = updateEnvironmentFile(true, proxySettings)
		if err == nil {
			err = updateZshrc(true, exportProxySettings)
		}
		if err == nil {
			fmt.Printf("Proxy berhasil diaktifkan dengan IP %s dan port %s\n", *proxyIP, *proxyPort)
		}
	} else {
		err = updateEnvironmentFile(false, proxySettings)
		if err == nil {
			err = updateZshrc(false, exportProxySettings)
		}
		if err == nil {
			fmt.Println("Proxy berhasil dinonaktifkan")
		}
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
