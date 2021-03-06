package model

import (
	"github.com/leekchan/accounting"
)

type Item struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Title    string `json:"title"`
	Price    int    `json:"price"`
	ImageURL string `json:"imageUrl"`
}

func (i Item) GetFormattedPrice() string {
	a := accounting.Accounting{
		Symbol:         "ā½",
		Precision:      0,
		Thousand:       "ā",
		Decimal:        ",",
		Format:         "%vā%s",
	}
	return a.FormatMoney(i.Price)
}
