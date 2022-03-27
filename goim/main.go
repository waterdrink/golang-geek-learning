package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

// https://github.com/Terry-Mao/goim/blob/master/docs/protocol.png
// 以下为goim中数据包协议定义：
const (
	packageLength   = 4
	headerLength    = 2
	protocolVersion = 2
	operation       = 4
	sequenceId      = 4

	headerSize            = packageLength + headerLength + protocolVersion + operation + sequenceId
	headerLengthOffset    = packageLength
	protocolVersionOffset = headerLengthOffset + headerLength
	operationOffset       = protocolVersionOffset + protocolVersion
	sequenceIdOffset      = operationOffset + operation
)

func main() {
	errCh := make(chan error, 1)
	pr, pw := io.Pipe()
	go writeMsg("hello, world.", pw)
	go func() {
		err := readAndPrint(pr)
		errCh <- err

	}()
	err := <-errCh
	fmt.Println(err)

}

func readAndPrint(r io.Reader) error {
	headerBuf := make([]byte, headerSize)
	err := readToBuf(r, headerBuf, 1)
	if nil != err {
		return errors.Wrapf(err, "read header")
	}
	packageLen := binary.BigEndian.Uint32(headerBuf[0:headerLengthOffset])
	headerLen := binary.BigEndian.Uint16(headerBuf[headerLengthOffset:protocolVersionOffset])
	protocolVer := binary.BigEndian.Uint16(headerBuf[protocolVersionOffset:operationOffset])
	op := binary.BigEndian.Uint32(headerBuf[operationOffset:sequenceIdOffset])
	seqId := binary.BigEndian.Uint32(headerBuf[sequenceIdOffset:])
	if headerLen != headerSize {
		return errors.Errorf("bad head size [%v:%v]", headerLen, headerSize)
	}
	fmt.Println("packageLen: ", packageLen)
	fmt.Println("headerLen: ", headerLen)
	fmt.Println("protocolVer: ", protocolVer)
	fmt.Println("op: ", op)
	fmt.Println("seqId: ", seqId)

	bodyLen := packageLen - uint32(headerLen)
	if bodyLen <= 0 {
		return errors.New("empty body")
	}
	bodyBuf := make([]byte, bodyLen)
	err = readToBuf(r, bodyBuf, 1)
	if nil != err {
		return errors.Wrapf(err, "read body")
	}
	fmt.Println("body: ", string(bodyBuf))
	return nil
}

func readToBuf(r io.Reader, buf []byte, timeoutSed int) error {
	deadline := time.Now().Add(time.Duration(timeoutSed) * time.Second)
	readN := 0
	for {
		n, err := r.Read(buf[readN:])
		if n < 0 {
			return errors.Errorf("bad read count: %v", n)
		}

		readN += n
		if readN == len(buf) {
			return nil
		}
		if time.Now().After(deadline) {
			return errors.Errorf("not read full buf[%v:%v], err: %v", readN, len(buf), err)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func writeMsg(msg string, w io.Writer) {
	headerBuf := make([]byte, headerSize)
	binary.BigEndian.PutUint32(headerBuf[:headerLengthOffset], uint32(headerSize+len(msg)))
	binary.BigEndian.PutUint16(headerBuf[headerLengthOffset:protocolVersionOffset], headerSize)
	binary.BigEndian.PutUint16(headerBuf[protocolVersionOffset:operationOffset], 10)
	binary.BigEndian.PutUint32(headerBuf[operationOffset:sequenceIdOffset], 11)
	binary.BigEndian.PutUint32(headerBuf[sequenceIdOffset:], 12)
	w.Write(headerBuf)
	if len(msg) > 0 {
		w.Write([]byte(msg))
	}
	return
}
