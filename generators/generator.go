package generators

import "github.com/opencommand/schema"

// Generator 定义了代码生成器的接口
type Generator interface {
	// Generate 根据schema生成代码
	Generate(schema *schema.Schema) (string, error)

	// GetLanguage 返回生成器支持的语言
	GetLanguage() string

	// GetFileExtension 返回生成文件的后缀
	GetFileExtension() string
}

// BaseGenerator 提供基础生成器功能
type BaseGenerator struct {
	language      string
	fileExtension string
}

func (b *BaseGenerator) GetLanguage() string {
	return b.language
}

func (b *BaseGenerator) GetFileExtension() string {
	return b.fileExtension
}
