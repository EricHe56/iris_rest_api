# iris_rest_api

本项目在iris框架基础上开发完成，用于快速开发restful接口。
两步完成开发：
1. 直接在api目录仿照user和faq添加函数
2. 在app.go中添加新的接口类名称即可build运行。
系统会通过条件编译自动生成路由器代码，开发者可以仅关注开发业务逻辑。


### 具体步骤：
1. clone 项目
2. 在api目录中仿照user.go和faq.go创建自己的api代码
3. 在app.go中添加你需要暴露的接口类，如： api.UserApi{}, api.FaqApi{},
4. 运行build.bat脚本将会自动生成(路由器代码、接口文档、测试文档和可执行程序):
    -  Build.bat包含如下命令：
    -  (1). 运行命令生成api的路由代码（同时生成-路由器代码apiRoutes.go、接口文档doc/_apiDoc.html）：
        -  go run -tags=create_router app.go emptyLoadRouteHandlers.go
    -  (2). 运行命令测试生成的api路由代码(同时生成-测试文档doc/_testDoc.html)：
        -  go test -count=1
    -  (3). 运行命令生成可执行文件：
        -  go build
        -  或直接运行测试环境（加dev参数运行，无参数运行会读取conf/prod.conf请自行配置生产环境值）
            -  go run app.go apiRoutes.go dev

### 说明：
1. 自动创建接口文档在http://localhost:8080/_doc/_apiDoc.html
2. 自动生成测试文档在http://localhost:8080/_doc/_testDoc.html
3. 配置文件在conf目录中： 端口、跨域、mongo、redis等
4. 可以将数据库相关的结构放在models/db_struct.go中，仅用于接口的结构放在models/api_struct中
5. 单元测试可以写在app_test.go中（有范例）运行： 
   - go test -count=1 即可（会自动生成测试文档）
