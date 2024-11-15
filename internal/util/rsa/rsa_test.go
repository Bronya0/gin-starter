package rsa

import (
	"testing"
)

const (
	publicKey = `-----BEGIN PUBLIC KEY-----
xxx
-----END PUBLIC KEY-----`

	privateKey = `-----BEGIN RSA PRIVATE KEY-----
xxx
-----END RSA PRIVATE KEY-----`
)

func TestEncrypt(t *testing.T) {
	str, err := NewPublic(publicKey).Encrypt("123456")
	if err != nil {
		t.Error("rsa public encrypt error", err)
		return
	}

	t.Log(str)
}

func TestDecrypt(t *testing.T) {
	decryptStr := "KTKXckjkCLI6Vk_y_XROnY-a6nJpllruL-CX-v_2AFxfghA2tZ2nkQyS6d1-IIYMlgwm4ivwlzy-phLtaN9BB03htA5D9rwjA_JwYtqAG4iwuvgaDl2SiZ_H2ACv-aV1kNRgnyjh14hs0JiSt5VHEiJ3D2xYzOCKwtEzoo0WczJ-MYb3u_-bfcnm9YtvgtG5-y3Jy7WYr-IwXdBKqPO0E-jzrtY7m3Q1yC4znHdzjNpxCj0I6YRx4CZ362b706qNX7sl3E5KTJeSmYrsurB-SxQT1CaqGzVt7mshx1v2qGnv5NBNXpj7ZPKWGJbgaCUxcuxd1Mg0o81HnfbsGuSlFQ=="

	str, err := NewPrivate(privateKey).Decrypt(decryptStr)
	if err != nil {
		t.Error("rsa private decrypt error", err)
		return
	}

	t.Log(str)
}

func BenchmarkEncryptAndDecrypt(b *testing.B) {
	b.ResetTimer()
	rsaPublic := NewPublic(publicKey)
	rsaPrivate := NewPrivate(privateKey)
	for i := 0; i < b.N; i++ {
		encryptString, _ := rsaPublic.Encrypt("123456")
		_, err := rsaPrivate.Decrypt(encryptString)
		if err != nil {
			return
		}
	}
}
