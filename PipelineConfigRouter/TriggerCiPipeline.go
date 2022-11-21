package PipelineConfigRouter

import (
	"automation-suite/HelperRouter"
	Base "automation-suite/testUtils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/sjson"
	"log"
	"strconv"
	"time"
)

func (suite *PipelinesConfigRouterTestSuite) TestClassD3TriggerCiPipeline() {
	config, _ := GetEnvironmentConfigPipelineConfigRouter()
	var configId int
	log.Println("=== Here we are creating an App ===")
	createAppApiResponse := Base.CreateApp(suite.authToken).Result
	log.Println("=== Here we are printing AppName ===>", createAppApiResponse.AppName)
	log.Println("=== Here we are creating App Material ===")
	createAppMaterialRequestDto := GetAppMaterialRequestDto(createAppApiResponse.Id, 1, false)
	appMaterialByteValue, _ := json.Marshal(createAppMaterialRequestDto)
	createAppMaterialResponse := HitCreateAppMaterialApi(appMaterialByteValue, createAppApiResponse.Id, 1, false, suite.authToken)

	log.Println("=== Here we are saving docker build config ===")
	requestPayloadForSaveAppCiPipeline := GetRequestPayloadForSaveAppCiPipeline(createAppApiResponse.Id, config.DockerRegistry, config.DockerRegistry+"/test", config.DockerfilePath, config.DockerfileRepository, config.DockerfileRelativePath, createAppMaterialResponse.Result.Material[0].Id)
	byteValueOfSaveAppCiPipeline, _ := json.Marshal(requestPayloadForSaveAppCiPipeline)
	HitSaveAppCiPipeline(byteValueOfSaveAppCiPipeline, suite.authToken)

	log.Println("=== Here we are fetching latestChartReferenceId ===")
	time.Sleep(2 * time.Second)
	getChartReferenceResponse := HitGetChartReferenceViaAppId(strconv.Itoa(createAppApiResponse.Id), suite.authToken)
	latestChartRef := getChartReferenceResponse.Result.LatestChartRef

	log.Println("=== Here we are fetching Template using getAppTemplateAPI ===")
	getTemplateResponse := HitGetTemplateViaAppIdAndChartRefId(strconv.Itoa(createAppApiResponse.Id), strconv.Itoa(latestChartRef), suite.authToken)

	log.Println("=== Here we are fetching DefaultAppOverride from template response ===")
	defaultAppOverride := getTemplateResponse.Result.GlobalConfig.DefaultAppOverride

	log.Println("=== Here we are creating payload for SaveTemplate API ===")
	saveDeploymentTemplate := GetRequestPayloadForSaveDeploymentTemplate(createAppApiResponse.Id, latestChartRef, defaultAppOverride)
	byteValueOfSaveDeploymentTemplate, _ := json.Marshal(saveDeploymentTemplate)
	jsonOfSaveDeploymentTemp := string(byteValueOfSaveDeploymentTemplate)
	jsonWithMicroserviceToleration, _ := sjson.Set(jsonOfSaveDeploymentTemp, "valuesOverride.tolerations.0", map[string]interface{}{"effect": "NoSchedule", "key": "microservice", "operator": "Equal", "value": "true"})
	finalJson, _ := sjson.Set(jsonWithMicroserviceToleration, "valuesOverride.tolerations.1", map[string]interface{}{"effect": "NoSchedule", "key": "kubernetes.azure.com/scalesetpriority", "operator": "Equal", "value": "spot"})
	updatedByteValueOfSaveDeploymentTemplate := []byte(finalJson)

	log.Println("=== Here we are hitting SaveTemplate API ===")
	HitSaveDeploymentTemplateApi(updatedByteValueOfSaveDeploymentTemplate, suite.authToken)

	log.Println("=== Here we are saving Global Configmap ===")
	requestPayloadForConfigMap := HelperRouter.GetRequestPayloadForSecretOrConfig(0, "-config1", createAppApiResponse.Id, "environment", "kubernetes", false, false, false, false)
	byteValueOfSaverConfigMap, _ := json.Marshal(requestPayloadForConfigMap)
	globalConfigMap := HelperRouter.HitSaveGlobalConfigMap(byteValueOfSaverConfigMap, suite.authToken)
	configId = globalConfigMap.Result.Id

	log.Println("=== Here we are saving Global Secret ===")
	requestPayloadForSecret := HelperRouter.GetRequestPayloadForSecretOrConfig(configId, "-secret1", createAppApiResponse.Id, "environment", "kubernetes", false, false, true, false)
	byteValueOfSecret, _ := json.Marshal(requestPayloadForSecret)
	HelperRouter.HitSaveGlobalSecretApi(byteValueOfSecret, suite.authToken)

	log.Println("=== Here we are saving workflow with Pre/Post CI ===")
	workflowResponse := HitCreateWorkflowApiWithFullPayload(createAppApiResponse.Id, suite.authToken).Result

	preStageScript, _ := Base.GetByteArrayOfGivenJsonFile("../testdata/PipeLineConfigRouter/preStageScript.txt")
	postStageScript, _ := Base.GetByteArrayOfGivenJsonFile("../testdata/PipeLineConfigRouter/postStageScript.txt")

	log.Println("=== Here we are saving CD pipeline ===")
	payload := GetRequestPayloadForSaveCdPipelineApi(createAppApiResponse.Id, workflowResponse.AppWorkflowId, 1, workflowResponse.CiPipelines[0].Id, workflowResponse.CiPipelines[0].ParentCiPipeline, Automatic, string(preStageScript), string(postStageScript), Automatic)
	bytePayload, _ := json.Marshal(payload)
	savePipelineResponse := HitSaveCdPipelineApi(bytePayload, suite.authToken)
	time.Sleep(2 * time.Second)

	log.Println("=== Here we are getting pipeline material ===")
	pipelineMaterial := HitGetCiPipelineMaterial(workflowResponse.CiPipelines[0].Id, suite.authToken)

	//here we are hitting GetWorkFlow API 2 time one just after the triggerCiPipeline and one after 4 minutes of triggering
	suite.Run("A=1=TriggerCiPipelineWithValidPayload", func() {
		payloadForTriggerCiPipeline := CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.CiPipelines[0].Id, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		triggerCiPipelineResponse := HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		if triggerCiPipelineResponse.Result.AuthStatus != "allowed for all pipelines" {
			time.Sleep(2 * time.Second)
			triggerCiPipelineResponse = HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
			assert.Equal(suite.T(), "allowed for all pipelines", triggerCiPipelineResponse.Result.AuthStatus)
			assert.NotNil(suite.T(), triggerCiPipelineResponse.Result.ApiResponse)
		}
		time.Sleep(10 * time.Second)
		log.Println("=== Here we are getting workflow after triggering ===")
		workflowStatus := HitGetWorkflowStatus(createAppApiResponse.Id, suite.authToken)
		if workflowStatus.Result.CiWorkflowStatus[0].CiStatus == "Starting" {
			time.Sleep(5 * time.Second)
			workflowStatus = HitGetWorkflowStatus(createAppApiResponse.Id, suite.authToken)
			assert.Equal(suite.T(), "Running", workflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		} else {
			assert.Equal(suite.T(), "Running", workflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		}
		log.Println("=== Here we are getting workflow and verifying the status after triggering via poll function ===")
		assert.True(suite.T(), PollForGettingCdDeployStatusAfterTrigger(createAppApiResponse.Id, suite.authToken))
		updatedWorkflowStatus := HitGetWorkflowStatus(createAppApiResponse.Id, suite.authToken)
		assert.Equal(suite.T(), "Succeeded", updatedWorkflowStatus.Result.CiWorkflowStatus[0].CiStatus)
		assert.Equal(suite.T(), "Healthy", updatedWorkflowStatus.Result.CdWorkflowStatus[0].DeployStatus)
	})

	suite.Run("A=2=TriggerCiPipelineWithInvalidateCacheAsFalse", func() {
		payloadForTriggerCiPipeline := CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.CiPipelines[0].Id, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		triggerCiPipelineResponse := HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		if triggerCiPipelineResponse.Result.AuthStatus != "allowed for all pipelines" {
			time.Sleep(2 * time.Second)
			triggerCiPipelineResponse = HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
			assert.Equal(suite.T(), "allowed for all pipelines", triggerCiPipelineResponse.Result.AuthStatus)
			assert.NotNil(suite.T(), triggerCiPipelineResponse.Result.ApiResponse)
		}
	})

	suite.Run("A=3=TriggerCiPipelineWithInvalidPipelineId", func() {
		invalidPipeLineId := Base.GetRandomNumberOf9Digit()
		payloadForTriggerCiPipeline := CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, invalidPipeLineId, pipelineMaterial.Result[0].Id, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		triggerCiPipelineResponse := HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		assert.Equal(suite.T(), "pg: no rows in result set", triggerCiPipelineResponse.Errors[0].UserMessage)
	})

	suite.Run("A=4=TriggerCiPipelineWithInvalidPipelineId", func() {
		invalidMaterialId := Base.GetRandomNumberOf9Digit()
		payloadForTriggerCiPipeline := CreatePayloadForTriggerCiPipeline(pipelineMaterial.Result[0].History[0].Commit, workflowResponse.CiPipelines[0].Id, invalidMaterialId, true)
		bytePayloadForTriggerCiPipeline, _ := json.Marshal(payloadForTriggerCiPipeline)
		triggerCiPipelineResponse := HitTriggerCiPipelineApi(bytePayloadForTriggerCiPipeline, suite.authToken)
		assert.Equal(suite.T(), "[{pg: no rows in result set}]", triggerCiPipelineResponse.Errors[0].InternalMessage)
	})

	log.Println("=== Here we are Deleting the CD pipeline ===")
	deletePipelinePayload := GetPayloadForDeleteCdPipeline(createAppApiResponse.Id, savePipelineResponse.Result.Pipelines[0].Id)
	deletePipelineByteCode, _ := json.Marshal(deletePipelinePayload)
	HitForceDeleteCdPipelineApi(deletePipelineByteCode, suite.authToken)

	log.Println("=== Here we are Deleting the CI pipeline ===")
	DeleteCiPipeline(createAppApiResponse.Id, workflowResponse.CiPipelines[0].Id, suite.authToken)
	log.Println("=== Here we are Deleting CI Workflow ===")
	HitDeleteWorkflowApi(createAppApiResponse.Id, workflowResponse.AppWorkflowId, suite.authToken)
	log.Println("=== Here we are Deleting the app after all verifications ===")
	Base.DeleteApp(createAppApiResponse.Id, createAppApiResponse.AppName, createAppApiResponse.TeamId, createAppApiResponse.TemplateId, suite.authToken)
}

func PollForGettingCdDeployStatusAfterTrigger(id int, authToken string) bool {
	count := 0
	for {
		updatedWorkflowStatus := HitGetWorkflowStatus(id, authToken)
		deploymentStatus := updatedWorkflowStatus.Result.CdWorkflowStatus[0].DeployStatus
		time.Sleep(1 * time.Second)
		count = count + 1
		if deploymentStatus == "Healthy" || count >= 500 {
			break
		}
	}
	return true
}
