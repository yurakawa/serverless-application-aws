package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"log"
)

// CloudWatchのAlarmの構造
type MessageBody struct {
	AWSAccountId     string `json:"AWSAccountId"`
	AlarmDescription string `json:"AlarmDescription"`
	AlarmName        string `json:"AlarmName"`
	NewStateReason   string `json:"NewStateReason"`
	NewStateValue    string `json:"NewStateValue"`
	OldStateValue    string `json:"OldStateValue"`
	Region           string `json:"Region"`
	StateChangeTime  string `json:"StateChangeTime"`
	Trigger          struct {
		ComparisonOperator string `json:"ComparisonOperator"`
		Dimensions         []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"Dimensions"`
		EvaluationPeriods float64     `json:"EvaluationPeriods"`
		MetricName        string      `json:"MetricName"`
		Namespace         string      `json:"Namespace"`
		Period            float64     `json:"Period"`
		Statistic         string      `json:"Statistic"`
		Threshold         float64     `json:"Threshold"`
		Unit              interface{} `json:"Unit"`
	} `json:"Trigger"`
}

func handleRecord(record events.SNSEventRecord) error {
	message := record.SNS.Message
	log.Println(message)

	awsBody := &MessageBody{}
	err := json.Unmarshal([]byte(message), &awsBody)
	if err != nil {
		panic(err)
	}

	alarmName := awsBody.AlarmName
	streamName := awsBody.Trigger.Dimensions[0].Value

	s := session.Must(session.NewSessionWithOptions(session.Options{}))
	k := kinesis.New(s)
	w := cloudwatch.New(s)


	streamSummary, err := k.DescribeStreamSummary(&kinesis.DescribeStreamSummaryInput{
		StreamName: aws.String(streamName),

	})
	if err != nil {
		return err
	}
	currentCnt := aws.Int64Value(streamSummary.StreamDescriptionSummary.OpenShardCount) // 現在のオープンシャード数を取得

	// シャード数を2倍に変更
	targetShardCnt := currentCnt * 2
	_, err = k.UpdateShardCount(&kinesis.UpdateShardCountInput{
		StreamName: aws.String(streamName),
		TargetShardCount: aws.Int64(targetShardCnt),
		ScalingType: aws.String("UNIFORM_SCALING"),

	})
	if err != nil {
		return err
	}

	// 現在のアラームのしきい値をシャード数x1000の80%に設定
	threshold := float64(targetShardCnt * 1000) * 0.8
	_, err = w.PutMetricAlarm(&cloudwatch.PutMetricAlarmInput{
		AlarmName: aws.String(alarmName),
		MetricName:         aws.String("IncomingRecords"),
		Namespace:          aws.String("AWS/Kinesis"),
		Period:             aws.Int64(60),
		EvaluationPeriods:  aws.Int64(1),
		ComparisonOperator: aws.String(cloudwatch.ComparisonOperatorGreaterThanThreshold),
		Threshold:          aws.Float64(threshold),
		Statistic:          aws.String(cloudwatch.StatisticSum),
	})

	if err != nil {
		return err
	}

	return nil
}

type Response struct {
	Message string `json:"message"`
}

func Handler(ctx context.Context, event events.SNSEvent) (Response, error) {
	log.Println("start")

	for _, record := range event.Records {
		select {
		case <-ctx.Done():
			return Response{Message: ""}, ctx.Err()
		default:
			if err := handleRecord(record); err != nil {
				return Response{Message: ""}, err
			}
		}
	}

	log.Println("end")
	return Response{Message: "Success"}, nil
}

func main() {
	lambda.Start(Handler)
}
