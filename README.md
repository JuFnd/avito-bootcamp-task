# avito-bootcamp-task
Avito Bootcamp Task

### Запуск
```
docker-compose up
```

### Описание архитектуры:
   Реализована микросевисная архитектура, общение сервисов по gRPC.
   - Реализовано два микросервиса:
     
   - Авторизация:
        Авторизация реализована на основе сессий.
        ![UNXyORX2pQ8](https://github.com/JuFnd/avito-task/assets/109366718/0a8f1eaa-9af5-4eef-bfc2-df2969b1bc46)

        - Схема БД:

          ![изображение](https://github.com/JuFnd/avito-task/assets/109366718/a36e0419-5f02-4d8d-a069-87d5304ffafd)

        - СУБД: Postgresql
        - БД Кэширования: Redis
     
   - Квартиры:
        - Схема БД:
          ![изображение](https://github.com/user-attachments/assets/d209f8d0-fa68-47fe-9ae9-fb3ce51815c3)
          
        - СУБД: Postgresql
        

   - Примеры запросов:
     - Регистрация
       ![изображение](https://github.com/JuFnd/ozon-task/assets/109366718/064c2b64-97a3-4e4c-b8a1-47a316d25a20)

     - Авторизация
       ![изображение](https://github.com/JuFnd/ozon-task/assets/109366718/2737acbd-eebc-40e5-ac7f-fe9f17f8a9a6)

     - Выход из аккаунта
       ![изображение](https://github.com/JuFnd/ozon-task/assets/109366718/29320114-ac6d-4b80-a09b-f8e161a3d45a)

     - Получение объявлений(юзер/модератор)
       ![изображение](https://github.com/user-attachments/assets/c4461139-b4cf-4ce4-8225-16f6eb1e5d09)
       ![изображение](https://github.com/user-attachments/assets/4a5942bd-aff2-4a13-adf9-dd4ead87704a)

     - Обновление квартиры
       ![изображение](https://github.com/user-attachments/assets/68d6e609-8421-41cc-83a8-55fad5301592)

     - Создание квартиры
       ![изображение](https://github.com/user-attachments/assets/9e2b6f7e-14e4-44ec-a8c4-b4b98026fd96)


Запросы:

      localhost:8081/api/v1/house/1 POST
      localhost:8081/api/v1/house/2 POST
      localhost:8081/api/v1/house/3 POST

      localhost:8080/login localhost:8080/register POST
      {
          "login":"test",
          "password":"test"
      }

      localhost:8081/api/v1/flat/update POST - Необходимы права модератора
      {
          "apartment_number": 114,
          "price": 500000,
          "rooms": 4,
          "house_id": 1,
          "address": "123 Main St",
          "status": "approved"
      }

      localhost:8081/api/v1/flat/create
      {     
         "apartment_number": 114,
         "price": 500000,
         "rooms": 4,
         "house_id": 1
      }


