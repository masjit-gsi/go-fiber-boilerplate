package format

import (
	"strings"
	"time"

	"github.com/bojanz/currency"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

const (
	// parsedTime1, err := time.Parse(time.RFC3339, tanggalJadwal)
	// parsedTime2, err := time.Parse(time.RFC3339, tanggalLog)

	// // Menghilangkan informasi zona waktu
	// timeWithoutZone := parsedTime1.UTC()
	// timeWithoutZone2 := parsedTime2.UTC()

	// // Memformat waktu tanpa zona waktu ke dalam format yang diinginkan
	// formattedTime1 := timeWithoutZone.Format("2006-01-02 15:04:05")
	// formattedTime2 := timeWithoutZone2.Format("2006-01-02 15:04:05")
	// DefaultDateFormat digunakan untuk standar format input dari string ke format tanggal
	DateDMYFormat = "02/01/2006"

	// DefaultDateFormat digunakan untuk standar format input dari string ke format tanggal
	DefaultDateFormat = "2006-01-02"

	// DefaultDateFormat digunakan untuk standar format input dari string ke format tanggal
	DefaultDateYearFormat = "2006-01"

	// DefaultDateFormat digunakan untuk standar format input dari string ke format tanggal dan jam
	DefaultDateTimeFormat = "2006-01-02 15:04:05"

	// DefaultTimeFormat digunakan untuk standar format input dari string ke format jam
	DefaultTimeFormat = "15:04:05"
)

// DecimalToRupiah parse decimal to rupiah format
func DecimalToRupiah(value decimal.Decimal) string {
	ac := accounting.Accounting{Symbol: "Rp. ", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(value)
}

// RupiahToDecimal parse rupiah format to decimal
func RupiahToDecimal(rupiah string) (value decimal.Decimal, err error) {
	locale := currency.NewLocale("id")
	formatter := currency.NewFormatter(locale)
	amount, err := formatter.Parse(rupiah, "IDR")
	if err != nil {
		return
	}
	value, err = decimal.NewFromString(amount.Number())
	if err != nil {
		return
	}
	return
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func ParseSplitString(val string) string {
	var ids string
	resSplit := strings.Split(val, ",")
	lengId := len(resSplit)
	for i, s := range resSplit {
		join := "'" + strings.TrimSpace(s) + "'"
		if i == lengId-1 {
			ids += join
		} else {
			ids += join + ","
		}
	}

	return ids
}

func StringJoin(arr []string) string {
	var ids string
	lengId := len(arr)
	for i, s := range arr {
		join := "'" + strings.TrimSpace(s) + "'"
		if i == lengId-1 {
			ids += join
		} else {
			ids += join + ","
		}
	}

	return ids
}

func SplitString(val string, demiliter string) []string {
	resSplit := strings.Split(val, demiliter)
	return resSplit
}

func ParseString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func ParseInt(s *int) int {
	if s != nil {
		return *s
	}

	return 0
}

func ParseFloat64(s *float64) float64 {
	if s != nil {
		return *s
	}

	return 0
}

func NullString(val string) *string {
	if val != "" {
		return &val
	}

	return nil
}

func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

func ParseBool(s *bool) bool {
	if s != nil {
		return *s
	}

	return false
}
