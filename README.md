# CAdv API 

The **Cash Advisor API** is a comprehensive backend service designed to manage user profiles, authentication, analytics, and more. It offers a robust set of endpoints to interact with various functionalities of the CAdv application.

# Team
> CTO: [Mikhail Vakhrushin](https://github.com/wachrusz)
>
>Email: wachrusz@gmail.com
>
>>Developer: [Rostislav Zhukov](https://github.com/zhukovrost/)

# Stack
>**Database**
>>Postgres, Redis (in progress)
>
>**API**
>>Go, REST, JWT, Oauth2, Gorilla Mux, TSB RF API (Only exchange rates), RabbitMQ, etc.
>
>**Server**
>>Nginx, Ubuntu, Pgbadger
>
>**Build**
>>Docker, Compose
>
>**Documentation**
>>Swagger UI
>

# External Services
>[Rostislav Zhukov](https://github.com/zhukovrost/)
>>[cadv_email](https://github.com/zhukovrost/cadv_email)
>>
>>[cadv_logger](https://github.com/zhukovrost/cadv_logger)

# Error Handling

The API returns standard HTTP status codes to indicate the success or failure of the requests. Detailed error messages are provided in the response body for failed requests. (or in logs if you can see them **_(:^)_** )

# Security
All sensitive data is handled securely, and JWT tokens are used for authenticated endpoints. Ensure that tokens are kept confidential and are not exposed to unauthorized parties.

# Contact
For any inquiries or support, please [email me](wachrusz@gmail.com), and we'll handle everything together **:)**

# P.S.
This documentation provides a low-level overview of the Cash Advisor API. For detailed information on each endpoint, including parameters and response formats, please refer to the [Swagger UI](https://212.233.78.3:8080/swagger/index.html) or the [Swagger JSON file](https://github.com/wachrusz/Back-End-API/blob/main/docs/swagger.json). _And some time we'll even add our own wiki..._

