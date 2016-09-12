package encryption

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type FemotCrypto struct{}

func (c *FemotCrypto) Encrypt(in, iv []byte) ([]byte, error) {
	inputSize := len(in)
	ivSize := len(iv)

	if ivSize != 32 {
		return nil, errors.New("ivSize must be 32!")
	}

	arr2 := make([]byte, 256)
	arr3 := make([]byte, 256)

	roundedSize := inputSize + (256 - (inputSize % 256))
	totalSize := roundedSize + 32
	output := make([]byte, totalSize)

	for j := 0; j < 8; j++ {
		for i := 0; i < 32; i++ {
			arr2[32*j+i] = rotateLeft(iv[i], uint(j))
		}
	}

	for i := 0; i < 32; i++ {
		output[i] = iv[i]
	}
	for i := 0; i < inputSize; i++ {
		output[i+32] = in[i]
	}
	output[totalSize-1] = byte(256 - (inputSize % 256))

	for offset := 32; offset < totalSize; offset += 256 {
		for i := 0; i < 256; i++ {
			output[offset+i] ^= arr2[i]
		}
		slice := output[offset:]

		// sub_9E9D8
		pikachu(slice, &arr3)

		for j := 0; j < 256; j++ {
			arr2[j] = arr3[j]
			output[j+offset] = arr3[j]
		}
	}
	return output, nil
}

func (c *FemotCrypto) CreateIV() []byte {
	iv := make([]byte, 32)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//for i := 0; i < 32; i++ {
	//iv[i] = byte(r.Int() % 256)
	//}
	r.Read(iv)
	return iv
}

func (c *FemotCrypto) Enabled() bool {
	return true
}

func rotateLeft(x byte, n uint) byte {
	return ((x << n) | (x >> (8 - n))) & 255
}

func pSlice(s []uint32) {
	out := make([]string, len(s))
	for i, b := range s {
		out[i] = fmt.Sprintf("%d", b)
	}
	//fmt.Printf("[%s]\n", strings.Join(out, ","))
}

// pikachu is an alias for sub_9E9D8
func pikachu(slice []byte, arr3 *[]byte) {
	temp := make([]uint32, 812/4)
	temp2 := make([]uint32, 256/4)

	for i := 0; i < 64; i++ {
		temp2[i] = uint32(slice[4*i]) | uint32(slice[4*i+1])<<8 | uint32(slice[4*i+2])<<16 | uint32(slice[4*i+3])<<24
	}

	// 1
	pSlice(temp)
	bulbasaur(&temp, &temp2) // sub_87568
	// 2
	pSlice(temp)
	ivysaur(&temp) // sub_8930C
	// 3
	pSlice(temp)
	venusaur(&temp) // sub_8B2F4
	// 4
	pSlice(temp)
	charmander(&temp) // sub_8D114
	// 5
	pSlice(temp)
	charmeleon(&temp) // sub_8F0B0
	// 6
	pSlice(temp)
	charizard(&temp) // sub_910A8
	// 7
	pSlice(temp)
	squirtle(&temp) // sub_92E08
	// 8
	//pSlice(temp)
	wartortle(&temp) // sub_94BDC
	// 9
	//pSlice(temp)
	blastoise(&temp) // sub_96984
	// 10
	//pSlice(temp)
	caterpie(&temp) // sub_985E0
	// 11
	//pSlice(temp)
	metapod(&temp) // sub_9A490
	// 12
	//pSlice(temp)
	butterfree(&temp) // sub_9C42C
	// 13
	pSlice(temp)
	weedle(&temp, &temp2) // sub_9E1C4
	// 14
	pSlice(temp)

	for i := 0; i < 64; i++ {
		(*arr3)[4*i] = byte(temp2[i] & 255)
		(*arr3)[4*i+1] = byte(temp2[i] >> 8 & 255)
		(*arr3)[4*i+2] = byte(temp2[i] >> 16 & 255)
		(*arr3)[4*i+3] = byte(temp2[i] >> 24 & 255)
	}
}
