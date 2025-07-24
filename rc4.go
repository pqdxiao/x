package x

import (
	"crypto/cipher"
	"crypto/rc4"
	"fmt"
	"io"
)

// 使用RC4流加密
type RC4Cipher struct {
	// 密钥长度，范围为1-256字节
	Key []byte
	// rc4加密器
	Cipher *rc4.Cipher
}

type RC4Writer struct {
	RC4Cipher
	Writer io.Writer
}

type RC4Reader struct {
	RC4Cipher
	Reader io.Reader
}

func NewRC4Writer(key []byte, writer io.Writer) (*RC4Writer, error) {
	if len(key) < 1 || len(key) > 256 {
		return nil, fmt.Errorf("invalid rc4 key length")
	}
	return &RC4Writer{RC4Cipher{Key: key}, writer}, nil
}

func NewRC4Reader(key []byte, reader io.Reader) (*RC4Reader, error) {
	if len(key) < 1 || len(key) > 256 {
		return nil, fmt.Errorf("invalid rc4 key length")
	}
	return &RC4Reader{RC4Cipher{Key: key}, reader}, nil
}

func (r *RC4Writer) Write(p []byte) (int, error) {
	if r.Cipher == nil {
		var cipher *rc4.Cipher
		cipher, err := rc4.NewCipher(r.Key)
		if err != nil {
			return 0, err
		}
		r.Cipher = cipher
	}
	streamWriter := &cipher.StreamWriter{S: r.Cipher, W: r.Writer}
	return streamWriter.Write(p)
}

func (r *RC4Reader) Read(p []byte) (int, error) {
	if r.Cipher == nil {
		var cipher *rc4.Cipher
		cipher, err := rc4.NewCipher(r.Key)
		if err != nil {
			return 0, err
		}
		r.Cipher = cipher
	}
	streamReader := &cipher.StreamReader{S: r.Cipher, R: r.Reader}
	return streamReader.Read(p)
}
