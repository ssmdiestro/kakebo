package repository

import "os"

var (
	DatabaseName       = os.Getenv("KAKEBO_DB")
	RecordCollection   = os.Getenv("RECORD_TABLE")
	CategoryCollection = os.Getenv("CATEGORY_TABLE")
)

func GetrecordCollection() string {
	return RecordCollection
}

func GetcategoryCollection() string {
	return CategoryCollection
}

func GetdatabaseName() string {
	return DatabaseName
}
