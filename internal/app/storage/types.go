package storage

type Metric struct {
	TimeStamp   string   `json:"TimeStamp"`
	IsApp       uint8    `json:"IsApp"`
	IsAuth      uint8    `json:"IsAuth"`
	IsNew       uint8    `json:"IsNew"`
	ResWidth    uint16   `json:"ResWidth"`
	ResHeight   uint16   `json:"ResHeight"`
	UserAgent   string   `json:"UserAgent"`
	UserID      string   `json:"UserID"`
	SessionID   string   `json:"SessionID"`
	DeviceType  string   `json:"DeviceType"`
	Reffer      string   `json:"Reffer"`
	Stage       string   `json:"Stage"`
	Action      string   `json:"Action"`
	ExtraKeys   []string `json:"ExtraKeys"`
	ExtraValues []string `json:"ExtraValues"`
}
