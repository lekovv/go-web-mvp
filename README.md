# Go-Checkup MVP

MVP веб-приложения онлайн записи к врачам

*данный проект является "игровым" (упрощенным), в "бою" он бы выглядел иначе

### В проекте реализовано: 
* Полноценная Jwt авторизация
* Token blacklist механизм для Logout
* Горутина для очистки истекших токенов по расписанию
* HMAC-SHA256 для хэширования токенов
* Bcrypt для хэширования паролей
* Role-based authorization
* DI с помощью интерфейсов
* Многослойная архитектура
* Middleware error handler
* Миграции БД
* Структурный валидатор
* Graceful shutdown
* Роуты для администраторов, пациентов, докторов

### Технологический стэк:
* REST API
* Golang
* Fiber
* PostgreSQL (GORM)
* Viper
* go-playground/validator
* golang-migrate

### Описание эндпойнтов:

Роуты администраторов (роль `admin` проверяется в `Jwt Claims`):

1) `POST /api/admin/create-doctor` - метод для создания доктора в системе

Тело запроса: 
```json
{
    "email": "doctor@mail.ru",
    "password": "0000",
    "gender": "male",
    "first_name": "Доктор",
    "last_name": "Хаус",
    "specialization": "Стоматолог",
    "phone_number": "89991112233",
    "bio": "супер доктор, обращайтесь",
    "experience_years": 5,
    "price": 5000
}
```
2) `PATCH /api/admin/update-user/:id` - метод для частичного обновления данных пользователя. Можно обновить как все поля сразу, так и одно из них

Тело запроса: 
```json
{
   "first_name": "John",
   "last_name": "Doe",
   "middle_name": "Vladimirovich",
   "is_active": false
}
```

3) `DELETE /api/admin/delete-user/:id` - метод для полного удаления пользователя из системы

Auth роуты:

1) `POST /api/auth/registration` - метод для самостоятельной регистрации пациента

Тело запроса:
```json
{
    "email": "patient@mail.ru",
    "password": "0000",
    "gender": "male",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "89991112233",
    "birth_date": "1994-12-28"
}
```

2) `POST /api/auth/login` - метод для авторизации в системе. В ответе: `jwt token`

Тело запроса: 
```json
{
  "email": "example@mail.ru",
  "password": "0000"
}
```

3) `POST api/auth/logout` - метод для выхода из системы. Вытаскивает токен, хэширует его и кладет в `black list`

Пользовательские роуты:

1) `GET api/user/get-user` - универсальный метод для вывода основной информации о пользователе (докторе или пациенте). Вытаскивает `user_id` из `Jwt Claims`

В ответе (если авторизован доктор): 
```json
{
    "data": {
        "id": "f50f3cc8-d3b8-46d2-85e6-aa1de4f383f9",
        "email": "doctor@mail.ru",
        "gender": "male",
        "first_name": "Доктор",
        "last_name": "Хаус",
        "is_active": true,
        "specialization": "Стоматолог",
        "phone_number": "89991112233",
        "bio": "супер доктор, обращайтесь",
        "experience_years": 5,
        "price": 5000
    },
    "status": "success"
}
```

В ответе (если авторизован пациент):
```json
{
    "data": {
        "id": "459158e2-5489-47d0-a0e3-3444af626ca2",
        "email": "patient@mail.ru",
        "gender": "male",
        "first_name": "John",
        "last_name": "Doe",
        "middle_name": "Vladimirovich",
        "is_active": true,
        "phone_number": "89991112233",
        "birth_date": "1994-12-28"
    },
    "status": "success"
}
```




