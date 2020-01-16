package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	usr "github.com/mainmast/iam-models/pkg/user"
)

var accessChr = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var secretChr = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123457890/&$-_!}{|")

var randInt = []rune("abcdezABCDEF123457890")

var hashSalt = os.Getenv("HASH_SALT")

// Generate Access Key
func GenerateAccesKey(orgName string) string {
	rand.Seed(time.Now().UnixNano())

	r := make([]rune, 20)

	for i := range r {
		r[i] = accessChr[rand.Intn(len(accessChr))]
	}

	rs := string(r)
	encodedOrg := base64.StdEncoding.EncodeToString([]byte(orgName))
	cleaned := strings.TrimSuffix(strings.TrimSuffix(encodedOrg, "="), "=")
	
	return rs[:10] + cleaned + rs[10:20]
}

// Generate Random String
func GenerateRandomString(nChar int) string {

	rand.Seed(time.Now().UnixNano())

	r := make([]rune, nChar)
	
	for i := range r {
		r[i] = randInt[rand.Intn(len(randInt))]
	}

	return string(r)
}

// Generate Secret Key
func GenerateSecretKey() string {
	rand.Seed(time.Now().UnixNano())

	r := make([]rune, 40)

	for i := range r {
		r[i] = secretChr[rand.Intn(len(secretChr))]
	}
	
	return string(r)
}

// Hashes a plaintext secret
func SecretHash(secret string) (string, error) {
	hs, err := bcrypt.GenerateFromPassword([]byte(hashSalt+secret), bcrypt.MinCost)

	if err != nil {
		fmt.Println("error creating hash secret", err)
		return "", errors.New("error creating hash secret")
	}

	return string(hs), nil
}

// Validates the given string with a hashed value.
func CompareHash(hash string, secret string) bool {

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(hashSalt+secret)); err != nil {
		fmt.Println("Secret doesn't match previously hashed value.", err)
		return false
	}

	return true
}

// NormaliseOrganisationName
func NormalizeName(orgName string, randChars bool) string {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {
		fmt.Print(err)

		if randChars {
			return "default_schema_" + GenerateRandString(5)
		}

		return "default_account_" + GenerateRandString(5)
	}

	orgName = reg.ReplaceAllString(orgName, "")

	if randChars {
		orgName = orgName + "_" + GenerateRandString(4)
	}

	orgName = strings.ToLower(orgName)
	return orgName
}
