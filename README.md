[![Build Status](https://travis-ci.org/anthill-com/ImageProcessorService.svg?branch=main)](https://travis-ci.org/anthill-com/ImageProcessorService)
# ImageProcessorService
Rest service with single post method for image processing

* Сервер поддерживает graceful shutdown
* Файла содержат Dockerfile и docker-compose.yml
* Имеется интеграция TravisCI запускающая go тесты

Загрузка:
1. По переданным изображениям генерируется хэш
2. В общем каталоге папок создается папка с хэш именем
3. В базе данных сохраняется информация об изображении с хэшем
4. На запрос передается id записи в базе данных
 
Выгрузка:
1. По переданному id достается хэш из базы
2. По хэшу берется изображение из хэш-папки в общей папки изображений
3. Если необходимо превью - делается превью и созраняется в папку превью в папке хэш
4. Возвращается строка с данными изображения

Сервис выдает ответы в json формате:

```
{ "message": "сообщение", "resCode": 0}
```

*message - ответ от сервера (строка с ошибкой или json)
*resCode:
- 0 - ок
- 1 - ошибка сервера
- 2 - ошибка запроса

Все загружаемые изображения проходят валидацию расширения по возможным из конфига

Запросы передаются в utf8 кодировке

# Список запросов и способы их получения

## json с base64 кодированным изображением

Запрос обязательно должен содержать заголовки: 
* Content-type = application/json
* Req-type = BASE64

В запросе может быть несколько изображений

### Формат запроса:

```
{
	"images": 
	[
		{
            "name": "ff14", 
            "extension": "png", 
            "Data": "base64 строка" 
        },
        {
            "name": "lilil", 
            "extension": "jpg", 
            "Data": "base64 строка" 
        }	
	]
}
```

* name - имя файла
* extension - расширение файла
* Data - base64 кодированная строка

### Формат ответа:

```
{
  "message": "{"file":[{"id":1,"name":"ff14","extension":"png","statis":1,"resMessage":""},{"file":[{"id":1,"name":"lilil","extension":"jpg","statis":1,"resMessage":""}]}",
  "resCode": 0
}
```

В массиве file лежит массив данных отдельно по каждому изображению

* id - id в базе
* name - название изображения
* extension - расширение изображения
* statis - статус обработки 0 - ошибка, 1 - ок
* resMessage - сообщение поясняющее ошибку

## json с URL адресом изображения

Запрос обязательно должен содержать заголовки: 
* Content-type = application/json
* Req-type = URL-LOAD

### Формат запроса:

```
{
	"url": "http://pngimg.com/uploads/peacock/peacock_PNG42.png"
}
```

url - адрес изображения

### Формат ответа:

```
{
  "message": "{"id":3,"name":"peacock_PNG42","extension":"png","statis":1,"resMessage":""}",
  "resCode": 0
}
```

* id - id в базе
* name - название изображения
* extension - расширение изображения
* statis - статус обработки 0 - ошибка, 1 -ок
* resMessage - сообщение поясняющее ошибку

## Загрузка Multipart изображений

Запрос обязательно должен содержать заголовок: 
* Content-type = multipart/form-data

### Формат запроса:

* Обычный Multipart запрос. Обязательно должен содержать filename

### Формат ответа:

```
{
  "message": "{"file":[{"id":1,"name":"Новый документ","extension":"jpg","statis":1,"resMessage":""},{"id":2,"name":"Новый текстовый документ","extension":"png","statis":1,"resMessage":""}]}",
  "resCode": 0
}
```

* file - массив данных по каждой фотографии
* id - id в базе
* name - название изображения
* extension - расширение изображения
* statis - статус обработки 0 - ошибка, 1 -ок
* resMessage - сообщение поясняющее ошибку

## Восстановление изображений

Запрос обязательно должен содержать заголовки: 
* Content-type = application/json
* Req-type = RESTORE

### Формат запроса:

```
{
	"id": 1
}
```

* id - id изображения в базе

### Формат ответа:

```
{
  "message": "{"name":"ff14","extension":"png","data":"строка с изображением"}",
  "resCode": 0
}
```

* name - название изображения
* extension - расширение изображения
* data - строка с данными из файла изображения

## Восстановление превью изображений

Запрос обязательно должен содержать заголовки: 
* Content-type = application/json
* Req-type = RESTORE-PREVIEW

### Формат запроса:

```
{
	"id": 1
}
```

* id - id изображения в базе

### Формат ответа:

```
{
  "message": "{"name":"ff14","extension":"png","data":"строка с изображением"}",
  "resCode": 0
}
```

* name - название изображения
* extension - расширение изображения
* data - строка с данными из файла изображения

## Конфигурация

Приложение конфигурируется по средствам конфигурационного файла config.toml

Список параметров:
* LogFilePath - путь до логов
* Port - порт на котором работает сервер
* ServedURL - адрес по которому принимаются запросы
* ReadTimeout - таймаут на чтение
* WriteTimeout - таймаут на запись
* FileSaveExtensionList - валидные изображения на запись
* ScaledImageRestoreExtension - валидные изображения на скейлинг
* ScaledImageH - размер скейла по высоте
* ScaledImageW - резмер скейла по ширине
* DataBasePath - путь до базы данный (sqlite)
* FileSavePath - папка для сохранения файлов
* PreviewFileFolder - название папки для сохранения превью

# Для сборки необходимо:

* Иметь установленный go
* Перейти в папку с файлом main.go
* Выполнить go build main.go Server.go
* На выходе будет main.exe. Для запуска необходим файл конфига в папке запуска

* Так же можно запустить без сборки приложения go run main.go Server.go

* Если будет нехватать библиотек, необходимо выполнить команды ниже
- go get github.com/mattn/go-sqlite3
- go get github.com/nfnt/resize
- go get github.com/pelletier/go-toml

# Запуск в Docker контейнере

Единственная особенность это то, что необходимо при запуске открывать порт указанный в конфиге для работы приложения

# Запуск DockerCompose

Неоходимо соответствие порта в docker-compose.yml и в файле конфига