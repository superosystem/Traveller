package helper

import (
	"crypto/rand"
	"io"
)

var numbers = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateOTP(length int) string {
	otp := make([]byte, length)

	n, err := io.ReadAtLeast(rand.Reader, otp, length)

	if n != length {
		panic(err)
	}

	for i := range otp {
		otp[i] = numbers[int(otp[i])%len(numbers)]
	}

	return string(otp)
}
