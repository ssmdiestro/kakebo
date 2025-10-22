package dto

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
	Date          Date            `json:"date,omitempty" bson:"date,omitempty"`
	Supervivencia CategorySummary `json:"supervivencia,omitempty" bson:"supervivencia,omitempty"`
	OcioYVicio    CategorySummary `json:"ocioyvicio,omitempty" bson:"ocioyvicio,omitempty"`
	Compras       CategorySummary `json:"compras,omitempty" bson:"compras,omitempty"`
	Total         float64         `json:"total" bson:"total"`
}

type MonthSummary struct {
	Month              int                `json:"month" bson:"month"`
	MonthName          string             `json:"monthName" bson:"monthName"`
	Year               int                `json:"year" bson:"year"`
	WeekSums           map[int]WeekSums   `json:"weekSums" bson:"weekSums"`
	SupervivenciaTotal map[string]float64 `json:"supervivenciaTotal" bson:"supervivenciaTotal"`
	OcioYVicioTotal    map[string]float64 `json:"ocioyvicioTotal" bson:"ocioyvicioTotal"`
	ComprasTotal       map[string]float64 `json:"comprasTotal" bson:"comprasTotal"`
	Total              float64            `json:"total" bson:"total"`
}

type WeekSums struct {
	Week             int                `json:"week" bson:"week"`
	SupervivenciaSum map[string]float64 `json:"supervivenciaSum" bson:"supervivenciaSum"`
	OcioYVicioSum    map[string]float64 `json:"ocioyvicioSum" bson:"ocioyvicioSum"`
	ComprasSum       map[string]float64 `json:"comprasSum" bson:"comprasSum"`
	Total            float64            `json:"total" bson:"total"`
}
type CategorySummary struct {
	Description string               `json:"description,omitempty" bson:"description,omitempty"`
	Subcategory []SubCategorySummary `json:"subcategory,omitempty" bson:"subcategory,omitempty"`
	Sum         float64              `json:"sum,omitempty" bson:"sum,omitempty"`
}

type SubCategorySummary struct {
	Description string      `json:"description,omitempty" bson:"description,omitempty"`
	Records     []RecordDTO `json:"records,omitempty" bson:"records,omitempty"`
	Sum         float64     `json:"sum,omitempty" bson:"sum,omitempty"`
}

type RecordDTO struct {
	Description string  `json:"description,omitempty" bson:"description,omitempty"`
	Amount      float64 `json:"amount,omitempty" bson:"amount,omitempty"`
	Notes       string  `json:"notes,omitempty" bson:"notes,omitempty"`
}

type WeekLimits struct {
	Week      int    `json:"week" bson:"week"`
	Month     int    `json:"month" bson:"month"`
	MonthName string `json:"monthName" bson:"monthName"`
	StartDate string `json:"startDate" bson:"startDate"`
	EndDate   string `json:"endDate" bson:"endDate"`
}
