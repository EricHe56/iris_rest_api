# iris_rest_api

English | [简体中文](./README-zh.md)

This project is developed on the basis of the iris framework for rapid development of restful interfaces.
Complete development in two steps:
1. Add functions directly in the api directory following user and faq
2. Add a new interface class name in app.go to build and run.
The system will automatically generate router code through conditional compilation, and developers can focus only on developing business logic.


### Specific steps:
1. Clone project
2. Create your own api code in the api directory following user.go and faq.go
3. Add the interface classes you need to expose in app.go, such as: api.UserApi {}, api.FaqApi {},
4. Running the build.bat script will automatically generate (router code, interface documentation, test documentation, and executable programs):
   -  Build.bat contains the following commands:
   -  (1). Run the command to generate the routing code of the api (also generate-router code apiRoutes.go, interface document ./doc/_apiDoc.html):
        - go run -tags=create_router app.go emptyLoadRouteHandlers.go
   -  (2). Run the command to test the generated api routing code (simultaneously generate-test document ./doc/_testDoc.html):
        - go test -count=1
   -  (3). Run the command to generate the executable file:
        - go build
        - Or directly run the test environment (add dev parameters to run, run without parameters will read ./conf/prod.conf, please configure the production environment value)
            - go run app.go apiRoutes.go dev

### Explanation:
1. Automatically create interface documentation at http://localhost:8080/_doc/_apiDoc.html
2. Automatically generate test documentation at http://localhost:8080/_doc/_testDoc.html
3. The configuration file is in the conf directory: port, cross domain, mongo, redis, etc.
4. Database related structures can be placed in ./models/db_struct.go, and structures used only for interfaces are placed in ./models/api_struct
5. Unit tests can be written in app_test.go (with examples) and run: 
    - go test -count = 1 (the test document will be automatically generated)


