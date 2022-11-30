package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	seed := time.Now().UnixNano()
	randomStr := getRandomString(seed, 6)
	t.Log(randomStr)
	assert.Len(t, randomStr, 6)
}
func TestRandomInt(t *testing.T) {
	seed := time.Now().UnixNano()
	randomInt := getRandomInt(seed)
	t.Log(randomInt)
	assert.Contains(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, randomInt)
}
func TestRandomURL(t *testing.T) {
	randomURL := GetRandomURL(5, 8)
	t.Log(randomURL)
	assert.NotEqual(t, "", randomURL)
}
func TestRandomURLMany(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	for i := 0; i < 10; i++ {
		TestRandomURL(t)
		time.Sleep(time.Nanosecond)
	}
}

func TestRandomUrlDuplicate(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	n := 1000
	urls := []string{}
	for i := 0; i < n; i++ {
		url := GetRandomURL(5, 8)
		urls = append(urls, url)
		time.Sleep(time.Nanosecond)
	}

	urlMap := make(map[string]int)
	for _, url := range urls {
		if _, ok := urlMap[url]; ok {
			t.Log("Duplicate entries")
			t.FailNow()
		}
		urlMap[url] = 0
	}
}
