package command

import "gitlab.bbinfra.net/3estack/alexandria/model"

type Command struct {
	inputAddresses  []string
	outputAddresses []string
	command         *model.Command
}

func ApplyModelCommand(command *model.Command, signerKey string, ba BlockchainAccess) {

}
