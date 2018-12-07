# pmoclient

[![Build Status](https://cloud.drone.io/api/badges/vistrcm/pmoclient/status.svg)](https://cloud.drone.io/vistrcm/pmoclient)

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
    "Spreadsheet": {
        "SpreadsheetID": "spreadsheet which store information of users",
        "SecretFile": "secrets file from google cloud"
    },
    "filterUsers": [
        "user55",
        "anotheruser"
    ]
}
```

## Command line options
pmoclient has single command line option ```-spreadsheet``. If specified tool will be using Spreadsheet to get list of users to filter and will update 'AutofillFromPMO' sheet in this document.
