package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 2
	saltLength  = 16
	keyLength   = 32
)

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func Hash(password string) (string, error) {
	salt, err := GenerateSalt(saltLength)
	if err != nil {
		return "", err
	}

	// Tính toán băm Argon2
	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	// Mã hóa kết quả thành chuỗi base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Kết hợp các thông tin thành một chuỗi duy nhất
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// Hàm kiểm tra mật khẩu
func Verify(password, encodedHash string) (bool, error) {
	// Phân tích chuỗi băm đã mã hóa
	parts := SplitEncodedHash(encodedHash)
	if parts == nil {
		return false, fmt.Errorf("invalid hash format")
	}

	// Giải mã salt và hash từ base64
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Tính toán băm mới từ mật khẩu đầu vào
	newHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))

	// So sánh hai băm
	if subtle.ConstantTimeCompare(hash, newHash) == 1 {
		return true, nil
	}
	return false, nil
}

// Hàm phân tích chuỗi băm đã mã hóa
func SplitEncodedHash(encodedHash string) []string {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil
	}
	return parts
}
