package loader

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "os"
    "strings"

    "github.com/joho/godotenv"
    "gopkg.in/yaml.v3"

    "github.com/joerucci/laravel-env-lint-cli/pkg/schema"
)

type laravelEnvelope struct {
    IV    string `json:"iv"`
    Value string `json:"value"`
    MAC   string `json:"mac"`
}

// LoadEnv reads a plaintext .env file.
func LoadEnv(path string) (map[string]string, error) {
    return godotenv.Read(path)
}

// LoadEncryptedEnv decrypts a Laravel-encrypted .env.encrypted file and parses it.
func LoadEncryptedEnv(path string, key []byte) (map[string]string, error) {
    payload, err := readAndDecode(path)
    if err != nil {
        return nil, err
    }

    iv, ct, err := parseEnvelope(payload)
    if err != nil {
        return nil, err
    }

    pt, err := decryptCBC(iv, ct, key)
    if err != nil {
        return nil, err
    }
    
    pt = unwrapPHP(pt)

    pt = ensureTrailingNewline(pt)
    return godotenv.Parse(strings.NewReader(string(pt)))
}

// readAndDecode reads the file and Base64-decodes its contents.
func readAndDecode(path string) ([]byte, error) {
    dataB64, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    return base64.StdEncoding.DecodeString(string(dataB64))
}

// parseEnvelope extracts IV and ciphertext from Laravel's JSON envelope.
// Falls back to raw IV||ciphertext if JSON unmarshal fails.
func parseEnvelope(payload []byte) (iv, ciphertext []byte, err error) {
    trimmed := strings.TrimLeft(string(payload), " \t\r\n")
    if len(trimmed) > 0 && trimmed[0] == '{' {
        var env laravelEnvelope
        if err = json.Unmarshal(payload, &env); err != nil {
            return
        }
        // Decode iv and ciphertext
        iv, err = base64.StdEncoding.DecodeString(env.IV)
        if err != nil {
            return
        }
        ciphertext, err = base64.StdEncoding.DecodeString(env.Value)
        return
    }
    // Raw IV||ciphertext
    if len(payload) < aes.BlockSize {
        err = fmt.Errorf("ciphertext too short")
        return
    }
    iv = payload[:aes.BlockSize]
    ciphertext = payload[aes.BlockSize:]
    return
}

// unwrapPHP checks for PHP's s:<len>:"..."; wrapper and returns the inner text.
func unwrapPHP(pt []byte) []byte {
    s := string(pt)
    // look for s:<number>:"..."; at the start
    if strings.HasPrefix(s, "s:") {
        // first quote
        if fq := strings.Index(s, `"`); fq >= 0 {
            // last occurrence of `";`
            if lq := strings.LastIndex(s, `";`); lq > fq {
                return []byte(s[fq+1 : lq])
            }
        }
    }
    return pt
}

// decryptCBC decrypts AES-CBC ciphertext and strips PKCS#7 padding.
func decryptCBC(iv, ct, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("invalid key: %w", err)
    }
    if len(ct)%aes.BlockSize != 0 {
        return nil, fmt.Errorf("ciphertext not multiple of block size")
    }
    pt := make([]byte, len(ct))
    cipher.NewCBCDecrypter(block, iv).CryptBlocks(pt, ct)

    padLen := int(pt[len(pt)-1])
    if padLen <= 0 || padLen > aes.BlockSize {
        return nil, fmt.Errorf("invalid padding")
    }
    for i := len(pt) - padLen; i < len(pt); i++ {
        if pt[i] != byte(padLen) {
            return nil, fmt.Errorf("padding mismatch at %d", i)
        }
    }
    return pt[:len(pt)-padLen], nil
}

// ensureTrailingNewline appends a newline if the buffer doesn't end with one.
func ensureTrailingNewline(pt []byte) []byte {
    if len(pt) > 0 && pt[len(pt)-1] != '\n' {
        return append(pt, '\n')
    }
    return pt
}

// LoadSchema reads and parses a YAML schema file.
func LoadSchema(path string) (schema.Schema, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    s := make(schema.Schema)
    if err := yaml.Unmarshal(data, &s); err != nil {
        return nil, err
    }
    return s, nil
}
