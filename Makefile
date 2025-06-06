run:
	go run cmd/comments-system/main.go --config=./configs/config.yaml

gqlgen:
	gqlgen generate --config ./internal/graph/gqlgen.yml