package model

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/google/uuid"
	"strings"
	"time"
)

const (
	EV_KEY_ID                  = "id"
	EV_KEY_TRANSACTION_ID      = "transactionId"
	EV_KEY_TIMESTAMP           = "timestamp"
	EV_KEY_EVENT_SEQ           = "eventSeq"
	EV_KEY_NUM_EVENTS          = "numEvents"
	SAWTOOTH_CURRENT_BLOCK_ID  = "block_id"
	SAWTOOTH_PREVIOUS_BLOCK_ID = "previous_block_id"
)

const (
	SAWTOOTH_BLOCK_COMMIT       = "sawtooth/block-commit"
	EV_TYPE_TRANSACTION_CONTROL = "evTransactionControl"
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

const expectedAddressLength = 6 + 64

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
