package ChartRepositoryRouter

import "github.com/stretchr/testify/assert"

func (suite *ChartRepoTestSuite) TestClassC6TriggerChartSyncManual() {

	suite.Run("A=1=TriggerChartSyncManual", func() {
		triggerChartSyncApiResponse := HitTriggerChartSyncManualApi(suite.authToken)
		assert.Equal(suite.T(), "ok", triggerChartSyncApiResponse.Result.Status)
	})
}
