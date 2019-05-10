package model

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/google/uuid"
	"strings"
	"time"
)

const ID = "id"
const PDF = "PDF"

const (
	TRANSACTION_ID             = "transactionId"
	TIMESTAMP                  = "timestamp"
	EVENT_SEQ                  = "eventSeq"
	NUM_EVENTS                 = "numEvents"
	SAWTOOTH_CURRENT_BLOCK_ID  = "block_id"
	SAWTOOTH_PREVIOUS_BLOCK_ID = "previous_block_id"
)

const (
	EV_SAWTOOTH_BLOCK_COMMIT = "sawtooth/block-commit"
	EV_TRANSACTION_CONTROL   = "evTransactionControl"
)

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
