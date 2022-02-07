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
