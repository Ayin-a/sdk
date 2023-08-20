package http

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"hk4e_sdk/pkg/logger"
	"os"
)

// 改index死妈
func (s *Secret) LoadSecret(mode bool) error {
	//var writeToPem bool = false

	s.PayPrivateKey, s.PayPublicKey = LoadPemData("data/keys/pay_private_key.pem", "data/keys/pay_public_key.pem")
	s.PasswordPrivateKey, s.PasswordPublicKey = LoadPemData("data/keys/password_private_key.pem", "data/keys/password_public_key.pem")

	return nil
}

// 改index死妈
func LoadPemData(priPemPath, pubPemPath string) (pri *PrivateKey, pub *PublicKey) {
	loadPrivateKey := func(path string) (*rsa.PrivateKey, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode(data)
		var privateKey *rsa.PrivateKey
		switch block.Type {
		case "PRIVATE KEY":
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			privateKey = key.(*rsa.PrivateKey)
		case "RSA PRIVATE KEY":
			privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		}
		return privateKey, nil
	}

	loadPublicKey := func(path string) (*rsa.PublicKey, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode(data)
		var publicKey *rsa.PublicKey
		switch block.Type {
		case "PUBLIC KEY":
			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			publicKey = key.(*rsa.PublicKey)
		case "RSA PUBLIC KEY":
			publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		}
		return publicKey, nil
	}

	sk, err := loadPrivateKey(priPemPath)
	if err != nil {
		logger.Error("read private key err")
	}

	pk, err := loadPublicKey(pubPemPath)
	if err != nil {
		logger.Error("read public key err")
	}

	return &PrivateKey{PrivateKey: sk}, &PublicKey{PublicKey: pk}
}
