package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
)

// EncryptionService handles encryption and decryption of sensitive payroll data
type EncryptionService struct {
	key []byte
}

// NewEncryptionService creates a new encryption service instance
func NewEncryptionService() (*EncryptionService, error) {
	// Get encryption key from config or use a fallback
	key := getEncryptionKey()
	
	// Validate key length for AES-256
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes for AES-256")
	}
	
	return &EncryptionService{key: key}, nil
}

// getEncryptionKey retrieves the encryption key from config or uses a fallback
func getEncryptionKey() []byte {
	if config.AppConfig != nil && config.AppConfig.Encryption.Key != "" {
		// Use configured key (pad or truncate to 32 bytes)
		key := []byte(config.AppConfig.Encryption.Key)
		return padKeyTo32Bytes(key)
	}
	
	// Fallback key for development (should be overridden in production)
	fallbackKey := "hrms-payroll-encryption-key-2024!"
	return []byte(fallbackKey)[:32] // Ensure 32 bytes
}

// padKeyTo32Bytes ensures the key is exactly 32 bytes
func padKeyTo32Bytes(key []byte) []byte {
	if len(key) == 32 {
		return key
	}
	
	// Pad with zeros if shorter
	if len(key) < 32 {
		padded := make([]byte, 32)
		copy(padded, key)
		return padded
	}
	
	// Truncate if longer
	return key[:32]
}

// EncryptPayrollData encrypts sensitive payroll data using employee-specific key derivation
func (s *EncryptionService) EncryptPayrollData(employeeID uint, data interface{}) (string, error) {
	// Convert data to JSON bytes
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}
	
	// Derive employee-specific key
	employeeKey := s.deriveEmployeeKey(employeeID)
	
	// Encrypt the data
	encrypted, err := s.encryptBytes(dataBytes, employeeKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %v", err)
	}
	
	// Return base64 encoded encrypted data
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptPayrollData decrypts payroll data using employee-specific key derivation
func (s *EncryptionService) DecryptPayrollData(employeeID uint, encryptedData string, result interface{}) error {
	// Decode base64
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to decode base64: %v", err)
	}
	
	// Derive employee-specific key
	employeeKey := s.deriveEmployeeKey(employeeID)
	
	// Decrypt the data
	decryptedBytes, err := s.decryptBytes(encryptedBytes, employeeKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %v", err)
	}
	
	// Unmarshal JSON back to result
	if err := json.Unmarshal(decryptedBytes, result); err != nil {
		return fmt.Errorf("failed to unmarshal decrypted data: %v", err)
	}
	
	return nil
}

// deriveEmployeeKey creates a unique encryption key for each employee
func (s *EncryptionService) deriveEmployeeKey(employeeID uint) []byte {
	// Combine master key with employee ID to create unique key
	employeeKey := fmt.Sprintf("%s-%d", string(s.key), employeeID)
	
	// Hash to get consistent 32-byte key
	return padKeyTo32Bytes([]byte(employeeKey))
}

// encryptBytes performs AES-GCM encryption
func (s *EncryptionService) encryptBytes(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// decryptBytes performs AES-GCM decryption
func (s *EncryptionService) decryptBytes(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// EncryptField encrypts a single field value
func (s *EncryptionService) EncryptField(employeeID uint, value string) (string, error) {
	if value == "" {
		return "", nil
	}
	
	employeeKey := s.deriveEmployeeKey(employeeID)
	
	// For short fields, we can encrypt directly
	if len(value) <= 100 {
		encrypted, err := s.encryptBytes([]byte(value), employeeKey)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(encrypted), nil
	}
	
	// For longer fields, use structured encryption
	fieldData := map[string]string{"value": value}
	return s.EncryptPayrollData(employeeID, fieldData)
}

// DecryptField decrypts a single field value
func (s *EncryptionService) DecryptField(employeeID uint, encryptedValue string) (string, error) {
	if encryptedValue == "" {
		return "", nil
	}
	
	// Try to decrypt as simple field first
	if !strings.Contains(encryptedValue, "{") {
		decryptedBytes, err := s.decryptFieldBytes(employeeID, encryptedValue)
		if err == nil {
			return string(decryptedBytes), nil
		}
	}
	
	// Fall back to structured decryption
	var result map[string]string
	if err := s.DecryptPayrollData(employeeID, encryptedValue, &result); err != nil {
		return "", err
	}
	
	return result["value"], nil
}

// decryptFieldBytes helper for simple field decryption
func (s *EncryptionService) decryptFieldBytes(employeeID uint, encryptedValue string) ([]byte, error) {
	employeeKey := s.deriveEmployeeKey(employeeID)
	
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return nil, err
	}
	
	return s.decryptBytes(encryptedBytes, employeeKey)
}

// IsEncrypted checks if a value appears to be encrypted
func (s *EncryptionService) IsEncrypted(value string) bool {
	// Simple heuristic: encrypted data is base64 encoded and doesn't look like plain text
	if _, err := base64.StdEncoding.DecodeString(value); err != nil {
		return false
	}
	
	// Additional checks could be added here
	return len(value) > 20 && !strings.Contains(value, " ") && strings.Contains(value, "=")
}