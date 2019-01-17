- permission.json をグループのインラインポリシーとして作成

- Amazon Kinesis   DataStreams の ストリーム を 作成 する

aws kinesis create-stream --stream-name sample --shard-count 1

- 作成されたことを確認する

aws kinesis list-streams

- permission.json を更新して sns の仕様を許可する
                "sns:*"
                
- アラームの通知先となるSNSトピックを作成

aws sns create-topic --name sample

- cloudwatch のアラームを設定する

aws cloudwatch put-metric-alarm \
  --alarm-name kinesis-mon --metric-name IncomingRecords \
  --namespace AWS/Kinesis --statistic Sum --period 60 \
  --threshold 10 --comparison-operator GreaterThanThreshold \
  --dimensions Name=StreamName,Value=sample \
  --evaluation-periods 1 \
  --alarm-actions arn:aws:sns:ap-northeast-1:780132170115:sample


- set-alarm-stateコマンドを使用してアラーム状態を変更してアラームをテストする

aws cloudwatch set-alarm-state --alarm-name kinesis-mon \
  --state-reason 'initializing' --state-value ALARM


- resharding-functionに紐付ける IAMロールを作成
aws iam create-role --role-name resharding_function_role --assume-role-policy-document file://trustpolicy.json

# 大まかな手順
- put-records.go で 15件kinesisに書き込む
- それを検知して シャーディングし直す


