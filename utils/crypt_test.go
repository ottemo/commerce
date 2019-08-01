package utils

import (
	"crypto/rand"
	"sync"
	"testing"
)

func TestEncryptString(t *testing.T) {
	unencrypted := "Some Data"
	encrypted := EncryptString(unencrypted)
	if encrypted == unencrypted {
		t.Errorf("Encryption error")
	}
	decrypted := DecryptString(encrypted)
	if decrypted != unencrypted {
		t.Error("decryption error: " + decrypted + " != " + unencrypted)
	}
}
func TestPasswordCheck(t *testing.T) {
	enc := PasswordEncode("pass", "salt")
	if PasswordCheck("pass", enc) {
		t.Error("password improperly encoded")
	}
}

//tests get and set key
func TestKey(t *testing.T) {
	cryptKey := make([]byte, 16)
	rand.Read(cryptKey)
	err := SetKey(cryptKey)
	if err != nil {
		t.Error("GetKey failiure")
	}
	if string(GetKey()) != string(cryptKey) {
		t.Error("Getkey failiure 2: " + string(cryptKey) + " != " + string(GetKey()))
	}
}
func TestDecryptData(t *testing.T) {
	unencrypted := "un"
	encrypted, _ := EncryptData([]byte(unencrypted))
	if string(encrypted) == unencrypted {
		t.Errorf("Encryption error")
	}
	decrypted, _ := DecryptData(encrypted)
	if string(decrypted) != unencrypted {
		t.Error("decryption error: " + string(decrypted) + " != " + unencrypted)
	}
}
func TestDecryptURLString(t *testing.T) {
	tester := "URL String"
	encrypted := CryptAsURLString(tester)
	if value, err := DecryptURLString(encrypted); err != nil || value != tester {
		t.Error("URL decryption error")
	}
}
func TestAsynchronous(t *testing.T) {
	var mutex sync.Mutex
	f := func(i int) {
		mutex.Lock()
		TestKey(t)
		TestEncryptString(t)
		TestDecryptData(t)
		TestDecryptURLString(t)
		mutex.Unlock()
	}
	var wg sync.WaitGroup
	for i := 0; i < 9999; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f(i)
		}()
	}
	wg.Wait()
}
