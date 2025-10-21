package service

import (
	"context"
	"fmt"
	"kakebo/internal/dto"
	"kakebo/internal/repository"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Service) NewRecord(ctx context.Context, recordRequest dto.RecordRequest) error {
	newRecord, err := buildRecordFromRequest(recordRequest)
	if err != nil {
		return fmt.Errorf("error al construir el registro: %w", err)
	}
	coll := s.DB.Collection(repository.GetdatabaseName(), repository.GetrecordCollection())
	_, err = repository.InsertOne(ctx, coll, newRecord)
	if err != nil {
		return fmt.Errorf("error al insertar el registro: %w", err)
	}
	return nil
}

func (s *Service) GetRecordByDate(ctx context.Context, date string) ([]dto.Record, error) {
	filter := bson.M{"date.realDate": date}
	coll := s.DB.Collection(repository.GetdatabaseName(), repository.GetrecordCollection())
	resultList, err := repository.Find[dto.Record](ctx, coll, filter)
	if err != nil {
		fmt.Print(err)
		return []dto.Record{}, err
	}
	if len(resultList) > 0 {
		return resultList, nil
	}
	return []dto.Record{}, nil
}

func (s *Service) GetDayReport(ctx context.Context, date string) (dto.DaySummary, error) {
	filter := bson.M{"date.realDate": date}
	coll := s.DB.Collection(repository.GetdatabaseName(), repository.GetrecordCollection())
	resultList, err := repository.Find[dto.Record](ctx, coll, filter)
	if err != nil {
		fmt.Print(err)
		return dto.DaySummary{}, err
	}
	daySummary := dto.DaySummary{}
	for _, result := range resultList {
		daySummary.Date = result.Date
		daySummary.Supervivencia, _ = s.getCategorySummary(resultList, dto.Supervivencia)
		daySummary.OcioYVicio, _ = s.getCategorySummary(resultList, dto.OcioYVicio)
		daySummary.Compras, _ = s.getCategorySummary(resultList, dto.Compras)
		daySummary.Total = daySummary.OcioYVicio.Sum + daySummary.Supervivencia.Sum + daySummary.Compras.Sum
	}
	return daySummary, nil
}

func (s *Service) GetWeekReport(ctx context.Context, week, month, year int) (dto.WeekSummary, error) {
	weekDays, err := WeekDaysInCustomMonth(year, month, week, time.Local)
	if err != nil {
		fmt.Print(err)
		return dto.WeekSummary{}, err
	}
	dayMap := make(map[int]dto.DaySummary)
	firstDay := 6
	lastDay := 0
	for k, day := range weekDays {
		daySummary, err := s.GetDayReport(ctx, day.Format("2006-01-02"))
		if err != nil {
			return dto.WeekSummary{}, err
		}
		if daySummary.Total > 0 {
			dayMap[k] = daySummary
		} else {
			if day.Year() == 2999 {
				dayMap[k] = dto.DaySummary{
					Date: dto.Date{
						RealDate:      day.Format("2006-01-02"),
						ContableMonth: month,
						Day:           day.Day(),
						DayOfWeek:     time.Weekday((k + 1) % 7).String(),
						Year:          day.Year(),
						WeekOfYear:    0,
						WeekOfMonth:   0,
					},
				}
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
		}
		//get the lower key
		if k < firstDay {
			firstDay = k
		}
		//get the higher key
		if k > lastDay && day.Year() != 2999 {
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
	return weekSummary, nil
}

func (s *Service) getCategorySummary(records []dto.Record, category dto.Category) (dto.CategorySummary, error) {
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
	return categorySummary, nil
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

func buildRecordFromRequest(recordRequest dto.RecordRequest) (dto.Record, error) {
	parsed, err := time.Parse("2006-01-02", recordRequest.Date)
	if err != nil {
		log.Println("fecha inválida:", err)
		return dto.Record{}, err
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
	return newRecord, nil
}
