package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// func EncryptString(data string, key string) (string, error) {
// 	result, err := EncryptData([]byte(data), []byte(key))
// 	return string(result), err
// }
//
// func DecryptString(data string, key string) (string, error) {
// 	result, err := EncryptData([]byte(data), []byte(key))
// 	return string(result), err
// }

// EncryptData encrypts given text with usage of given key
func EncryptData(data []byte, key []byte) ([]byte, error) {
	// CBC mode works on blocks so data may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the data is already of the correct length.
	if len(data)%aes.BlockSize != 0 {
		// return nil, errors.New("data is not a multiple of the block size")
		diff := len(data) - aes.BlockSize
		for diff < 0 {
			data = append(data, 0)
			diff++
		}
	}
	if len(key)%aes.BlockSize != 0 {
		// return nil, errors.New("data is not a multiple of the block size")
		diff := len(key) - aes.BlockSize
		for diff < 0 {
			key = append(key, 0)
			diff++
		}
	}

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(cipherBlock, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], data)

	return cipherText, nil
}

// DecryptData decrypts given text with usage of given key
func DecryptData(encodedData []byte, key []byte) ([]byte, error) {

	if len(key)%aes.BlockSize != 0 {
		// return nil, errors.New("data is not a multiple of the block size")
		diff := len(key) - aes.BlockSize
		for diff < 0 {
			key = append(key, 0)
			diff++
		}
	}

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipher text.
	if len(encodedData) < aes.BlockSize {
		return nil, errors.New("data too short")
	}
	iv := encodedData[:aes.BlockSize]
	encodedData = encodedData[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(encodedData)%aes.BlockSize != 0 {
		panic("data is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(cipherBlock, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	result := make([]byte, len(encodedData))
	mode.CryptBlocks(result, encodedData)

	return result, nil
}

func testEncript() {
	data := "some text"
	key := "test"
	encrypted, err := EncryptData([]byte(data), []byte(key))
	if err != nil {
		fmt.Println(err.Error())
	}
	decrypted, err := DecryptData(encrypted, []byte(key))
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%x\n", encrypted)
	fmt.Printf("%s\n", decrypted)
}
