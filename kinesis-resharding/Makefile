.PHONY: deps clean build

functionName = resharding-function

deps:
	go get -u "github.com/aws/aws-lambda-go/events"
	go get -u "github.com/aws/aws-lambda-go/lambda"
	go get -u "github.com/aws/aws-sdk-go/aws"
	go get -u "github.com/aws/aws-sdk-go/aws/session"
	go get -u "github.com/aws/aws-sdk-go/service/cloudwatch"
	go get -u "github.com/aws/aws-sdk-go/service/kinesis"

clean:
	rm -rf ./build/$(functionName)

#exec:
# 	sam local generate-event s3 put --bucket elemental-upload --key sample.mp4 | sam local invoke
#	make build
#	chmod 755 build/create-job
#	sam local invoke CreateJob --event put_event.json

build:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/$(functionName) ./$(functionName)

package:
	sam package --template-file template.yml \
				--output-template-file packaged.yml \
				--s3-prefix $(functionName) \
				--s3-bucket yurakawa-sam-packages


deploy:
	make package
	sam deploy --template-file packaged.yml \
				--stack-name $(functionName) \
				--capabilities CAPABILITY_IAM \
				--region ap-northeast-1

