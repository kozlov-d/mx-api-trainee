## Запуск
```
docker-compose up
```
* Развертывает два контейнера - с Postgres и самим приложением
* В образе приложения есть скрипт для ожидания готовности базы данных
* Dockerfile приложения создает образ в несколько этапов
* БД инициируется скриптом для создания таблиц
* Устанавливает необходимые переменные среды для контейнеров из файла .env

## Хранение данных
Merchant | Offer
------------ | -------------
(pk) MerchantID | OfferID 
| | OfferName
| | Price
| | Quantity
| | (pk) OfferID, MerchantID

В рамках выполнения задания разделение на две таблицы было излишним и добавило больше работы. Статистика по задачам на обработку файлов хранится в простом кэше

## Описание API
* После запуска доступен Swagger http://localhost:3000/swagger/docs/index.html

## Использованы

* Роутер [gorilla/mux](http://github.com/gorilla/mux)
* Сервер net/http
* Excel lib [tealeg/xlsx](http://github.com/tealeg/xlsx)
* Драйвер Postgres [jackc/pgx](http://github.com/jackc/pgx)

### Пример запросов с помощью curl

* Для тестирования создавались файлы на 10, 50k и 200k строк со случайным наполнением с помощью скрипта generator.py. Файл с 200k строк размером >6mb обрабатывается примерно за 50 секунд
* Отдельно запускался сервер, хранящий файлы. При этом главное приложение при получении ссылки скачивает файл с помощью Get-запроса

``` curl -X PUT -i http://localhost:3000/merchants/20 -d "{\"link\":\"http://localhost:3001/offers.xlsx\"}" ```

![Put request](https://github.com/kozlov-d/mx-api-trainee/blob/main/docs/put.png)

```curl -X GET -i -G localhost:3000/tasks/1```

После окончания обработки:

![Get tasks request](https://github.com/kozlov-d/mx-api-trainee/blob/main/docs/get_task_completed.png)

```curl -X GET -i -G localhost:3000/offers? --data-urlencode "sub=ABCT"```

![Put request](https://github.com/kozlov-d/mx-api-trainee/blob/main/docs/get_offers.png)

## Усложнения
- [x] Асинхронная работа
- [x] Нагрузочное тестирование
- [ ] Юнит-тесты

## Тесты k6

На следующих графиках отображены два теста для PUT запроса, с max VUs 120 и 250 соответственно. VU посылает запрос раз в секунду. Используется .xlsx небольшого размера на 10 строк. Использование файлов с большим размером крайне затруднительно, ввиду тестирования на хост-машине.   
  
С помощью pprof изучим объекты в куче. Выделив топ-20, можно прийти к выводу, что большая часть размещается библиотекой чтения .xlsx, а также запросом к БД. Это также видно исходя из графа (здесь не отображен), полученного командой web. Некоторые из запросов к БД возвращают указатель на структуру, отсюда и появление их в топе pprof.

![top_20](https://github.com/kozlov-d/mx-api-trainee/blob/main/docs/top_20_alloc.png)

  
Из графиков Grafana видно, что минимальная длительность запроса повысилась с увеличением VUs, а RPS, хоть и возрос, но видны заметные скачки. Для первого теста mean req_duration ~6ms, а для второго - уже 130ms

![put_test](https://github.com/kozlov-d/mx-api-trainee/blob/main/docs/grafana_120_250.png)

Для GET запросов RPS >1k, с заметными перепадами, а max длительность запроса >1.5sec.

![get_test](https://github.com/kozlov-d/mx-api-trainee/blob/main/docs/staged_get_test.png)

Собрать статистику по обработке файлов можно было бы с помощью pandas, предварительно получив с сервера массив всех Tasks, и, например, определить среднее затраченное время на обработку, зависимость времени обработки от индекса задачи

## Что стоит сделать

* Использовать batch-запросы к базе
* Согласно api спецификации, сделать поддержку query params типа ?offerId=1,2 etc
* Переделать все под Clean Architecture
