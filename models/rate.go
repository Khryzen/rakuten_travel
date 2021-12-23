package models

import (
	"encoding/xml"
	"io"
)

type TableRate struct {
	Date string `xml:"time,attr"`
	Rate []Rate
}

type Rate struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

type Rates struct {
	Rates RateList `xml:"Cube>Cube"`
}

type RateList []TableRate

func (ls *RateList) UnmarshalXML(decode *xml.Decoder, start xml.StartElement) error {
	var rate TableRate

	date := start.Attr[0].Value
	rate.Date = date

	for {
		token, err := decode.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if startElement, ok := token.(xml.StartElement); ok {
			rate.Rate = append(rate.Rate, Rate{Currency: startElement.Attr[0].Value, Rate: startElement.Attr[1].Value})
			if err := decode.DecodeElement(&rate, &startElement); err != nil {
				return err
			}

			*ls = append(*ls, rate)
		}
	}
}
