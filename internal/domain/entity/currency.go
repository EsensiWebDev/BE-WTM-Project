package entity

type Currency struct {
	ID         uint
	ExternalID string
	Code       string
	Name       string
	Symbol     string
	IsActive   bool
}
