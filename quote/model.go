package quote

import (
	"github.com/jinzhu/gorm"
)

// Model for the Database table Quote
type QuoteModel struct {
	gorm.Model
	Author string
	Quote  string
}

// struct to use with Create and Update Requests
type QuoteRequest struct {
	Author string `json:author`
	Quote  string `json:quote`
}
