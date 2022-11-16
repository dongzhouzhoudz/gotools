### 使用说明
##### 使用范例
```go
 type OutPutData struct {
NameTitle    string `type:"title" title_column:"A1-C3" title_value:"姓名"`
BarcodeTitle string `type:"title" title_column:"D1-F2" title_value:"代码说明"`
OriginTitle  string `type:"title" title_column:"D3" title_value:"源代码"`
AuthorTitle  string `type:"title" title_column:"E3-F3" title_value:"作者"`
SexTitle     string `type:"title" title_column:"G1-G3" title_value:"测试重复"`
NameBody     string `type:"data"  data_column:"A-C"`
BarcodeBody  string `type:"data"   data_column:"D-F"`
SexBody      string `type:"data"   data_column:"G"`
}

type OutPutTwo struct {
NameTitle    string `type:"title" title_column:"H1-H3" title_value:"新增加列1"`
BarcodeTitle string `type:"title" title_column:"I1-I3" title_value:"新增加列2"`
SexTitle     string `type:"title" title_column:"J1-K3" title_value:"新增加列3"`
NameBody     string `type:"data"  data_column:"H"`
BarcodeBody  string `type:"data"   data_column:"I"`
SexBody      string `type:"data"   data_column:"J-K"`
}

func OutPutExcelDemo() {
var appendData []interface{}
for i := 0; i < 1000; i++ {
dataOne := OutPutData{
NameBody:    "One",
BarcodeBody: "11111111",
SexBody:     "1",
}

appendData = append(appendData, dataOne)

dataTwo := OutPutData{
NameBody:    "Two",
BarcodeBody: "22222222",
SexBody:     "1",
}

appendData = append(appendData, dataTwo)

}

for i := 0; i < 1000; i++ {
dataThree := OutPutTwo{
NameBody:    "Three",
BarcodeBody: "10011",
SexBody:     "22",
}
appendData = append(appendData, dataThree)
}

style := `{
      "alignment":{"horizontal":"center","Vertical":"center"},
      "font":{"bold":true,"color":"#FFFFFF","size":13},
      "border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"right","color":"000000","style":1},{"type":"bottom","color":"000000","style":1}],
      "fill":{"type":"pattern","color":["#4472C4"],"pattern":1,"shading":3}
       }`
sheetStyle := excel.ProduceSheetStyle("TestSheet", "A1", "j3", style)
//先配置头部信息
excel.SaveExcelFromStructListWithStyle("~/Download/test.xlsx", "TestSheet", appendData, 4, sheetStyle)
//再次配置头部样式，保留当前已有样式

}

```


### 导出结构体定义

```go
 type OutPutData struct {
NameTitle    string `type:"title" title_column:"A1-C3" title_value:"姓名"`
BarcodeTitle string `type:"title" title_column:"D1-F2" title_value:"代码说明"`
OriginTitle  string `type:"title" title_column:"D3" title_value:"源代码"`
AuthorTitle  string `type:"title" title_column:"E3-F3" title_value:"作者"`
SexTitle     string `type:"title" title_column:"G1-G3" title_value:"测试重复"`
NameBody     string `type:"data"  data_column:"A-C"`
BarcodeBody  string `type:"data"   data_column:"D-F"`
SexBody      string `type:"data"   data_column:"G"`
}
```
##### 结构体标签说明
###### 表头标签配置说明
type 导出这个结构体中该字段类型是 title 表头 <br>
title_column 该表头的区域 A1-C3 表示这个表头占用 A1代表矩形左上角点 C3 代表右下角点，该矩形内单元格会合并，如果是单个单元格可以配置为 D3 这种单格形式<br>
title_value  该表头显示的具体值

##### 表中具体值配置说明
type  导出这个结构体中该字段类型是 data 值<br>
data_column 该值需要填充的单元格范围，仅支持单行数据 多列合并，A-C 表示 某个值需要再某行合并A-C 单元格显示，如果无需合并单元格 可以配置为 D 这种单格形式 <br>



### 设置单元格格式

```go
style := `{
      "alignment":{"horizontal":"center","Vertical":"center"},
      "font":{"bold":true,"color":"#FFFFFF","size":13},
      "border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"right","color":"000000","style":1},{"type":"bottom","color":"000000","style":1}],
      "fill":{"type":"pattern","color":["#4472C4"],"pattern":1,"shading":3}
       }`
// TestSheet sheet 名称
// "A1" 起始单元格
// "j3" 结束单元格
// style 需要适配的样式
sheetStyle := excel.ProduceSheetStyle("TestSheet", "A1", "j3", style)
//先配置头部信息
excel.SaveExcelFromStructListWithStyle("~/Download/test.xlsx", "TestSheet", appendData, 4, sheetStyle)
//再次配置头部样式，保留当前已有样式

```



