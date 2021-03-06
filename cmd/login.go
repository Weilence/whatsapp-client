package cmd

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"go.mau.fi/whatsmeow"
	"log"
	"strings"
	"whatsapp-client/pkg/utils"
	"whatsapp-client/pkg/whatsapp"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login by phone",
	Args:  cobra.ExactArgs(1),
	Run: newRun(func(str string) (stop bool) {
		result, msg, _ := strings.Cut(str, ",")

		if result == "code" {
			code, err := qrcode.New(msg, qrcode.Low)
			utils.NoError(err)
			fmt.Println(code.ToSmallString(false))
		} else {
			fmt.Println(msg)
			return true
		}
		return
	}),
}

func (ctx *Context) login() {
	var client *whatsapp.Client
	jid := ctx.args[1]
	if c, _ := whatsapp.GetClient(jid); c != nil {
		ctx.Write("error,当前已登录，请先退出")
		return
	}

	client = whatsapp.NewClient(jid)
	qrChan := client.Login()

	if qrChan == nil {
		ctx.Write("success,直接登录成功")
		return
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			ctx.Write("code," + evt.Code)
		} else if evt == whatsmeow.QRChannelSuccess {
			ctx.Write("success,扫码登录成功")
			break
		} else if evt == whatsmeow.QRChannelScannedWithoutMultidevice {
			ctx.Write("error,请开启多设备测试版")
			break
		} else {
			log.Println(evt)
			ctx.Write("error,扫码登录失败")
			break
		}
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
