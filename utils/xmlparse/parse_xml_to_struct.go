package xmlparse

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"reflect"
	"strings"
)

var ElementCallFunction = map[string]func(parentElement interface{}) (interface{}, error){}
var StructFiledMap = map[string]map[string]int{}

// SetElementCallFunction 设置不同节点处理数据的方法
func SetElementCallFunction(elementNodeName string, dealWithFunction func(parentElement interface{}) (interface{}, error)) {
	ElementCallFunction[elementNodeName] = dealWithFunction
}

// CallElementFunction 处理不同节点的特殊方法
func callElementFunction(elementNodeName string, childData interface{}) (interface{}, error) {
	if _, ok := ElementCallFunction[elementNodeName]; ok {
		return ElementCallFunction[elementNodeName](childData)
	} else {
		return nil, errors.New(fmt.Sprintf("XML节点<%s> 没有找到对应的处理方法,当前系统不做任何处理", elementNodeName))
	}
}

// GetEachChildElement 获取不同子节点对应的数据
func GetEachChildElement(elements []*etree.Element, childData interface{}) {
	if len(elements) > 0 {
		for _, element := range elements {
			if len(element.ChildElements()) > 0 {
				method, err := callElementFunction(element.Tag, childData)
				if err != nil {
					GetEachChildElement(element.ChildElements(), childData)
				} else {
					GetEachChildElement(element.ChildElements(), method)
				}
			} else { //如果是当个节点 制动获取相关数据
				tag := element.Tag
				text := element.Text()
				typeOfChildData := reflect.TypeOf(childData).Elem()
				valueOfChildData := reflect.ValueOf(childData).Elem()
				structName := typeOfChildData.Name()
				if _, ok := StructFiledMap[structName]; !ok {
					oneFiledMap := map[string]int{}
					for i := 0; i < typeOfChildData.NumField(); i++ {
						get := typeOfChildData.Field(i).Tag.Get("key")
						split := strings.Split(get, ",")
						if len(split) > 0 {
							for _, oneData := range split {
								oneFiledMap[oneData] = i
							}
						}
					}
					StructFiledMap[structName] = oneFiledMap
				}
				if _, ok := StructFiledMap[structName]; ok {
					if _, valueOk := StructFiledMap[structName][tag]; valueOk {
						i := StructFiledMap[structName][tag]
						valueOfChildData.Field(i).Set(reflect.ValueOf(text))
					}
				}
			}
		}
	} else {
		fmt.Println("解析过程中发现该数据节点没有任何子节点，可以考虑跳过")
	}
}

// ParseXmlToStruct 将XML 解析到结构体中
func ParseXmlToStruct(filePath string, topStructPtr interface{}) interface{} {
	doc := etree.NewDocument()
	err := doc.ReadFromFile(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		elements := doc.ChildElements()
		GetEachChildElement(elements, topStructPtr)
	}
	return topStructPtr
}
