package atwa

//var Client *whatsmeow.Client
//var V interface{}

type MediaType int

const (
	_              = iota
	MediaTypeImage = MediaType(iota)
	MediaTypeVideo
	MediaTypeAudio
	MediaTypeDocument
)

var (
	AppInfo = map[MediaType]string{
		MediaTypeImage:    "WhatsApp Image Keys",
		MediaTypeVideo:    "WhatsApp Video Keys",
		MediaTypeAudio:    "WhatsApp Audio Keys",
		MediaTypeDocument: "WhatsApp Document Keys",
	}
)

type IteungMessage struct {
	Phone_number string  `json:"phone_number"`
	Group_name   string  `json:"group_name"`
	Alias_name   string  `json:"alias_name"`
	Messages     string  `json:"messages"`
	Is_group     string  `json:"is_group"`
	Filename     string  `json:"filename"`
	Filedata     string  `json:"filedata"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Api_key      string  `json:"api_key"`
}

type IteungRespon struct {
	Message string `json:"message"`
}

type Response struct {
	Response string `json:"response"`
}

type QRScan struct {
	QR      string `json:"qr"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Chat struct {
	Phone_number string `json:"phone_number"`
	Messages     string `json:"messages"`
}

type Notif struct {
	User     string `json:"user"`
	Server   string `json:"server"`
	Messages string `json:"messages"`
}

type NotifButton struct {
	User    string         `json:"user"`
	Server  string         `json:"server"`
	Message ButtonsMessage `json:"button_messages"`
}

type WaButton struct {
	ButtonId    string `json:"button_id,omitempty"`
	DisplayText string `json:"display_text,omitempty"`
}

type WaButtonsMessage struct {
	HeaderText  string `json:"header_text,omitempty"`
	ContentText string `json:"content_text,omitempty"`
	FooterText  string `json:"footer_text,omitempty"`
}

type ButtonsMessage struct {
	Message WaButtonsMessage `json:"message,omitempty"`
	Buttons []WaButton       `json:"buttons,omitempty"`
}

type ListMessage struct {
	Title       string
	Description string
	FooterText  string
	ButtonText  string
	Sections    []WaListSection
}

type WaListSection struct {
	Title string
	Rows  []WaListRow
}

type WaListRow struct {
	Title       string
	Description string
	RowId       string
}

type PhoneList struct {
	Phones []string `json:"phones,omitempty"`
}
