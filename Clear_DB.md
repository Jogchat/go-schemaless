# Clearing script migrate to 
```
https://github.com/Jogchat/schemaless_clean
```

# Connect database;
```
mysql -u root -p
Umiuni_jogchat_schemales_2018@
use jogchat0;
```
# Clearing Schemaless database commands:

```
DELETE from cell;
DELETE from index_users_id;
DELETE from index_users_username;
DELETE from index_users_email;
DELETE from index_users_phone;
DELETE from index_users_password;
DELETE from index_users_activate;
DELETE from index_users_token;
```

# Clearing companies
```
DELETE from index_companies_id;
DELETE from index_companies_category;
DELETE from index_companies_domain;
DELETE from index_companies_name;
```

# Clearing schools
```
DELETE from index_schools_id;
DELETE from index_schools_category;
DELETE from index_schools_domain;
DELETE from index_schools_name;
```

# Test users:
```
test0
password0

test1
password1

test2
password2
```

# Clear remote tables using scripts:
```
mysql -h "165.227.25.43" -u "root" "-pUmiuni_jogchat_schemales_2018@" "jogchat0" < "mysql_clean.sql"
```
```
mysql -h "138.197.103.33" -u "root" "-pUmiuni_jogchat_schemales_2018@" "jogchat1" < "mysql_clean.sql"
```
