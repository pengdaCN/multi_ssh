package m_terminal

import (
	"fmt"
	"io"
	"multi_ssh/model"
	"strings"
)

const sudoPrefix = "[sudo] password for %s: "

func sudo(t *Terminal, in []byte, out io.WriteCloser) error {
	line := string(in)
	beenMatched := fmt.Sprintf(sudoPrefix, t.user.User())
	if strings.Contains(beenMatched, line) {
		u, _ := t.user.(*model.SSHUserByPassphrase)
		_, err := out.Write([]byte(u.Password + "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
