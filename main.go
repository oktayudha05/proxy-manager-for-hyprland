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
    proxyIP   = "192.168.92.209"
    proxyPort = "8183"
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
    proxySettings  = []string{
        fmt.Sprintf("http_proxy=\"http://%s:%s\"", proxyIP, proxyPort),
        fmt.Sprintf("https_proxy=\"http://%s:%s\"", proxyIP, proxyPort),
        "no_proxy=\"localhost,127.0.0.1\"",
    }
    exportProxySettings = []string{
        fmt.Sprintf("export http_proxy=\"http://%s:%s\"", proxyIP, proxyPort),
        fmt.Sprintf("export https_proxy=\"http://%s:%s\"", proxyIP, proxyPort),
        "export no_proxy=\"localhost,127.0.0.1\"",
    }
)

func updateEnvironmentFile(activate bool) error {
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

func updateZshrc(activate bool) error {
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
    newLines = append(newLines, "")  // Add single trailing newline

    return os.WriteFile(zshrcFile, []byte(strings.Join(newLines, "\n")), 0644)
}

func main() {
    activate := flag.Bool("aktif", false, "Aktifkan proxy")
    deactivate := flag.Bool("mati", false, "Matikan proxy")
    flag.Parse()

    if !*activate && !*deactivate {
        fmt.Println("Gunakan flag -aktif untuk mengaktifkan proxy atau -mati untuk menonaktifkan proxy")
        return
    }

    if *activate && *deactivate {
        fmt.Println("Tidak bisa menggunakan kedua flag sekaligus")
        return
    }

    var err error
    if *activate {
        err = updateEnvironmentFile(true)
        if err == nil {
            err = updateZshrc(true)
        }
        if err == nil {
            fmt.Println("Proxy berhasil diaktifkan")
        }
    } else {
        err = updateEnvironmentFile(false)
        if err == nil {
            err = updateZshrc(false)
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