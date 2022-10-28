.PHONY: submit test
tags=logging callback metrics scheduler
cgo=0
goos=linux
goarch=amd64
#export TAGS
#export CGO_ENABLED
#export GOARCH
# Set next value at compile time for fallback log-level
ldXflags=github.com/mongodb/mongodbatlas-cloudformation-resources/util.defaultLogLevel=info
ldXflagsD=github.com/mongodb/mongodbatlas-cloudformation-resources/util.defaultLogLevel=debug

# example for resource: MongoDB::Atlas::Teams
# example for go path: github.com/mongodb/mongodbatlas-cloudformation-resources/invites
# example for go path: github.com/mongodb/mongodbatlas-cloudformation-resources/indexes

create:
	cd .. && \
	mkdir $(resource) && cp autogen/configs/mongodb-atlas-$(schemaName).json $(resource)/mongodb-atlas-$(schemaName).json && \
	cd $(resource) && \
 	cfn init -t Mongodb::Atlas::$(schemaName) -a r && cfn generate
	cd ../autogen && \
	go run main.go -res ../$(resource) -schemaName $(schemaName)
	cd ../$(resource) && rm go.mod


schema:
	cd schema-gen && go build
	./schema-gen/schema-gen
