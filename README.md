# SOCI API - Now in Go!

This small program runs the backend API for the SCOI Web App.

# Requirements

You need to have a few environment variables set in order for the go app to run. If you don't have these set, it'll make a few (sensible?) default assumptions:

```
APP_PORT=9000
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=socidb
DB_USER=dbuser
DB_PASSWORD=password
```

There is one environment variable that you'll need to set that doesn't have a default. `APP_KEY` is what is used to sign the JWT tokens, and it should be a nice random string between 32 and 64 chars. If you ever change this, the JWTs that are signed with one key won't be readable with another key. So, on my local computer, I start the linux binary with this command from the project root:

`APP_KEY=asdf DB_USER=root DB_PASSWORD=secret dist/socid`

Next, you'll probably need to make sure that database exists and that your db user has the correct permissions to work with that database.

If so, you can get the app database up to date by running the migrations inside the `migrations/` folder. You'll need to have goose (https://github.com/pressly/goose) added to your path.

If that's all good to go, go ahead and run this command from inside the migrations folder:

```
goose mysql "dbuser:dbpass@tcp(dbhost:dbport)/dbname" up
```

Here's what I run on my local dev machine:

```
goose mysql "root:password@tcp(127.0.0.1:3306)/socidb" up
```

## Building the app

If all is well up to this point, you can build the binary. The included `build.sh` bash script will try and build the go code and place two different binary files into the `dist/` folder. There's one for linux, and one for OSX, so jump into that folder and run whichever one makes sense for you.
