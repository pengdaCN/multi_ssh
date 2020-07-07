package cmd

import (
	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"multi_ssh/m_terminal"
	"multi_ssh/tools"
	"os"
	"sync"
	"time"
)

var (
	mode string
	uid  int
	gid  int
)

func init() {
	rootCmd.AddCommand(&copyCmd)
	copyCmd.Flags().StringVar(&mode, "mode", "", "设置文件上传后的权限")
	copyCmd.Flags().IntVar(&uid, "uid", -1, "设置上传后文件uid")
	copyCmd.Flags().IntVar(&gid, "gid", -1, "设置上传后文件的gid")
}

var copyCmd = cobra.Command{
	Use:   "copy src dst",
	Short: "copy 命令将本地的文件拷贝到远端",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		srcPaths := args[:len(args)-1]
		dstPath := args[len(args)-1]
		ch := make(chan *commandResult, 0)
		go func() {
			for m := range ch {
				outputByFormat(outFormat, m, os.Stdout)
			}
		}()
		var w sync.WaitGroup
		for _, t := range terminals {
			w.Add(1)
			go func(term *m_terminal.Terminal) {
				defer w.Done()
				err := term.SftpUpdates(srcPaths, dstPath, func(file *sftp.File) error {
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
		w.Wait()
		time.Sleep(time.Second)
		close(ch)
	},
}
