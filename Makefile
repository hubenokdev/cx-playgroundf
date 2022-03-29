export GO111MODULE=on

dep: ## Update go vendor
	go mod vendor
	go mod verify
	go mod tidy

docker-run:
	docker build -t cxplayground .
	docker run -p 5336:5336 -tid cxplayground