package dict

type Word struct {
	Name      string     `json:"name"`
	Trans     []Tran     `json:"trans"`
	USPhone   string     `json:"us_phone"`
	UKPhone   string     `json:"uk_phone"`
	Sentences []Sentence `json:"sentences"`

	Mistake bool `json:"-"`
	Degree  int  `json:"-"`
	IsNew   bool `json:"-"`
}

type Tran struct {
	Pos       string `json:"pos"`
	TranCn    string `json:"tran_cn"`
	TranOther string `json:"tran_other"`
	DescCn    string `json:"desc_cn"`
	DescOther string `json:"desc_other"`
}

type Sentence struct {
	Content string `json:"content"`
	Tran    string `json:"tran"`
}
