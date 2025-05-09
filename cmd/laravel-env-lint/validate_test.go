package main

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

// helper to run rootCmd with args and capture output+error
func runCommand(args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	_, cmdErr := rootCmd.ExecuteC()    // returns error instead of exiting
    // clear out any commands left behind so subsequent tests start fresh
    rootCmd.SetArgs([]string{})
    return buf.String(), cmdErr
}

func TestValidate_Help(t *testing.T) {
	out, err := runCommand("validate", "--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(out, "Validate a .env file") {
		t.Errorf("expected help text, got %q", out)
	}
}

func TestNormalizeKey_RawAndBase64(t *testing.T) {
	// Raw 32-char key
	raw := strings.Repeat("a", 32)
	if b, err := normalizeKey(raw); err != nil || len(b) != 32 {
		t.Errorf("normalizeKey(raw) = %v, %v; want len=32 no error", b, err)
	}

	// Base64 key
	enc := "base64:" + base64.StdEncoding.EncodeToString([]byte(raw))
	if b, err := normalizeKey(enc); err != nil || string(b) != raw {
		t.Errorf("normalizeKey(base64) = %v, %v; want %q, nil", b, err, raw)
	}

	// Too-short key
	if _, err := normalizeKey("short"); err == nil {
		t.Errorf("normalizeKey(short) did not error")
	}
}