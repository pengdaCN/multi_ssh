package m_terminal

func ExpandCmd(t *Terminal) {
	t.pressCmd(t.currentRawCmd)
}

func TrimSudo(t *Terminal) {

}
