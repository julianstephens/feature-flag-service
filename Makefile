
.PHONY: help apigen

help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

apigen:   ## Generate gRPC API code
	@protoc --go_out=gen/go/grpc/v1 --go-grpc_out=gen/go/grpc/v1 api/grpc/v1/*.proto
