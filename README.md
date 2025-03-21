# Blog Aggregator
Blog Aggregator is, as the name implies, an RSS feed aggregator written in Go. With Blog Aggregator, you can:
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of aggregated posts in the terminal, with a link to the full post.

## Requirements
You need to have installed the Go toolchain. For this project, I used 1.22.2. [Click here](https://go.dev/doc/install) for information on how to install Go on Linux.
You also need to have install PostgreSQL. For this project, I used 16.8. I ran `sudo apt install postgresql postgresql-contrib` in the terminal to install it. After you have installed it, run `sudo passwd postgres`  to update the password for the user `postgres`. Make sure to remember this password! 
Then, run `sudo service postgresql start` to start the server. I'm using the `psql` client to connect to the server. To connect to the server, I run `sudo -u postgres psql`. Now you should be in the database server!

## Installation
To install Blog Aggregator, download the source code to a specified folder. In your terminal, go to the specified folder and run `go install`. This will install the program and make it accessible in your terminal if you start your command with `blog-aggregator`.

## Config
I'm using a JSON file to keep track of the configuration. It has two important details:
1. Who is currently logged in
2. The connection credentials for the PostgreSQL database.

To ensure the configuration is correct, go into your home directory and create a file called `.gatorconfig.json`. Next, you need to find the connection link to your PostgreSQL database. It will have the following structure:
```
"postgres://<username>:<password>@localhost:5432/gator"
```
In this case, `<username>` is `postgres` if you ran the commands above, and `<password>` is whatever password you enter after running `sudo passwd postgres`.  

If you have your connection link in `.gatorconfig.json`, you should be all set!

## Commands
Blog Aggregator supports the following commands:
- `register <name>` lets you register a new user to the database, and will automatically log this user in.
- `login <name>` lets you log in to a specified username, as long as this user has been registered before.
- `users` lets you view all the users that are registered in the database.
- `addfeed <feedname> <url>` lets the current user add a new feed to the database using the provided url.
- `feeds` lets you view all the feeds that are registered in the database.
- `follow <url>` lets the current user follow the feed using the provided url.
- `unfollow <url>` lets the current user unfollow a feed using the provided url.
- `following` lets you view all the feeds all users are following.
- `agg <time_between_requests>` starts a loop that will continuously scrape posts from the registered feeds, with a provided delay.
- `browse <limit>` lets you view the posts of the feeds that the current users is following. The default limit is 2.
- `reset` resets the database.
