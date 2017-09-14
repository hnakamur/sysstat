package sysstat

func nextField(s []byte) (start, end int) {
	i := 0
	for isSpace(s[i]) {
		i++
	}
	j := i + 1
	for !isSpace(s[j]) {
		j++
	}
	return i, j
}

var asciiSpace = [256]bool{'\t': true, '\n': true, '\v': true, '\f': true, '\r': true, ' ': true}

// isSpace returns true if b is one of '\t', '\n', '\v', '\f', '\r', or ' '.
// It returns false otherwise.
func isSpace(b byte) bool {
	return asciiSpace[b]
}
