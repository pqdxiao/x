package x

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// #define PLM_STRING_ENCRYPT_KEY		_T("_STRING__ENCRYPT_")
const PLM_STRING_ENCRYPT_KEY = "_STRING__ENCRYPT_"

// #define PLM_STRING_ENCRYPT_HEADER	_T("_HEADER_")
const PLM_STRING_ENCRYPT_HEADER = "_HEADER_"

type RC4 struct {
	state [256]byte
	x, y  byte
}

func NewRC4WithState(state [256]byte, x, y byte) *RC4 {
	return &RC4{
		state: state,
		x:     x,
		y:     y,
	}
}

func (rc4 *RC4) Crypt(data []byte) {
	for i := range data {
		rc4.x++
		rc4.y += rc4.state[rc4.x]
		rc4.state[rc4.x], rc4.state[rc4.y] = rc4.state[rc4.y], rc4.state[rc4.x]
		data[i] ^= rc4.state[(rc4.state[rc4.x]+rc4.state[rc4.y])&0xFF]
	}
}

func unencryptBase64(base64Str string, rc4 *RC4) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}

	rc4.Crypt(decodedBytes)

	return decodedBytes, nil
}

func RC4Decrypt(base64Str string) (string, error) {
	// Provided RC4 state and initial values for x and y
	state := [256]byte{
		189, 47, 79, 98, 16, 145, 45, 161, 123, 9, 12, 3, 13, 227, 11, 7, 150, 53, 63, 196, 76, 70, 102, 122, 228, 52, 136, 43, 155, 198, 234, 211,
		86, 177, 179, 154, 170, 32, 28, 131, 137, 166, 158, 24, 56, 219, 95, 65, 238, 77, 104, 94, 109, 62, 126, 20, 160, 169, 172, 182, 66, 134, 138,
		195, 78, 205, 176, 124, 116, 135, 163, 21, 71, 220, 67, 22, 100, 197, 192, 149, 185, 4, 162, 125, 225, 64, 36, 105, 72, 184, 246, 92, 110,
		191, 96, 217, 233, 26, 113, 250, 237, 93, 39, 35, 183, 147, 87, 142, 85, 119, 40, 130, 143, 224, 152, 221, 41, 60, 243, 25, 159, 139, 186,
		218, 140, 230, 226, 114, 146, 29, 129, 180, 69, 255, 209, 229, 101, 37, 171, 156, 188, 151, 115, 249, 210, 23, 244, 187, 203, 178, 19, 201,
		54, 121, 51, 106, 0, 55, 42, 6, 193, 81, 254, 144, 208, 236, 148, 17, 46, 232, 128, 59, 118, 48, 242, 49, 190, 239, 31, 223, 202, 34, 213,
		111, 248, 199, 50, 80, 57, 231, 84, 133, 168, 103, 222, 164, 27, 68, 8, 75, 44, 241, 181, 14, 127, 10, 141, 165, 2, 112, 153, 235, 82, 204,
		73, 214, 245, 117, 5, 200, 167, 247, 97, 194, 33, 88, 61, 207, 83, 90, 212, 157, 174, 132, 173, 74, 91, 18, 215, 108, 206, 15, 251, 58, 216,
		252, 240, 38, 99, 253, 175, 30, 89, 1, 107, 120,
	}
	x := byte(0) // Assuming x is initialized to 0
	y := byte(0) // Assuming y is initialized to 0

	rc4 := NewRC4WithState(state, x, y)

	decryptedData, err := unencryptBase64(base64Str, rc4)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	decryptedString := strings.TrimPrefix(string(decryptedData), PLM_STRING_ENCRYPT_HEADER)
	return decryptedString, nil
}
