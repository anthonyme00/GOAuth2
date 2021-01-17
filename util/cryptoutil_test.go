package util

import (
	"math/rand"
	"testing"
)

func TestEncryptionDecryption(t *testing.T) {
	input := []byte("inputtest")
	key := []byte("secretkey")

	encrypted, err := Encrypt(input, key)
	if err != nil {
		t.Errorf("Unable to encrypt data")
		t.FailNow()
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Errorf("Unable to decrypt data")
		t.FailNow()
	}

	if string(decrypted) != string(input) {
		t.Errorf("Failed ! Expected %s got %s", string(input), string(decrypted))
		t.FailNow()
	}

	t.Log("Success !")
}

func TestEncryptionKeyLimit(t *testing.T) {
	var (
		encrypt []byte
		decrypt []byte
		err     error
	)

	input := []byte("dummyinput")
	key24bytes := []byte("thiskeyislessthan32bytes")
	encrypt, err = Encrypt(input, key24bytes)
	if err != nil || encrypt == nil {
		t.Errorf("Failed ! Should be able to encrypt with key of length less than 32 bytes")
		t.FailNow()
	}
	decrypt, err = Decrypt(encrypt, key24bytes)
	if err != nil || decrypt == nil {
		t.Errorf("Failed ! Should be able to decrypt with key of length less than 32 bytes")
		t.FailNow()
	}

	key32bytes := []byte("thiskeyisexactlyequalsto32bytes!")
	encrypt, err = Encrypt(input, key32bytes)
	if err != nil || encrypt == nil {
		t.Errorf("Failed ! Should be able to encrypt with key of length equals to 32 bytes")
		t.FailNow()
	}
	decrypt, err = Decrypt(encrypt, key32bytes)
	if err != nil || decrypt == nil {
		t.Errorf("Failed ! Should be able to decrypt with key of length equals to 32 bytes")
		t.FailNow()
	}

	key42bytes := []byte("thiskeyiswaylongerthan32bytes...shouldfail")
	encrypt, err = Encrypt(input, key42bytes)
	if err == nil || encrypt != nil {
		t.Errorf("Failed ! Should not be able to encrypt with key of length more than to 32 bytes")
		t.FailNow()
	}
	decrypt, err = Decrypt(encrypt, key42bytes)
	if err == nil || decrypt != nil {
		t.Errorf("Failed ! Should not be able to decrypt with key of length more than 32 bytes")
		t.FailNow()
	}

	t.Log("Success !")
}

func TestBase64Gen(t *testing.T) {
	for i := 0; i < 10; i++ {
		length := rand.Intn(64)
		base64 := GenerateBase64URLNoPadding(uint32(length))

		if len(base64) < length {
			t.Errorf("Failed ! Length of generated base64 number is not equal to the requested length\nRequested: %v , got %v", length, len(base64))
			t.FailNow()
		}
	}

	t.Log("Success !")
}
