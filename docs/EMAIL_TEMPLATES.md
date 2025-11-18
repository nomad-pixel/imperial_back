# HTML Email Templates

## Обзор

Система отправки email использует красивые HTML шаблоны для всех писем. Шаблоны встроены в приложение через `embed.FS` и автоматически загружаются при инициализации.

## Структура

```
internal/infrastructure/email/
├── templates/
│   ├── verification_code.html    # Шаблон для кода верификации
│   └── password_reset.html       # Шаблон для сброса пароля
├── template_manager.go           # Менеджер шаблонов
├── smtp_email_service.go        # SMTP сервис с HTML поддержкой
└── console_email_service.go     # Консольный сервис (для разработки)
```

## Как это работает

### 1. Template Manager

`TemplateManager` загружает все HTML шаблоны из директории `templates/` при инициализации:

```go
tm, err := NewTemplateManager()
// Загружает все .html файлы из templates/
```

### 2. Рендеринг шаблонов

Шаблоны используют стандартный Go `html/template`:

```go
htmlBody, err := tm.Render("verification_code.html", TemplateData{Code: "123456"})
```

### 3. Multipart Email

SMTP сервис отправляет email в формате `multipart/alternative`:
- **text/plain** - простая текстовая версия (fallback)
- **text/html** - красивая HTML версия

Клиенты email автоматически выберут подходящую версию.

## Шаблоны

### Verification Code Template

**Файл:** `templates/verification_code.html`

**Переменные:**
- `{{.Code}}` - код верификации (6 цифр)

**Особенности:**
- Градиентный заголовок (фиолетовый)
- Большой выделенный код
- Предупреждение о времени действия
- Адаптивный дизайн

### Password Reset Template

**Файл:** `templates/password_reset.html`

**Переменные:**
- `{{.Code}}` - код для сброса пароля (6 цифр)

**Особенности:**
- Градиентный заголовок (розовый)
- Большой выделенный код
- Предупреждение о безопасности
- Адаптивный дизайн

## Кастомизация шаблонов

### Добавление новых переменных

1. Обновите `TemplateData` в `template_manager.go`:

```go
type TemplateData struct {
    Code     string
    UserName string  // новая переменная
    ExpiresIn int    // новая переменная
}
```

2. Обновите шаблон:

```html
<p>Здравствуйте, {{.UserName}}!</p>
<p>Код действителен {{.ExpiresIn}} минут.</p>
```

3. Обновите вызов рендеринга:

```go
data := TemplateData{
    Code: code,
    UserName: user.Name,
    ExpiresIn: 15,
}
htmlBody, err := tm.Render("verification_code.html", data)
```

### Создание нового шаблона

1. Создайте новый `.html` файл в `templates/`:

```html
<!-- templates/welcome.html -->
<!DOCTYPE html>
<html>
<head>
    <title>Добро пожаловать</title>
</head>
<body>
    <h1>Добро пожаловать, {{.UserName}}!</h1>
</body>
</html>
```

2. `TemplateManager` автоматически загрузит его при следующей инициализации.

## Тестирование

### Локальное тестирование (Console)

При использовании `EMAIL_PROVIDER=console`, HTML не рендерится, выводится только код в консоль.

### Тестирование HTML шаблонов

1. Установите `EMAIL_PROVIDER=smtp`
2. Настройте SMTP (можно использовать тестовый сервис)
3. Отправьте email
4. Проверьте в почтовом клиенте

### Предпросмотр шаблонов

Можно создать простой тестовый скрипт:

```go
package main

import (
    "fmt"
    "github.com/nomad-pixel/imperial/internal/infrastructure/email"
)

func main() {
    tm, _ := email.NewTemplateManager()
    html, _ := tm.Render("verification_code.html", email.TemplateData{Code: "123456"})
    fmt.Println(html)
    // Сохраните вывод в .html файл и откройте в браузере
}
```

## Best Practices

1. **Всегда включайте plain text версию** - для клиентов без HTML поддержки
2. **Используйте inline CSS** - многие клиенты игнорируют `<style>` теги
3. **Тестируйте на разных клиентах** - Gmail, Outlook, Apple Mail
4. **Оптимизируйте изображения** - используйте внешние URL или встраивайте base64
5. **Проверяйте мобильную версию** - большинство пользователей читают email на телефоне

## Безопасность

⚠️ **Важно:**
- Шаблоны автоматически экранируют HTML (`html/template`)
- Не используйте `text/template` для пользовательского контента
- Все переменные безопасно экранируются

## Примеры

### Базовый пример

```go
// В SMTPEmailService
htmlBody, err := s.templateManager.Render("verification_code.html", TemplateData{Code: code})
if err != nil {
    return err
}
```

### Расширенный пример с дополнительными данными

```go
data := TemplateData{
    Code: code,
    // Добавьте больше полей при необходимости
}
htmlBody, err := tm.Render("verification_code.html", data)
```

## Troubleshooting

### Шаблон не найден

Убедитесь, что файл находится в `templates/` и имеет расширение `.html`.

### Ошибка рендеринга

Проверьте синтаксис Go template:
- `{{.Field}}` - правильный синтаксис
- `{{Field}}` - неправильный (нет точки)

### HTML не отображается

1. Проверьте, что используется `multipart/alternative`
2. Убедитесь, что `Content-Type: text/html` установлен
3. Проверьте логи SMTP сервера

