# ecommerce

This is a backend of an Ecommerce purely implemented using Golang.
The database used in MongoDB

## Approach

The code is written using the following packages:

    github.com/dgrijalva/jwt-go v3.2.0+incompatible
    github.com/gin-gonic/gin v1.8.1
    github.com/go-playground/validator/v10 v10.11.1
    go.mongodb.org/mongo-driver v1.10.2
    golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d

## How to run the program

    * to run the program we need to run the docker-compose yaml which will setup a local mongodb server in the local container. To run the it use the following command:
    `docker-compose up -d`
    This will setup the mongo server in the standard 27017 port
    * To run the main program just run the main file using the command
    `go run main.go`
    This will run the Program in `8000` port
    * use Postman app to test the commands with the required payloads

# Endpoints

    * Sign Up (Post request)
    URL: `http://localhost:8000/users/signup`

    ```
    {
    "first_name": "aflah",
    "last_name": "ahamed",
    "email": "aflah@test.com",
    "password": "aflahaflah",
    "phone": "+91740654657"
    }
    ```

    * Log In (Post request)
    URL:`http://localhost:8000/users/login`

    ```
    {
    "email":"aflah@test.com",
    "password":"aflahaflah"
    }
    ```

    * Add Product (POst Request)
    URL:`http://localhost:8000/admin/addproduct`

    ```
    {
    "product_name": "ASUS laptop",
    "price": 50000,
    "rating": 7,
    "image": "asus.png"
    }
    ```

    * All products view (Get request)
    URL:`http://localhost:8000/users/productview`

    * Regex (Get request)
    URL:`http://localhost:8000/users/search?name=asus`

    * Add to cart (Get request)
    URL:`http://localhost:8000/addtocart?id=6345ba523cacbf5b04ebde4c&userID=63445d653066839ebb492085`

    *Remove from cart (Get request)
    URL:`http://localhost:8000/removeitem?id=63446135e8156f9f963eecbc&userID=63445d653066839ebb492085`

    * Add adress (Post request)
    URL:`http://localhost:8000/addaddress?id=63445d653066839ebb492085`

    ```
    {
    "house_name": "black house",
    "street_name": "black street",
    "city_name": "mangalore",
    "pin_code": "54564186"
    }
    ```

    * Edit work address (Put request)
    URL:`http://localhost:8000/editworkaddress?id=63445d653066839ebb492085`

    ```
    {
    "house_name": "blue house",
    "street_name": "lol street",
    "city_name": "lalalla",
    "pin_code": "58568"
    }
    ```

    * Edit Home address (Put request)
    URL:`http://localhost:8000/editworkaddress?id=63445d653066839ebb492085`

    ```
    {
    "house_name": "blue house",
    "street_name": "lol street",
    "city_name": "lalalla",
    "pin_code": "58568"
    }
    ```

    * Delete both address (Get request)
    URL:`http://localhost:8000/deleteaddresses?id=63445d653066839ebb492085`

    * Cart checkout (Get request)
    URL:`http://localhost:8000/cartcheckout?id=63445d653066839ebb492085`

    * Instant checkout (Getrequest)
    URL:`http://localhost:8000/instantbuy?userID=63445d653066839ebb492085&id=6345ba523cacbf5b04ebde4c`
