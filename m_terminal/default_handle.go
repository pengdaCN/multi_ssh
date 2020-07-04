package m_terminal

//import (
//	"log"
//	"multi_ssh/model"
//)

//func ExpectHandleSudo(out string, t *Terminal, s *TermSession) {
//	u, ok := t.user.(*model.SSHUserByPassphrase)
//	if !ok {
//		log.Fatalln("ERROR user password not found")
//	}
//	_, err := s.stdin.Write([]byte(u.Password + "\n"))
//	if err != nil {
//		log.Fatalln("ERROR 执行sudo时失败")
//	}
//}
