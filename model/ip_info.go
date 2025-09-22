package model

type IPInfo struct {
    Country     string  `json:"country"`
    CountryCode string  `json:"countryCode"`
    Region      string  `json:"region"`
    RegionName  string  `json:"regionName"`
    City        string  `json:"city"`
    Zip         string  `json:"zip"`
    Lat         float64 `json:"lat"`
    Lon         float64 `json:"lon"`
    Timezone    string  `json:"timezone"`
    ISP         string  `json:"isp"`
    Org         string  `json:"org"`
    AS          string  `json:"as"`
    IP          string  `json:"query"`
}
