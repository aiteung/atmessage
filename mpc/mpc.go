package mpc

import (
	"fmt"
	"strings"

	"github.com/aiteung/atmessage"
	"github.com/aiteung/atmessage/autoiteung"
	"github.com/aiteung/atmessage/iteung"
	"github.com/aiteung/atmessage/mediadecrypt"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func IteungV1(IteungIPAddress string, apikey string, Info *types.MessageInfo, Message *waProto.Message, waclient *whatsmeow.Client) {
	var im atmessage.IteungMessage
	im.Phone_number = Info.Sender.User
	im.Group_name = Info.Sender.User
	im.Alias_name = Info.PushName
	im.Messages = Message.GetConversation()
	im.Is_group = "false"
	im.Filename = ""
	im.Filedata = ""
	im.Latitude = 0.0
	im.Longitude = 0.0
	im.Api_key = apikey
	if Info.Chat.Server == "g.us" {
		groupInfo, err := waclient.GetGroupInfo(Info.Chat)
		fmt.Println("cek err : ", err)
		if groupInfo != nil {
			im.Group_name = groupInfo.GroupName.Name + "@" + Info.Chat.User
		} else {
			fmt.Println("groupInfo : ", groupInfo)
		}
		im.Is_group = "true"
		if strings.Contains(Message.GetConversation(), "teung") || strings.Contains(Message.GetConversation(), "Teung") {
			go waclient.SendChatPresence(Info.Chat, "composing", "")
		}
	} else {
		go waclient.SendChatPresence(Info.Chat, "composing", "")
	}
	MessageEvent(IteungIPAddress, Info, Message, &im, waclient)
}

func MessageEvent(IteungIPAddress string, Info *types.MessageInfo, Message *waProto.Message, im *atmessage.IteungMessage, waclient *whatsmeow.Client) {
	if Info.MediaType == "livelocation" {
		LiveLoc(Message, Info, im, waclient)
	}
	if Message.ExtendedTextMessage != nil {
		Extended(Message, im)

	}
	if Info.MediaType == "document" {
		Document(Message, im)
	}
	if im.Messages != "" {
		iteung.Send(IteungIPAddress, im, Info.Chat, waclient)
	}
}

func LiveLoc(Message *waProto.Message, Info *types.MessageInfo, im *atmessage.IteungMessage, waclient *whatsmeow.Client) {
	if Message.LiveLocationMessage != nil {
		im.Latitude = *Message.LiveLocationMessage.DegreesLatitude
		im.Longitude = *Message.LiveLocationMessage.DegreesLongitude
		if im.Is_group == "true" {
			im.Messages = autoiteung.BukaKelas(im.Group_name)
		} else {
			LiveLocinPrivateMessage(Message, Info, im, waclient)
		}
	} else {
		fmt.Println("LiveLocationMessage : ", Message.LiveLocationMessage)
	}
}

func Extended(Message *waProto.Message, im *atmessage.IteungMessage) {
	im.Messages = *Message.ExtendedTextMessage.Text
	fmt.Println(Message)
	if Message.ExtendedTextMessage.ContextInfo != nil {
		if Message.ExtendedTextMessage.ContextInfo.Participant != nil {
			im.Phone_number = strings.Split(*Message.ExtendedTextMessage.ContextInfo.Participant, "@")[0]
			if Message.ExtendedTextMessage.ContextInfo.QuotedMessage.LiveLocationMessage != nil {
				im.Latitude = *Message.ExtendedTextMessage.ContextInfo.QuotedMessage.LiveLocationMessage.DegreesLatitude
				im.Longitude = *Message.ExtendedTextMessage.ContextInfo.QuotedMessage.LiveLocationMessage.DegreesLongitude
			}
			if Message.ExtendedTextMessage.ContextInfo.QuotedMessage.DocumentMessage != nil {
				im.Filename = *Message.ExtendedTextMessage.ContextInfo.QuotedMessage.DocumentMessage.DirectPath
				im.Filedata = mediadecrypt.GetBase64Filedata(Message.ExtendedTextMessage.ContextInfo.QuotedMessage.DocumentMessage.Url, Message.ExtendedTextMessage.ContextInfo.QuotedMessage.DocumentMessage.MediaKey)
			}
			if Message.ExtendedTextMessage.ContextInfo.QuotedMessage.DocumentWithCaptionMessage != nil {
				im.Filename = *Message.ExtendedTextMessage.ContextInfo.QuotedMessage.DocumentWithCaptionMessage.Message.DocumentMessage.DirectPath
				im.Filedata = mediadecrypt.GetBase64Filedata(Message.ExtendedTextMessage.ContextInfo.QuotedMessage.DocumentWithCaptionMessage.Message.DocumentMessage.Url, Message.ExtendedTextMessage.ContextInfo.QuotedMessage.DocumentWithCaptionMessage.Message.DocumentMessage.MediaKey)
			}
		} else {
			fmt.Println("ContextInfo", Message.ExtendedTextMessage.ContextInfo)
		}
	} else {
		fmt.Println("ExtendedTextMessage", Message.ExtendedTextMessage)
	}
}

func Document(Message *waProto.Message, im *atmessage.IteungMessage) {
	if Message.DocumentMessage == nil {
		return
	}

	switch {
	case Message.DocumentMessage.Title != nil:
		im.Filename = *Message.DocumentMessage.Title
	case Message.DocumentMessage.FileName != nil:
		im.Filename = *Message.DocumentMessage.FileName

	}

	if Message.DocumentMessage.Caption != nil {
		im.Messages = *Message.DocumentMessage.Caption
	}

	if im.Messages != "" {
		im.Filedata = mediadecrypt.GetBase64Filedata(Message.DocumentMessage.Url, Message.DocumentMessage.MediaKey)
	}
}
