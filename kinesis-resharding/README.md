## 概要

Amazon KinesisのPutレコードの量をAmazon CloudWatchを使ってモニタリングし、一定の値(10回/min)を超えてアラームが発生したらシャード数を増減させる(リシャーディング)

アラームの通知先はAmazon SNSのトピック。アラームが発生したらAamazon SNSのトピックへ通知し、そのトピックをサブスクライブしているLambda関数が呼び出され、Lambda関数を用いてAmazon KinesisのリシャーディングAPiを実行する。

## 実行までに必要な手順

- デプロイ

```
make deps
make build
make deploy
```

- kinesis へ書き込みをPutしアラームを発生させる

```
go run put-records/put-records.go 
```

- シャード数を確認する

aws kinesis describe-stream --stream-name kinesis-resharding-sample-stream | jq -r '.StreamDescription | .Shards[] | "\(.ShardId) \(.ParentShardId) \( .AdjacentParentShardId) \(.HashKeyRange | .StartingHashKey)"' | sed s/shardId-//g |awk 'BEGIN{i=0;print}{hashKey[i]=$4;state[i]="Open";if($2 != "null"){state[$2/1]="Close"};if($3 != "null"){state[$3/1]="Close"};i++}END{for(j=0;j<i;j++) printf("shardId-%012d %5s start: %s\n",j,state[j],hashKey[j]);}'

参考: https://qiita.com/daikumatan/items/d58fb628d5c03cc8a710
