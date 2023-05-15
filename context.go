package journify

import (
	"net"
)

type Context struct {
	App       AppInfo      `json:"app,omitempty"`
	Campaign  CampaignInfo `json:"campaign,omitempty"`
	Device    DeviceInfo   `json:"device,omitempty"`
	Library   LibraryInfo  `json:"library,omitempty"`
	Location  LocationInfo `json:"location,omitempty"`
	Network   NetworkInfo  `json:"network,omitempty"`
	OS        OSInfo       `json:"os,omitempty"`
	Page      PageInfo     `json:"page,omitempty"`
	Referrer  ReferrerInfo `json:"referrer,omitempty"`
	IP        net.IP       `json:"ip,omitempty"`
	Direct    bool         `json:"direct,omitempty"`
	Locale    string       `json:"locale,omitempty"`
	Timezone  string       `json:"timezone,omitempty"`
	UserAgent string       `json:"userAgent,omitempty"`
	Traits    Traits       `json:"traits,omitempty"`
}

type AppInfo struct {
	Name      string `json:"name,omitempty"`
	Version   string `json:"version,omitempty"`
	Build     string `json:"build,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type CampaignInfo struct {
	Name    string `json:"name,omitempty"`
	Source  string `json:"source,omitempty"`
	Medium  string `json:"medium,omitempty"`
	Term    string `json:"term,omitempty"`
	Content string `json:"content,omitempty"`
}

type DeviceInfo struct {
	Id            string `json:"id,omitempty"`
	Manufacturer  string `json:"manufacturer,omitempty"`
	Model         string `json:"model,omitempty"`
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	Version       string `json:"version,omitempty"`
	AdvertisingID string `json:"advertisingId,omitempty"`
}

type LibraryInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type LocationInfo struct {
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
	Region    string  `json:"region,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Speed     float64 `json:"speed,omitempty"`
}

type NetworkInfo struct {
	Bluetooth bool   `json:"bluetooth,omitempty"`
	Cellular  bool   `json:"cellular,omitempty"`
	WIFI      bool   `json:"wifi,omitempty"`
	Carrier   string `json:"carrier,omitempty"`
}

type OSInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type PageInfo struct {
	Hash     string `json:"hash,omitempty"`
	Path     string `json:"path,omitempty"`
	Referrer string `json:"referrer,omitempty"`
	Search   string `json:"search,omitempty"`
	Title    string `json:"title,omitempty"`
	URL      string `json:"url,omitempty"`
}

type ReferrerInfo struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
	Link string `json:"link,omitempty"`
}
