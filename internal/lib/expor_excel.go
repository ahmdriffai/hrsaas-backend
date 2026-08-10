package lib

import (
	"fmt"
	"hr-sas/internal/model"
	"time"

	"github.com/xuri/excelize/v2"
)

func ExportAbsensiToExcel(sheets []model.AttendanceSheet, periodeData string) (*excelize.File, error) {
	f := excelize.NewFile()

	// f.SetSheetName("Sheet1", "Hallo")
	for i, sheet := range sheets {
		sheetName := sheet.Name
		if i == 0 {
			f.SetSheetName("Sheet1", sheetName)
		} else {
			f.NewSheet(sheetName)
		}

		headStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true, Size: 14, Color: "#000000ff"},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#cfcfcfff"}, Pattern: 1},
		})
		f.SetCellStyle(sheetName, "A1", "AA", headStyle)

		// Style Judul
		titleStyle, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Size: 14, Color: "#1D1D1Dff"},
			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"#d4d4d4ff"}, Pattern: 1},
		})

		// Style Header
		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"#18254F"}, Pattern: 1},
			Border: []excelize.Border{
				{Type: "left", Color: "656565", Style: 1},
				{Type: "right", Color: "656565", Style: 1},
				{Type: "top", Color: "656565", Style: 1},
				{Type: "bottom", Color: "656565", Style: 1},
			},
		})

		// Style Default Isi
		cellStyle, _ := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
			Border: []excelize.Border{
				{Type: "left", Color: "656565", Style: 1},
				{Type: "right", Color: "656565", Style: 1},
				{Type: "top", Color: "656565", Style: 1},
				{Type: "bottom", Color: "656565", Style: 1},
			},
		})

		// Judul
		f.SetCellValue(sheetName, "B1", "REKAP KEHADIRAN HARIAN")
		f.MergeCell(sheetName, "B1", "H1")
		f.SetCellStyle(sheetName, "A1", "H1", titleStyle)

		// Detail Personalia
		f.SetCellValue(sheetName, "B3", "PT BPR BANK WONOSOBO (PERSERODA)")
		f.MergeCell(sheetName, "B3", "E3")
		// Style Default Isi
		styleTitleBawon, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Size: 16, Color: "#ffffffff"},
			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"#18254F"}, Pattern: 1},
		})

		f.SetCellStyle(sheetName, "B3", "D3", styleTitleBawon)

		// Style Default Isi
		styleTitleDetail, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Size: 10, Color: "#1D1D1Dff"},
			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"#dfe5f7ff"}, Pattern: 1},
		})

		f.SetCellValue(sheetName, "B4", "Karyawan")
		f.SetCellValue(sheetName, "C4", sheet.Name)
		f.SetCellValue(sheetName, "B5", "Periode")
		f.SetCellValue(sheetName, "C5", periodeData)
		f.SetCellValue(sheetName, "B6", "Total Hadir")
		f.SetCellValue(sheetName, "C6", sheet.Total)
		f.SetCellStyle(sheetName, "B4", "B6", styleTitleDetail)

		// Header Tabel
		headers := []string{"No", "Tanggal", "Status", "Jam Masuk", "Jam Keluar", "Terlambat", "Pulang Lebih Awal", "Catatan"}
		for col, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(col+2, 8)
			f.SetCellValue(sheetName, cell, h)
		}
		f.SetCellStyle(sheetName, "B8", "I8", headerStyle)

		no := 1
		// Data
		for rowIdx, row := range sheet.Data {
			r := rowIdx + 9
			f.SetRowHeight(sheetName, r, 20)

			f.SetCellValue(sheetName, fmt.Sprintf("B%d", r), no)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", r), formatTimeMilli(row.Date, "2006-01-02"))
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", r), row.Status)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", r), formatTimeMilli(row.CheckInTime, "15:04:05"))
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", r), formatTimeMilli(row.CheckOutTime, "15:04:05"))
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", r), row.LateCheckIn)
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", r), row.LateCheckOut)
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", r), row.Note)

			f.SetCellStyle(sheetName, fmt.Sprintf("B%d", r), fmt.Sprintf("I%d", r), cellStyle)

			// Warna kuning untuk "Bukan Hari Kerja"
			if row.Status == "Bukan Hari Kerja" {
				style, _ := f.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
					Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFF00"}, Pattern: 1},
				})
				f.SetCellStyle(sheetName, fmt.Sprintf("D%d", r), fmt.Sprintf("D%d", r), style)
			}

			if row.Status == "Belum Ada Status (TS)" {
				style, _ := f.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
					Fill:      excelize.Fill{Type: "pattern", Color: []string{"#8a8a8aff"}, Pattern: 1},
				})
				f.SetCellStyle(sheetName, fmt.Sprintf("D%d", r), fmt.Sprintf("D%d", r), style)
			}
			no++
		}

		// Atur lebar kolom
		f.SetColWidth(sheetName, "A", "A", 5)
		f.SetColWidth(sheetName, "B", "B", 10)
		f.SetColWidth(sheetName, "C", "C", 22)
		f.SetColWidth(sheetName, "D", "D", 30)
		f.SetColWidth(sheetName, "E", "F", 20)
		f.SetColWidth(sheetName, "G", "G", 20)
		f.SetColWidth(sheetName, "H", "H", 20)
		f.SetColWidth(sheetName, "I", "I", 20)
		f.SetRowHeight(sheetName, 3, 30)
		f.SetRowHeight(sheetName, 1, 30)
		f.SetRowHeight(sheetName, 8, 30)

	}

	return f, nil
}

func ExportEmployeeToExcel(employees []model.EmployeeResponse) (*excelize.File, error) {
	return nil, nil
}

func ExportVisitToExcel(visits []model.VisitResponse) (*excelize.File, error) {
	return nil, nil
}

func ExportCreditCollectionToExcel(remidialVisits []model.RemidialVisitResponse) (*excelize.File, error) {
	return nil, nil
}

func formatTimeMilli(milli int64, format string) string {
	if milli <= 0 {
		return "-"
	}
	return time.UnixMilli(milli).Format(format)
}
