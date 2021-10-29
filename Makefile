.PHONY: build

build:
	sam build
create-table:
	aws dynamodb create-table --cli-input-json file://sleepRecord_table.json --endpoint-url http://localhost:8000
input-data:
	aws dynamodb batch-write-item --request-items file://dynamodb/sleepRecord_table_data.json --endpoint-url http://localhost:8000 
scan-data:
	aws dynamodb scan --table-name SleepRecord --endpoint-url http://localhost:8000  
create-backet:
	aws s3 mb --profile test-for-s3 s3://static --endpoint-url http://localhost:9000
check-backet:
	aws s3 ls s3://static --profile test-for-s3 --endpoint-url http://127.0.0.1:9000

