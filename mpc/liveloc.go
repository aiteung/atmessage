package mpc

import (
	"fmt"

	"github.com/aiteung/atmessage"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func LiveLocinPrivateMessage(Message *waProto.Message, Info *types.MessageInfo, im *atmessage.IteungMessage, waclient *whatsmeow.Client) {
	msg := "Hai, " + Info.PushName + ". Nomor kamu " + Info.Sender.User + " belum terdaftar di layanan kami. Silakan gunakan nomor whatsapp yang terdaftar. Terima kasih"
	fmt.Println(msg)
	fmt.Println(Info.Sender)
	atmessage.SendMessage(msg, Info.Sender, waclient)
	print("hai")
}
