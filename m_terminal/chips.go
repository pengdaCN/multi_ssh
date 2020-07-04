package m_terminal

import (
	"fmt"
	"multi_ssh/model"
	"strings"
)

const sudoPrefix = "[sudo] password for %s: "

// hookAfterExec
func sudo(t *Terminal, s *TermSession, rst []byte) {
	line := t.termStdoutCache.getLast()
	beenMatched := fmt.Sprintf(sudoPrefix, t.user.User())
	if strings.Contains(beenMatched, line) {
		u, _ := t.user.(*model.SSHUserByPassphrase)
		_, err := s.TermStdin.Write([]byte(u.Password + "\n"))
		if err != nil {
			panic(err)
		}
	}
}
