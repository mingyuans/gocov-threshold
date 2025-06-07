package pr

import "github.com/mingyuans/gocov-threshold/cmd/threshold/arg"

type Service struct {
	env      Environment
	info     GitHubPRInfo
	inputArg arg.Arg
}

func NewService(inputArg arg.Arg) *Service {
	env := getPREnvironment()
	info, getInfoErr := gettPRInfo(env, inputArg)
	if getInfoErr != nil {
		panic("Failed to get PR info: " + getInfoErr.Error())
	}
	return &Service{
		env:      env,
		info:     info,
		inputArg: inputArg,
	}
}
