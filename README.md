# Budget Control Application

## Tech Stack

- Go [http://go.dev/](http://go.dev)
- PostgreSQL [https://www.postgresql.org/](https://www.postgresql.org/)
- HTMX [https://htmx.org/](https://htmx.org/)

## Description

The objective of this application is to manage the budget of a construction company, to achieve this it will be able to:

- Manage Suppliers
- Manage Budget Items
- Create Budgets for each project
- Create Invoices
- Each invoice will decrease the budgeted data
- Several reports to be defined

## Endpoints

In order to achieve the project's description, the application will have both public and protected Endpoints

### Public Endpoints

***/login***<br />
This endpoint will allow you to access all the protected routes, as a user you will have to provide your credentials (email and password) and if authorized you will be able to read the protected routes.

### Protected Endpoints

***/api/v1/***<br />
This endpoint will show the logged in application home page, with a collapsable menu in the left side of the screen.  All interactions should be displayed on the content part of the HTML
for navigating and interacting with the application, we will be using HTMX

