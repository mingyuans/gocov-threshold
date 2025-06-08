package pr

import "github.com/mingyuans/gocov-threshold/cmd/threshold/model"

type Service struct {
	env      Environment
	info     GitHubPRInfo
	inputArg model.Arg
}

func NewService(inputArg model.Arg) *Service {
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
