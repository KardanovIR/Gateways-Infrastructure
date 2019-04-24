package models

type NetworkType string

func (n NetworkType) Schema() byte {
	return byte(n[0])
}
