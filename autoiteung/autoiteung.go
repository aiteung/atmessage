package autoiteung

import (
	"strconv"
	"strings"
)

func BukaKelas(group_name string) string {
	messages := ""
	listgroup := strings.Split(group_name, "-")
	_, err := strconv.Atoi(strings.TrimSpace(listgroup[0]))
	if (len(listgroup) >= 3) && (err == nil) {
		tokendosenpengganti := strings.Split(listgroup[2], "|")
		token := strings.Split(tokendosenpengganti[1], "@")
		if len(token) > 1 {
			messages = "iteung kelas luring dosen pengganti mulai mode tm passcode " + strings.TrimSpace(token[0])
		} else {
			messages = "iteung kelas luring mulai mode tm"
		}

	}
	return messages
}
