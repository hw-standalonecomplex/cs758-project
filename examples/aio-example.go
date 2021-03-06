package main

// package main must be the first thing in this program for the test to work.

import "fmt"
import "os"
import "errors"
import "syscall"
import aio "github.com/spwilson2/cs758-project/libaio"

const VALID_STRING string = "package main"
const TESTFILE string = "aio-example.go"

var fail = errors.New("")

const BUFSIZE = 5000000

func chk_err(err error) {
	if err != nil {
		fmt.Printf("(FAILED) %s\n", os.Args[0])
		panic(err) //os.Exit(-1)
	}
}

func main() {

	var ctx syscall.AioContext_t
	var err error

	// Create our AioContext with support for 1 inflight aio.
	chk_err(syscall.IoSetup(1, &ctx))

	defer func() {
		chk_err(syscall.IoDestroy(ctx))
	}()

	// Open our file using the systemcall
	fd, err := syscall.Open(TESTFILE, syscall.O_RDONLY, 0)
	chk_err(err)

	var iocb syscall.Iocb
	var iocbp = &iocb

	// Create our buffer for reading into
	var buffer = make([]byte, BUFSIZE, BUFSIZE)
	aio.PrepPread(iocbp, fd, buffer, len(buffer), 0)

	// Submit our request for AIO.
	chk_err(syscall.IoSubmit(ctx, 1, &iocbp))

	var event syscall.IoEvent
	var timeout syscall.Timespec
	timeout.Sec = 1

	// Get back the result from our AIO call.
	events := syscall.IoGetevents(ctx, 1, 1, &event, &timeout)
	if events == 0 {
		chk_err(fail)
	} else if events < 0 {
		chk_err(fail)
	}

	// Check the result..
	if string(buffer[:len(VALID_STRING)]) != VALID_STRING {
		fmt.Printf("Expected:%s, Found:%s\n", VALID_STRING, buffer)
		chk_err(fail)
	}

	fmt.Printf("(OK) %s\n", os.Args[0])
}
