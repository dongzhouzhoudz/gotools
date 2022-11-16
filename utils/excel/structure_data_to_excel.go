package excel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

// SaveExcelFromStructListWithStyle 将结构化数据写入到excel 支持多种类型对象，多类型
func SaveExcelFromStructListWithStyle(saveFilePath string, sheetName string, saveData []interface{}, valueStartLineCount int32, styleList ...SheetStyle) {
	// 打开文件句柄
	file := excelize.NewFile()
	// 配置sheet 展示样式
	setSheetStyle(file, styleList)
	// 将数据写入到excel 中
	loadStructToFile(file, sheetName, saveData, valueStartLineCount)
	file.SaveAs(saveFilePath)
}

// SaveExcelFromStructList 将结构化数据写入到excel 支持多种类型对象，多类型
func SaveExcelFromStructList(saveFilePath string, sheetName string, saveData []interface{}, valueStartLineCount int32) {
	file := excelize.NewFile()
	loadStructToFile(file, sheetName, saveData, valueStartLineCount)
	file.SaveAs(saveFilePath)
}

// loadStructToFile 将每一行数据写入到excel 中并且支持多类型重复定义
func loadStructToFile(f *excelize.File, sheetName string, loadData []interface{}, valueStartLine int32) *excelize.File {
	var sheetTitleMap = new(sync.Map)
	activeSheet := f.NewSheet(sheetName)
	if len(loadData) > 1 {
		filePointer, _ := setSheetTitle(f, sheetName, loadData[0])
		writeCellLineCount := valueStartLine
		for _, value := range loadData {
			titleName := reflect.TypeOf(value).Name()
			_, ok := sheetTitleMap.Load(titleName)
			if !ok {
				setSheetTitle(f, sheetName, value)
				sheetTitleMap.Store(titleName, "store")
			}
			setSheetBodyData(filePointer, sheetName, value, &writeCellLineCount)
		}
	} else {
		panic("导出的数据结构不能为空，至少要保证数据结构长度>1")
	}
	f.SetActiveSheet(activeSheet)
	return f
}

// setSheetBodyData 将结构体表表格中的内容 写入到excel中
func setSheetBodyData(f *excelize.File, sheetName string, structInfo interface{}, lineNo *int32) {
	refType := reflect.TypeOf(structInfo)
	refValue := reflect.ValueOf(structInfo)
	fieldCount := refType.NumField()
	for i := 0; i < fieldCount; i++ {
		fileTmp := refType.Field(i)
		columnType := fileTmp.Tag.Get("type")
		columnTitleRange := fileTmp.Tag.Get("data_column")
		field := refValue.Field(i)
		//只有类型是title 部分才能生成excel 表头
		if columnType == "data" {
			if strings.Contains(columnTitleRange, "-") {
				if splitArray := strings.Split(columnTitleRange, "-"); len(splitArray) == 2 {
					startColumn := fmt.Sprintf("%s%d", splitArray[0], *lineNo)
					endColumn := fmt.Sprintf("%s%d", splitArray[1], *lineNo)
					f.MergeCell(sheetName, startColumn, endColumn)
					f.SetCellValue(sheetName, startColumn, field)
				} else {
					panic(fmt.Sprintf("%s结构中的%s字段标签配置有误，|配置样例单个单元格 A | 合并单元格 A-C|", refType.Name(), fileTmp.Name))
				}
			} else {
				cellAddr := fmt.Sprintf("%s%d", columnTitleRange, *lineNo)
				f.SetCellValue(sheetName, cellAddr, field)
			}
		}
	}
	// 写入cell 自增操作
	*lineNo = atomic.AddInt32(lineNo, 1)
}

// setSheetTitle 获取某个对象的参数名称，对应的数据展示出来 //
func setSheetTitle(f *excelize.File, sheetName string, structInfo interface{}) (*excelize.File, error) {
	refType := reflect.TypeOf(structInfo)
	fieldCount := refType.NumField()
	for i := 0; i < fieldCount; i++ {
		fileTmp := refType.Field(i)
		columnType := fileTmp.Tag.Get("type")
		columnTitleRange := fileTmp.Tag.Get("title_column")
		columnTitleValue := fileTmp.Tag.Get("title_value")
		//只有类型是title 部分才能生成excel 表头
		if columnType == "title" {
			if strings.Contains(columnTitleRange, "-") {
				if splitArray := strings.Split(columnTitleRange, "-"); len(splitArray) == 2 {
					f.MergeCell(sheetName, splitArray[0], splitArray[1])
					f.SetCellValue(sheetName, splitArray[0], columnTitleValue)
				} else {
					panic(fmt.Sprintf("%s结构中的%s字段标签配置有误，|配置样例单个单元格 A1 | 合并单元格 A1-C2 |", refType.Name(), fileTmp.Name))
				}
			} else {
				f.SetCellValue(sheetName, columnTitleRange, columnTitleValue)
			}
		}
	}
	return f, nil
}

type SheetStyle struct {
	SheetName   string
	HCellName   string
	VCellName   string
	StyleString string
}

// ProduceSheetStyle 设置excel样式
func ProduceSheetStyle(sheetName string, hCellName string, vCellName string, style string) SheetStyle {
	return SheetStyle{
		sheetName,
		hCellName,
		vCellName,
		style,
	}
}

// SetSheetStyle 设置excel 单元格格式 //
func setSheetStyle(f *excelize.File, styleList []SheetStyle) bool {
	if len(styleList) >= 1 {
		for _, oneStyle := range styleList {
			activeSheet := f.NewSheet(oneStyle.SheetName)
			f.SetActiveSheet(activeSheet)
			style, err := f.NewStyle(oneStyle.StyleString)
			if err != nil {
				panic(err)
			} else {
				f.SetCellStyle(oneStyle.SheetName, oneStyle.HCellName, oneStyle.VCellName, style)
				return true
			}

		}

	} else {
		panic("配置单元格样式不能为空")
	}
	return true
}
