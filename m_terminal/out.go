package m_terminal

type out chan []byte

func (o out) Write(src []byte) (n int, err error) {
	o <- src
	return len(src), nil
}
