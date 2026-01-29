***Mirocservice-productscatalog***

This document outlines the system design for a scalable


**Prerequisites**
* Docker
* Go

**How to test?**

Run application 

1. if you want to create a mysql database for test run docker-compo


*Architecture*

This microservice follows a Hexagonal Architecture (also known as Ports and Adapters) to ensure a clear separation of concerns and make the application independent of external services.
Domain Layer: This layer contains the core business logic and entities.

Application Layer (Services): This layer contains the use cases and interfaces that define the input and output ports of the system.
Infrastructure Layer: This layer holds the implementations of the output ports, such as the Product and Order repository. This is also where you would add implementations for message queues, middlewares, and other external integrations.


1. *Business Logic*

The core of the system revolves around products that can be purchased through an order.  It also provides for a CRUD that applies to products.
It is considered that the service will be part of a core ecommerce, this indicates that the system will be of heavy reading, with a large number of users analyzing products and a smaller percentage making purchase orders. 
How is a test case only contemplates the purchase of a product by order (subject to be discussed in the technical interview).

2. *Requirements*
   * Product CRUD: Users can create, read, update and delete products
   * Create Order: Users can simulate a purchase of a product. Specifying the amount they want and viewing the total money in the order
   * Get Order: Users can view all the orders generated previously 

2.1 *Functional Requirements*

   * Product CRUD: Users can CREATE, UPDATE, DELETE AND FETCH products.
   * Create Order: Users can generate purchase orders.
   * Get Orders: Users can view a timeline displaying orders.
   * Authentication Assumption: All users are considered valid. There is no need to implement a sign-in module or session management.

2.2 *Non-Functional Requirements*

   * The solution must be able to scale up to hundreds of orders generated simultaneously.
   * The application must be optimized for reads.


3. **Data Model**

*Relational Model (MySQL)*

*Product Table*
* id (uuid, v4)
* description (string)
* price (float64)
* stock (int)
* created_at TIMESTAMP (date),
* updated_at TIMESTAMP (date)



*Order Table*
* id (uuid, v4)
* product_id (uuid, v4)
* quantity (int)
* total (int)
* date (date)



4. *API Endpoint Design*

*GET*

/api/products

Request:
GET /api/products?limit=10

Success Response:
Code: 200 

Content:
[
{
"id": "f4691a93-f2c0-4480-8172-39f5a9b0105e",
"name": "Gopher",
"price": 12.21,
"stock": 50,
"created_at": "2023-09-24T15:30:00Z"
},
{
"id": "a8b9c123-d456-7890-1234-56789abcdef0",
"name": "Gadget",
"price": 23.50,
"stock": 20,
"created_at": "2023-09-20T12:00:00Z"
}
]

Response Code Errors:
400	Bad Request
500	Internal Server Error


*POST*

/api/products

Request Body:
{
"name": "Nuevo Producto",
"price": 19.99,
"stock": 10
}
* Success Response:
Código: 201 Created

* Response Code Errors:
400	Bad Request
500	Internal Server Error


*DELETE* 
/api/products/:id

* Success Response:
204 No Content
* 
* Response Code Errors:
404	Not Found
500	Internal Server Error


*PUT*
/api/products/:id

* Request Body:
{
"name": "Gopher Pro",
"price": 15.50,
"stock": 40
}
* 
* Success Response:
Code: 204 Not Content

* Response Code Errors:
400	Bad Request
404	Not Found
500	Internal Server Error

  

6. *Next Iterations & Discution Points:*

*Decouple services:* Service responsibilities can be separated by implementing REST or gRPC flames applying patterns such as Circuit Breaker

*Caching:* A cache could be implemented to temporality store requested resources like orders or products. Reducing the number
of database queries 

*Metrics (CPU, RAM, Request per Second, Go Routines):* A metrics collection system, such as Prometheus, could be integrated to monitor
the microservice's performance, including CPU usage, RAM, and request per second.

*Author:*
Lautaro Olmedo

