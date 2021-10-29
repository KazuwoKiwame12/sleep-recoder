.PHONY: build

build:
	sam build

create-table:
	aws dynamodb create-table --cli-input-json file://sleepRecord_table.json --endpoint-url http://localhost:8000
input-data:
	aws dynamodb batch-write-item --request-items file://dynamodb/sleepRecord_table_data.json --endpoint-url http://localhost:8000 
scan-data:
	aws dynamodb scan --table-name SleepRecord --endpoint-url http://localhost:8000  