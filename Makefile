start:
	docker-compose up
restart:
	docker-compose up --build
stop:
	docker-compose down
test:
	go test ./...