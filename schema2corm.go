package main

import (
	"fmt"
	"os"
	"schema2corm/generators"
	"github.com/opencommand/schema"

	"gopkg.in/yaml.v3"
)

func main() {
	// 读取schema文件
	data, err := os.ReadFile("docker.schema.yaml")
	if err != nil {
		fmt.Printf("Error reading schema file: %v\n", err)
		os.Exit(1)
	}

	// 解析schema
	var schemaData schema.Schema
	err = yaml.Unmarshal(data, &schemaData)
	if err != nil {
		fmt.Printf("Error parsing schema: %v\n", err)
		os.Exit(1)
	}

	// 创建生成器管理器
	manager := generators.NewManager()

	// 生成Go代码
	goCode, err := manager.Generate(&schemaData, "go")
	if err != nil {
		fmt.Printf("Error generating Go code: %v\n", err)
		os.Exit(1)
	}

	// 输出生成的代码
	fmt.Println("Generated Go code:")
	fmt.Println("==================")
	fmt.Println(goCode)

	// 显示支持的语言
	fmt.Println("\nSupported languages:")
	for _, lang := range manager.GetSupportedLanguages() {
		fmt.Printf("- %s\n", lang)
	}
}
