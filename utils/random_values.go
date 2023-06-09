package utils

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// generates a random integer number within given range including the 0
func RandomInt(max int) int {
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(max)
}

// generates a random password of 12 characters long with random uppercase and lowercase letters
// and random numbers
func GeneratePassword() string {
	letters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	numbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	result := ""
	var index int
	turned := false
	rand.NewSource(time.Now().UnixNano())

	// add random characters that are either numbers or letters
	for i := 0; i < 12; i++ {
		letterOrNumber := RandomInt(2)
		possibleLetter := RandomInt(len(letters) - 1)
		possibleNumber := RandomInt(len(numbers) - 1)
		choose := letterOrNumber%2 == 0
		value := ""
		if choose {
			value = letters[possibleLetter]
		} else {
			value = numbers[possibleNumber]
		}
		result += value
	}

	// iterate until at least one letter is turned uppercase
	index = RandomInt(len(result) - 1)
	for !turned {
		isNumeric := unicode.IsDigit(rune(result[index]))
		if !isNumeric {
			upperized := strings.ToUpper(string(result[index]))
			result = strings.Replace(result, string(result[index]), upperized, 1)
			turned = true
		}
		index = RandomInt(len(result) - 1)
	}

	return result
}

// make http request and get a random SVG profile picture from given randomness seed,
// returns empty string if an error occurs.
// The name is auto-generated by the function, on the path must go only the name of the
// user directory,
func GetRandomPfp(seed, path string) (string, error) {
	pfpApi := fmt.Sprintf("%v/svg?seed=%v", GlobalEnv.URLs.PfpApi, seed)
	response, err := http.Get(pfpApi)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	filename := "pfp.svg"
	filepath := filepath.Join(GlobalEnv.GitRoot, path, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("%v/%v/%v/%v/%v", GlobalEnv.URLs.Project, "media", path, "picture", filename)
	return result, nil
}
