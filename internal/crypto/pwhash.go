// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = fmt.Errorf("argon2 invalid hash")
	ErrIncompatibleVersion = fmt.Errorf("argon2 incompatible version")
)

type argonParams struct {
	iterations  uint32 // time
	memory      uint32
	keyLength   uint32
	saltLength  uint32
	parallelism uint8
}

func generateHash(rawPw []byte, p *argonParams) ([]byte, []byte, *argonParams, error) {
	if p == nil {
		p = &argonParams{
			iterations:  3,
			memory:      64 * 1024,
			keyLength:   32,
			saltLength:  16,
			parallelism: 2,
		}
	}

	salt := make([]byte, p.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error generating random: %w", err)
	}

	hash := argon2.IDKey(rawPw, salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	return hash, salt, p, nil
}

func encodeHash(hash, salt []byte, p *argonParams) []byte {
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	// "standard" encoding: https://github.com/P-H-C/phc-winner-argon2#command-line-utility
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.memory, p.iterations, p.parallelism, b64Salt, b64Hash,
	)

	return []byte(encodedHash)
}

func decodePWHash(pwHash []byte) (hash, salt []byte, p *argonParams, err error) {
	vals := strings.Split(string(pwHash), "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &argonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return hash, salt, p, nil
}

func HashPW(rawPw []byte) ([]byte, error) {
	hash, salt, params, err := generateHash(rawPw, nil)
	if err != nil {
		return nil, fmt.Errorf("error generating pw hash: %w", err)
	}
	pwHash := encodeHash(hash, salt, params)
	return pwHash, nil
}

func CheckPWHash(pwHash, rawPw []byte) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	hash, salt, p, err := decodePWHash(pwHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(rawPw), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}
