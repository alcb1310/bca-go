# Budget Control Application

## Tech Stack

- Go [![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](http://go.dev)
- PostgreSQL [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
- HTMX [https://htmx.org/](https://htmx.org/)

## Socials

- [![Twitter Badge](https://img.shields.io/twitter/follow/username.svg?style=social&label=Follow)](https://twitter.com/alcb1310)
- [![Github Badg](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/alcb1310)

## Description

The objective of this application is to manage the budget of a construction company, to achieve this it will be able to:

- Manage Suppliers
- Manage Budget Items
- Create Budgets for each project
- Create Invoices
- Each invoice will decrease the budgeted data
- Several reports to be defined

## Table of contents

- [Screens](#screens)
    - [Login](#login)
    - [Home](#home)
    - [Users](#users)

## Screens

In order to achieve the project's description, the application will have both public and protected Endpoints

### Login

The application will start at the `/login` route which allows a user to login to it by providing their email and password, then the
server will validate their credentials, and on success it will go to the protected routes and if it didn't succeed, it will display
a message indicating `invalid credentials`

### Home

All the application will reside inside the `/bca/` route, and display the menu at the left pane and on its top displays an exit button,
which will logout the user and in its main pane it will display a greeting to the logged in user

### Users

Within the **Usuarios** menu, there will be 2 options

- `Editar usuario`
This option will allow a user with admin privileges to create, update and delete users from the application

- `Cambiar contraseaña`
This option will allow any user to change their password to access de application
