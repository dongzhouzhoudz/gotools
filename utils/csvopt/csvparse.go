package csvopt

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

type callBackFunction func(filePath string) (bool, error)

// 扩展回调解析方式，除了定义的哪些文件格式，还可以自定义处理格式
var extendParseFunction map[string]func(originBytes []byte, replaceBomBytes []byte) (io.Reader, error)

// 初始化扩展调用方式
func init() {
	extendParseFunction = make(map[string]func(originBytes []byte, replaceBomBytes []byte) (io.Reader, error))
	extendParseFunction["ISO-8859-1"] = func(originBytes []byte, replaceBomBytes []byte) (io.Reader, error) {
		return simplifiedchinese.GBK.NewDecoder().Reader(bytes.NewReader(replaceBomBytes)), nil
	}

	extendParseFunction["UTF-16LE"] = func(originBytes []byte, replaceBomBytes []byte) (io.Reader, error) {
		transformReader := transform.NewReader(bytes.NewReader(replaceBomBytes), unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder())
		all, errorOne := io.ReadAll(transformReader)
		if errorOne != nil {
			return nil, errorOne
		} else {
			return bytes.NewReader([]byte(strings.ReplaceAll(string(all), "=", ""))), nil
		}
	}

	extendParseFunction["UTF-8"] = func(originBytes []byte, replaceBomBytes []byte) (io.Reader, error) {
		return bytes.NewReader([]byte(strings.ReplaceAll(string(replaceBomBytes), "=", ""))), nil
	}

}

// ParseCsvFileForOneRecord 解析单条文件数据,支持callBack方法
func ParseCsvFileForOneRecord(filePath string, resultStruct interface{}, function callBackFunction) (bool, error) {
	exits, err := checkFileExitsAndExt(filePath, ".csv")
	if !exits {
		panic(err)
	}
	csvData, dealFileError := dealWithDifferentEncodingCsvFile(filePath)
	if dealFileError != nil {
		panic(dealFileError)
	}

	if len(csvData) > 0 {

		typeOfChildData := reflect.TypeOf(resultStruct).Elem()
		valueOfChildData := reflect.ValueOf(resultStruct).Elem()
		structName := typeOfChildData.Name()

		var csvTitleToField = map[string]int{}
		var csvFiledToPosition = map[int]int{}
		for i := 0; i < typeOfChildData.NumField(); i++ {
			tag := typeOfChildData.Field(i).Tag
			structFiledName := typeOfChildData.Field(i).Name
			dataKey := tag.Get("title")

			split := strings.Split(dataKey, ",")

			if len(split) == 1 {
				csvTitleToField[split[0]] = i
			} else if len(split) >= 2 {
				for _, splitValue := range split {
					csvTitleToField[splitValue] = i
				}
			} else {
				panic(fmt.Sprintf("%s 结构体标签 title 设置不合法,当前设置为 %s", structName, dataKey))
			}
			dataNumber := tag.Get("column")
			csvPositionId, convertError := strconv.Atoi(dataNumber)
			if convertError != nil {
				panic(fmt.Sprintf("%s 结构体元素 %s 的标签 column 不能设置成非数字类型,当前设置为 %s", structName, structFiledName, dataNumber))
			} else {
				csvFiledToPosition[i] = csvPositionId
			}
		}

		if len(csvData) == 1 {
			// 单行数据结构 需要通过列号匹配
			for key, value := range csvFiledToPosition {
				valueOfChildData.Field(key).Set(reflect.ValueOf(csvData[0][value]))
			}
			return true, nil
		} else if len(csvData) >= 2 {
			// 多行数据结构 需要通过title 匹配字段
			for titlePosition := 0; titlePosition < len(csvData[0]); titlePosition++ {
				titleName := csvData[0][titlePosition]
				if _, ok := csvTitleToField[titleName]; ok {
					valueOfChildData.Field(csvTitleToField[titleName]).Set(reflect.ValueOf(csvData[1][titlePosition]))
				}
			}

			return true, nil
		} else {
			return false, errors.New("文件内容行数有误，实际行数超出解析范围")
		}
	} else {
		return false, errors.New("文件内容有误当前解析结果为0行，缺少必要的数据无法解析")
	}

}

// SetOtherFileEncodingParseFunction 对外暴露不同文件类型解析方法
func SetOtherFileEncodingParseFunction(encodingType string, newFunction func(originBytes []byte, replaceBomBytes []byte) (io.Reader, error)) {
	extendParseFunction[encodingType] = newFunction
}

// dealWithDifferentEncodingCsvFile 处理不同格式csv 文件解析
func dealWithDifferentEncodingCsvFile(filePath string) (csvData [][]string, error error) {
	file, openFileError := os.Open(filePath)
	originBytes, _ := io.ReadAll(file)
	fileData, et := Skip(bytes.NewReader(originBytes))
	fmt.Println("CSV-BOM 文件检测--当前文件数据类型为--", et)
	if openFileError != nil {
		panic(openFileError)
	}

	fileBytes, err := io.ReadAll(fileData)
	if err != nil {
		return nil, err
	}

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(originBytes)
	if err != nil {
		return nil, err
	}
	fmt.Println("CSV-ENCODING 文件编码格式检测--当前文件编码格式为--", result.Charset)
	if _, ok := extendParseFunction[result.Charset]; ok {
		parseReader, parseError := extendParseFunction[result.Charset](originBytes, fileBytes)
		if parseError != nil {
			return nil, parseError
		} else {
			csvReader := csv.NewReader(parseReader)
			csvReader.LazyQuotes = true
			csvData, readError := csvReader.ReadAll()
			if readError != nil {
				return nil, readError
			} else {
				return csvData, nil
			}

		}
	} else {
		return nil, errors.New(fmt.Sprintf("当前文件编码格式为%s，并没有对应的解析方式，请调用csvopt.SetOtherFileEncodingParseFunction 配置相关解析方案", result.Charset))
	}

}

// checkFileExitsAndExt 判断文件是否存在 并且判断是否是合法后缀
func checkFileExitsAndExt(filePath string, extString string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		ext := path.Ext(filePath)
		if ext != extString {
			return false, errors.New(fmt.Sprintf("文件%s格式不合法,要求文件后缀为%s", filePath, extString))
		} else {
			return true, nil
		}
	}
	if os.IsNotExist(err) {
		return false, errors.New(fmt.Sprintf("%s 文件不存在", filePath))
	}
	return false, err
}
