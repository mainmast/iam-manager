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

//GenerateAPILogin ...
func GenerateAPILogin(input string) string {

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if re.MatchString(input) {

		idx := strings.Index(input, "@")
		base := input[:idx]
		return "root_" + base + "_api"
	}

	return "root_" + input + "_api"
}

//GenerateAccesKey ...
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

//GenerateRandString ...
func GenerateRandString(nChar int) string {

	rand.Seed(time.Now().UnixNano())

	r := make([]rune, nChar)
	for i := range r {
		r[i] = randInt[rand.Intn(len(randInt))]
	}
	return string(r)
}

//GenerateSecretKey ...
func GenerateSecretKey() string {

	rand.Seed(time.Now().UnixNano())

	r := make([]rune, 40)
	for i := range r {
		r[i] = secretChr[rand.Intn(len(secretChr))]
	}
	return string(r)

}

//SecretHash ...
func SecretHash(secret string) (string, error) {

	hs, err := bcrypt.GenerateFromPassword([]byte(hashSalt+secret), bcrypt.MinCost)

	if err != nil {

		fmt.Println("error creating hash secret", err)
		return "", errors.New("error creating hash secret")
	}

	return string(hs), nil
}

//CompareHash ...
func CompareHash(hash string, secret string) bool {

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(hashSalt+secret)); err != nil {

		fmt.Println("Password does not match with hash", err)
		return false
	}
	return true
}

//CheckUserType checks user and hash secret
func CheckUserType(user *usr.User) error {

	if user.UserType == usr.Plaftorm {

		hash, err := SecretHash(user.UserSecret)
		if err != nil {
			return err
		}
		user.UserSecret = hash
	} else {

		user.UserType = usr.API
		user.UserSecret = ""
	}

	return nil
}

//NormalizeName ...
func NormalizeName(orgName string, randChars bool) string {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {

		fmt.Print(err)
		if randChars {
			return "defaul_schema_" + GenerateRandString(5)
		}
		return "defaul_acc_" + GenerateRandString(5)
	}

	orgName = reg.ReplaceAllString(orgName, "")
	if randChars {
		orgName = orgName + "_" + GenerateRandString(4)
	}
	orgName = strings.ToLower(orgName)

	return orgName
}
