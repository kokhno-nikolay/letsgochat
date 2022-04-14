/* Package hasher provides hashing and password hash comparison. */
package hasher

import (
	"crypto/aes"
	"encoding/hex"
)

type Hasher interface {
	HashPassword(pass string) (string, error)
	CheckHashPassword(pass, hash string) bool
}

type Client struct {
	secretKey string
}

func NewHasher(secretKey string) Hasher {
	return &Client{
		secretKey: secretKey,
	}
}

/*
	HashPassword returns the AES hash of the password from 32-bit encryption key.
	The Advanced Encryption Standard (AES) aka Rijndael is an encryption algorithm
	created in 2001 by NIST. It uses 128-bit blocks of data to encrypt and is a
	symmetric block cipher.
*/
func (c *Client) HashPassword(pass string) (string, error) {
	cr, err := aes.NewCipher([]byte(c.secretKey))
	if err != nil {
		return "", err
	}

	out := make([]byte, len(pass))
	cr.Encrypt(out, []byte(pass))

	return hex.EncodeToString(out), nil
}

/*
	CheckHashPassword method compares password with hash and
	returns a boolean value.
*/
func (c *Client) CheckHashPassword(pass, hash string) bool {
	ciphertext, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}

	cr, err := aes.NewCipher([]byte(c.secretKey))
	if err != nil {
		return false
	}

	pt := make([]byte, len(ciphertext))
	cr.Decrypt(pt, ciphertext)

	s := string(pt[:])
	if pass != s {
		return true
	}

	return false
}
