package dto

type RecordRequest struct {
	Description string      `json:"description" bson:"description"`
	Date        string      `json:"date" binding:"required"` // <- string
	Subcategory Subcategory `json:"subcategory" bson:"subcategory"`
	Amount      float64     `json:"amount" bson:"amount"`
	Notes       string      `json:"notes" bson:"notes"`
}

type Record struct {
	Description string      `json:"description" bson:"description"`
	Date        Date        `json:"date" bson:"date"`
	Subcategory Subcategory `json:"subcategory" bson:"subcategory"`
	Amount      float64     `json:"amount" bson:"amount"`
	Notes       string      `json:"notes" bson:"notes"`
}

type Date struct {
	RealDate      string `json:"realDate" bson:"realDate"`
	ContableMonth int    `json:"contableMonth" bson:"contableMonth"`
	Day           int    `json:"day" bson:"day"`
	DayOfWeek     string `json:"dayofweek" bson:"dayofweek"`
	Year          int    `json:"year" bson:"year"`
	WeekOfYear    int    `json:"weekOfYear" bson:"weekOfYear"`
	WeekOfMonth   int    `json:"weekOfMonth" bson:"weekOfMonth"`
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

type WeekSummary struct {
	Week          int             `json:"week" bson:"week"`
	StartDate     string          `json:"startDate" bson:"startDate"`
	EndDate       string          `json:"endDate" bson:"endDate"`
	Supervivencia CategorySummary `json:"supervivencia" bson:"supervivencia"`
	OcioYVicio    CategorySummary `json:"ocioyvicio" bson:"ocioyvicio"`
	Compras       CategorySummary `json:"compras" bson:"compras"`
	Total         float64         `json:"total" bson:"total"`
}

type DaySummary struct {
	Date          Date            `json:"date" bson:"date"`
	Supervivencia CategorySummary `json:"supervivencia" bson:"supervivencia"`
	OcioYVicio    CategorySummary `json:"ocioyvicio" bson:"ocioyvicio"`
	Compras       CategorySummary `json:"compras" bson:"compras"`
	Total         float64         `json:"total" bson:"total"`
}

type CategorySummary struct {
	Description string               `json:"description" bson:"description"`
	Subcategory []SubCategorySummary `json:"subcategory" bson:"subcategory"`
	Sum         float64              `json:"sum" bson:"sum"`
}

type SubCategorySummary struct {
	Description string      `json:"description" bson:"description"`
	Records     []RecordDTO `json:"records" bson:"records"`
	Sum         float64     `json:"sum" bson:"sum"`
}

type RecordDTO struct {
	Description string  `json:"description" bson:"description"`
	Amount      float64 `json:"amount" bson:"amount"`
	Notes       string  `json:"notes" bson:"notes"`
}

type WeekLimits struct {
	Week      int    `json:"week" bson:"week"`
	Month     int    `json:"month" bson:"month"`
	StartDate string `json:"startDate" bson:"startDate"`
	EndDate   string `json:"endDate" bson:"endDate"`
}
