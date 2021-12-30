# Software
| Name           | Version                 |
| -------------- | ----------------------- |
| Go Lang        | 1.15                    |
| MySQL          | 14.14 Distrib 5.7.36    |

# Installation of GO Lang
```bash
dev@dev: wget https://golang.org/dl/go1.15.linux-amd64.tar.gz
dev@dev: sudo tar -C /usr/local -xzf go1.15.linux-amd64.tar.gz
```

Go Configuration in **.profile**
```bash
dev@dev: sudo vim .profile


# GO configuration
PATH="$HOME/bin:$HOME/.local/bin:$PATH"
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

# Installation of MySQL

### Step 1 — Installing MySQL

Update the package index on your server with apt.
```bash
dev@dev:~$ sudo apt update
```

Install the default package.
```bash
dev@dev:~$ sudo apt install mysql-server
```
### Step 2 — Configuring MySQL

Run the security script.
```bash
dev@dev:~$ sudo mysql_secure_installation
```

### Step 3 — (Optional) Adjusting User Authentication and Privileges

In order to use a password to connect to MySQL as root, you will need to switch its authentication method from `auth_socket` to `mysql_native_password`.

```bash
dev@dev:~$ sudo mysql

mysql > SELECT user,authentication_string,plugin,host FROM mysql.user;
mysql > ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'password';
mysql > FLUSH PRIVILEGES;
```

# How to run the code?

```bash
# Linux
dev@dev:~/rakuten_travel$ clear; go build; ./rakuten_travel

# Windows
C:\workspace\go\rakuten_travel> cls && go build && .\rakuten_travel
```


# Reference
* [How to set up Go for Windows](https://www.freecodecamp.org/news/setting-up-go-programming-language-on-windows-f02c8c14e2f/)
* [How To Install MySQL on Ubuntu 18.04](https://www.digitalocean.com/community/tutorials/how-to-install-mysql-on-ubuntu-18-04)