package util

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
	"strconv"
)

func makeTable() map[string]int {
	hashTable := make(map[string]int)
	data1 := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	getKey := func(i, j int) string {
		return data1[i] + data1[j]
	}

	num := 1
	
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			hashTable[getKey(i, j)] = num
			num++
		}
	}
	return hashTable
}

func MakeHashRandomRange(blockhash []byte, diy []byte, height int64, min, max *big.Int) sdk.Dec {
	if min.Cmp(max) == 0 {
		return sdk.NewDecFromBigInt(min)
	}

	if min.Cmp(max) == 1 {
		return sdk.NewDec(0)
	}
	heightStr := strconv.FormatInt(height, 10)
	mysha512 := sha512.New()
	mysha512.Write(blockhash)
	mysha512.Write([]byte(heightStr))
	if diy != nil {
		mysha512.Write(diy)
	}
	diff := new(big.Int).Sub(max, min)
	diff = new(big.Int).Add(diff, big.NewInt(1))
	randBigInt, err := rand.Int(bytes.NewReader(mysha512.Sum(nil)), diff)
	if err != nil {
		return sdk.NewDec(0)
	}
	return sdk.NewDecFromBigInt(new(big.Int).Add(randBigInt, min))
}


func MakeRandomRate(hash string) int {
	
	if len(hash) < 2 {
		return 0
	}
	table := makeTable()
	
	
	
	
	key := hash[:2]
	
	keyNumber := table[key] 
	
	
	index := 0

	
	if keyNumber <= 112 {
		index = keyNumber / 7
	} else {
		index = (keyNumber-112)/6 + 16
	}
	
	return index + 80
}


func GenerateSecureRandomString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	return randomString[:length], nil
}
