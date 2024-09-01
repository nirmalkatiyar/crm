# CRM
### CRM Backend System

This project is a backend system for a CRM application developed using the Go programming language and the Gin framework, with MongoDB for data storage. The system handles essential CRM functionalities, including user and customer management, interaction tracking, and advanced analytics, with a strong emphasis on security best practices such as proper encryption and password storage. Additionally, it supports role-based access control, activity notifications, and optional features like email integration and reporting. The application is containerized using Docker and is designed for easy deployment on cloud platforms like AWS, GCP, or Heroku. This repository includes comprehensive API documentation, test cases, and a detailed README file with database schema design, system architecture, and setup instructions.

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
