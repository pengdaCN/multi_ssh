package m_terminal

type out chan []byte

func (o out) Write(src []byte) (n int, err error) {
	o <- src
	return len(src), nil
}

func readFromRemote(t *Terminal, s *TermSession)  {
	stdout, ok := s.Stdout.(out)
	if ! ok {
		panic("stdout 不是out类型")
	}
	stderr, ok := s.Stderr.(out)
	if ! ok {
		panic("stderr 不是out类型")
	}
	for {
		if stdout == nil && stderr == nil {
			break
		}
		select {
		case o, ok := <- stdout:
			if ! ok {
				stdout = nil
				continue
			}
			_, _ = t.termStdoutCache.Write(o)
			s.rst = append(s.rst, o...)
		case o2, ok := <- stderr:
			if ! ok {
				stdout = nil
				continue
			}
			_, _ = t.termStderrCache.Write(o2)
			s.rst = append(s.rst, o2...)
		}
	}
}

//import (
//	"fmt"
//	"regexp"
//)
//
//func out(t *Terminal, s *TermSession) {
//	for {
//		select {
//		case m, ok := <-s.stdout:
//			if !ok {
//				s.stdout = nil
//				continue
//			}
//			str := string(m)
//			fmt.Print(str)
//			_, _ = t.termStdoutCache.Write(m)
//			_, _ = t.termCache.Write(m)
//		case m2, ok := <-s.stderr:
//			if !ok {
//				s.stderr = nil
//				continue
//			}
//			fmt.Print(string(m2))
//			_, _ = t.termStdoutCache.Write(m2)
//			_, _ = t.termCache.Write(m2)
//		}
//	}
//}
//
//func trigger(msg string, t *Terminal, s *TermSession) {
//	for k, v := range t.handle {
//		if ok, err := regexp.MatchString(k, msg); ok && err == nil {
//			v(msg, t, s)
//		}
//	}
//}
