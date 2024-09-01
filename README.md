# CRM
### CRM Backend System

This project is a backend system for a CRM application developed using the Go programming language and the Gin framework, with MongoDB for data storage. The system handles essential CRM functionalities, including user and customer management, interaction tracking, and advanced analytics, with a strong emphasis on security best practices such as proper encryption and password storage. Additionally, it supports role-based access control, activity notifications, and optional features like email integration and reporting. The application is containerized using Docker and is designed for easy deployment on cloud platforms like AWS, GCP, or Heroku. This repository includes comprehensive API documentation, test cases, and a detailed README file with database schema design, system architecture, and setup instructions.


### Setup and Installation
**Clone the repo**
```
git clone <github directory url >
cd <project directory>
```
**Install dipendency**
```
go mod tidy
```
**Run Application**
```
go run main.go
```


### Folder Structure
 ```
/crm     
  |
  |-- /database
  |   |-- dbConnection.go           # DB connection management (e.g., loading environment variables)
  |
  |-- /controllers
  |   |-- userController.go         # Handler functions for user-related operations
  |   |-- customerController.go     # Handler functions for customer-related operations
  |   |-- interactionController.go  # Handler functions for interaction management
  |   |-- ticketController.go       # Handler functions for ticket
  |   |-- expImportController.go    # Handler functions for export import data
  |
  |-- /models
  |   |-- user.go                    # User model definition
  |   |-- customer.go                # Customer model definition
  |   |-- interaction.go             # Interaction model definition
  |   |-- ticket.go                  # ticket-related models
  |
  |-- /routes
  |   |-- userRoutes.go             # Routes related to user operations
  |   |-- customerRoutes.go         # Routes related to customer operations
  |   |-- authRoutes.go             # Routes related to authentication
  |   |-- exportImportDataRoutes.go # Routes related to analytics and reporting
  |
  |-- /middleware
  |   |-- middleware.go             # Middleware for authentication
  |   |-- rateLimit.go              # Middleware for rate limiting
  |
  |-- /helpers
  |   |-- auth.go                    # Helper function for authentication operations
  |   |-- customer.go                # Helper function for customer operations
  |   |-- user.go                    # Helper function for user operations
  |
  |-- /utils
  |   |-- constant.go               # Utility functions for JWT handling
  |   |-- email.go                  # Utility functions for email notification
  |
  |-- README.md                      # Project overview and setup instructions
  |-- Dockerfile                     # Dockerfile for containerizing the application
  |-- docker-compose.yml             # Docker Compose file for multi-container setup
  |-- go.mod                         # Go module file
  |-- go.sum                         # Go checksum file
  |-- main.go                        # Entry point for the application
  
```

## Features

- **User and Customer Management:**
  - Comprehensive CRUD operations for users and customers.
  - JWT-based authentication for secure user and customer access.
  - Role-based access control (RBAC) to ensure appropriate authorization levels.

- **Interaction Management:**
  - Efficient management of user-customer interactions, including meetings, tasks, and follow-ups.
  - Secure handling and storage of sensitive interaction data.

- **Security Measures:**
  - Role-based access control (RBAC) for precise authorization.
  - Robust password encryption and secure data storage practices.

- **Email Notifications:**
  - Automated email notifications for various interactions and updates.
  - Customizable SMTP settings for email service integration.

- **Data Import/Export:**
  - Support for importing and exporting data in CSV and JSON formats.
  - Role-based permissions for controlling data import and export access.

- **Rate Limiting:**
  - Implemented rate limiting to guard against DOS/DDOS attacks and prevent abuse by controlling request rates.

 
### API Endpoints

### Import/Export Data Routes
 - ExportcData (CSV/JSON): GET /export/customer_data
 - Import Data:            POST import/customer_data

### Customer Routes
 - Get Customers:           GET /customers
 - Get Customer by ID:      GET /customers/:cust_id
 - Update Customer:         PATCH /customers/:cust_id
 - Delete Customer:         DELETE /customers/:cust_id
   
**Auth Routes**
 - Register User:           POST customer/signup
 - Login User:              POST /customer/signin
 - Register Customer:       POST user/signup
 - Login Customer:          POST user/signin

### User Routes
 - Get Users:               GET /users
 - Get User by ID:          GET /users/:user_id
 - Update User:             PATCH /users/:user_id
 - Delete User:             DELETE /users/:user_id

### Ticket Routes
 - Get All Tickets:         GET /customers/tickets/
 - Create Ticket :          POST /customers/ticket/:interaction_id
 - Update Ticket:           PATCH /customers/ticket/:ticket_id
 - Delete Ticket:           DELETE /customers/ticket/:ticket_id

### Interaction Routes
 - Get Interactions:                  GET /users/meetings/
 - Get Interactions by User ID:       GET /interactions/user/:user_id
 - Get Interactions by CustomerId ID: GET /interactions/user/:cust_id
 - Get Interaction by Interaction ID: GET /interactions/:interaction_id
 - Create Interaction:                POST /interactions
 - Update Interaction:                PATCH /interactions/:interaction_id
 - Delete Interaction:                DELETE /interactions/:interaction_id

### System Diagram
```  +-------------------+
  | Client Application|
  +--------+----------+
           |
           v
  +--------+----------+
  |     API Gateway   |--------------------+
  +--------+----------+                    |     
           |                               |
           v                               |
+----------------------------+             |
|  Authentication &          |             |
|  Authorization (JWT)       |             |
+----------------------------+             |
           |                               |
           v                               v
  +--------+----------+----     -+--------------------+
  |        |          |           |                  |
  v        v          v           v                  v
+---------+   +------+-------+ +-------+------+ +---------------+
| User    |   | Customer     | | Interaction |  | Export/Import |
| Service |   | Service      | | Service     |  | Data Service  |
+---------+   +--------------+ +-------------+  +------+---------
    |             |                  |                 |
    |             v                  v                 |
    |     +-------------+ +---------------+            |
    |     | Import/Export|  | Notification|            |
    |     | Service      |  | Service     |            |
    |     +-------------+ +---------------+            |
    |                |                                 |
    |                v                                 |
    |     +--------------------------+                 |
    |   |        Mongo DB            |                 |
    |---|        Database            |-----------------|
        +----------------------------+
```
## APIs
**customer/signup**
- Curl :
```
    curl --location 'localhost:8080/customer/signup' \
 --header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJYVxETlARgxBtGbUdMsBuwJltNsTbH41k' \
 --header 'Content-Type: multipart/form-data' \
 --data-raw '{
         "Name": "ramlal",
         "Email": "ram@gmail.com",
         "Password": "$2a$15$lShiva",
         "Company": "TATA",
         "Phone": "+91 35636621762"
     }'
```
- Response :
```
   {
     "insertId": {
         "InsertedID": "66d3ccc9e71590f28320f639"
     },
     "message": "Customer created successfully"
 }
```
**customer/signin**
  -Curl : 
  ```
 curl --location 'localhost:8080/customer/signin' \
 --header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1ZGYyMWQ3MWIzIiwiUm9sZSI6IkFETUlOIiwiZtGbUdMsBuwJltNsTbH41k' \
 --header 'Content-Type: multipart/form-data' \
 --data-raw '{
         "Email": "ram@gmail.com",
         "Password": "$2a$15$lShiva"
     }'
```
  -Response:
  ```
 {
     "customer": {
         "customer_id": "66d3ccc9e71590f28320f639",
         "id": "66d3ccc9e71590f28320f639",
         "name": "ramlal",
         "email": "ram@gmail.com",
         "password": "$2a$15$ngAkPv3vGF82A.z/VAIz8.QW/zaqGG0Ytu4TbOEOsXL0NBnjIItFW",
         "company": "TATA",
         "phone": "+91 35636621762",
         "token": "token",
         "created_at": "2024-09-01T02:09:13Z",
         "updated_at": "2024-09-01T02:11:17Z"
     },
     "message": "Customer logged in successfully"
 }
```

**user/signup**
-Curl:  
```
 curl --location 'localhost:8080/user/signup' \
 --header 'token: eyJhbGcik' \
 --header 'Content-Type: multipart/form-data' \
 --data-raw '{
         "Email": "sita@gmail.com",
         "Password": "$15$lShiva",
         "name":"sita",
         "role":"ADMIN"
     }'
```
-Response:
     ```
     {
     "message": "User created successfully",
     "user": {
         "id": "66d3cdf3e71590f28320f63a",
         "user_id": "66d3cdf3e71590f28320f63a",
         "name": "sita",
         "password": "$2a$15$hhJdAdiMr2Kk11/iqnf3EuO9aqFI/ZyVAm8GuucDsa4SU1EdxRmiC",
         "email": "sita@gmail.com",
         "role": "ADMIN",
         "token": "eyJyODMyMGY2M2EiLCJSb2xE3M6o",
         "created_at": "2024-09-01T07:44:11+05:30",
         "updated_at": "2024-09-01T07:44:11+05:30"
     }
 }```
 
**user/signin**
    -Curl: 
     ```curl --location 'localhost:8080/user/signin' \
 --header 'token: eyJhbGciOiTHGvAgYVxETlARgxBtGbUdMsBuwJltNsTbH41k' \
 --header 'Content-Type: multipart/form-data' \
 --data-raw '{
         "Email": "sita@gmail.com",
         "Password": "$15$lShiva"
     }' ```
-Response:
     ```{
     "msg": "User logged in successfully",
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InNpdGFAZ21haWwuY29tIiwiTmFtZSI6InNpdGEiLCJVaWQiOiI2NmQzY2RmM2U3MTU5MGYyODMyMGY2M2EiLCJSb2xlIjoiQURNSU4iLCJleHAiOjE3MjUyNDM0MDh9.eqae6DvxivdkH4kuZnlT1Dw3CuwalJbE_TFKM4giU20",
     "user": {
         "id": "66d3cdf3e71590f28320f63a",
         "user_id": "66d3cdf3e71590f28320f63a",
         "name": "sita",
         "password": "$2a$15$hhJdAdiMr2Kk11/iqnf3EuO9aqFI/ZyVAm8GuucDsa4SU1EdxRmiC",
         "email": "sita@gmail.com",
         "role": "ADMIN",
         "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InNpdGFAZ21haWwuY29tIiwiTmFtZSI6InNpdGEiLCJVaWQiOiI2NmQzY2RmM2U3MTU5MGYyODMyMGY2M2EiLCJSb2xlIjoiQURNSU4iLCJleHAiOjE3MjUyNDM0MDh9.eqae6DvxivdkH4kuZnlT1Dw3CuwalJbE_TFKM4giU20",
         "created_at": "2024-09-01T02:14:11Z",
         "updated_at": "2024-09-01T02:16:48Z"
     }
}```

**users**
-Curll : 
```
  curl --location 'localhost:8080/users' \
  --header 'token: 1ZGYyMWQ3MWIzIiwiUm9sZSI6IkFETUlOIiwiZXhwIjoxNzI1MTk5MjIxfQ.VbHxfTHGvAgYVxETlARgxBtGbUdMsBuwJltNsTbH41k' \
  --header 'Content-Type: multipart/form-data' \
  --data ''
```

Reaponse: 
```
[
    {
        "id": "66d321f59940e15df21d71b3",
        "user_id": "66d321f59940e15df21d71b3",
        "name": "Nirmal",
        "password": "$2a$15$KAKnHqDbrgn/4lZIxJ1Unu6RyRHn4CvmWgAAe3c3thV8IE0UPbkmu",
        "email": "abc@gmail.com",
        "role": "USER",
        "token": "eyJhbGciOiJIUzIZTE1ZGYyMWQ3MWIzIiwiUm9sZSI6IkFETUlOIiwiZXhwIjoxNzI1MTk5MjIxfQ.VbHxfTHGvAgYVxETlARgxBtGbUdMsBuwJltNsTbH41k",
        "created_at": "2024-08-31T14:00:21Z",
        "updated_at": "2024-08-31T14:00:21Z"
    },
    {
        "id": "66d3233ae3d3bbd29769c153",
        "user_id": "66d3233ae3d3bbd29769c153",
        "name": "temp",
        "password": "$2a$15$Jg/2lUsl4Cts/mDPtVMt0eQhE3OTa5kELElvUyj1JoSL39ApZPeuO",
        "email": "abcd@gmail.com",
        "role": "USER",
        "token": "eyJhbGciOiJIUzI1NiIs\MiLCJSb2xlIjoiVVNFUiIsImV4cCI6MTcyNTE5OTc4OH0.a3PKeom_j9LWtCW1YZszaBBVV-VAPnj9fjZaEfYSBYA",
        "created_at": "2024-08-31T14:05:46Z",
        "updated_at": "2024-08-31T14:09:48Z"
    },
    {
        "id": "66d3cdf3e71590f28320f63a",
        "user_id": "66d3cdf3e71590f28320f63a",
        "name": "sita",
        "password": "$2a$15$hhJdAdiMr2Kk11/iqnf3EuO9aqFI/ZyVAm8GuucDsa4SU1EdxRmiC",
        "email": "sita@gmail.com",
        "role": "ADMIN",
        "token": "eyJhbGciOiJIUzI1NlIjoiQURNSU4iLCJleHAiOjE3MjUyNDM0MDh9.eqae6DvxivdkH4kuZnlT1Dw3CuwalJbE_TFKM4giU20",
        "created_at": "2024-09-01T02:14:11Z",
        "updated_at": "2024-09-01T02:16:48Z"
    }
]

```

**users/:user_id**
-Curl : 
```
  curl --location 'localhost:8080/users/66d321f59940e15df21d71b3' \
  --header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6Ikp1ZGYyMWQ3MWIzIiwiUm9sZSI6IkFETUlOIiwiZXhwIjoxNzI1MTk5MjIxfQ.VbHxfTHGvAgYVxETlARgxBtGbUdMsBuwJltNsTbH41k' \
  --header 'Content-Type: multipart/form-data' \
  --data ''
```

-Response: 
```
  {
      "id": "66d321f59940e15df21d71b3",
      "user_id": "66d321f59940e15df21d71b3",
      "name": "Nirmal",
      "password": "$2a$15$KAKnHqDbrgn/4lZIxJ1Unu6RyRHn4CvmWgAAe3c3thV8IE0UPbkmu",
      "email": "abc@gmail.com",
      "role": "USER",
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbTUlOIiwiZXhwIjoxNzItNsTbH41k",
      "created_at": "2024-08-31T14:00:21Z",
      "updated_at": "2024-08-31T14:00:21Z"
  }
```

**user/:user_id**
-Curl: 
```
 curl --location --request DELETE 'localhost:8080/users/66d321f59940e15df21d71b3' \
 --header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFjU5OTQwZTE1ZGYyMWQ3MWIzIiwiUm9sZSI6IkFETUlOIiwiZXhwIjoxNzI1MTk5MjIxfQ.VbHxfTHGvAgYVxETlARgxBtGbUdMsBuwJltNsTbH41k' \
 --header 'Content-Type: multipart/form-data' \
 --data ''```

-Response: 
```
 {
     "message": "User deleted successfully"
 }
```

**user/:user_id**
-Curl: 
```
 curl --location --request PUT 'localhost:8080/users/66d3233ae3d3bbd29769c153' \
 --header 'token: eyJhbGciOiJIUzI1NiIsInRU5OTQwZTE1ZGYyMWQ3MWIzIiwiUm9sZSI6IkFETUlOIifTHGvAgYVxETlARgxBtGbUdMsBuwJltNsTbH41k' \
 --header 'Content-Type: multipart/form-data' \
 --data '{
     "name":"sahil"
 }'
```
-Response: 
```
 {
     "message": "user updated successfully"
 }
```

**users/meetings/:customer_id**
 -Curl: 
 ```
  curl --location 'localhost:8080/users/meetings/66d324c3e3d3bbd29769c154' \
   --header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpiwiZXhwIjoxNzI1MTk5MjIxfQ.VbHxfTHGvAgYVxETlARgxBtGbUdMsBuwJltNsTbH41k' \
   --header 'Content-Type: multipart/form-data' \
   --data '{}'
```

-Response:
```
 {
     "InsertedID": "66d3d83afff6681e0956ca2d"
 }```
