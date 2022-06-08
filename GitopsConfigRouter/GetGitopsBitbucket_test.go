package GitopsConfigRouter

import (
	"encoding/json"
	"log"

	"github.com/stretchr/testify/assert"
)

func (suite *GitOpsRouterTestSuite) TestClassA2FetchAllGitopsConfig() {

	suite.Run("A=1=FetchAllGitopsConfig", func() {
		log.Println("Hitting GET api for Gitops config")
		fetchAllLinkResponseDto := HitFetchAllGitopsConfigApi(suite.authToken)
		noOfGitopsConfig := len(fetchAllLinkResponseDto.Result)

		log.Println("Hitting the 'Save Gitops Config' Api for creating a new entry")
		gitopsConfig, _ := GetGitopsConfig()

		createGitopsConfigRequestDto := GetGitopsConfigRequestDto(gitopsConfig.Provider, gitopsConfig.Username, gitopsConfig.Host, gitopsConfig.Token, gitopsConfig.GitHubOrgId)
		byteValueOfCreateGitopsConfig, _ := json.Marshal(createGitopsConfigRequestDto)

		log.Println("Hitting The post gitops config API")
		HitCreateGitopsConfigApi(byteValueOfCreateGitopsConfig, gitopsConfig.Provider, gitopsConfig.Username, gitopsConfig.Host, gitopsConfig.Token, gitopsConfig.GitHubOrgId, suite.authToken)

		log.Println("Hitting the HitFetchAllGitopsConfigApi again for verifying the functionality of it")
		fetchAllLinkResponseDto = HitFetchAllGitopsConfigApi(suite.authToken)

		log.Println("Validating the response of FetchAllLink API")

		// as response is not sending id or any parameter we are using if else using return code
		assert.Equal(suite.T(), noOfGitopsConfig+1, len(fetchAllLinkResponseDto.Result))

	})
}
