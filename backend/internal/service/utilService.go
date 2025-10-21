package service

import (
	"errors"
	"fmt"
	"time"
)

// WeekNumberInCustomMonth devuelve el número de semana (1..N) dentro del “mes contable”
// cuyo rango es [24 del mes anterior, 23 del mes actual]. Las semanas son ISO (lunes–domingo).
// También devuelve el inicio y fin (acotados a la ventana 24–23).
func WeekNumberInCustomMonth(dateStr string, loc *time.Location) (week, month int, spanStart, spanEnd time.Time, err error) {
	if loc == nil {
		loc = time.Local
	}

	// 1) Parseo de fecha
	t, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return 0, 0, time.Time{}, time.Time{}, err
	}
	t = atMidnight(t, loc)

	// 2) Determinar “mes contable” al que pertenece t
	cYear, cMonth := customMonthOf(t)

	// 3) Ventana del mes contable [24/mes-1, 23/mes]
	winStart, winEnd := customMonthWindow(cYear, cMonth, loc)

	// Validación básica
	if t.Before(winStart) || t.After(winEnd) {
		return 0, 0, time.Time{}, time.Time{}, errors.New("la fecha no cae dentro de la ventana del mes contable")
	}

	// 4) Lunes–domingo de la semana ISO que contiene t
	weekStart := mondayOfWeek(t, loc)
	weekEnd := weekStart.AddDate(0, 0, 6)

	// 5) Lunes de la semana 1 del mes contable (la semana que contiene winStart)
	firstWeekStart := mondayOfWeek(winStart, loc)

	// 6) Nº de semana = 1 + semanas completas desde firstWeekStart hasta weekStart
	days := int(weekStart.Sub(firstWeekStart).Hours() / 24)
	week = (days/7 + 1)

	// 7) Acotar el span al interior de la ventana 24–23
	spanStart = maxTime(weekStart, winStart)
	spanEnd = minTime(weekEnd, winEnd)

	month = int(t.Month())
	if t.Day() >= 24 {
		month = month + 1
	}
	if month > 12 {
		month = 1
	}
	return week, month, spanStart, spanEnd, nil
}

// --- Helpers ---

func customMonthOf(t time.Time) (year int, month time.Month) {
	// Si el día es >= 24, pertenece al mes contable del mes+1 (con cambio de año si toca).
	if t.Day() >= 24 {
		n := t.AddDate(0, 1, 0)
		return n.Year(), n.Month()
	}
	// Si es 1..23, pertenece al mes contable del mes actual.
	return t.Year(), t.Month()
}

func customMonthWindow(year int, month time.Month, loc *time.Location) (start, end time.Time) {
	// Ventana: [24 del mes-1, 23 del mes]
	prev := time.Date(year, month, 1, 0, 0, 0, 0, loc).AddDate(0, -1, 0)
	start = time.Date(prev.Year(), prev.Month(), 24, 0, 0, 0, 0, loc)
	end = time.Date(year, month, 23, 23, 59, 59, 999999999, loc)
	return
}

func mondayOfWeek(t time.Time, loc *time.Location) time.Time {
	t = atMidnight(t, loc)
	// Go: Sunday=0, Monday=1, ... Saturday=6
	wd := int(t.Weekday())
	// Convertimos a índice con Monday=0
	offset := (wd + 6) % 7
	return t.AddDate(0, 0, -offset)
}

func atMidnight(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

// WeekDaysInCustomMonth devuelve los días (nombre→fecha) de la semana 'week'
// del mes contable 'month' (1..12) en 'year'.
// El mes contable va del 24 del mes anterior al 23 del mes indicado.
// Las semanas son ISO (lunes→domingo) y se intersectan con la ventana 24–23.
// Ej.: (2025, 11, 1) -> {"Viernes":"2025-10-24","Sábado":"2025-10-25","Domingo":"2025-10-26"}
func WeekDaysInCustomMonth(year, month, week int, loc *time.Location) (map[int]time.Time, error) {
	if loc == nil {
		loc = time.Local
	}
	if month < 1 || month > 12 || week < 1 {
		return nil, fmt.Errorf("parámetros inválidos (month=1..12, week>=1)")
	}

	// Ventana del mes contable: [24 del mes-1, 23 del mes]
	winStart, winEnd := customMonthWindow(year, time.Month(month), loc)

	// Lunes de la semana 1 (la que contiene el 24 del mes-1)
	firstWeekStart := mondayOfWeek(winStart, loc)

	// Lunes y domingo de la semana solicitada
	weekStart := firstWeekStart.AddDate(0, 0, (week-1)*7)
	weekEnd := weekStart.AddDate(0, 0, 6)

	// Si la semana está completamente antes/después de la ventana, no existe
	if weekStart.After(winEnd) || weekEnd.Before(winStart) {
		return nil, fmt.Errorf("la semana %d no existe en el mes contable %02d/%d", week, month, year)
	}

	// Intersección real con la ventana
	start := maxTime(weekStart, winStart)
	end := minTime(weekEnd, winEnd)

	esWeekday := func(w time.Weekday) int {
		switch w {
		case time.Monday:
			return 0
		case time.Tuesday:
			return 1
		case time.Wednesday:
			return 2
		case time.Thursday:
			return 3
		case time.Friday:
			return 4
		case time.Saturday:
			return 5
		default:
			return 6
		}
	}

	out := make(map[int]time.Time)
	for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {
		key := esWeekday(day.Weekday())
		out[key] = day
	}
	return out, nil
}
