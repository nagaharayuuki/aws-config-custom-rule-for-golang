package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
)

type AWSConfigService interface {
	PutEvaluations(*configservice.PutEvaluationsInput) (*configservice.PutEvaluationsOutput, error)
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, configEvent events.ConfigEvent) error {
	cSession := session.Must(session.NewSession())
	svc := configservice.New(cSession)
	err := handleRequestWithConfigService(ctx, configEvent, svc)
	if err != nil {
		return err
	}
	return nil
}

func handleRequestWithConfigService(ctx context.Context, configEvent events.ConfigEvent, svc AWSConfigService) error {
	var status string
	var invokingEvent InvokingEvent
	var configurationItem ConfigurationItem

    // 設定変更の呼び出しイベントを取得（https://docs.aws.amazon.com/ja_jp/config/latest/developerguide/evaluate-config_develop-rules_example-events.html）
	invokingEvent, err := getInvokingEvent([]byte(configEvent.InvokingEvent))

	if err != nil {
		panic(err)
	}

	configurationItem = invokingEvent.ConfigurationItem

	if params := getParams(configEvent, "excludeBuckets"); params != nil {
		for _, v := range params {
			if v == configurationItem.ResourceName {
				fmt.Println("Skipping over Compliance check for resource", v, "Params: excludeBuckets")
				status = "NOT_APPLICABLE"
			}
		}
	}

    // 評価対象リソースのイベントかを判定
	if isApplicable(configurationItem, configEvent) && status == "" {
        // 評価結果のステータス
		status = evaluateCompliance(configurationItem)
	} else {
        // 評価対象外のステータス
		status = "NOT_APPLICABLE"
	}

	var evaluations []*configservice.Evaluation

	evaluation := &configservice.Evaluation{
		ComplianceResourceId:   aws.String(configurationItem.ResourceID),
		ComplianceResourceType: aws.String(configurationItem.ResourceType),
		ComplianceType:         aws.String(status),
		OrderingTimestamp:      aws.Time(configurationItem.ConfigurationItemCaptureTime),
	}

	evaluations = append(evaluations, evaluation)
	putEvaluations := &configservice.PutEvaluationsInput{
		Evaluations: evaluations,
		ResultToken: aws.String(configEvent.ResultToken),
		TestMode:    aws.Bool(false),
	}

    // ConfigService へ評価結果を送信
	_, err = svc.PutEvaluations(putEvaluations)

	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	return nil
}

func evaluateCompliance(c ConfigurationItem) string {
	if c.ResourceType != "AWS::S3::Bucket" {
		return "NOT_APPLICABLE"
	}

    // ライフサイクル設定を取得
	bucketLifecycle := c.SupplementaryConfiguration.BucketLifecycleConfiguration

    // ライフサイクルが設定されているか
	if bucketLifecycle.Rules != nil {
		return "COMPLIANT"
	}
	return "NON_COMPLIANT"
}

func getInvokingEvent(event []byte) (InvokingEvent, error) {
	var result InvokingEvent
	err := json.Unmarshal(event, &result)
	if err != nil {
		fmt.Println("Error:", err)
		return result, err
	}

	return result, nil
}

func isApplicable(c ConfigurationItem, e events.ConfigEvent) bool {
	status := c.ConfigurationItemStatus
	if e.EventLeftScope == false && status == "OK" || status == "ResourceDiscovered" {
		return true
	}
	return false

}

func getParams(p events.ConfigEvent, param string) []string {
	if p.RuleParameters != "" {
		var result map[string]string
		err := json.Unmarshal([]byte(p.RuleParameters), &result)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		if _, ok := result[param]; ok {
			value := strings.Split(result[param], ",")
			return value
		}
		return nil
	}
	return nil
}
