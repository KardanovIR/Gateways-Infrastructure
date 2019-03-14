package models

type NetworkType string

const (
	MainNet NetworkType = "MainNet"
	TestNet NetworkType = "TestNet"
	DevNet  NetworkType = "DevNet"
)

func (n NetworkType) Schema() byte {
	switch n {
	case MainNet:
		return 'W'
	case TestNet:
		return 'T'
	case DevNet:
		return 'D'
	}
	return 0
}
