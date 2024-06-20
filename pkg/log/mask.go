package log

import "strings"

func MaskSessionID(input string) string {
	const sessionID = "SESSIONID="
	const sessionIDLen = 10 // len(sessionID)

	out := []rune(input)
	count := strings.Count(strings.ToUpper(input), sessionID)

	var offset int
	for k := 0; k < count; k++ {
		subURI := input[offset:]
		if index := strings.Index(strings.ToUpper(subURI), sessionID); index > -1 {
			j := 0
			for i := index + sessionIDLen; i < len(subURI); i++ {
				if subURI[i] == '&' || subURI[i] == '.' {
					offset += i
					break
				}
				if j > sessionIDLen+1 {
					out[i+offset] = '*'
				}
				j++
			}
		}
	}
	return string(out)
}
