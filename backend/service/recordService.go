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
	filter := bson.M{"date.realDate": date}
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
	filter := bson.M{"date.realDate": date}
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

func GetWeekReport(week, month, year int) dto.WeekSummary {
	weekDays, err := WeekDaysInCustomMonth(year, month, week, time.Local)
	fmt.Println(weekDays)
	if err != nil {
		fmt.Print(err)
		return dto.WeekSummary{}
	}
	dayMap := make(map[int]dto.DaySummary)
	firstDay := 6
	lastDay := 0
	for k, day := range weekDays {
		daySummary := GetDayReport(day.Format("2006-01-02"))
		if daySummary.Total > 0 {
			dayMap[k] = daySummary
		} else {
			_, weekISO := day.ISOWeek()
			dayMap[k] = dto.DaySummary{
				Date: dto.Date{
					RealDate:      day.Format("2006-01-02"),
					ContableMonth: month,
					Day:           day.Day(),
					DayOfWeek:     day.Weekday().String(),
					Year:          day.Year(),
					WeekOfYear:    weekISO,
					WeekOfMonth:   week,
				},
			}
		}
		//get the lower key
		if k < firstDay {
			firstDay = k
		}
		//get the higher key
		if k > lastDay {
			lastDay = k
		}
	}
	weekSummary := dto.WeekSummary{
		Week:             week,
		StartDate:        weekDays[firstDay].Format("2006-01-02"),
		EndDate:          weekDays[lastDay].Format("2006-01-02"),
		DaySummary:       dayMap,
		SupervivenciaSum: GetCategoriesSum(dayMap, dto.Supervivencia),
		OcioYVicioSum:    GetCategoriesSum(dayMap, dto.OcioYVicio),
		ComprasSum:       GetCategoriesSum(dayMap, dto.Compras),
	}
	weekSummary.Total = weekSummary.SupervivenciaSum["total"] + weekSummary.OcioYVicioSum["total"] + weekSummary.ComprasSum["total"]
	return weekSummary
}

func getCategorySummary(records []dto.Record, category dto.Category) dto.CategorySummary {
	categorySummary := dto.CategorySummary{
		Description: "",
		Subcategory: []dto.SubCategorySummary{},
		Sum:         0,
	}
	subCategoryMap := make(map[string][]dto.RecordDTO)
	for _, record := range records {
		if record.Subcategory.Category == category {
			categorySummary.Sum += record.Amount
			categorySummary.Description = string(category)
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

func GetCategoriesSum(dayMap map[int]dto.DaySummary, category dto.Category) map[string]float64 {
	sum := make(map[string]float64)
	sum[string("total")] = 0
	for _, day := range dayMap {
		if day.Supervivencia.Description == string(category) {
			for _, v := range day.Supervivencia.Subcategory {
				sum[v.Description] += v.Sum
				sum[string("total")] += v.Sum
			}
		}
		if day.OcioYVicio.Description == string(category) {
			for _, v := range day.OcioYVicio.Subcategory {
				sum[v.Description] += v.Sum
				sum[string("total")] += v.Sum
			}
		}
		if day.Compras.Description == string(category) {
			for _, v := range day.Compras.Subcategory {
				sum[v.Description] += v.Sum
				sum[string("total")] += v.Sum
			}
		}
	}
	return sum
}
