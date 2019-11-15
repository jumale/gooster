package command

func Complete(cmd string, completion string) string {
	for i := len(cmd) - 1; i >= 0; i -= 1 {
		if rune(cmd[i]) == space && (i-1 > 0 || rune(cmd[i-1]) != escape) {
			return cmd[:i+1] + completion + " "
		}
	}
	return completion
}
