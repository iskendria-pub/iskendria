package blockchain

import (
	"os"
)

const envVarIp = "ALEXANDRIA_IP"

func getIp() string {
	return os.Getenv(envVarIp)
}
