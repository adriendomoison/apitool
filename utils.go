package apitool

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// ExtractQueryParams format arguments for gorm from a map[string]interface{}
func ExtractQueryParams(queryParams map[string]interface{}) (query string, args []interface{}) {
	for key, value := range queryParams {
		if query == "" {
			query += key
		} else {
			query += " AND " + key
		}
		args = append(args, value)
	}
	return
}

func ReadSecret(fileName string) (string, error) {
	b, err := ioutil.ReadFile("/run/secrets/" + fileName)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(b), "\n"), nil
}

// GenerateRandomString create a random string of the requested length using the hexadecimal symbols
func GenerateRandomString(strLen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "ABCDEF0123456789"
	result := make([]byte, strLen)
	for i := 0; i < strLen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// GenerateCleanString take a string with symbols and extra spaces and return an alpha numerical string with no accents
func GenerateCleanString(originalString string) (string, error) {

	// Remove accents
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	accentLessString, _, _ := transform.String(t, originalString)

	// Remove characters others than alphanumericals
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	alphanumericalString := reg.ReplaceAllString(accentLessString, " ")

	// Remove extra spaces
	cleanString := strings.Join(strings.Fields(alphanumericalString), " ")

	return cleanString, nil
}
