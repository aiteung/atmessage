package atwa

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

//Client *whatsmeow.Client

func GetLiveLoc(Message *waProto.Message) (lat float64, long float64) {
	if Message.LiveLocationMessage != nil {
		lat = *Message.LiveLocationMessage.DegreesLatitude
		long = *Message.LiveLocationMessage.DegreesLongitude
	} else {
		fmt.Println("LiveLocationMessage : ", Message.LiveLocationMessage)
	}
	return lat, long
}

func SendMessage(msg string, toJID types.JID) (resp whatsmeow.SendResponse, err error) {
	resp, err = model.Client.SendMessage(context.Background(), toJID, "", &waProto.Message{
		Conversation: proto.String(msg),
	})
	return resp, err
}

func SendListMessage(lstmsg model.ListMessage, toJID types.JID) (resp whatsmeow.SendResponse, err error) {
	var lms []*waProto.ListMessage_Section
	for _, sec := range lstmsg.Sections {

		var lmr []*waProto.ListMessage_Row
		for _, lst := range sec.Rows {
			tmplst := waProto.ListMessage_Row{
				Title:       proto.String(lst.Title),
				Description: proto.String(lst.Description),
				RowId:       proto.String(lst.RowId),
			}
			lmr = append(lmr, &tmplst)
		}

		tmpsec := waProto.ListMessage_Section{
			Title: proto.String(sec.Title),
			Rows:  lmr,
		}
		lms = append(lms, &tmpsec)
	}

	message := &waProto.Message{
		ListMessage: &waProto.ListMessage{
			Title:       proto.String(lstmsg.Title),
			Description: proto.String(lstmsg.Description),
			FooterText:  proto.String(lstmsg.FooterText),
			ButtonText:  proto.String(lstmsg.ButtonText),
			ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
			Sections:    lms,
		},
	}
	viewOnce := &waProto.Message{
		ViewOnceMessage: &waProto.FutureProofMessage{
			Message: message,
		},
	}
	resp, err = model.Client.SendMessage(context.Background(), toJID, "", viewOnce)
	return resp, err

}

func SendButtonMessage(btnmsg model.ButtonsMessage, toJID types.JID) (resp whatsmeow.SendResponse, err error) {
	var buttons []*waProto.ButtonsMessage_Button
	for _, btn := range btnmsg.Buttons {
		tmpbtn := waProto.ButtonsMessage_Button{
			ButtonId: proto.String(btn.ButtonId),
			ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
				DisplayText: proto.String(btn.DisplayText),
			},
			Type: waProto.ButtonsMessage_Button_RESPONSE.Enum(),
		}
		buttons = append(buttons, &tmpbtn)
	}
	this_message := &waProto.Message{
		ButtonsMessage: &waProto.ButtonsMessage{
			ContentText: proto.String(btnmsg.Message.ContentText),
			FooterText:  proto.String(btnmsg.Message.FooterText),
			Buttons:     buttons,
			HeaderType:  waProto.ButtonsMessage_TEXT.Enum(),
			Header: &waProto.ButtonsMessage_Text{
				Text: btnmsg.Message.HeaderText,
			},
		},
	}
	viewOnce := &waProto.Message{
		ViewOnceMessage: &waProto.FutureProofMessage{
			Message: this_message,
		},
	}
	resp, err = model.Client.SendMessage(context.Background(), toJID, "", viewOnce)
	return resp, err
}
