Тестовое задание для KODE.
Rest API на языке golang.
Для запуска проекта нужно:
1. Поднять контейнер - docker-compose -f docker-compose.yml up
2. Подключиться к локальной базе данных:
   Зайти на localhost:5051
   Авторизоваться логин - root@root.com пароль- root
   Добавить сервер - name любой, hostname - dbnotes, user - root, password root.
Сервис работает через RestAPI, передача данных с помощью в формате json.

Библиотека логгера - log. Библиотека web-сервера - gorilla.
Сервис предоставляет возможность создавать заметки, выводить их и удалять их.
В проекте реализована аутентификация и авторизация. Пользователи хранятся в файле user.txt внутри проекта.
Сами заметки поподают в БД - postgresql. Также перед добавлением они проверяются с помощью сервиса yandex speller.

После запуска проекта слушается порт 8080. Для авторизации - зайти на 127.0.0.1:8080/login
В postman отправляем запрос POST
{
	"username": "Nikita"
}
После этого вы получите jwt токен.
Можно добавить в файле user.txt любых других пользователей. Логины считываются из этого файле.

Добавить запись на 127.0.0.1:8080/notes с помощью запроса POST
{
	"content": "Купить картошку"
}
Запрос автоматически провериться с помощью yandex speller и добавить в локальную БД.

Получить записи можно зайдя на 127.0.0.1:8080/notes и отправив запрос GET.
Удалить записи можно зайдя на 127.0.0.1:8080/notes и отправив запрос DELETE.





