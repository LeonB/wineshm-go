all: assets/shmwrapper1.exe assets/shmwrapper2.bin bindata.go

CC=gcc
WINECC=i686-w64-mingw32-gcc
CFLAGS=-Wall -Os -g

assets/shmwrapper2.bin: shmwrapper2.c
	$(CC) $< $(CFLAGS) -o $@

assets/shmwrapper1.exe: shmwrapper1.c
	$(WINECC) $< $(CFLAGS) -mconsole -o $@

bindata.go: assets/shmwrapper1.exe assets/shmwrapper2.bin
	go-bindata -pkg=wineshm -ignore=.gitkeep assets/

clean:
	rm -f assets/shmwrapper1.exe assets/shmwrapper2.bin bindata.go

# vim: syntax=make ts=4 sw=4 sts=4 sr noet
