package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mingalevme/avito/internal/model"
	log "github.com/mingalevme/gologger"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
	"strings"
)

type Parser struct {
	HTMLDocumentGetter HTMLDocumentGetter
	Logger             log.Logger
}

func (p Parser) Parse(sourceUrl string) ([]model.Item, error) {
	u, err := url.Parse(sourceUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid source url")
	}
	doc, err := p.HTMLDocumentGetter.Get(sourceUrl)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting HTML doc")
	}
	items := make([]model.Item, 0)
	doc.Find("div[data-marker=\"item\"]").Each(func(i int, div *goquery.Selection) {
		item, err := p.parseItem(div, u)
		if err != nil {
			p.Logger.WithError(err).Error("error while parsing item")
		} else {
			items = append(items, item)
		}
	})
	return items, err
}

func (p Parser) parseItem(div *goquery.Selection, u *url.URL) (model.Item, error) {
	id, err := parseId(div)
	if err != nil {
		return model.Item{}, errors.Wrap(err, "error while getting item id")
	}
	title, err := parseTitleFromAnchor(div)
	if err != nil {
		p.Logger.WithError(err).Error("error while getting item title")
	}
	if title == "" {
		title, err = parseTitleFromHeader(div)
		if err != nil {
			p.Logger.WithError(err).Error("error while getting item title")
			title = fmt.Sprintf("#%d", id)
		}
	}
	linkPath, err := parseLinkPath(div)
	if err != nil {
		p.Logger.WithError(err).Error("error while getting item link path")
		title = fmt.Sprintf("#%d", id)
	}
	link := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, linkPath)
	price, err := parsePrice(div)
	if err != nil {
		p.Logger.WithError(err).Error("error while getting item price")
	}
	imageUrl, err := parseImageUrl(div)
	if err != nil {
		p.Logger.Infof("%d %s", id, title)
		p.Logger.WithError(err).Error("error while getting item image url")
	}
	return model.Item{
		ID:       id,
		URL:      link,
		Title:    title,
		Price:    price,
		ImageURL: imageUrl,
	}, nil
}

func parseId(div *goquery.Selection) (int, error) {
	idStr, ok := div.Attr("data-item-id")
	if !ok {
		return 0, errors.New("attr does not exist")
	}
	idStr = strings.TrimSpace(idStr)
	if idStr == "" {
		return 0, errors.New("attr is empty")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("error while converting to int")
	}
	if id < 1 {
		return 0, errors.New("id < 1")
	}
	return id, nil
}

func parseTitleFromHeader(div *goquery.Selection) (string, error) {
	h3 := div.Find("h3")
	if h3.Length() < 1 {
		return "", errors.New("h3 element is not found")
	}
	title := strings.TrimSpace(h3.Text())
	if title == "" {
		return "", errors.New("element is empty")
	}
	return title, nil
}

func parseTitleFromAnchor(div *goquery.Selection) (string, error) {
	a := div.Find("a")
	if a.Length() < 2 {
		return "", nil
	}
	a = a.Eq(1)
	title, ok := a.Attr("title")
	if !ok {
		return "", nil
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return "", nil
	}
	return title, nil
}

func parseLinkPath(div *goquery.Selection) (string, error) {
	a := div.Find("a")
	if a.Length() < 1 {
		return "", errors.New("\"a\" is not found")
	}
	href, ok := a.Attr("href")
	if !ok {
		return "", errors.New("href-attr does not exist")
	}
	href = strings.TrimSpace(href)
	if href == "" {
		return "", errors.New("href-attr is empty")
	}
	return href, nil
}

func parsePrice(div *goquery.Selection) (int, error) {
	meta := div.Find("meta[itemprop=\"price\"]")
	if meta.Length() < 1 {
		return 0, errors.New("meta-element does not exist")
	}
	content, ok := meta.Attr("content")
	if !ok {
		return 0, errors.New("meta element does not contain content-attr")
	}
	priceStr := strings.TrimSpace(content)
	if priceStr == "" {
		return 0, errors.New("content-attr is empty")
	}
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return 0, errors.New("error while converting to int")
	}
	if price < 1 {
		return 0, errors.New("price < 1")
	}
	return price, nil
}

func parseImageUrl(div *goquery.Selection) (string, error) {
	img := div.Find("img")
	if img.Length() < 1 {
		return "", nil
	}
	src, ok := img.Attr("src")
	if !ok {
		return "", errors.New("src-attr does not exist")
	}
	src = strings.TrimSpace(src)
	if src == "" {
		return "", errors.New("src-attr is empty")
	}
	return src, nil
}
