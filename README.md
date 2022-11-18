# Avito_internship
Тестовое задание на позицию стажера-бэкендера

Реализовано:
1. Метод получения баланса пользователя
2. Метод начисления средств на баланс пользователя
3. Метод резервирования средств  
4. Метод отмены резервирования средств  
5. Метод перевода средств от одного пользователя к другому
6. Метод создания заказов
7. Первое дополнительно задание про создания отчета и открытия его
8. Выгрузка истории транзакций пользователя с возможностью сортровки 

Доступ к базе данных:
- Database: avito
- User: postgres
- Password: 123

Реализация работы программы:
- База данных: postgres
- Драйвер для работы: pgx
- Фреймворк: gin
- Есть дополнитильная валидация для баланса пользователя

Запуск:
- К сожалению, у меня возникли проблемы с docker, так как я столкнулся с ним впервые. На данный момент выкладываю без него, но не оставляю попыток разобраться с ним, чтобы отправить работу согласно требованиям

Примеры методов:

(GET) /test
Тестовый запрос 

Пример:
![image](https://user-images.githubusercontent.com/80826818/202691380-ab7541e2-7a63-43b0-8c4a-ec5099b72db7.png)


(GET) /balance/get
Запрос на состояние баланса пользователя 

Пример:
![image](https://user-images.githubusercontent.com/80826818/202691648-45ef375c-2c61-4641-b046-52ec12911849.png)

(POST) /balance/add
Запрос на добавление баланса для пользователя  

Пример:

![image](https://user-images.githubusercontent.com/80826818/202691831-53613bf4-0241-40a0-ad61-f5bab016e6dd.png)

(POST) /reserve/make
Запрос на покупку услуги 

Должны быть соблюдены следующие условия:
- аккаунт существует
- сумма не больше счета пользователя

У пользователя достаточно средств:
Пример:
![image](https://user-images.githubusercontent.com/80826818/202692032-b87b6a86-a818-4b35-99e2-60806e25daef.png)

У пользователя не достаточно средств:
Пример:
![image](https://user-images.githubusercontent.com/80826818/202692083-540487f6-7340-4cca-9837-9b2987a71090.png)

(POST) /reserve/accept

Запрос на подтверждение или отмену операции 
Пример:
![image](https://user-images.githubusercontent.com/80826818/202692195-4dc3f4f3-84e9-4f19-b7b5-0025c416d095.png)

(POST) /balance/transfer

Запрос на перевод средств от одного пользователя к другому 
Пример:
![image](https://user-images.githubusercontent.com/80826818/202692475-9c0ad002-3a6f-44fa-b5e4-7047060def22.png)

(GET) /report

Запрос на создание отчета 
Пример:
![image](https://user-images.githubusercontent.com/80826818/202692575-d379583c-595e-440b-ad83-dca61e58a0cb.png)

(GET) /reserve/csv/:file

Запрос на открытие отчета 
Пример:
![image](https://user-images.githubusercontent.com/80826818/202692632-e5ba6b0d-cb33-4b85-8768-79d188ae16bf.png)

(GET) /balance/list

Запрос на историю транзакций пользователя 
Пример:
![image](https://user-images.githubusercontent.com/80826818/202692758-1d785e85-06bf-451d-84d5-b7c570b6ffd0.png)
