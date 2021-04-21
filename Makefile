.PHONY: all
all: swagger

.PHONY: swagger
swagger:
	swagger generate spec -o ./docs/swagger.json

.PHONY: docker
docker:
	docker build -t todoapi .

.PHONY: docker-run
docker-run:
	docker run -p 8081:8081 todoapi