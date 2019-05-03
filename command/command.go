package command

import "gitlab.bbinfra.net/3estack/alexandria/model"

type Command struct {
	inputAddresses  []string
	outputAddresses []string
	command         *model.Command
}
