# Примечание перед началом работы

Тесты в данном проекте на данный момент не являются полностью доделанными.  Я признаю в этом свою ошибку ведь я поверхностно прочитал документацию
и пошел делать тесты не качеством - а количеством не проверяя тесты на их работоспособность. 
До конца законченные тесты отмечены //go: build unit. Также проект не является  хорошо масштабируемым. Это доказывает сложность написания тестов.


# Содержание
1. [Установка компонентов](#установка)
2. [как создать таблицы в mysql](#создание-таблиц-mysql)
3.  [как запустить проект](#запуск-проекта)
4. [примеры curl запросов](#curl-запросы)
5. [Как работает проект (e2e объяснение)](#Архитектура)
5. [информация о важных пакетах](#информация-о-пакетахсамых-основных)

# Установка

## Установка проекта
```bash
git clone git@github.com:bust6k/lms-project-final.git
```
или если не работает попробуйте
```bash
git clone https://github.com/bust6k/lms-project-final.git
```
## Установка зависимостей

#### установка go gin
```bash
go get -u github.com/gin-gonic/gin

```
#### Установка zap
```bash
go get -u go.uber.org/zap
```

#### Установка sonic

```bash
go get -u github.com/bytedance/sonic
```
#### установка gRPC кода
```bash
go get github.com/bust6k/protoLMS
```

#### установка драйвера mysql 

```bash
go get github.com/go-sql-driver/mysql
```

#### установка  protobuf  библиотеки
```bash
go get -u google.golang.org/protobuf
```

#### установка библиотеки gRPC
```bash
go get google.golang.org/grpc
```
#### Установка crypto
```bash
go get -u golang.org/x/crypto
```
#### Установка testify
```bash
go get github.com/stretchr/testify
```

#### Установка sqlx
```bash
go get github.com/jmoiron/sqlx 
```
#### установка jwt/v5
```bash
go get github.com/golang-jwt/jwt/v5
```

#### Установка uuid
```bash
go get github.com/google/uuid 
```

###  Или если не хочется скачивать все отдельно выполните
```bash
make -f Makefile.addiction
```

### предварительно  перейдя по каталогам
```bash
cd  lms-project-final
cd 'project yandex lms'

```

## Установка mysql
### для linux/mac os
Обновление репозитория
```bash
sudo apt-get update
```
Установка сервера
```bash
sudo apt-get install mysql-server
```
Установка клиента
```bash
sudo apt-get install mysql-client
```
зайти в mysql
```bash
sudo mysql -u root
```
Изменить на пароль используемый в коде
```bash
ALTER USER 'root'@'localhost' IDENTIFIED BY 'pass12345';
```

###  Для Windows
[скачайте mysql](https://dev.mysql.com/downloads/installer/)

Выберите MySQL Installer (MSI) → нажмите Download
(версия mysql-installer-community-*.msi)



   Запустите скачанный .msi файл от имени администратора.

   Выберите тип установки:
   Developer Default (все необходимое для разработки).

   Нажмите Execute → дождитесь загрузки компонентов → Next.


   В разделе High Availability выберите:
   Standalone MySQL Server / Classic MySQL Replication.

   В Type and Networking оставьте настройки по умолчанию → Next.

   Укажите пароль root-пользователя (запомните его!) → Next.



   Нажмите Execute → дождитесь установки → Finish.

   
В конце установки поставьте галочку "Start MySQL Server at Startup".


Открыть MySQL Command Line Client



После установки через MySQL Installer найдите в меню Пуск:
    MySQL → MySQL Command Line Client
    (или запустите cmd и введите: mysql -u root -p)

Войти в MySQL

   Введите пароль, который задали при установке.

Изменить пароль root


# Создание таблиц mysql
Сначала создадите базу данных LMS
```sql
CREATE DATABASE LMS;
```
Зайдите в эту базу данные 
```sql
USE LMS;
```

создадите таблицу Users
```sql
CREATE TABLE Users (
                       id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                       Login VARCHAR(20) UNIQUE,
                       Password VARCHAR(100),
                       User_id VARCHAR(100) UNIQUE
);
```
создадите таблицу ProcessedExpressions

```sql
CREATE TABLE ProcessedExpressions (
                                      Id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                                      Status VARCHAR(20),
                                      Result FLOAT,
                                      user_id VARCHAR(100) NOT NULL,
                                      CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES Users(User_id)
);
```


# Запуск проекта

### Запуск приложения
```bash
make  all #находясь в директории 'project yandex lms'
```
###  Запуск тестов
```bash
go test -v -tags=unit ./... #находясь в директории  'project yandex lms'
```

# curl запросы

### регистрация
Отправить корректные данные для регистрации
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"ivan":"123"}'
```
Отправить данные для регистрации в виде xml(некорректные данные)
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/xml" \
  -d '<user><ivan>123</ivan></user>'
```
Отправить синтаксически неправильные данные
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{ivan:"123"}'
```

### логин
Отправить корректные данные для логина 
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"ivan":"123"}' \
  -c cookies.txt
```
Отправить несуществующего пользователя

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"ruslan":"12345"}' \

```

### вычисления
Отправить корректное выражение с куками
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression":"2+2*2"}' \
-b cookies.txt
```
отправить корректное выражение без кук
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression":"2+2*2"}' \
```
Отправить некорректное выражение с куками
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression":"2в"}' \
-b cookies.txt
```
отправить некорректное выражение без кук
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression":"2в"}' \

```
### Получение выражений
Получить все выражения
```bash
curl -X GET http://localhost:8080/api/v1/expressions \
-H "Content-Type: application/json" \
 \
-b cookies.txt
```
Получить конкретное выражение
```bash
curl -X GET http://localhost:8080/api/v1/expressions/id \
-H "Content-Type: application/json" \
 \
-b cookies.txt
```
*Замените поле id на существующий id*




# Архитектура

### Регистрация 

юзер отправляет данные для регистрации на сервер. Сервер  в свою очередь берет имя и пароль, генерирует уникальный user_id для пользователя и сохраняет юзера в **СУБД**

### Логин

юзер отправляет данные для  аутентификации в систему. Сервер в свою очередь вытягивает user_id  зная логин и пароль,тем самым ещё и аутентифицируя юзера(нарушение srp). после вы вытягивания user_id,сервер создаёт accsess  и  refresh токен с полем user_id.

### Создание выражения 

юзер аутентифицировался, теперь хочет вычислить выражение. Для этого он идёт на /api/v1/calculate. сервер проматывает его через цепочку миддлвейров(Аутентификации и авторизации) и только после этого он берет user_id из цепочки и создаёт выражение в системе вызывая Агента  для подсчёта выражения и парсинга выражений на задачи отправляя их в /internal как gRPC stub.

#### Миддлвейры

 ##### Миддлвейр Аутенцификации

миддлвейр берет Из Кук Accsess токен,если куки нет на месте значит accsess просрочен, сразу отправляем на Миддлвейр авторизации,пусть сам там разберется.
Если же Кука на места  то мы проверяем юзера в **СУБД** по его user_id которое мы передаём в токенах.  Если такого юзера нет то тогда кидаем ошибку 403-юзер не Аутентифицирован и передаём наш токен с user_id

##### Миддлвейр авторизации

миддлвейр сначала проверяет является ли accsess токен инвалидным. Если это так то  мы используем refresh токен. если refresh токена нет значит юзер не авторизован. Если же есть то мы берём из него user_id и обновляем токены на основе user_id. если accsess валидный,то мы берём user_id и токен из цепочки Миддлвейров  и проверяем что токен валидный(просрочился exp, не начался nbf,ошибка подписи токеа и др. ошибки) если это так то передаём в цепочку миддлвейров user_id. иначе-401 юзер не авторизован.

### Просмотр выражений

#### просмотр всех выражений

юзер хочет посмотреть все свои выражения. Он идёт на /api/v1/expressions. Сервер ведёт его по цепочке Миддлвейров см    [Миддлверы](#Архитектура#Миддлверы). после цепочки Миддлвейров сервер берет user_id Из цепочки Миддлвейров, подключается как gRPC stub к  gRPC  серверу processed expressions вызывая метод для получения всех выражений. дальше происходит преобразование массива  ListPorcessedExlressionRequest в processed expressions и преобразованный слайс отправляется в виде json как ответ.

#### просмотр конкретного выражения

сервер прокручивает юзера через цепочку миддлвейров см [Миддлверы](#Архитектура#Миддлверы) дальше сервер получает user_id Из цепочки Миддлвейров а также получает id задачи из url. дальше сервер вызывает метод получения отдельного выражения по id. сервер отправляет полученный результат 

#  Информация о пакетах(самых основных)



| Имя пакета   | предназначение |
|--------------|----------------|
| pkg          |  содержит код  который может быть пере использован в следующих проектах
| calc         | содержит код калькулятора из 0 спринта    
| config       |  содержит структуры-конфиги и их конструкторы с дефолтными значениями
| database     | содержит основные CRUD функции для работы с **БД**. не путайте с repository! database имеет более низкоуровневую логику
| entities      | содержит общие структуры которые используются другими пакетами
| grpc         |  содержит код g**RPC** серверов и main.go файлы для их запуска
| models       |  содержит структуры представляющие  таблицы в **БД**
| variables    |  содержит глобальные переменные  используемые всей программой. но на данный момент почти не используется,  оставлен для обратной совместимости со старой версией.
| application  |  содержит код для запуска оркестратора и агента
| Orchestrator | Ядро системы. содержит основные компоненты оркестратора, алгоритмы на которых он работает,  и web компоненты
| proto        | содержит proto-контракт которому придерживаются все g**RPC** серверы из пакета *grpc*
| agent        | содержит код самой главной рабочей лошадки - агента







