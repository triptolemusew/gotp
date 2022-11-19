package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	KEY_LENGTH  = 32
	SALT_LENGTH = 8
)

func Encrypt(passphrase, data []byte) ([]byte, error) {
	computedKey, salt := DeriveKey(passphrase, nil)

	blockCipher, err := aes.NewCipher(computedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

func Decrypt(key, data []byte) ([]byte, error) {
	salt, computedData := data[len(data)-SALT_LENGTH:], data[:len(data)-SALT_LENGTH]

	key, _ = DeriveKey(key, salt)

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := computedData[:gcm.NonceSize()], computedData[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func DeriveKey(password, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		rand.Read(salt)
	}
	return pbkdf2.Key([]byte(password), salt, 10000, KEY_LENGTH, sha256.New), salt
}
