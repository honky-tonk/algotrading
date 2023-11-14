package asset

type Asset interface {
	Get_Price(int) error
}

const (
	Daily = iota + 1
	Weekly
	Monthly
)
