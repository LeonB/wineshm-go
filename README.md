# Wine Shm grabber in go

This package retrieves a Wine shared memory map file descriptor and makes it
available in Linux.

## Installation

Because of the dependency on C / compiled binary files the best way is to copy
this repo and include it in your own project.

## Example

``` go
// Get wine file descriptor
shmfd, err := GetWineShm("Local\\IRSDKMemMapFileName", FILE_MAP_READ)
if err != nil {
  log.Fatal(err)
}

// Turn file descriptor into os.File
file := os.NewFile(shmfd, "Local\\IRSDKMemMapFileName")
```

## How it works

Go launches a windows binary through wine and passes a unix socket to the new
wine process (shmwrapper1) by attaching the socket to the subprocess' stdout.

The wine process then opens the shared memory map and sets the file handle as
it's stdin. Wine then launches a new linux process (shmwrapper2) which inherits
the stdin & stdout from the parent process (shmwrapper1).

So now shmwrapper2 is a native Linux process with a unix socket as stdout and
the windows shared memory map handle as it's stdin. It then uses a unix feature
to send a file descriptor (stdin) over a unix socket (stdout) to the original go
process initiating the wine process.

## Thanks

Very special thanks to [jspenguin](https://github.com/jspenguin/) for creating
the c code and thinking of a way to share file descriptors between Linux and Wine.
