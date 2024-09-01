# CRM
### CRM Backend System

This project is a backend system for a CRM application developed using the Go programming language and the Gin framework, with MongoDB for data storage. The system handles essential CRM functionalities, including user and customer management, interaction tracking, and advanced analytics, with a strong emphasis on security best practices such as proper encryption and password storage. Additionally, it supports role-based access control, activity notifications, and optional features like email integration and reporting. The application is containerized using Docker and is designed for easy deployment on cloud platforms like AWS, GCP, or Heroku. This repository includes comprehensive API documentation, test cases, and a detailed README file with database schema design, system architecture, and setup instructions.


### Folder Structure
 ``` /crm     
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

