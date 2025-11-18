package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"path/filepath"
)

//go:embed templates/*.html
var templatesFS embed.FS

type TemplateData struct {
	Code string
}

type TemplateManager struct {
	templates map[string]*template.Template
}

func NewTemplateManager() (*TemplateManager, error) {
	tm := &TemplateManager{
		templates: make(map[string]*template.Template),
	}

	templateFiles := []string{
		"templates/verification_code.html",
		"templates/password_reset.html",
	}

	for _, tmplFile := range templateFiles {
		tmpl, err := template.ParseFS(templatesFS, tmplFile)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", tmplFile, err)
		}

		templateName := filepath.Base(tmplFile)
		tm.templates[templateName] = tmpl
	}

	return tm, nil
}

func (tm *TemplateManager) Render(templateName string, data TemplateData) (string, error) {
	tmpl, exists := tm.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

func (tm *TemplateManager) GetPlainText(code string) string {
	return fmt.Sprintf(`
Здравствуйте!

Ваш код: %s

Этот код действителен в течение 15 минут.
Пожалуйста, не сообщайте этот код никому.

С уважением,
Команда Imperial
`, code)
}
