package main

import (
	"fmt"
	"log"
	"os"
	"unsafe"

	mmap "github.com/edsrzf/mmap-go"
	wineshm "github.com/leonb/wineshm-go"
)

type StatusField int32

const (
	MAX_BUFS = 4
)

type Header struct {
	Ver      int32       // api version 1 for now
	Status   StatusField // bitfield using StatusField
	TickRate int32       // ticks per second (60 or 360 etc)

	// session information, updated periodicaly
	SessionInfoUpdate int32 // Incremented when session info changes
	SessionInfoLen    int32 // Length in bytes of session info string
	SessionInfoOffset int32 // Session info, encoded in YAML format

	// State data, output at tickRate
	NumVars         int32 // length of array pointed to by varHeaderOffset
	VarHeaderOffset int32 // offset to VarHeader[numVars] array, Describes the variables recieved in varBuf

	NumBuf int32    // <= MAX_BUFS (3 for now)
	BufLen int32    // length in bytes for one line
	Pad1   [2]int32 // (16 byte align)
	VarBuf [MAX_BUFS]VarBuf
}

type VarBuf struct {
	TickCount int32    // used to detect changes in data
	BufOffset int32    // offset from header
	Pad       [2]int32 // (16 byte align)
}

func main() {
	// Get wine file descriptor
	wineshm.WineCmd = []string{"/opt/iracing/bin/wine", "--bottle", "default"}
	shmfd, err := wineshm.GetWineShm("Local\\IRSDKMemMapFileName", wineshm.FILE_MAP_READ)
	if err != nil {
		log.Fatal(err)
	}

	// Turn file descriptor into os.File
	file := os.NewFile(shmfd, "Local\\IRSDKMemMapFileName")

	mmap, err := mmap.Map(file, mmap.RDONLY, 0)
	if err != nil {
		log.Fatalf("error mapping: %s", err)
	}

	fmt.Println(len(mmap))
	header := (*Header)(unsafe.Pointer(&mmap))
	fmt.Printf("%+v", header)
}
