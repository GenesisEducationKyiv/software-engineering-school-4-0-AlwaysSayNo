# Genesis Currency API

This project implements a service with an API that allows users to:

- Get the current exchange rate of USD to UAH.
- Subscribe an email to receive updates on exchange rate changes.
- Send daily emails with the current exchange rate to all subscribed users.

## Features

- **HTTP Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: PostgreSQL for storing user data (emails).
- **Database Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **3rd Party API**: [PrivatBank API](https://api.privatbank.ua/#p24/exchange)
- **Scheduler**: [github.com/robfig/cron](https://github.com/robfig/cron) for scheduling tasks.

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Installation

1. **Clone the repository**:
    ```sh
    git clone https://github.com/AlwaysSayNo/genesis-currency-api.git
    cd genesis-currency-api
    ```

2. **Set up environment variables**:
The __./pkg/common/envs/.env__ file stores all public variables, such as:
    ```env
    APP_PORT=:8080 // the port on which the application will be launched;
    SMTP_HOST=smtp.gmail.com // the host of the mail server;
    SMTP_PORT=587 // the port of the mail server;
    EMAIL_SUBJECT="Subject: Currency Rate Update" // the subject line of the email sent to users;
    THIRD_PARTY_API="https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5" // a link to the 3rd party API for obtaining currency data.
    ```


All other variables are considered secret and must be provided as environment variables. 
In order to run the project with docker-compose, you must create the following two files in the project root: 

- __app-env-local__ - responsible the variables needed for the application to work:
    ```app-env-local
       DB_NAME=currency-api //example: a name of the database schema (for migration);
       DB_PASSWORD=root //example: database user password (for migration);
       DB_USER=root //example: database user login (for migration);
       SMTP_PASSWORD=zvai mawk wdbn aamd //example: 16 character app password; 
       SMTP_USER=test@gmail.com //example: sender gmail address.
    ```

- __db-env-local__ - responsible the variables needed for the database to work:
    ```db-env-local 
    POSTGRES_USER=root //example: a name of the database schema;
    POSTGRES_PASSWORD=root //example: database user password;
    POSTGRES_DB=currency-api //example: database user login;
    ```

3. **Build and run the Docker containers**:
    ```sh
    docker-compose up --build
    ```

## Main API Endpoints

### Get Current Exchange Rate

- **Endpoint**: `/api/rate`
- **Method**: `GET`
- **Description**: Retrieves the current exchange rate of USD to UAH.
- **Note**: from the 3rd party API we can fetch both sale and buy rates. But since we have to return one value, I've decided to return only sale.
- **Response**: number (rate)

### Subscribe Email

- **Endpoint**: `/api/subscribe`
- **Method**: `POST`
- **Description**: Subscribes an email to receive daily updates on the exchange rate.
- **Request**: "application/x-www-form-urlencoded": "email"
- **Response**: string (message)


## Utility API Endpoints Emails are successfully sent

## Get All Emails

- **Endpoint**: `/api/util/emails`
- **Method**: `GET`
- **Description**: Retrieves data about all subscribed users.
- **Response**:
    ```json
    [
      {
        "id":  "number",
        "email": "string"
      }   
    ]
    ```

## Send Emails

- **Endpoint**: `/api/util/emails/send`
- **Method**: `POST`
- **Description**: Sends emails to subscribed users.
- **Note**: from the 3rd party API we can fetch both sale and buy rates. In the email we use both buy and sale rates.
- **Response**: string


## Scheduler

- **Hourly**: The service sends a request to the PrivatBank API to fetch the latest exchange rate every hour.
- **Daily**: Emails with the current exchange rate are sent to all subscribed users every day at 9:00 AM UTC.