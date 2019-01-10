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

	stream_name := "sample"
	partition_key := uuid.New().String()
	data := time.Now().Unix()

	for i := 0; i < 15; i++ {
		putOutput, err := k.PutRecord(&kinesis.PutRecordInput{
			Data: []byte{byte(data)},
			StreamName: aws.String(stream_name),
			PartitionKey: aws.String(partition_key),
		})

		if err != nil {
			panic(err)
		}

		fmt.Println(putOutput)
	}
	fmt.Println("end...")
}

