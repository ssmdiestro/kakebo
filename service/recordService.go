package service

import (
	"fmt"
	"kakebo/dto"
	"kakebo/repository"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewRecord(recordRequest dto.RecordRequest) error {
	parsed, err := time.Parse("2006-01-02", recordRequest.Date)
	if err != nil {
		log.Println("fecha inválida:", err)
		return err
	}
	date := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
	_, week := date.ISOWeek()
	w, m, _, _, _ := WeekNumberInCustomMonth(date.Format("2006-01-02"), time.Local)
	dateObject := dto.Date{
		RealDate:      date.Format("2006-01-02"),
		WeekOfYear:    week,
		WeekOfMonth:   w,
		ContableMonth: m,
		Day:           date.Day(),
		DayOfWeek:     date.Weekday().String(),
		Year:          date.Year(),
	}
	newRecord := dto.Record{
		Description: recordRequest.Description,
		Date:        dateObject,
		Subcategory: recordRequest.Subcategory,
		Amount:      recordRequest.Amount,
		Notes:       recordRequest.Notes,
	}

	repository.Insert([]dto.Record{newRecord}, os.Getenv("RECORD_TABLE"), os.Getenv("KAKEBO_DB"))
	return nil
}

func GetRecordById(id string) dto.Record {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("GetRecordById - primitive.ObjectIDFromHex(): ", err)
		return dto.Record{}
	}
	filter := bson.M{"_id": objID}
	resultList, err := repository.SelectWithError([]dto.Record{}, filter, os.Getenv("RECORD_TABLE"), os.Getenv("KAKEBO_DB"))
	if err != nil {
		return dto.Record{}
	}
	if len(resultList) > 0 {
		return resultList[0]
	}
	return dto.Record{}
}

func GetRecordByDate(date string) []dto.Record {
	filter := bson.M{"date.date": date}
	resultList, err := repository.SelectWithError([]dto.Record{}, filter, os.Getenv("RECORD_TABLE"), os.Getenv("KAKEBO_DB"))
	if err != nil {
		fmt.Print(err)
		return []dto.Record{}
	}
	if len(resultList) > 0 {
		return resultList
	}
	return []dto.Record{}
}

func GetDayReport(date string) dto.DaySummary {
	filter := bson.M{"date.date": date}
	resultList, err := repository.SelectWithError([]dto.Record{}, filter, os.Getenv("RECORD_TABLE"), os.Getenv("KAKEBO_DB"))
	if err != nil {
		fmt.Print(err)
		return dto.DaySummary{}
	}
	daySummary := dto.DaySummary{}
	for _, result := range resultList {
		daySummary.Date = result.Date
		daySummary.Supervivencia = getCategorySummary(resultList, dto.Supervivencia)
		daySummary.OcioYVicio = getCategorySummary(resultList, dto.OcioYVicio)
		daySummary.Compras = getCategorySummary(resultList, dto.Compras)
		daySummary.Total = daySummary.OcioYVicio.Sum + daySummary.Supervivencia.Sum + daySummary.Compras.Sum
	}
	return daySummary
}

func getCategorySummary(records []dto.Record, category dto.Category) dto.CategorySummary {
	categorySummary := dto.CategorySummary{
		Description: string(category),
		Subcategory: []dto.SubCategorySummary{},
		Sum:         0,
	}
	subCategoryMap := make(map[string][]dto.RecordDTO)
	for _, record := range records {
		if record.Subcategory.Category == category {
			categorySummary.Sum += record.Amount
			newRecord := dto.RecordDTO{
				Description: record.Description,
				Amount:      record.Amount,
				Notes:       record.Notes,
			}
			subCategoryMap[record.Subcategory.Description] = append(subCategoryMap[record.Subcategory.Description], newRecord)
		}
	}
	for k, v := range subCategoryMap {
		newSubCategorySummary := dto.SubCategorySummary{
			Description: k,
			Records:     v,
			Sum:         getSum(v),
		}
		categorySummary.Subcategory = append(categorySummary.Subcategory, newSubCategorySummary)
	}
	return categorySummary
}

func getSum(records []dto.RecordDTO) float64 {
	var sum float64 = 0
	for _, record := range records {
		sum = +record.Amount
	}
	return sum
}
