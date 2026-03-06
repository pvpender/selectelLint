# SelectelLint

Кастомный линтер для Go с расширяемыми правилами проверки кода.

## Установка

### Вариант 1: Установка бинарника

1. Скачать последникй bin [раздела Releases](https://github.com/pvpender/selectelLint/releases).

2. Запустить

```bash
chmod +x custom-gcl
./selectellint-main --help
```

### Вариант 2: Собрать из исходников

**Требования:**

Необходимо наличие **Go** и **golangci-lint**

1. Клониировать репоизиторий
```bash
git clone https://github.com/pvpender/selectelLint
```
2. Перейти в папку
```bash
cd selectelLint
```
3. Изменить в `.custom-gci.yml` **version** на номер версии вашего **golangci-lint**. Формат: `version: v{номер версии}`

4. Запустить

```bash
golangci-lint custom -v
```

## Конфигурация

Линтер использует YAML конфигурацию для настройки правил:

```yaml
version: "2"

linters:
  default: none
  enable:
    - sclint

  settings:
    custom:
      sclint:
        type: module
        settings:
          capitalLetter: true          # Проверка заглавных букв
          englishLetter: true          # Проверка английских букв
          specialLetters: true          # Проверка специальных символов
          sensitiveData: true           # Проверка чувствительных данных
          enableCustomRules: true       # Включение кастомных правил
          rules:
            - name: "Test"
              description: "Digits restricted!"
              pattern: '\d+'             # Regex паттерн для поиска
```

## Использование

```bash
./selectellint-main run
```

## Кастомные правила

Вы можете добавлять свои правила через конфигурацию:

```yaml
rules:
- name: "NoDebug"
  description: "Debug statements are not allowed"
  pattern: 'fmt\.Println\('

- name: "TodoComments"
  description: "TODO comments should be addressed"
  pattern: '//\s*TODO'
```

## Разработка

### Структура проекта

```
├── analyzers/           # Анализаторы кода
│   └── selectelLint/    # Основной анализатор
│       ├── analyzer.go
│       └── testdata/     # Тестовые данные
├── cmd/                  # Точки входа
│   └── selectelLint/     
│       └── main.go
├── config/               # Работа с конфигурацией
│   └── config.go
└── plugin/               # Плагин для golangci-lint
    └── selectellint.go
```

### Тестирование

Тесты запускаются в папке `/analyzers/selectelLint`

```bash
go test
```

## Использование ИИ

В проекте использовался ИИ в основном для консультирования по некоторым вопросам и для форматировании документации