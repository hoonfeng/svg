package io

import (
	"os"

	"github.com/hoonfeng/svg/types"
)

// SaveSVG 将SVG文档保存为文件
func SaveSVG(doc *types.Document, filename string) error {
	// 将文档转换为XML字符串
	xml := doc.ToXML()

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入文件
	_, err = file.WriteString(xml)
	if err != nil {
		return err
	}

	return nil
}
