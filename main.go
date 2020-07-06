package main

//var (
//	hostsPath string
//	cmd       string
//	users string
//	path string
//	f         *flag.FlagSet
//)
//
//func init() {
//	f = flag.NewFlagSet("multi_ssh", flag.ContinueOnError)
//	f.StringVar(&hostsPath, "hosts", "./hosts", "使用multi—ssh要处理的机器")
//	f.StringVar(&cmd, "cmd", "", "需要执行的命令")
//	f.StringVar(&users, "users", "", "可手动添加的用户，用`,`分隔")
//	f.StringVar(&path, "copy", "", "拷贝文件")
//}
//
//func main() {
//	if err := f.Parse(os.Args[1:]); err != nil {
//		log.Fatalln("参数解析失败:", err.Error())
//	}
//	terms, err := config.ReadHosts(hostsPath)
//	if err != nil {
//		log.Fatalln("读取hosts文件失败", err.Error())
//	}
//	candes := make([]*candy, 0)
//	for _, i := range terms {
//		c := &candy{
//			u:   i,
//			cmd: cmd,
//		}
//		candes = append(candes, c)
//	}
//	out := make(chan pack, 0)
//	go exec2(out, candes...)
//	func() {
//		for o := range out {
//			fmt.Printf("%s\n\t%s", o.userInfo.Host(), o.msg)
//		}
//	}()
//}
