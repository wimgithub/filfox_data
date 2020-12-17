package model

type Data struct {
	ID      uint `gorm:"primary_key"`
	Time    int64
	FilFrom string
	Height  int64
	Message string
	FilTo   string
	Type    string // send:发出  burn-fee:销毁手续费  miner-fee：矿工手续费
	Value   string
}

type Resp struct {
	TotalCount int64             `json:"totalCount"`
	Transfers  []*FilFoxResponse `json:"transfers"`
}
type FilFoxResponse struct {
	Height    int64  `json:"height"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Type      string `json:"type"`
}

type PageData struct {
	Page []int64 `json:"page"`
}
