
# Cash Advisor Back-end API

Документация находится на
<http://localhost:8080/swagger/index.html>

## Восстановление базы данных

1. Загрузите и установите необходимую схему базы данных из SQL-файла.

    ```bash
    psql -U ваше_имя_пользователя -d имя_вашей_базы_данных -h хост -p порт -f ваш_файл.sql
    ```

    Например:

    ```bash
    psql -U postgres -d backendapi -h localhost -p 5432 -f backup.sql
    ```

2. Установите переменные окружения для вашего проекта, чтобы обеспечить подключение к базе данных. Пример:

    ```bash
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER=ваш_пользователь
    export DB_PASSWORD=ваш_пароль
    export DB_NAME=ваша_база_данных
    ```

3. Запустите ваш проект.

    ```bash
    go run main.go
    ```

Теперь ваш проект должен быть подключен к базе данных PostgreSQL с использованием указанных переменных окружения.
