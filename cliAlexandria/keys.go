package cliAlexandria

import (
	"encoding/hex"
	"errors"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"io/ioutil"
	"os"
	"strings"
)

func createKeyPair(publicKeyFile, privateKeyFile string) error {
	context := signing.NewSecp256k1Context()
	privateKey := context.NewRandomPrivateKey()
	publicKey := context.GetPublicKey(privateKey)
	privateKeyString := privateKey.AsHex()
	publicKeyString := publicKey.AsHex()
	modeRwRR := os.FileMode(0664)
	err := ioutil.WriteFile(publicKeyFile, []byte(publicKeyString), modeRwRR)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(privateKeyFile, []byte(privateKeyString), modeRwRR)
	if err != nil {
		return err
	}
	return nil
}

var LoggedInPublicKeyStr string
var LoggedInPublicKey signing.PublicKey
var LoggedInPrivateKey signing.PrivateKey

func IsLoggedIn() bool {
	return LoggedInPublicKeyStr != ""
}

func login(publicKeyFile, privateKeyFile string) error {
	publicKey, publicKeyAsString, err := readPublicKeyFile(publicKeyFile)
	if err != nil {
		return err
	}
	privateKey, err := readPrivateKeyFile(privateKeyFile)
	if err != nil {
		return err
	}
	err = signAndVerifyChallengeString(publicKey, privateKey)
	if err != nil {
		return err
	}
	LoggedInPublicKeyStr = publicKeyAsString
	LoggedInPublicKey = publicKey
	LoggedInPrivateKey = privateKey
	return nil
}

func readPublicKeyFile(publicKeyFile string) (signing.PublicKey, string, error) {
	publicKeyAsString, publicKeyBytes, err := readAndDecode(publicKeyFile)
	if err != nil {
		return nil, "", err
	}
	publicKey := signing.NewSecp256k1PublicKey(publicKeyBytes)
	return publicKey, publicKeyAsString, err
}

func readAndDecode(hexEncodedFile string) (string, []byte, error) {
	hexEncodedContents, err := ioutil.ReadFile(hexEncodedFile)
	if err != nil {
		return "", nil, err
	}
	decoded, err := decodeKey(hexEncodedContents)
	asString := string(hexEncodedContents)
	asString = strings.ToLower(asString)
	if err != nil {
		return asString, nil, err
	}
	return asString, decoded, nil
}

func decodeKey(encoded []byte) ([]byte, error) {
	trimmed := []byte(strings.TrimSpace(string(encoded)))
	if len(trimmed)%2 != 0 {
		return nil, errors.New("Keys should have an even number of bytes")
	}
	result := make([]byte, len(trimmed)/2)
	_, err := hex.Decode(result, trimmed)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func readPrivateKeyFile(privateKeyFile string) (signing.PrivateKey, error) {
	_, privateKeyBytes, err := readAndDecode(privateKeyFile)
	if err != nil {
		return nil, err
	}
	privateKey := signing.NewSecp256k1PrivateKey(privateKeyBytes)
	return privateKey, nil
}

func signAndVerifyChallengeString(publicKey signing.PublicKey, privateKey signing.PrivateKey) error {
	context := signing.CreateContext(privateKey.GetAlgorithmName())
	msg := []byte("some string")
	signature := context.Sign(msg, privateKey)
	if !context.Verify(signature, msg, publicKey) {
		return errors.New("ERROR: Invalid key pair")
	}
	return nil
}

func logout() error {
	LoggedInPublicKeyStr = ""
	LoggedInPublicKey = nil
	LoggedInPrivateKey = nil
	return nil
}
