package generators

import (
	"fmt"
	"github.com/opencommand/schema"
)

// Manager 管理所有代码生成器
type Manager struct {
	generators map[string]Generator
}

// NewManager 创建新的生成器管理器
func NewManager() *Manager {
	manager := &Manager{
		generators: make(map[string]Generator),
	}

	// 注册默认生成器
	manager.Register(NewGoGenerator())

	return manager
}

// Register 注册生成器
func (m *Manager) Register(generator Generator) {
	m.generators[generator.GetLanguage()] = generator
}

// GetGenerator 获取指定语言的生成器
func (m *Manager) GetGenerator(language string) (Generator, error) {
	generator, exists := m.generators[language]
	if !exists {
		return nil, fmt.Errorf("generator for language '%s' not found", language)
	}
	return generator, nil
}

// Generate 使用指定语言生成代码
func (m *Manager) Generate(schema *schema.Schema, language string) (string, error) {
	generator, err := m.GetGenerator(language)
	if err != nil {
		return "", err
	}

	return generator.Generate(schema)
}

// GetSupportedLanguages 获取支持的语言列表
func (m *Manager) GetSupportedLanguages() []string {
	languages := make([]string, 0, len(m.generators))
	for lang := range m.generators {
		languages = append(languages, lang)
	}
	return languages
}
