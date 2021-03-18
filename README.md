# processor-webapp
Processor Web App GitHub Repository Setup

Test locally:
```.env
export DBHostName=rds
export DBUserName=username
export DBPassword=pwd
export DBSchema=db
export KafkaServer=localhost
export ESAddr=http://localhost:9200
```

## processor-webapp with Docker

1. Build the image

    ```
    docker build -t processor-webapp:1.0 .
    ```

2. Run processor-webapp in local

   Make sure you have MySQL database set up in your local.
   You will need to pass in your MySQL `DBHostName`, `DBUserName`, `DBPassword`, `DBSchema`, `KafkaServer` and `ESAddr`

    ```
    docker run -it \
        -e DBHostName=endpoint \
        -e DBUserName=usrname \
        -e DBPassword=pwd \
        -e DBSchema=stories \
        -e KafkaServer=kafka \
        -e ESAddr=hostname:port \
        processor-webapp:1.0
    ```

    