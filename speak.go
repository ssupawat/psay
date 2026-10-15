package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var sentenceRe = regexp.MustCompile(`[^.!?…]+[.!?…]+["')\]]*\s+`)

// SplitSentences breaks announce text into speakable chunks while keeping
// version numbers ("v1.2.3") and file names ("auth_service.go") intact.
// Fragments of fewer than two words fold into a neighbor so sentences
// like "Hi. Deploy is done." stay in one chunk.
func SplitSentences(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	var chunks []string
	prevEnd := 0
	for _, m := range sentenceRe.FindAllStringIndex(text, -1) {
		chunk := strings.TrimSpace(text[prevEnd:m[1]])
		if chunk != "" {
			chunks = append(chunks, chunk)
		}
		prevEnd = m[1]
	}
	if rest := strings.TrimSpace(text[prevEnd:]); rest != "" {
		chunks = append(chunks, rest)
	}

	merged := make([]string, 0, len(chunks))
	for _, c := range chunks {
		if len(strings.Fields(c)) < 2 && len(merged) > 0 {
			merged[len(merged)-1] += " " + c
			continue
		}
		merged = append(merged, c)
	}
	if len(merged) > 1 && len(strings.Fields(merged[0])) < 2 {
		merged[1] = merged[0] + " " + merged[1]
		merged = merged[1:]
	}
	return merged
}

// CacheKey identifies a synthesized announce by voice, speed, and text.
func CacheKey(voice string, speed float64, text string) string {
	sum := sha256.Sum256([]byte(voice + "\x00" + fmt.Sprintf("%g", speed) + "\x00" + text))
	return hex.EncodeToString(sum[:])
}

func cachedWav(voice string, speed float64, text string) (string, bool) {
	psayDir, err := defaultDir()
	if err != nil {
		return "", false
	}
	wav := filepath.Join(psayDir, "cache", CacheKey(voice, speed, text)+".wav")
	info, err := os.Stat(wav)
	if err != nil || info.Size() == 0 {
		return "", false
	}
	return wav, true
}

type speakRequest struct {
	Chunks []string `json:"chunks"`
	Voice  string   `json:"voice"`
	Speed  float64  `json:"speed"`
	Cache  string   `json:"cache"`
	Dir    string   `json:"dir"`
}

type speakEvent struct {
	Event   string `json:"event"`
	Message string `json:"message,omitempty"`
}

func daemonSocketPath() (string, error) {
	psayDir, err := defaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(psayDir, "daemon.sock"), nil
}

// speakViaDaemon streams the request to a running daemon, starting one if
// needed. played is false when no daemon could serve, so the caller can
// fall back to a cold one-shot synth.
func speakViaDaemon(req speakRequest) (played bool, err error) {
	sock, err := daemonSocketPath()
	if err != nil {
		return false, err
	}

	conn, err := net.Dial("unix", sock)
	if err != nil {
		if err := startDaemon(sock); err != nil {
			return false, err
		}
		conn, err = dialPoll(sock, 40, 100*time.Millisecond)
		if err != nil {
			return false, err
		}
	}
	defer conn.Close()

	payload, err := json.Marshal(req)
	if err != nil {
		return true, err
	}
	if _, err := conn.Write(append(payload, '\n')); err != nil {
		return true, err
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var ev speakEvent
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			return true, err
		}
		switch ev.Event {
		case "done":
			return true, nil
		case "error":
			return true, fmt.Errorf("daemon: %s", ev.Message)
		}
	}
	return true, scanner.Err()
}

func dialPoll(sock string, attempts int, interval time.Duration) (net.Conn, error) {
	for i := 0; i < attempts; i++ {
		conn, err := net.Dial("unix", sock)
		if err == nil {
			return conn, nil
		}
		time.Sleep(interval)
	}
	return nil, fmt.Errorf("daemon did not accept on %s", sock)
}

func startDaemon(sock string) error {
	psayDir, err := defaultDir()
	if err != nil {
		return err
	}
	python := filepath.Join(psayDir, "venv", "bin", "python")
	if _, err := os.Stat(python); err != nil {
		return fmt.Errorf("no psay venv at %s", python)
	}

	os.Remove(sock)
	cmd := exec.Command(python, "-c", speakScript,
		"daemon", sock, filepath.Join(psayDir, "kokoro"))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}

// announce plays the text: cached wav if present, else daemon synth with a
// cold one-shot fallback.
func announce(voice string, speed float64, text string) error {
	psayDir, err := defaultDir()
	if err != nil {
		return err
	}
	dir, err := os.MkdirTemp("", "psay-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if wav, ok := cachedWav(voice, speed, text); ok {
		_, err := run(exec.Command("afplay", wav))
		return err
	}

	req := speakRequest{
		Chunks: SplitSentences(text),
		Voice:  voice,
		Speed:  speed,
		Cache:  filepath.Join(psayDir, "cache", CacheKey(voice, speed, text)+".wav"),
		Dir:    dir,
	}

	played, err := speakViaDaemon(req)
	if played {
		return err
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return speakOneshot(req)
}

// speakOneshot synthesizes and plays in a fresh python process, no daemon.
func speakOneshot(req speakRequest) error {
	psayDir, err := defaultDir()
	if err != nil {
		return err
	}
	python := filepath.Join(psayDir, "venv", "bin", "python")
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	cmd := exec.Command(python, "-c", speakScript, "oneshot", filepath.Join(psayDir, "kokoro"))
	cmd.Stdin = strings.NewReader(string(payload) + "\n")
	_, err = run(cmd)
	return err
}
