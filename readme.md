# Github / Mail / Google auth example with Postgresql, Redis and Memcached
![image](https://user-images.githubusercontent.com/61962654/169509440-485a95f0-e2de-4d79-b7a6-bf953536c634.png)
![image](https://user-images.githubusercontent.com/61962654/169366937-a5472d37-2c9f-463e-8193-64824b3938b6.png)
![image](https://user-images.githubusercontent.com/61962654/169367192-9967579b-b7c9-47d8-bf88-434ddd48e190.png)



## Requirements

#### - Go: `1.18`
#### - PostgreSQL: `14.2`
#### - Memcached: `1.6.15`
#### - Redis: `6.2.6`
#### - `Google` and `Github` keys

---

## Setup OAuth

### - Github
1. Go to the **[Developer settings](https://github.com/settings/apps)**
2. Create **[Application](https://github.com/settings/apps)**
3. Enable `User permissions` -> `Email addresses` -> `Read Only` in the **[Permissions](https://github.com/settings/apps/permissions)**
4. Generate **secret token**
5. Enter the `URIs` that are allowed to be redirect-URIs (e.g. `https://localhost:8080/oauth/github/callback`)
6. Paste both `Client ID` and `Client Secret` to the `github_secret.json`

### - Google
1. Go to the **[Google Cloud Console](https://console.cloud.google.com/projectselector2/apis/credentials)**
2. Create project (add content to the consent screen like title and logo) or use existing
3. `Credentials` -> `Create Credentials` -> `Oauth Client ID`
4. Choose the `Web Application` type and give it a name
5. Enter the `URIs` that are allowed to be redirect-URIs (e.g. `https://localhost:8080/oauth/google/callback`)
6. Paste both `Client ID` and `Client Secret` to the `google_secret.json`

---

## docker-compose
Server is ready immediately after containers are up
```shell
docker-compose up
```

It is possible to additionally configure the app using environment variables
```yaml
environment:
  POSTGRES_IP: 127.0.0.1 # connect to local database
  HOST_PORT: 8082 # change server port
```

---

## Setup PostgreSQL
```shell
migrate -database ${POSTGRESQL_URL} -path migrate/ up
```
### Down
```shell
migrate -database ${POSTGRESQL_URL} -path migrate/ down
```

## Build / Run

```shell
git clone https://github.com/illiafox/go-auth.git auth
cd auth

make build
make run # /cmd/server/bin
```

### Run arguments

#### HTTP mode

```shell
server -http
```


#### With non-standard config and log file paths

```shell
server -config config.toml -log log.txt
```

#### With reading from `environment`:

Available keys can be found in **[config structure](https://github.com/illiafox/go-auth/blob/master/utils/config/struct.go)** tags

```shell
POSTGRES_PORT=4585 server -env
```

---

## Logs
In addition to the terminal output, logs are also written to the file
```shell
# Terminal
20/05/2022 10:50:20 |   info    Initializing repository
20/05/2022 10:50:20 |   info    Done    {"time": 0.012092004}
20/05/2022 10:50:20 |   info    Server started at 0.0.0.0:8080

```

```json5
// File (default log.txt)
{"level":"info","ts":"Sun, 01 May 2022 14:54:38 EEST","msg":"Initializing database"}
```
---

## Endpoints

### `/` Main Page

### `/register` Register

### `/login` Login

### `/logout` Logout

### `/verify` Mail verify

#### `/oauth/github/login`  `/oauth/github/callback` Github OAuth

#### `/oauth/google/login`  `/oauth/google/callback` Google OAuth
