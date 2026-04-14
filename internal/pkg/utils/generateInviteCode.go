package utils

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"
)

const words = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func GenerateInviteCode(lenCode int) string {
	slice := strings.Split(words, "")
	var sliceIC []string
	slog.Debug("slice words", "slice: ", slice)
	for i := 0; i < lenCode; i++ {
		randW := rand.IntN(62)
		slog.Debug("Random Num", "randW: ", randW)
		tempIC := slice[randW]
		sliceIC = append(sliceIC, tempIC)
	}
	return strings.Join(sliceIC, "")
}

func GenerateInviteCodeURL(inviteCode string, botName string) string {
	return fmt.Sprintf("t.me/%s?start=%s", botName, inviteCode)
}
