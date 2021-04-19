.PHONY: all
all: swagger

.PHONY: swagger
swagger:
	swagger generate spec -o ./docs/swagger.json
