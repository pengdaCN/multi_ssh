package cmd

import (
	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"multi_ssh/m_terminal"
	"multi_ssh/tools"
	"os"
	"sync"
)

var (
	mode       string
	uid        int
	gid        int
	copySudo   bool
	copyExists bool
)

func init() {
	rootCmd.AddCommand(&copyCmd)
	copyCmd.Flags().StringVar(&mode, "mode", "", "设置文件上传后的权限")
	copyCmd.Flags().IntVar(&uid, "uid", -1, "设置上传后文件uid")
	copyCmd.Flags().IntVar(&gid, "gid", -1, "设置上传后文件的gid")
	copyCmd.Flags().BoolVarP(&copySudo, "sudo", "S", false, "可将本地文件无限制的拷贝到远端")
	copyCmd.Flags().BoolVarP(&copyExists, "exists", "e", false, "当远程目录不存在则创建")
}

var copyCmd = cobra.Command{
	Use:   "copy src dst",
	Short: "copy 命令将本地的文件拷贝到远端",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		srcPaths := args[:len(args)-1]
		dstPath := args[len(args)-1]
		ch := make(chan *commandResult, 0)
		var w sync.WaitGroup
		w.Add(1)
		go func() {
			w.Done()
			for m := range ch {
				outputByFormat(outFormat, m, os.Stdout)
			}
		}()
		var w2 sync.WaitGroup
		for _, t := range terminals {
			w2.Add(1)
			go func(term *m_terminal.Terminal) {
				defer w2.Done()
				err := term.Copy(copyExists, copySudo, srcPaths, dstPath, func(file *sftp.File) error {
					if mode != "" {
						m, err := tools.String2FileMode(mode)
						if err != nil {
							return err
						}
						if err := file.Chmod(m); err != nil {
							return err
						}
					}
					if uid != -1 && gid != -1 {
						if err := file.Chown(uid, gid); err != nil {
							return err
						}
					}
					return nil
				})
				if err != nil {
					ch <- &commandResult{
						u:   term.GetUser(),
						msg: []byte(err.Error()),
					}
				} else {
					ch <- &commandResult{
						u:   term.GetUser(),
						msg: []byte("OK"),
					}
				}
			}(t)
		}
		w2.Wait()
		close(ch)
		w.Wait()
	},
}
