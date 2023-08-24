BINARY_NAME=agent
 
all: ${BINARY_NAME} test
 
${BINARY_NAME}:
	go build -tags netgo -ldflags '-extldflags "-static"' -o ${BINARY_NAME} main.go
 
run: ${BINARY_NAME}
	./${BINARY_NAME} -interval 0
 
clean:
	go clean
	rm -f ${BINARY_NAME}
