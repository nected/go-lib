package date

import (
	"time"
)

type DateFormatEnum string

const (
	EMPTY_FORMAT   DateFormatEnum = ""
	IN_DATE_FORMAT DateFormatEnum = "in"
	US_DATE_FORMAT DateFormatEnum = "us"
)

func (e DateFormatEnum) String() string {
	return string(e)
}

type DateLayout string

const (
	Blank          DateLayout = ""
	ISO            DateLayout = "2006-01-02T15:04:05Z"
	ISOMilisecond  DateLayout = "2006-01-02T15:04:05.000Z"
	IsoMicrosecond DateLayout = "2006-01-02T15:04:05.000000Z"

	CommonDatetime1 DateLayout = "2006-01-02 15:04:05"
	CommonDatetime2 DateLayout = "2006/01/02 15:04:05"
	CommonDate1     DateLayout = "2006-01-02"
	CommonDate2     DateLayout = "2006/01/02"

	InDatetime1 DateLayout = "02/01/2006 15:04:05"
	InDatetime2 DateLayout = "02-01-2006 15:04:05"
	InDate1     DateLayout = "02/01/2006"
	InDate2     DateLayout = "02-01-2006"

	Datetime1 DateLayout = "01/02/2006 15:04:05"
	Datetime2 DateLayout = "01-02-2006 15:04:05"
	Date1     DateLayout = "01/02/2006"
	Date2     DateLayout = "01-02-2006"
)

type DateParseOption struct {
	DateFormat DateFormatEnum `json:"dateFormat"`
	Timezone   string         `json:"timezone"`
}

var (
	commonDatetimeLayouts = []DateLayout{
		time.RFC3339,
		ISO, ISOMilisecond, IsoMicrosecond,
		CommonDatetime1, CommonDatetime2,
	}
	usDateLayouts = []DateLayout{CommonDate1, CommonDate2, Date1, Date2}
	inDateLayouts = []DateLayout{CommonDate1, CommonDate2, InDate1, InDate2}

	usDatetimeLayouts = append([]DateLayout{Datetime1, Datetime2}, commonDatetimeLayouts...)
	inDatetimeLayouts = append([]DateLayout{InDatetime1, InDatetime2}, commonDatetimeLayouts...)
)

func (s DateLayout) String() string {
	return string(s)
}

func InspectDateFormat(val string, dateFromat DateFormatEnum) (DateLayout, bool) {
	switch dateFromat {
	case EMPTY_FORMAT, IN_DATE_FORMAT:
		for _, layout := range inDateLayouts {
			if _, err := time.Parse(layout.String(), val); err == nil {
				return layout, true
			}
		}
	case US_DATE_FORMAT:
		for _, layout := range usDateLayouts {
			if _, err := time.Parse(layout.String(), val); err == nil {
				return layout, true
			}
		}
	}
	return Blank, false
}

func InspectDateTimeFormat(val string, dateFormat DateFormatEnum) (DateLayout, bool) {
	switch dateFormat {
	case EMPTY_FORMAT, IN_DATE_FORMAT:
		for _, layout := range inDatetimeLayouts {
			if _, err := time.Parse(layout.String(), val); err == nil {
				return layout, true
			}
		}
	case US_DATE_FORMAT:
		for _, layout := range usDatetimeLayouts {
			if _, err := time.Parse(layout.String(), val); err == nil {
				return layout, true
			}
		}
	}

	return Blank, false
}

func InspectAllTimeFormat(val string, dateFormat DateFormatEnum) (DateLayout, bool) {
	if layout, ok := InspectDateTimeFormat(val, dateFormat); ok {
		return layout, ok
	}
	return InspectDateFormat(val, dateFormat)
}

func GetAllTimeLayout(dateFormat DateFormatEnum) []DateLayout {
	switch dateFormat {
	case EMPTY_FORMAT, IN_DATE_FORMAT:
		return append(inDateLayouts, inDatetimeLayouts...)
	case US_DATE_FORMAT:
		return append(usDateLayouts, usDatetimeLayouts...)
	default:
		return []DateLayout{}
	}
}
