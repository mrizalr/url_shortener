package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func getRandomString(seed int64, length int) string {
	token := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	nt := len(token)
	rand.Seed(seed)

	result := ""
	for i := 0; i < length; i++ {
		result += string(token[rand.Intn(nt)])
	}

	return result
}

func getRandomInt(seed int64) int {
	rand.Seed(seed)
	return rand.Intn(10)
}

func GetRandomURL(min, max int) string {
	result := ""
	rand.Seed(time.Now().UnixNano())
	randomLength := rand.Intn(max-min) + min

	for i := 0; i < randomLength; i++ {
		seed := int64(time.Now().UnixNano() * rand.Int63())
		rand.Seed(seed)
		randomType := rand.Intn(3000)

		if randomType%2 == 0 {
			result += getRandomString(seed, 1)
		} else {
			result += strconv.Itoa(getRandomInt(seed))
		}
	}

	return result
}
