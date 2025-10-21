package dto

type RecordRequest struct {
	Description string      `json:"description" bson:"description"`
	Date        string      `json:"date" binding:"required"` // <- string
	Subcategory Subcategory `json:"subcategory" bson:"subcategory"`
	Amount      float64     `json:"amount" bson:"amount"`
	Notes       string      `json:"notes" bson:"notes"`
}

type Subcategory struct {
	Description string   `json:"description" bson:"description"`
	Category    Category `json:"category" bson:"category"`
}

type Category string

const (
	Supervivencia Category = "Supervivencia"
	OcioYVicio    Category = "Ocio y Vicio"
	Compras       Category = "Compras"
)
