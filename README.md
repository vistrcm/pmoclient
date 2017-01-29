# pmoclient

very basic script to simplify work with 'PMO' tool. 

This scrip logs in to PMO, get list of engineers and prints table filtered by _filterUsers_ from config.


## configuration file
Config file located in ```~/.config/pmoclient.json```.
Sample config:
```json
{
    "loginUrl": "https://pmoserver/login",
    "peopleListUrl": "https://pmoserver/people",
    "username": "superuser",
    "password": "verylongpassword",
    "filterUsers": [
        "user55",
        "anotheruser"
    ]
}
```
