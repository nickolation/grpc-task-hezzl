down-system:
	docker-compose down 

remove-cache:
	sudo rm -R shared/ch-data/

run-system:
	MY_IP=192.168.0.104 docker-compose up  

up-db-migrations:
	cd ./.. && go run storage/test_up_dbmigrations/main.go

run-grpc-server:
	cd ./.. && go run grpc/grpc_server/server.go

run-grpc-client:
	cd ./.. && go run grpc/test_grpc_apiclient/main.go 

create-logs-topic:
	docker run --net=host --rm confluentinc/cp-kafka:latest kafka-topics --create --topic logs --partitions 4 --replication-factor 2 --if-not-exists --bootstrap-server localhost:19092

consume-logs:
	kafkacat -C -b localhost:19092,localhost:29092,localhost:39092 -t logs 

