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

const INCRSLEEP = "MULTI_SSH_EXIT_CODE=`echo $?`;sleep 1;exit $MULTI_SSH_EXIT_CODE"

func autoIncrSleep(t *Terminal) {
	t.currentCmd = strings.TrimSpace(t.currentCmd)
	switch t.currentCmd[len(t.currentCmd)-1] {
	case '&':
		t.currentCmd += " " + INCRSLEEP
	case ';':
		t.currentCmd += INCRSLEEP
	default:
		t.currentCmd += ";" + INCRSLEEP
	}
}
