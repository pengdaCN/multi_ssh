package m_terminal

import (
	"fmt"
	"strings"
)

func ExpandCmd(t *Terminal) {
	t.pressCmd(t.currentRawCmd)
}

func TrimSudo(t *Terminal) {
	r := t.content.result
	r.msg = strings.TrimSpace(r.msg)
	beenMatched := fmt.Sprintf(sudoPrefix, t.GetUser().User())
	if strings.HasPrefix(r.msg, beenMatched) {
		r.msg = strings.Replace(r.msg, beenMatched, "", 1)
	}
	r.msg = strings.TrimSpace(r.msg)
}
