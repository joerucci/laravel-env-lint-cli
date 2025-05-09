package loader

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
)

func TestEnsureTrailingNewline(t *testing.T) {
    cases := []struct {
        input []byte
        want  []byte
    }{
        {[]byte("foo"), []byte("foo\n")},
        {[]byte("bar\n"), []byte("bar\n")},
    }
    for _, c := range cases {
        got := ensureTrailingNewline(c.input)
        if string(got) != string(c.want) {
            t.Errorf("ensureTrailingNewline(%q) = %q; want %q", c.input, got, c.want)
        }
    }
}

func TestReadAndDecode(t *testing.T) {
    // create temp file
    dir := t.TempDir()
    path := filepath.Join(dir, "data.b64")
    original := "hello world"
    enc := base64.StdEncoding.EncodeToString([]byte(original))
    if err := os.WriteFile(path, []byte(enc), 0644); err != nil {
        t.Fatalf("failed to write temp file: %v", err)
    }
    // read and decode
    got, err := readAndDecode(path)
    if err != nil {
        t.Fatalf("readAndDecode error: %v", err)
    }
    if string(got) != original {
        t.Errorf("readAndDecode = %q; want %q", got, original)
    }
}

func TestParseEnvelopeRaw(t *testing.T) {
    // prepare raw payload: 16-byte IV + ciphertext
    iv := make([]byte, aes.BlockSize)
    for i := range iv {
        iv[i] = byte(i)
    }
    ct := []byte("ciphertext-data")
    payload := append(iv, ct...)
    gotIV, gotCT, err := parseEnvelope(payload)
    if err != nil {
        t.Fatalf("parseEnvelope raw error: %v", err)
    }
    if !equal(gotIV, iv) {
        t.Errorf("parseEnvelope raw IV = %v; want %v", gotIV, iv)
    }
    if !equal(gotCT, ct) {
        t.Errorf("parseEnvelope raw ciphertext = %v; want %v", gotCT, ct)
    }
}

func TestParseEnvelopeJSON(t *testing.T) {
    // create envelope
    iv := []byte("1234567890123456")
    ct := []byte("encrypted-data")
    env := laravelEnvelope{
        IV:    base64.StdEncoding.EncodeToString(iv),
        Value: base64.StdEncoding.EncodeToString(ct),
        MAC:   "",
    }
    j, err := json.Marshal(env)
    if err != nil {
        t.Fatalf("json.Marshal error: %v", err)
    }
    payload := j
    gotIV, gotCT, err := parseEnvelope(payload)
    if err != nil {
        t.Fatalf("parseEnvelope json error: %v", err)
    }
    if !equal(gotIV, iv) {
        t.Errorf("parseEnvelope json IV = %v; want %v", gotIV, iv)
    }
    if !equal(gotCT, ct) {
        t.Errorf("parseEnvelope json ciphertext = %v; want %v", gotCT, ct)
    }
}

func TestDecryptCBC(t *testing.T) {
    key := []byte("abcdefghijklmnopqrstuvwxzy012345") // 32 bytes
    iv := []byte("1234567890123456")                  // 16 bytes
    plaintext := []byte("secret message!")
    // pad PKCS#7
    padLen := aes.BlockSize - len(plaintext)%aes.BlockSize
    pad := bytesRepeat(byte(padLen), padLen)
    pt := append(plaintext, pad...)
    // encrypt
    block, err := aes.NewCipher(key)
    if err != nil {
        t.Fatalf("aes.NewCipher error: %v", err)
    }
    ct := make([]byte, len(pt))
    cipher.NewCBCEncrypter(block, iv).CryptBlocks(ct, pt)
    // decrypt
    got, err := decryptCBC(iv, ct, key)
    if err != nil {
        t.Fatalf("decryptCBC error: %v", err)
    }
    if string(got) != string(plaintext) {
        t.Errorf("decryptCBC = %q; want %q", got, plaintext)
    }
}

func TestLoadEnv(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "test.env")
    content := "FOO=bar\nBAZ=qux\n"
    if err := os.WriteFile(path, []byte(content), 0644); err != nil {
        t.Fatalf("WriteFile error: %v", err)
    }
    got, err := LoadEnv(path)
    if err != nil {
        t.Fatalf("LoadEnv error: %v", err)
    }
    if got["FOO"] != "bar" || got["BAZ"] != "qux" {
        t.Errorf("LoadEnv map = %v; want FOO=bar, BAZ=qux", got)
    }
}

func TestLoadSchema(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "schema.yaml")
    yamlContent := `FOO:
  type: string
  required: true
`  
    if err := os.WriteFile(path, []byte(yamlContent), 0644); err != nil {
        t.Fatalf("WriteFile error: %v", err)
    }
    got, err := LoadSchema(path)
    if err != nil {
        t.Fatalf("LoadSchema error: %v", err)
    }
    spec, ok := got["FOO"]
    if !ok || spec.Type != "string" || !spec.Required {
        t.Errorf("LoadSchema map = %v; want 'FOO' with type=string, required=true", got)
    }
}

// equal is a simple byte-slice comparator
func equal(a, b []byte) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

// bytesRepeat returns a slice of length n with each element set to b.
func bytesRepeat(b byte, n int) []byte {
    out := make([]byte, n)
    for i := range out {
        out[i] = b
    }
    return out
}
