# RTUITLAB_backend

# Introduction

This is a 3 level project of purchases, shops and fabric services, created by Daniil Rozanov from IVBO-05-19 group. Each of service connected with own postgres database and one general RabbitMQ. Each service divided on 3 levels: handlers, middleware, repository, and communicate with each other by interfaces.

Actually services do this stuff:
- fabric produces products and sending it to shops service.
- shops provides to user possibility to view shops and products list and buying it, creating receipt. Receipt sends to purchases service.
- purchases can register users and give them list of their receipts.


# Install and Run services

To install services, clone this repository in any folder, go into created folder by ```cd``` and run
```
  docker-compose pull
  docker-compose up
```
Servises run by pre-compiled binary file, so if you want to do any changes and up it, you need to execute ```CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd``` in each directory of ```/shops```, ```/purchases``` and ```/factory```. Also you can use Makefile to rebuild binary main files by ```make rebuild_all_docker``` command, executed from project root directory.

Important path: before doing any changes, you need to correct ```docker-compose.yaml``` file setting docker image source from ```image: ddzzan/*``` to ```build: ./*``` or your own docker hub image. Also if you want to use command ```make``` to repush your custom image, you need to set your own docker hub destination address in each project's Makefile. 

# API documentation

## Swagger

To test all available functions, you can use swagger in purchases and shops service. To use it, go to ```http://localhost:8081/shops/swagger/index.html#/``` or ```http://localhost:8081/purchases/swagger/index.html#/``` while services is running.

## Shops service

## Auth request group
---
## POST /signin
Enter user's data and gain JWT token or error. Returns code 200 if success and 404, 500 on error.

#### Example
Request example
```json
{
  "name": "string",
  "password": "string"
}
```
response example
```json
{
  "token": "XXX.YYY.ZZZ"
}
```
## Default request group
Request if this groups doesn't need authentification.

---
## GET /products
Get products list. Returns code 200 if success and 500 on error.

#### Example
Response example
```json
{
  "data": [
    {
      "id": 2,
      "title": "Уголок для плитки",
      "description": "Уголок для плитки 9 мм МАК №25 горох, наружный",
      "cost": 300,
      "category": "Стеновые покрытия",
      "code": 972,
      "map": [
        {
          "shop_id": 2,
          "quantity": 66200
        },
        {
          "shop_id": 1,
          "quantity": 33100
        }
      ]
    },
    {
      "id": 1,
      "title": "Профиль перегородочный стоечный",
      "description": "Профиль перегородочный стоечный CW 100х50 (PR ПС 100) N",
      "cost": 300,
      "category": "Профиль укороченный",
      "code": 933,
      "map": [
        {
          "shop_id": 2,
          "quantity": 22590
        },
        {
          "shop_id": 1,
          "quantity": 14650
        }
      ]
    }
  ]
}
```

## GET /shops
Get shops list. Returns code 200 if success and 500 on error.

#### Example
Response example

```json
{
  "data": [
    {
      "id": 1,
      "title": "Тысяча мелочей",
      "address": "Самарская 134",
      "number": "+79948156309"
    },
    {
      "id": 2,
      "title": "Крепёж",
      "address": "Грибоедова 21",
      "number": "+78849137751"
    }
  ]
}
```

## Logged request group
To send this requests, you need to set Authorization field in header of HTTP request as "Bearer XXX.YYY.ZZZ", where XXX.YYY.ZZZ is your JWT token from auth.

---
## GET /carts
Get list of user's carts. Returns code 200 if success and 401, 500 on error.

#### Example
Response example
```json
{
  "data": [
    {
      "shop": {
        "id": 1,
        "title": "Тысяча мелочей",
        "address": "Самарская 134",
        "number": "+79948156309"
      },
      "products": [
        {
          "id": 1,
          "title": "Профиль перегородочный стоечный",
          "cost": 300,
          "quantity": 1,
          "entire_cost": 300,
          "category": "Профиль укороченный"
        }
      ],
      "summary_cost": 300
    }
  ]
}
````

## POST /products
Add to cart any product. Quantity and category fields is optional. Returns code 200 if success and 404, 401, 500 on error.

#### Example
Request example
```json
{
  "category":"mycat",
  "product_id": 1,
  "quantity": 0,
  "shop_id": 1
}
```
response example
```json
{
  "status": "success"
}
```

## DELETE /carts
Delete product from cart. Quantity field is optional. If quantity = 0 or ommited, delete all of this product. Returns code 200 if success and 404, 401, 500 on error.

#### Example
Request example
```json
{
  "product_id": 1,
  "quantity": 1,
  "shop_id": 1
}
```
Response type
```json
{
  "status": "success"
}
```

## POST /carts
Create receipt and returning id. Returns code 200 if success and 404, 401, 500 on error.

#### Example
Request example
```json
{
  "payoption": 1,
  "shop_id": 1
}
```
Response example
```json
{
  "id": 1
}
```

## GET /receipts
Get receipts list. Returns code 200 if success and 401, 500 on error.

#### Example
Response example
```json
{
  "data": [
    {
      "shop": {
        "id": 1,
        "title": "Тысяча мелочей",
        "address": "Самарская 134",
        "number": "+79948156309"
      },
      "products": [
        {
          "id": 1,
          "title": "Профиль перегородочный стоечный",
          "cost": 300,
          "quantity": 24,
          "entire_cost": 7200,
          "category": "Профиль укороченный"
        }
      ],
      "summary_cost": 7200,
      "pay_option": "VISA",
      "CreatedTime": "2021-03-14T13:19:49.187412Z"
    }
  ]
}
```
---

## Purchases service

## Auth request group
---
## POST /signup
Sign up. Returns code 200 if success and 404, 500 on error.
#### Example
Request example
```json
{
  "name": "string",
  "password": "string"
}
```
Resonse example
```json
{
  "id":1
}
```

## POST /signin
Sign in. Returns code 200 if success and 404, 500 on error.
#### Example
Request example
```json
{
  "name": "string",
  "password": "string"
}
```
Resonse example
```json
{
  "token":"XXX.YYY.ZZZ"
}
```

## POST /cheques
Get cheques. Returns code 200 if success and 404, 401, 500 on error.
#### Example
Response example
```json
{
  "data": [
    {
      "shop": {
        "id": 1,
        "title": "Тысяча мелочей",
        "address": "Самарская 134",
        "number": "+79948156309"
      },
      "products": [
        {
          "id": 1,
          "title": "Профиль перегородочный стоечный",
          "cost": 300,
          "quantity": 24,
          "entire_cost": 7200,
          "category": "Профиль укороченный"
        }
      ],
      "summary_cost": 7200,
      "pay_option": "VISA",
      "CreatedTime": "2021-03-14T13:19:49.187412Z"
    }
  ]
}
```
---

## Fabric service

Only you should to do is set fabric is configurate ```products.json``` file.

#### Example

Example of ```products.json``` file. ```period``` fiels is time is seconds that sets period of sending products to shops service. ```power``` field sets quantity of product that will be produced by one period. Map of ```quantity``` and ```shop_id``` fields sets how much items fabric needs to sent to every shop.

```json
{
  "products": [
    {
      "product": {
        "title": "Уголок для плитки",
        "description": "Уголок для плитки 9 мм МАК №25 горох, наружный",
        "cost": 300,
        "category": "Стеновые покрытия",
        "code": 972
      },
      "map": [
        {
          "shop_id": 2,
          "quantity": 100
        },
        {
          "shop_id": 1,
          "quantity": 50
        }
      ],
      "power": 100
    },
    {
      "product": {
        "title": "Профиль перегородочный стоечный",
        "description": "Профиль перегородочный стоечный CW 100х50 (PR ПС 100) N",
        "cost": 300,
        "category": "Профиль укороченный",
        "code": 933
      },
      "map": [
        {
          "shop_id": 2,
          "quantity": 30
        },
        {
          "shop_id": 1,
          "quantity": 20
        }
      ],
      "power": 40
    }
  ],
  "period": 7
}

```
