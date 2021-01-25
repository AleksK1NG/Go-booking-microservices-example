.PHONY:

linter:
	echo "Starting linters"
	cd main && echo 'cool ' && golangci-lint run ./...
	cd ..
	cd user && golangci-lint run ./...
	cd ..
	cd sessions && golangci-lint run ./...
	cd ..

