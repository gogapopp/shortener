# go-url-shortener

how to build application: docker-compose up  

POSTMAN request example:  
 - POST http://localhost:8080  
    BODY: https://www.notion.so/  
    RESPONSE: http://localhost:8080/geoY7c  

 - GET http://localhost:8080/geoY7c  
    RESPONSE: regirecting to the https://www.notion.so/  

 - POST http://localhost:8080/api/shorten  
    BODY:  
    {  
        "url":"https://www.notion.so/"  
    }  
    RESPONSE:  
    {  
        "result": "http://localhost:8080/wA2s8h"  
    }  

 - GET http://localhost:8080/ping  
    RESPONSE: http status OK(200)  

 - POST http://localhost:8080/api/shorten/batch  
    BODY:  
    [  
        {  
            "correlation_id": "abc",  
            "original_url": "https://yandex.ru/"  
        },  
        {  
            "correlation_id": "a2c",  
            "original_url": "https://notion.ru/"  
        }  
    ]  
    RESPONSE:  
    [   
        {  
            "correlation_id": "abc",  
            "short_url": "http://localhost:8080/bog8IN"  
        },  
        {  
            "correlation_id": "a2c",  
            "short_url": "http://localhost:8080/gkOINk"  
        }  
    ]  

 - GET http://localhost:8080/api/user/urls  
    RESPONSE:  
    [  
        {  
        "original_url": "https://yandex.ru/",  
        "short_url": "http://localhost:8080/bog8IN"  
        },  
        {  
            "original_url": "https://notion.ru/",  
            "short_url": "http://localhost:8080/gkOINk"  
        }  
    ]  

 - DELETE http://localhost:8080/api/user/urls  
    BODY:  
    [  
        "bog8IN",  
        "/gkOINk",  
    ]  
    RESPONSE: http status Accepted(202)  

the project can be started using env, flags and config.json  
the settings are set as follows: env, if not set, then the flags and then config.json  
 -с config.json file path | env:"CONFIG"  
 -a address and port to run server | env:"SERVER_ADDRESS  
 -b base url | env:"BASE_URL  
 -f file storage path | env:"FILE_STORAGE_PATH  
 -d database path | env:"DATABASE_DSN  
 -s https enable (true/false) | env:"ENABLE_HTTPS  
 -t trusted subnet | env:"TRUSTED_SUBNET  