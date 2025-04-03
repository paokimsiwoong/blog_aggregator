# blog_aggregator
> ## BOOT.DEV guided project 6
> * ### CLI RSS feeds collector
***

### blog_aggregator automatically connects and collects rss feeds' data and let you browse the data  

***
<details>
<summary> <h2> Prerequisites </h2> </summary>
<div markdown="1">

### 1. Install go v1.24 or later
```bash
curl -sS https://webi.sh/golang | sh
```

### 2. Install Postgres v15 or later
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```
#### then 
```bash
sudo passwd postgres
# set a password for user postgres
```
#### Finally, start the Postgres server in the background
```bash
sudo service postgresql start
```
</div>
</details>

***

> ## How to Install
```bash
go install https://github.com/paokimsiwoong/blog_aggregator@latest
```
### then create .gatorconfig.json file at your $HOME directory 
### the json's `db_url` field must contain the connection string to your sql database.
```json
{
    "db_url": "postgres://<username>:<password>@localhost:5432/<dbname>?sslmode=disable"
}
```

> ## How to Use
```bash
blog_aggregator <command name> \[<argument>\]
```
### commands:
* ### `register <userName>`
> #### register and login as `<userName>`
* ### `login <userName>`
> #### login as `<userName>`
* ### `users`
> #### list stored users
* ### `addfeed <feedName> <feedURL>`
> #### save feed url and follow the feed
* ### `feeds`
> #### list stored feeds
* ### `follow <feedURL>`
> #### (current user) follow a feed. the feed must have been stored beforehand
* ### `unfollow <feedURL>`
> #### (current user) unfollow a feed. the feed must have been stored beforehand
* ### `following`
> #### list feeds current user following
* ### `agg <timeBetweenReqs>`
> #### connect and collect feeds' data every `<timeBetweenReqs>`. Ctrl+C to abort 
* ### `browse [<numPosts>]`
> #### list posts from feeds current user following
* ### `reset`
> #### reset database

