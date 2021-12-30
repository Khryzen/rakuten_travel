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


# Reference
* [How To Install MySQL on Ubuntu 18.04](https://www.digitalocean.com/community/tutorials/how-to-install-mysql-on-ubuntu-18-04)