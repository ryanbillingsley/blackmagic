package blackmagic

import "fmt"

type WUndergroundResponse struct {
	Forecast forecast `json:"forecast"`
}

type forecast struct {
	SimpleForecast forecastSimple `json:"simpleforecast"`
}

type forecastSimple struct {
	ForecastDays []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Date       forecastDate
	Period     int
	High       forecastTemp
	Low        forecastTemp
	Conditions string
	Icon       string
	IconUrl    string `json:"icon_url"`
	Skyicon    string
	Pop        int

	QpfAllDay forecastPrecip `json:"qpf_allday"`
	QpfDay    forecastPrecip `json:"qpf_day"`
	QpfNight  forecastPrecip `json:"qpf_night"`

	SnowAllDay forecastPrecip `json:"snow_allday"`
	SnowDay    forecastPrecip `json:"snow_day"`
	SnowNight  forecastPrecip `json:"snow_night"`

	MaxWind forecastWind `json:"maxwind"`
	AvgWind forecastWind `json:"avewind"`

	AvgHumidity int `json:"avehumidity"`
	MaxHumidity int `json:"maxhumidity"`
	MinHumidity int `json:"minhumidity"`
}

type forecastDate struct {
	Epoch          string `json:"epoch"`
	Pretty         string `json:"pretty"`
	Day            int
	Month          int
	Year           int
	Yday           int
	Hour           int
	Min            string
	Sec            int
	IsDST          string `json:"isdst"`
	MonthName      string `json:"monthname"`
	MonthNameShort string `json:"monthname_short"`
	WeekdayShort   string `json:"weekday_short"`
	Weekday        string
	AmPm           string `json:"ampm"`
	TimeZoneShort  string `json:"tz_short"`
	TimeZoneLong   string `json:"tz_long"`
}

func (date *forecastDate) StandardFormat() string {
	return fmt.Sprintf("%v %v %v %v:%v %v %v", date.WeekdayShort, date.MonthNameShort, date.Day, date.Hour, date.Min, date.TimeZoneShort, date.Year)
}

type forecastTemp struct {
	Fahrenheit string
	Celsius    string
}

type forecastPrecip struct {
	Inches      float64 `json:"in"`
	Millimeters int     `json:"mm"`
}

type forecastWind struct {
	Mph       int    `json:"mpg"`
	Kph       int    `json:"kph"`
	Direction string `json:"dir"`
	Degrees   int    `json:"degrees"`
}
