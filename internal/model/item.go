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
	a := accounting.Accounting{Symbol: "₽", Precision: 0, Thousand: ".", Decimal: ","}
	return a.FormatMoney(i.Price)
}
