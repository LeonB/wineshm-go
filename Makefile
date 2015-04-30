all: shmwrapper1.exe shmwrapper2.bin

CC=gcc
WINECC=i686-w64-mingw32-gcc
CFLAGS=-Wall -Os -g

shmwrapper2.bin: shmwrapper2.c
	$(CC) $< $(CFLAGS) -o $@

shmwrapper1.exe: shmwrapper1.c
	$(WINECC) $< $(CFLAGS) -mconsole -o $@

clean:
	rm -f shmwrapper1.exe shmwrapper2.bin
