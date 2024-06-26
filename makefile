#generate a basic makefile for a go project
# Generate a basic makefile for a Go project

# Variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GORUN = $(GOCMD) run
GOTEST = $(GOCMD) test
BINARY_NAME = dist-key-value-storage
MAIN_FILE = main.go

# Build task
build:
	$(GOBUILD) -o /tmp/bin/$(BINARY_NAME) $(MAIN_FILE)

build-race:
	$(GOBUILD) -race -o /tmp/bin/$(BINARY_NAME) $(MAIN_FILE)

# Clean task
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Test task
test:
	$(GOTEST) -v ./...

# Run task
run:
	$(GORUN) $(MAIN_FILE)
