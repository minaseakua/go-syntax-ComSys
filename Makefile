OBJ = $(wildcard *.go)

all: server

server: $(OBJ)
	go build -o server $(OBJ)

clean:
	rm server