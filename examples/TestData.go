package examples

import "encoding/xml"

type Name struct {
	First string
	Last  string
}

type table struct {
	XMLName xml.Name `xml:"table"`
	Ths     []string `xml:"thead>tr>th"`
	Trs     []tr     `xml:"tbody>tr"`
}

func newTable(keys []string) *table {
	return &table{
		Ths: keys,
	}
}

func (t *table) addRow(row tr) {
	t.Trs = append(t.Trs, row)
}

type tr struct {
	XMLName xml.Name `xml:"tr"`
	Tds     []string `xml:"td"`
}
