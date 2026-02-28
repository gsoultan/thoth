package objects

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rc4"
	"crypto/sha256"
	"fmt"
	"io"
)

// EncryptionContext holds keys and data for PDF encryption.
type EncryptionContext struct {
	O          []byte // Owner key
	U          []byte // User key
	P          int32  // Permissions
	FileID     []byte // ID from trailer
	EncryptKey []byte // Global encryption key
	Algorithm  int    // 1 (RC4), 2 (RC4 128), 4 (AES 128), 5 (AES 256)
	Revision   int    // 2 (40-bit), 3 (128-bit), 4 (Acrobat 7), 5 (Acrobat 9), 6 (Acrobat X)
}

// NewEncryptionContext creates a new EncryptionContext from a password and file ID.
func NewEncryptionContext(password string, fileID []byte) *EncryptionContext {
	ec := &EncryptionContext{
		P:         -4, // Default permissions (all allowed)
		FileID:    fileID,
		Algorithm: 5, // AES 256
		Revision:  5, // Revision 5 (Acrobat 9)
	}

	// 1. Derivation of Encryption Key (Standard Security Handler Algorithm 3.2 for R5)
	if ec.Revision >= 5 {
		// Simplified Revision 5 implementation
		h := sha256.New()
		h.Write([]byte(password))
		h.Write(fileID)
		ec.EncryptKey = h.Sum(nil) // 32 bytes for AES-256

		// In real R5, O and U are much more complex.
		// We'll use dummy O/U for now, but the encryption logic will be real AES.
		ec.O = make([]byte, 32)
		ec.U = make([]byte, 32)
		copy(ec.O, []byte("OWNERKEYPADDINGOWNERKEYPADDINGOK"))
		copy(ec.U, []byte("USERKEYPADDINGUSERKEYPADDINGUSER"))
	} else {
		// Standard Security Handler Algorithm 3.2 (Revision 2/3)
		h := md5.New()
		h.Write([]byte(password))
		h.Write(fileID)
		ec.EncryptKey = h.Sum(nil)

		ec.O = make([]byte, 32)
		ec.U = make([]byte, 32)
		copy(ec.O, []byte("OWNERKEYPADDINGOWNERKEYPADDINGOK"))
		copy(ec.U, []byte("USERKEYPADDINGUSERKEYPADDINGUSER"))
	}

	return ec
}

// Encrypt encrypts data using the object's ID and generation.
func (ec *EncryptionContext) Encrypt(data []byte, objNum, objGen int) []byte {
	if ec == nil {
		return data
	}

	// 1. Create object-specific key
	key := make([]byte, len(ec.EncryptKey)+5)
	copy(key, ec.EncryptKey)
	key[len(ec.EncryptKey)] = byte(objNum)
	key[len(ec.EncryptKey)+1] = byte(objNum >> 8)
	key[len(ec.EncryptKey)+2] = byte(objNum >> 16)
	key[len(ec.EncryptKey)+3] = byte(objGen)
	key[len(ec.EncryptKey)+4] = byte(objGen >> 8)

	h := md5.New()
	h.Write(key)
	objKey := h.Sum(nil)

	keyLen := len(ec.EncryptKey) + 5
	if keyLen > 16 && ec.Algorithm < 5 {
		keyLen = 16
	}
	finalKey := objKey
	if keyLen < len(objKey) {
		finalKey = objKey[:keyLen]
	}

	if ec.Algorithm >= 4 {
		// AES 128/256 (Algorithm 3.2 in PDF 1.7)
		block, err := aes.NewCipher(finalKey)
		if err != nil {
			return data
		}

		// PKCS7 padding
		blockSize := block.BlockSize()
		padding := blockSize - len(data)%blockSize
		padtext := make([]byte, padding)
		for i := range padding {
			padtext[i] = byte(padding)
		}
		paddedData := append(data, padtext...)

		// 16 bytes IV
		iv := make([]byte, blockSize)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return data
		}

		mode := cipher.NewCBCEncrypter(block, iv)
		dst := make([]byte, len(paddedData))
		mode.CryptBlocks(dst, paddedData)

		// Result is IV + encrypted data
		return append(iv, dst...)
	}

	// RC4
	cipher, _ := rc4.NewCipher(finalKey)
	dst := make([]byte, len(data))
	cipher.XORKeyStream(dst, data)
	return dst
}

// WriteEncrypted writes data as a hex string or literal string, encrypted if necessary.
func (ec *EncryptionContext) WriteEncrypted(w io.Writer, data []byte, objNum, objGen int) (int64, error) {
	encrypted := ec.Encrypt(data, objNum, objGen)
	// For simplicity, we write as a hex string <...> if encrypted, or let PDFString handle it.
	// Actually, PDFString handles literal ( ... ).
	// If encrypted, hex string is often safer.
	n, err := fmt.Fprintf(w, "<%x>", encrypted)
	return int64(n), err
}
