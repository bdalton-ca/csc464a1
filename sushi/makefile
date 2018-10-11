# Variables
CC=g++
STD=c++11
CFLAGS= -std=$(STD) -Wno-write-strings -c -g #-O3 -static -msse -mfpmath=sse
LDFLAGS= -lpthread #-lSDL 
SRC= $(wildcard src/*.c)
HDR= $(wildcard src/*.h)
OBJ= $(patsubst src/%.c,obj/%.o,$(SRC)) 


# Rules


all: $(SRC) $(OBJ)
	$(CC) -o main $(OBJ) $(LDFLAGS)
	go fmt src/gomain.go
	go build src/gomain.go
	


obj/%.o: src/%.c $(HDR)
	$(CC) $(CFLAGS) $< -o $@
	

# Utilities

.PHONY: new clean arch export
export: all

new:
	mkdir src
	mkdir obj
	echo "#include <stdlib.h>\n" > src/main.c
	echo "int main(int argc, char * argv[])\n{\n\treturn 0;\n}\n" >> src/main.c
	
clean:
	@rm -r obj/*
	@rm main
	
arch:
	@mkdir -p $@
	@mkdir -p arch/`date +%y.%m.%d`
	@cp -r src arch/`date +%y.%m.%d`
