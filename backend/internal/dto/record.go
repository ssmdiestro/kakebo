package dto

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
