package common

import (
	"net"
	"io"
	"log"
	"errors"
)

type ConnInfo struct {
	Client net.Conn
	Server net.Conn
	Infer chan bool
}

type ConnChan chan net.TCPConn


var ErrShortWrite = errors.New("short write")

// errInvalidWrite means that a write returned an impossible count.
var errInvalidWrite = errors.New("invalid write result")

// ErrShortBuffer means that a read required a longer buffer than was provided.
var ErrShortBuffer = errors.New("short buffer")

// EOF is the error returned by Read when no more input is available.
// (Read must return EOF itself, not an error wrapping EOF,
// because callers will test for EOF using ==.)
// Functions should return EOF only to signal a graceful end of input.
// If the EOF occurs unexpectedly in a structured data stream,
// the appropriate error is either ErrUnexpectedEOF or some other error
// giving more detail.
var EOF = errors.New("EOF")

// ErrUnexpectedEOF means that EOF was encountered in the
// middle of reading a fixed-size block or data structure.
var ErrUnexpectedEOF = errors.New("unexpected EOF")

func Transfer(from *net.TCPConn, to *net.TCPConn, ch chan bool, director int){
	if director == 0 {
		log.Println("will close from automate")
		defer from.Close()
	}else if director == 1{
		log.Println("will close to automate")
		defer to.Close()
	}else{
		log.Println("will not close socket")
	}
	buf := make([]byte,1024)
	for{
		log.Println("start read data from",from.RemoteAddr().String())
		nr,er := io.ReadAtLeast(from, buf, 1)
		log.Println("has read data:",nr)
		if nr > 0 {
			log.Println("start write data to ")
			nw, ew := to.Write(buf[0:nr])
			log.Println("has write data",nw)
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
					log.Println("write data err",ew)
				}
			}
			if nr != nw {
				log.Println("write data err",ErrShortWrite)
				if director < 2 {
					ch <- true
					break
				}

			}
		}
		if er != nil {
			if er == io.EOF{
				log.Println("read has arrive last")
				if director < 2 {
					ch <- true
					break
				}

			}
			log.Println("err in read data,",er)
			if director < 2 {
				break
			}
		}
		log.Println("success transfer data ",nr)

	}
}