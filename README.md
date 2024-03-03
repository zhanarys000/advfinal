# Introduction
This project consists of a web application for managing books and user subscriptions. It includes backend services developed in Go (Golang) and a front-end interface built with Node.js and EJS.

# Creators:
- SE-2206
- Musabekova Zhansaya 
- Asan Zhanarys
- Mazhit Edige

# Components
The project is organized into several components:

## Backend (Go):
- Handlers: Handle HTTP requests, including user authentication, book management, and subscription management.
- Repository: Interacts with the database to perform CRUD operations on users, books, and subscriptions.
- Router: Defines the API routes and middleware using the Gin framework.
- Models: Structs representing user, book, subscriber, and user-book relationships.
- DB Package: Establishes and manages connections to the PostgreSQL database.
- Config: Holds configuration settings for the application, such as database credentials and server ports.
- Logger: Initializes and configures the logging functionality for the application.

## Frontend (Node.js and EJS):
- App.js: Main entry point for the Node.js application, sets up routes and views.
- Views (EJS Templates): Front-end views for displaying user interfaces.
- NPM: Manages project dependencies and script execution.

# Setup and Installation
To set up the project locally, follow these steps:

##  Backend (Go)
Ensure you have Go installed on your machine.
Navigate to the project directory and run go mod tidy to install dependencies.
Create a .env file in the root directory and configure environment variables (database connection, server ports, etc.).
Run go run . to start the Go server.

## Frontend (Node.js and EJS)
Ensure you have Node.js and npm installed on your machine.
Navigate to the front directory and run npm install to install dependencies.
Run nodemon app.js to start the Node.js server.

# Usage
Once the backend and frontend servers are running, you can access the application by visiting http://localhost:3000 in your web browser. From there, you can register/login users, manage books, and subscribe to book updates.

# Contributing
Contributions to this project are welcome. If you find any bugs or have suggestions for improvements, please open an issue or submit a pull request.


# Acknowledgements
Special thanks to the creators and maintainers of the Gin framework, Node.js, EJS, and other open-source tools used in this project. Their contributions make projects like this possible.