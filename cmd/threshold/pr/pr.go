package pr

type Service struct {
	env  Environment
	info GitHubPRInfo
}

func NewService() *Service {
	env := getPREnvironment()
	info, getInfoErr := gettPRInfo(env)
	if getInfoErr != nil {
		panic("Failed to get PR info: " + getInfoErr.Error())
	}
	return &Service{
		env:  env,
		info: info,
	}
}
