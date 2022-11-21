package UrlsRouter

import (
	helmRouter "automation-suite/HelmAppRouter"
	"automation-suite/testUtils"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
)

func (suite *UrlsTestSuite) TestGetUrlsForHelmApp() {
	envConf, _ := helmRouter.GetEnvironmentConfigForHelmApp()
	queryParams := map[string]string{"appId": envConf.HAppId}
	log.Println("Hitting Get urls API")
	resp := HitGetUrls(queryParams, suite.authToken)
	assert.Equal(suite.T(), 200, resp.Code)
}
func (suite *UrlsTestSuite) TestGetUrlsForHelmAppWithIncorrectAppId() {
	randomHAppId := testUtils.GetRandomNumberOf9Digit()
	queryParams := map[string]string{"appId": strconv.Itoa(randomHAppId)}
	resp := HitGetUrls(queryParams, suite.authToken)
	assert.Equal(suite.T(), 400, resp.Code)
	assert.Equal(suite.T(), "Bad Request", resp.Status)
}

func (suite *UrlsTestSuite) TestGetUrlsForDevtronApp() {
	envConf, _ := GetEnvironmentConfigForDevtronApp()
	queryParams := map[string]string{"appId": envConf.AppId, "envId": envConf.EnvId}
	log.Println("Hitting Get urls API")
	resp := HitGetUrls(queryParams, suite.authToken)
	assert.Equal(suite.T(), 200, resp.Code)
}

func (suite *UrlsTestSuite) TestGetUrlsForDevtronAppWithIncorrectAppId() {
	randomInstalledAppId := "installedAppId-1"
	randomEnvId := "envid-1"
	queryParams := map[string]string{"installedAppId": randomInstalledAppId, "envId": randomEnvId}
	resp := HitGetUrls(queryParams, suite.authToken)
	assert.Equal(suite.T(), 400, resp.Code)
	assert.Equal(suite.T(), "Bad Request", resp.Status)
}

func (suite *UrlsTestSuite) TestGetUrlsForInstalledApp() {
	envConf, _ := GetEnvironmentConfigForInstalledApp()
	queryParams := map[string]string{"appId": envConf.InstalledAppId, "envId": envConf.EnvId}
	log.Println("Hitting Get urls API")
	resp := HitGetUrls(queryParams, suite.authToken)
	assert.Equal(suite.T(), 200, resp.Code)
}

func (suite *UrlsTestSuite) TestGetUrlsForInstalledAppWithIncorrectAppId() {
	randomAppId := "appId-1"
	randomEnvId := "envid-1"
	queryParams := map[string]string{"appId": randomAppId, "envId": randomEnvId}
	resp := HitGetUrls(queryParams, suite.authToken)
	assert.Equal(suite.T(), 400, resp.Code)
	assert.Equal(suite.T(), "Bad Request", resp.Status)
}
