package mpc

import (
	"fmt"
	"strings"

	"github.com/aiteung/atmessage"
	"github.com/aiteung/atmessage/autoiteung"
	"github.com/aiteung/atmessage/mediadecrypt"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func LiveLoc(Message *waProto.Message, im *atmessage.IteungMessage) {
	if Message.LiveLocationMessage != nil {
		im.Latitude = *Message.LiveLocationMessage.DegreesLatitude
		im.Longitude = *Message.LiveLocationMessage.DegreesLongitude
		im.Messages = autoiteung.BukaKelas(im.Group_name)
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
	if Message.DocumentMessage != nil {
		im.Filename = *Message.DocumentMessage.Title
		im.Messages = *Message.DocumentMessage.Caption
		if im.Messages != "" {
			im.Filedata = mediadecrypt.GetBase64Filedata(Message.DocumentMessage.Url, Message.DocumentMessage.MediaKey)
		}
	}
}
