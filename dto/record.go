package dto

type RecordRequest struct {
	Description string      `json:"description" bson:"description"`
	Date        string      `json:"date" binding:"required"` // <- string
	Subcategory Subcategory `json:"subcategory" bson:"subcategory"`
	Amount      float64     `json:"amount" bson:"amount"`
	Notes       string      `json:"notes" bson:"notes"`
}

type Record struct {
	Description string      `json:"description,omitempty" bson:"description,omitempty"`
	Date        Date        `json:"date,omitempty" bson:"date,omitempty"`
	Subcategory Subcategory `json:"subcategory,omitempty" bson:"subcategory,omitempty"`
	Amount      float64     `json:"amount,omitempty" bson:"amount,omitempty"`
	Notes       string      `json:"notes,omitempty" bson:"notes,omitempty"`
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
	Week             int                `json:"week" bson:"week"`
	StartDate        string             `json:"startDate" bson:"startDate"`
	EndDate          string             `json:"endDate" bson:"endDate"`
	DaySummary       map[int]DaySummary `json:"daySummary" bson:"daySummary"`
	SupervivenciaSum map[string]float64 `json:"supervivenciaSum" bson:"supervivenciaSum"`
	OcioYVicioSum    map[string]float64 `json:"ocioyvicioSum" bson:"ocioyvicioSum"`
	ComprasSum       map[string]float64 `json:"comprasSum" bson:"comprasSum"`
	Total            float64            `json:"total" bson:"total"`
}

type DaySummary struct {
	//add omitempty
	Date          Date            `json:"date,omitempty" bson:"date,omitempty"`
	Supervivencia CategorySummary `json:"supervivencia,omitempty" bson:"supervivencia,omitempty"`
	OcioYVicio    CategorySummary `json:"ocioyvicio,omitempty" bson:"ocioyvicio,omitempty"`
	Compras       CategorySummary `json:"compras,omitempty" bson:"compras,omitempty"`
	Total         float64         `json:"total" bson:"total"`
}

type CategorySummary struct {
	//add omitempty
	Description string               `json:"description,omitempty" bson:"description,omitempty"`
	Subcategory []SubCategorySummary `json:"subcategory,omitempty" bson:"subcategory,omitempty"`
	Sum         float64              `json:"sum,omitempty" bson:"sum,omitempty"`
}

type SubCategorySummary struct {
	//add omitempty
	Description string      `json:"description,omitempty" bson:"description,omitempty"`
	Records     []RecordDTO `json:"records,omitempty" bson:"records,omitempty"`
	Sum         float64     `json:"sum,omitempty" bson:"sum,omitempty"`
}

type RecordDTO struct {
	//add omitempty
	Description string  `json:"description,omitempty" bson:"description,omitempty"`
	Amount      float64 `json:"amount,omitempty" bson:"amount,omitempty"`
	Notes       string  `json:"notes,omitempty" bson:"notes,omitempty"`
}

type WeekLimits struct {
	Week      int    `json:"week" bson:"week"`
	Month     int    `json:"month" bson:"month"`
	StartDate string `json:"startDate" bson:"startDate"`
	EndDate   string `json:"endDate" bson:"endDate"`
}
