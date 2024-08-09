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
        

   - Примеры запросов:
     - Регистрация
       ![изображение](https://github.com/JuFnd/ozon-task/assets/109366718/064c2b64-97a3-4e4c-b8a1-47a316d25a20)

     - Авторизация
       ![изображение](https://github.com/JuFnd/ozon-task/assets/109366718/2737acbd-eebc-40e5-ac7f-fe9f17f8a9a6)

     - Выход из аккаунта
       ![изображение](https://github.com/JuFnd/ozon-task/assets/109366718/29320114-ac6d-4b80-a09b-f8e161a3d45a)
