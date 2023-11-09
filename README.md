# Budget Control Application

## Tech Stack

- [![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](http://go.dev)
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
- [https://htmx.org/](https://htmx.org/)
- [![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)](https://jwt.io)

## Socials

[![Twitter Badge](https://img.shields.io/twitter/follow/username.svg?style=social&label=Follow)](https://twitter.com/alcb1310)
[![Github Badg](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/alcb1310)

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
    - [Settings](#settings)

- [Deployment](#deployment)

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

- `Cambiar contraseña`
This option will allow any user to change their password to access de application

### Settings

Within the **Parámetros** menu, there will be 3 options

- `Proyectos`
This option will allow a user to create and update projects, deletion of a project will not be allowed

- `Proveedores`
This option will allow a user to create and update suppliers, deletion of a supplier will not be allowed

## Deployment

In order to be able to deploy this application, the following is needed:

1. Clone the repository using the following command

```bash
git clone https://github.com/alcb1310/bca-go-w-test.git
```

2. Download all of the project's dependencies by running 

```bash
go mod tidy
```

3. At the root of the project directory, create a `.env` file with the following fields:

```.env
PORT=<Port number where the application will listen>

PGDATABASE=<Name of the postgres database>
PGHOST=<Host where the postgres database server is running>
PGPASSWORD=<Password that the postgres server uses>
PGPORT=<Port where the postgres server is listening>
PGUSER=<Username of the postgres server>

SECRET=<Secret used to generate the JWT Token>
```

4. To start the deployed version run the following command

```bash
go run main.go
```

