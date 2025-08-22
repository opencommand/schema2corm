package generators

import (
	"fmt"
	"github.com/opencommand/schema"
	"strings"
)

// GoGenerator 生成Go语言的corm代码
type GoGenerator struct {
	BaseGenerator
}

// NewGoGenerator 创建新的Go生成器
func NewGoGenerator() *GoGenerator {
	return &GoGenerator{
		BaseGenerator: BaseGenerator{
			language:      "go",
			fileExtension: ".go",
		},
	}
}

// Generate 生成Go代码
func (g *GoGenerator) Generate(schema *schema.Schema) (string, error) {
	var code strings.Builder

	// 生成包声明
	code.WriteString("package main\n\n")

	// 生成导入
	code.WriteString("import (\n")
	code.WriteString("\t\"fmt\"\n")
	code.WriteString("\t\"strings\"\n")
	code.WriteString("\t\"corm\"\n")
	code.WriteString(")\n\n")

	// 生成主命令结构体
	mainCmdName := g.toPascalCase(schema.Name) + "Command"
	code.WriteString(fmt.Sprintf("type %s struct {\n", mainCmdName))

	// 生成选项字段
	for optName, opt := range schema.Options {
		fieldName := "opt_" + optName
		fieldType := g.getGoType(opt.Type)
		code.WriteString(fmt.Sprintf("\t%s %s\n", fieldName, fieldType))
	}

	code.WriteString("\tcorm.BaseCommand\n")
	code.WriteString("}\n\n")

	// 生成主命令构造函数
	code.WriteString(fmt.Sprintf("func %s() *%s {\n", g.toPascalCase(schema.Name), mainCmdName))
	code.WriteString(fmt.Sprintf("\treturn &%s{}\n", mainCmdName))
	code.WriteString("}\n\n")

	// 生成Name方法
	code.WriteString(fmt.Sprintf("func (c *%s) Name() string {\n", mainCmdName))
	code.WriteString(fmt.Sprintf("\treturn \"%s\"\n", schema.Name))
	code.WriteString("}\n\n")

	// 生成选项设置方法
	for optName, opt := range schema.Options {
		methodName := g.toPascalCase(optName)
		paramType := g.getGoType(opt.Type)
		fieldName := "opt_" + optName

		code.WriteString(fmt.Sprintf("func (c *%s) %s(%s %s) *%s {\n",
			mainCmdName, methodName, optName, paramType, mainCmdName))
		code.WriteString(fmt.Sprintf("\tc.%s = %s\n", fieldName, optName))
		code.WriteString(fmt.Sprintf("\treturn c\n"))
		code.WriteString("}\n\n")
	}

	// 生成String方法
	code.WriteString(fmt.Sprintf("func (c *%s) String() string {\n", mainCmdName))
	code.WriteString("\targs := []string{}\n")

	for optName, opt := range schema.Options {
		fieldName := "opt_" + optName
		if opt.Type == "bool" {
			code.WriteString(fmt.Sprintf("\tif c.%s {\n", fieldName))
			code.WriteString(fmt.Sprintf("\t\targs = append(args, \"-%s\")\n", optName))
			code.WriteString("\t}\n")
		} else {
			code.WriteString(fmt.Sprintf("\tif c.%s != \"\" {\n", fieldName))
			code.WriteString(fmt.Sprintf("\t\targs = append(args, fmt.Sprintf(\"-%s=%s\", c.%s))\n", optName, "%s", fieldName))
			code.WriteString("\t}\n")
		}
	}

	code.WriteString(fmt.Sprintf("\treturn fmt.Sprintf(\"%s %%s\", strings.Join(args, \" \"))\n", schema.Name))
	code.WriteString("}\n\n")

	// 生成子命令
	for cmdName, cmd := range schema.SubCommands {
		subCmdName := g.toPascalCase(schema.Name) + g.toPascalCase(cmdName) + "Command"

		// 子命令结构体
		code.WriteString(fmt.Sprintf("type %s struct {\n", subCmdName))

		// 子命令选项字段
		for optName, opt := range cmd.Options {
			fieldName := "opt_" + optName
			fieldType := g.getGoType(opt.Type)
			code.WriteString(fmt.Sprintf("\t%s %s\n", fieldName, fieldType))
		}

		code.WriteString("\tcorm.BaseCommand\n")
		code.WriteString("}\n\n")

		// 子命令构造函数
		code.WriteString(fmt.Sprintf("func (c *%s) %s() *%s {\n",
			mainCmdName, g.toPascalCase(cmdName), subCmdName))
		code.WriteString(fmt.Sprintf("\treturn &%s{}\n", subCmdName))
		code.WriteString("}\n\n")

		// 子命令Name方法
		code.WriteString(fmt.Sprintf("func (c *%s) Name() string {\n", subCmdName))
		code.WriteString(fmt.Sprintf("\treturn \"%s %s\"\n", schema.Name, cmdName))
		code.WriteString("}\n\n")

		// 子命令选项设置方法
		for optName, opt := range cmd.Options {
			methodName := g.toPascalCase(optName)
			paramType := g.getGoType(opt.Type)
			fieldName := "opt_" + optName

			code.WriteString(fmt.Sprintf("func (c *%s) %s(%s %s) *%s {\n",
				subCmdName, methodName, optName, paramType, subCmdName))
			code.WriteString(fmt.Sprintf("\tc.%s = %s\n", fieldName, optName))
			code.WriteString(fmt.Sprintf("\treturn c\n"))
			code.WriteString("}\n\n")
		}

		// 子命令String方法
		code.WriteString(fmt.Sprintf("func (c *%s) String() string {\n", subCmdName))
		code.WriteString("\targs := []string{}\n")

		for optName, opt := range cmd.Options {
			fieldName := "opt_" + optName
			if opt.Type == "bool" {
				code.WriteString(fmt.Sprintf("\tif c.%s {\n", fieldName))
				code.WriteString(fmt.Sprintf("\t\targs = append(args, \"-%s\")\n", optName))
				code.WriteString("\t}\n")
			} else {
				code.WriteString(fmt.Sprintf("\tif c.%s != \"\" {\n", fieldName))
				code.WriteString(fmt.Sprintf("\t\targs = append(args, fmt.Sprintf(\"-%s=%s\", c.%s))\n", optName, "%s", fieldName))
				code.WriteString("\t}\n")
			}
		}

		code.WriteString(fmt.Sprintf("\treturn fmt.Sprintf(\"%s %s %%s\", strings.Join(args, \" \"))\n", schema.Name, cmdName))
		code.WriteString("}\n\n")
	}

	return code.String(), nil
}

// toPascalCase 转换为PascalCase
func (g *GoGenerator) toPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// getGoType 获取Go类型
func (g *GoGenerator) getGoType(schemaType string) string {
	switch schemaType {
	case "bool":
		return "bool"
	case "str", "string":
		return "string"
	case "int":
		return "int"
	case "float":
		return "float64"
	default:
		return "string"
	}
}
