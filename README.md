# blog_aggregator
> ## BOOT.DEV guided project 6
> * ### CLI RSS feeds collector
***

### blog_aggregator automatically connects and collects rss feeds' data and let you browse the data
### BOOT.DEV 백엔드 과정에서 진행한 간단한 Go, postgreSQL 토이 프로젝트로 사용자가 선택한 RSS 피드들을 수집, 갱신하여 출력하는 CLI 도구입니다. 
### 사용자가 직접 입력하거나 데이터베이스에서 저장되어 있는 사이트 피드 주소들을 선택하면 일정 주기마다 피드 데이터들을 내려받아 갱신하고 데이터베이스에 저장한 뒤 사용자가 원하면 출력하여 제목, 요약, 게시글url 등을 볼 수 있습니다.

***
<details>
<summary> <h2> Prerequisites </h2> </summary>
<div markdown="1">

### 1. Install go v1.24 or later
```bash
curl -sS https://webi.sh/golang | sh
```

### 2. Install Postgres v15 or later
#### 2-1. Install
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```
#### 2-2. Set a password for user postgres
```bash
sudo passwd postgres
# set a password for user postgres
```
#### 2-3. Start the Postgres server in the background
```bash
sudo service postgresql start
```
#### 2-4. Enter the `psql` shell
```bash
sudo -u postgres psql
# psql shell should show a new prompt : postgres=#
```
#### 2-5. Create a new database
```bash
# while in psql shell
CREATE DATABASE <db_name>;
# ex: CREATE DATABASE gator;
```
#### 2-6. Set the database user's password
```bash
# while in psql shell
# connect to the new database
\c <db_name>
# then psql shell should show a new prompt : <db_name>=#

# set the database user's password
ALTER USER postgres PASSWORD '<your_password>';
# this password is the one used in your connection string
``` 
### 3. Install goose and run up migrations in the project's sql/schema directory 
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
#### then download and `cd` into the project's sql/schema directory and run
```bash
goose postgres <connection_string> up
```
> #### your connection string should look like this
> ```
> "postgres://postgres:<database user's password>@localhost:5432/<database name>"
> ```
>> *Postgres' default port is :`5432`*

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

