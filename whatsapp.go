package atmessage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aiteung/musik"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func RunWA(handler whatsmeow.EventHandler) (waclient *whatsmeow.Client) {
	fmt.Println("Starting Whatsapp")
	dbLog := waLog.Stdout("Database", "ERROR", true)
	musik.CreateFolderifNotExist("./session/")
	container, err := sqlstore.New("sqlite3", "file:./session/gowa.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "ERROR", true)
	waclient = whatsmeow.NewClient(deviceStore, clientLog)
	waclient.AddEventHandler(handler)
	if waclient.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := waclient.GetQRChannel(context.Background())
		err = waclient.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = waclient.Connect()
		if err != nil {
			panic(err)
		}
		fmt.Println("Client Connected")
	}
	return

}

func GetLiveLoc(Message *waProto.Message) (lat float64, long float64) {
	if Message.LiveLocationMessage != nil {
		lat = *Message.LiveLocationMessage.DegreesLatitude
		long = *Message.LiveLocationMessage.DegreesLongitude
	} else {
		fmt.Println("LiveLocationMessage : ", Message.LiveLocationMessage)
	}
	return lat, long
}

func SendMessage(msg string, toJID types.JID, whatsapp *whatsmeow.Client) (resp whatsmeow.SendResponse, err error) {
	var wamsg waProto.Message
	wamsg.Conversation = proto.String(msg)
	resp, err = whatsapp.SendMessage(context.Background(), toJID, &wamsg)
	return resp, err
}

func SendListMessage(lstmsg ListMessage, toJID types.JID, whatsapp *whatsmeow.Client) (resp whatsmeow.SendResponse, err error) {
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
	resp, err = whatsapp.SendMessage(context.Background(), toJID, viewOnce)
	return resp, err

}

func SendButtonMessage(btnmsg ButtonsMessage, toJID types.JID, whatsapp *whatsmeow.Client) (resp whatsmeow.SendResponse, err error) {
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
	resp, err = whatsapp.SendMessage(context.Background(), toJID, viewOnce)
	return resp, err
}

func SendDocumentMessage(plaintext []byte, filename string, caption string, toJID types.JID, whatsapp *whatsmeow.Client) (resp whatsmeow.SendResponse, err error) {
	respupload, err := whatsapp.Upload(context.Background(), plaintext, whatsmeow.MediaDocument)
	if err != nil {
		msg := fmt.Sprintf("SendDocumentMessage to wa server : %s", err)
		fmt.Println(msg)
		SendMessage(msg, toJID, whatsapp)
	}
	docMsg := &waProto.DocumentMessage{
		Caption:       proto.String(caption),
		Mimetype:      proto.String(http.DetectContentType(plaintext)),
		FileName:      &filename,
		Url:           &respupload.URL,
		DirectPath:    &respupload.DirectPath,
		MediaKey:      respupload.MediaKey,
		FileEncSha256: respupload.FileEncSHA256,
		FileSha256:    respupload.FileSHA256,
		FileLength:    &respupload.FileLength,
	}
	docMessage := &waProto.Message{
		DocumentMessage: docMsg,
	}
	resp, err = whatsapp.SendMessage(context.Background(), toJID, docMessage)
	return resp, err
}

func SendImageMessage(plaintext []byte, caption string, toJID types.JID, whatsapp *whatsmeow.Client) (resp whatsmeow.SendResponse, err error) {
	respupload, err := whatsapp.Upload(context.Background(), plaintext, whatsmeow.MediaImage)
	if err != nil {
		msg := fmt.Sprintf("SendImageMessage to wa server : %s", err)
		fmt.Println(msg)
		SendMessage(msg, toJID, whatsapp)
	}
	imgMsg := &waProto.ImageMessage{
		Caption:       proto.String(caption),
		Url:           proto.String(respupload.URL),
		DirectPath:    proto.String(respupload.DirectPath),
		MediaKey:      respupload.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(plaintext)),
		FileEncSha256: respupload.FileEncSHA256,
		FileSha256:    respupload.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(plaintext))),
	}

	imgMessage := &waProto.Message{
		ImageMessage: imgMsg,
	}
	resp, err = whatsapp.SendMessage(context.Background(), toJID, imgMessage)
	return resp, err
}
