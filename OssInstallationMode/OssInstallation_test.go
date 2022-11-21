package OssInstallationMode

import (
	"automation-suite/AppLabelsRouter"
	"automation-suite/AppListingRouter"
	"automation-suite/AttributesRouter"
	"automation-suite/ChartRepositoryRouter"
	"automation-suite/PipelineConfigRouter"
	"automation-suite/TeamRouter"
	"automation-suite/UserRouter"
	"automation-suite/dockerRegRouter"
	"automation-suite/externalLinkoutRouter"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestAttributeRouterSuite(t *testing.T) {
	suite.Run(t, new(AttributesRouter.AttributeRouterTestSuite))
}

//disabling this as after extracting gitOps as integration I have to handle this in different way
/*func TestGitOpsRouterSuite(t *testing.T) {
	suite.Run(t, new(GitopsConfigRouter.GitOpsRouterTestSuite))
}*/

func TestTeamRouterSuite(t *testing.T) {
	suite.Run(t, new(TeamRouter.TeamTestSuite))
}

func TestUserRouterSuite(t *testing.T) {
	suite.Run(t, new(UserRouter.UserTestSuite))
}

func TestChartRepoRouterSuite(t *testing.T) {
	suite.Run(t, new(ChartRepositoryRouter.ChartRepoTestSuite))
}

func TestAppLabelsRouterSuite(t *testing.T) {
	suite.Run(t, new(AppLabelsRouter.AppLabelRouterTestSuite))
}

func TestAppListingRouterSuite(t *testing.T) {
	suite.Run(t, new(AppListingRouter.AppsListingRouterTestSuite))
}

func TestDockerRegRouterSuite(t *testing.T) {
	suite.Run(t, new(dockerRegRouter.DockersRegRouterTestSuite))
}

func TestLinkOutRouterSuite(t *testing.T) {
	suite.Run(t, new(externalLinkoutRouter.LinkOutRouterTestSuite))
}

func TestPipelineConfigSuite(t *testing.T) {
	suite.Run(t, new(PipelineConfigRouter.PipelinesConfigRouterTestSuite))
}
