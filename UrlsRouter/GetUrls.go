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
	resp := HitGetUrls(queryParams, suite.authToken, GetUrlsUrlHelm)
	assert.Equal(suite.T(), 200, resp.Code)
}
func (suite *UrlsTestSuite) TestGetUrlsForHelmAppWithIncorrectAppId() {
	randomHAppId := testUtils.GetRandomNumberOf9Digit()
	queryParams := map[string]string{"appId": strconv.Itoa(randomHAppId)}
	resp := HitGetUrls(queryParams, suite.authToken, GetUrlsUrlHelm)
	assert.Equal(suite.T(), 400, resp.Code)
	assert.Equal(suite.T(), "Bad Request", resp.Status)
}

func (suite *UrlsTestSuite) TestGetUrlsForDevtronApp() {
	envConf, _ := GetEnvironmentConfigForDevtronApp()
	queryParams := map[string]string{"appId": envConf.AppId, "envId": envConf.EnvId}
	log.Println("Hitting Get urls API")
	resp := HitGetUrls(queryParams, suite.authToken, GetUrlsUrl)
	assert.Equal(suite.T(), 200, resp.Code)
}

func (suite *UrlsTestSuite) TestGetUrlsForDevtronAppWithIncorrectAppId() {
	randomInstalledAppId := "installedAppId-1"
	randomEnvId := "envid-1"
	queryParams := map[string]string{"installedAppId": randomInstalledAppId, "envId": randomEnvId}
	resp := HitGetUrls(queryParams, suite.authToken, GetUrlsUrl)
	assert.Equal(suite.T(), 400, resp.Code)
	assert.Equal(suite.T(), "Bad Request", resp.Status)
}

func (suite *UrlsTestSuite) TestGetUrlsForInstalledApp() {
	envConf, _ := GetEnvironmentConfigForInstalledApp()
	queryParams := map[string]string{"installedAppId": envConf.InstalledAppId, "envId": envConf.EnvId}
	log.Println("Hitting Get urls API")
	resp := HitGetUrls(queryParams, suite.authToken, GetUrlsUrl)
	assert.Equal(suite.T(), 200, resp.Code)
}

func (suite *UrlsTestSuite) TestGetUrlsForInstalledAppWithIncorrectAppId() {
	randomAppId := "appId-1"
	randomEnvId := "envid-1"
	queryParams := map[string]string{"appId": randomAppId, "envId": randomEnvId}
	resp := HitGetUrls(queryParams, suite.authToken, GetUrlsUrl)
	assert.Equal(suite.T(), 400, resp.Code)
	assert.Equal(suite.T(), "Bad Request", resp.Status)
}

func (suite *UrlsTestSuite) TestGetUrlsdata() {
	expected := getTestExpectedUrlsData()
	envConf, _ := GetEnvironmentConfigForDevtronApp()
	queryParams := map[string]string{"appId": envConf.AppId, "envId": envConf.EnvId}
	log.Println("Hitting Get urls API")
	resp := HitGetUrls(queryParams, suite.authToken, GetUrlsUrl)
	assert.Equal(suite.T(), 200, resp.Code)
	assert.Equal(suite.T(), 3, len(resp.Result))
	for j, _ := range resp.Result {
		respData := resp.Result[j]
		assert.Equal(suite.T(), respData.Name, expected[j].Name)
		assert.Equal(suite.T(), respData.PointsTo, expected[j].PointsTo)
		assert.Equal(suite.T(), respData.Kind, expected[j].Kind)
		for i, url := range respData.Urls {
			assert.Equal(suite.T(), url, expected[j].Urls[i])
		}
	}
}

func getTestExpectedUrlsData() []UrlsResponse {
	res := make([]UrlsResponse, 0)
	res = append(res, UrlsResponse{
		Kind:     "Service",
		Name:     "ajay-test-devtron-demo-preview-service",
		PointsTo: "",
	})
	res = append(res, UrlsResponse{
		Kind:     "Service",
		Name:     "ajay-test-devtron-demo-service",
		PointsTo: "",
	})
	res = append(res, UrlsResponse{
		Kind:     "Ingress",
		Name:     "ajay-test-devtron-demo-ingress",
		PointsTo: "10.152.183.5",
		Urls:     []string{"chart-example1.local/example1", "chart-example2.local/example2", "chart-example2.local/example2/example2/healthz"},
	})

	return res
}
