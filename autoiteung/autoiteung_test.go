package autoiteung

import (
	"fmt"
	"testing"
)

func TestBukaKelas(t *testing.T) {
	nama_group := "21666-2A-PEMOGRAMAN III | TYGUSAD@62114545454"
	pesaniteung := BukaKelas(nama_group)
	fmt.Println(pesaniteung)

}
