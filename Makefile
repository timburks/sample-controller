
lite:
	go install ./...

all:	code
	go install ./...

code:
	./hack/update-codegen.sh
