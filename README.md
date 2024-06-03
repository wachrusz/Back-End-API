# CAdv API 

The Cash Advisor API is a comprehensive backend service designed to manage user profiles, authentication, analytics, and more. It provides a robust set of endpoints to interact with various functionalities of the CAdv application.
API Specification
# Base URL

    Host: 212.233.78.3:8080


# Security

    HTTPS: All API communications are secured via HTTPS.

    Authentication: Some endpoints require JWT tokens for authentication.

# Contact Information

    Name: Mikhail Vakhrushin

    Email: lstwrd@yandex.com

# API Version

    Version: 0.9.7.4

# Error Handling

The API returns standard HTTP status codes to indicate the success or failure of the requests. Detailed error messages are provided in the response body for failed requests. (or in logs if you can see them(:^))
# Security

All sensitive data is handled securely, and JWT tokens are used for authenticated endpoints. Ensure that tokens are kept confidential and are not exposed to unauthorized parties.
# Contact

For any inquiries or support, please email me, and we'll handle everything together:)

# P.S.
This documentation provides a low-level overview of the Cash Advisor API. For detailed information on each endpoint, including parameters and response formats, please refer to the [Swagger UI](https://212.233.78.3:8080/swagger/index.html) or the [Swagger JSON file](https://github.com/wachrusz/Back-End-API/blob/main/docs/swagger.json).

# Stack

```
  Database: Postgres, Redis (in progress)
  API: Go, REST, JWT, Oauth2, Gorilla Mux, TSB RF API (Only exchange rates)
  Server: Nginx, Ubuntu, Pgbadger
  Build: Docker (In past now R.I.P. due sanctions), Podman
  Documentation: Swagger UI
```
