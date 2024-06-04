
# goravel-demo-app

This is a demo app for the [Goravel](https://github.com/saalikmubeen/goravel) web framework.



![Goravel Demo App](/screenshots/browser.png)



## Features demonstrated in this demo app:

- Routing
- Handlers (For Laravel's Controllers)
- Middlewares
- Views
- Database (PostgreSQL)
- Migrations
- Models
- Upper/db ORM
- Validation
- Session Management
- User Authentication
- Cache Management
- How to send JSON and XML responses using Goravel
- Email Sending
- Password Reset
- Remember Me functionality using Cookies
- API Routes
- Crud Operations



## Installation

1. Clone the repository:

   ```bash
   git clone git@github.com:saalikmubeen/goravel-demo-app.git

    cd goravel-demo-app
    ```

2. Install the dependencies:

   ```bash
   go mod tidy
   ```


3. Fill the .env file with your database credentials:



4. Install the [Goravel Command Line Tool]:

   ```bash
   go install github.com/saalikmubeen/goravel/cmd/goravel@latest
   ```

Make sure you have the `$GOPATH/bin` directory in your `PATH` environment variable. If you don't have it, you can add it to your `~/.bashrc` or `~/.bash_profile` or `~/.zshrc` file:

```bash
export GOPATH="$HOME/go"
export PATH=$PATH:$GOPATH/bin
```

5. Run the migrations (from the root of the project)

   ```bash
    goravel migrate up
  ```

6. Start the server:

   ```bash
   go run ./*.go
   ```

7. Visit `http://localhost:4000` in your browser.