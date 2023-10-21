package asset

type Stock_Price struct{
	Open 	float64  	`json:"open"`
	Close 	float64		`json:"close"`
	High 	float64		`json:"high"`
	Low 	float64		`json:"low"`
	Volume 	int			`json:"volume"`
}

type Stocks struct{
	Price 	map[string]Stock_Price //[time]stock_proce
	Type 	int
	Name	string
}

func (s *Stocks)Get_Price()(err error, message *string){
	return nil, nil
}

