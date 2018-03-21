## Spreader
This is just a proof of concept of a tool that reads the binlog of your local MySQL database and publish them as events to Amazon Kinesis.
The data schema of the events is the following: 
```json
{
    "Database": "spreader",
    "Table": "user",
    "Type": "insert",  
    "OldValues": null,  
    "NewValues": {     
        "age": 20,
        "id": 18,
        "name": "testUser"
    }
}
```
* `"Type"`: Contains the type of the event: `{insert, update, delete}`
* `"OldValues"`: Contains the old row values. Only when the event type is `{update, delete}`
* `"NewValues"`: Contains the new row values. Only when the event type is `{update, insert}`

The repository contains the `Spreader` application, and two consumer examples: one in Go and another in Java. 
As you can see, everything is pretty basic, but it will grow little by little.

The goal is to be able to declare `filters`, `conditions`, and `transformers`. This way you can choose which events to publish, set some
correlations between them, and define the final schema of the data sent to Kinesis.

#### Prerequisites

1.- In the global mysql configuration file `/etc/my.cnf` or the `.my.cnf` file in your home directory, add the following lines:
```sh
[mysqld]
server-id=master
binlog-format=ROW
log-bin=~/.mysql-binlogs/binlog
```
2.- Put your AWS credentials in a file located in `~/.aws/credentials`. The contents of the file should be:
```sh
[default]
aws_access_key_id = <Your key>
aws_secret_access_key = <Your secret>
```
