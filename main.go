package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	lockFileName = "ente-sync.lock"
	timeout      = 6 * time.Hour
	ntfyURL      = "https://ntfy.sh"
	resultsFile  = "/srv/backups/ente-sync-results.log"
)

func main() {
	log.SetFlags(log.LstdFlags)

	lockFile := filepath.Join(os.TempDir(), lockFileName)

	exportDir := os.Getenv("EXPORT_DIR")
	secretsPath := os.Getenv("SECRETS_PATH")

	if exportDir == "" {
		log.Fatal("EXPORT_DIR must be set")
	}

	if secretsPath == "" {
		log.Fatal("SECRETS_PATH must be set")
	}

	// Prevent overlapping runs
	lockHandle, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("failed to open lock file: %v", err)
	}
	defer lockHandle.Close()

	if err := syscall.Flock(int(lockHandle.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		log.Println("another sync process is already running, exiting")
		os.Exit(0)
	}
	defer syscall.Flock(int(lockHandle.Fd()), syscall.LOCK_UN)

	if err := os.MkdirAll(filepath.Dir(secretsPath), 0700); err != nil {
		log.Fatalf("failed to create secrets directory: %v", err)
	}

	if err := os.MkdirAll(exportDir, 0755); err != nil {
		log.Fatalf("failed to create export directory: %v", err)
	}

	os.Setenv("ENTE_CLI_SECRETS_PATH", secretsPath)

	args := []string{"export", exportDir}

	if albums := os.Getenv("ALBUMS"); albums != "" {
		args = append(args, "--albums", albums)
	}

	if os.Getenv("INCLUDE_HIDDEN") == "true" {
		args = append(args, "--hidden=true")
	}

	log.Printf("starting ente export to %s", exportDir)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ente", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		recordResult("FAIL")
		if ctx.Err() == context.DeadlineExceeded {
			notify("ente export timed out after " + timeout.String())
			log.Fatalf("ente export timed out after %s", timeout)
		}
		notify(fmt.Sprintf("ente export failed: %v", err))
		log.Fatalf("ente export failed: %v", err)
	}

	recordResult("OK")
	log.Printf("export completed successfully")
}

func recordResult(result string) {
	f, err := os.OpenFile(resultsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to write result: %v", err)
		return
	}
	defer f.Close()
	fmt.Fprintln(f, result)

	if result != "OK" {
		return
	}

	data, err := os.ReadFile(resultsFile)
	if err != nil {
		return
	}

	ok, fail := 0, 0
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		switch strings.TrimSpace(line) {
		case "OK":
			ok++
		case "FAIL":
			fail++
		}
	}

	if ok < 24 {
		return
	}

	total := ok + fail
	if fail > 0 {
		notify(fmt.Sprintf("Daily ente summary: %d/%d runs succeeded, %d failed", ok, total, fail))
	} else {
		notify(fmt.Sprintf("Daily ente summary: all %d runs succeeded", total))
	}

	os.Truncate(resultsFile, 0)
}

func notify(msg string) {
	ntfyTopic := os.Getenv("NTFY_TOPIC")
	if ntfyURL == "" || ntfyTopic == "" {
		return
	}
	url := fmt.Sprintf("%s/%s", ntfyURL, ntfyTopic)
	resp, err := http.Post(url, "text/plain", strings.NewReader(msg))
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return
	}
	resp.Body.Close()
}
