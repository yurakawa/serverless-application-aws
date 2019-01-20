package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/google/uuid"
	"time"
)

// パーティションキーはUUIDを元にしたユニークな値。
// 実際に保存するデータは現在時刻
//
// 15件のデータを書き込む

func main(){
	s := session.Must(session.NewSessionWithOptions(session.Options{}))
	k := kinesis.New(
		s,
		aws.NewConfig().WithRegion("ap-northeast-1"),
	)

	streamName := "kinesis-resharding-sample-stream"
	partitionKey := uuid.New().String()
	data := time.Now().Unix()

	for i := 0; i < 15; i++ {
		putOutput, err := k.PutRecord(&kinesis.PutRecordInput{
			Data: []byte{byte(data)},
			StreamName: aws.String(streamName),
			PartitionKey: aws.String(partitionKey),
		})

		if err != nil {
			panic(err)
		}

		fmt.Println(putOutput)
	}
	fmt.Println("end...")
}

