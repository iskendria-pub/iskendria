package model

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/google/uuid"
	"strings"
	"time"
)

const CREATED_ON = "createdOn"
const MODIFIED_ON = "modifiedOn"
const SIGNER = "signer"
const ID = "id"
const UTF8 = "UTF8"

const FamilyName = "alexandria"
const FamilyVersion = "1.0"

var Namespace string = Hexdigest(FamilyName)[:6]

func Hexdigest(s string) string {
	hash := sha512.New()
	hash.Write([]byte([]byte(s)))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}

func hexdigestOfUuid(uuid uuid.UUID) string {
	hash := sha512.New()
	hash.Write(uuid[:])
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}

var expectedAddressLength = 6 + 64

func getAddressPrefixFromAddress(address string) string {
	if len(address) != expectedAddressLength {
		panic("Address of invalid length encountered: " + address)
	}
	return address[6:8]
}

func GetCurrentTime() int64 {
	now := time.Now()
	return now.Unix()
}
