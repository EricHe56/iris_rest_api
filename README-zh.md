#iris_rest_api

本项目在iris框架基础上开发完成，用于快速开发restful接口。
两步完成开发：
1. 直接在api目录仿照user和faq添加函数
2. 在app.go中添加新的接口类名称即可build运行。
系统会通过条件编译自动生成路由器代码，开发者可以仅关注开发业务逻辑。


###具体步骤：
1. clone 项目
2. 修改项目目录中的apiRoutes.go.template为apiRoutes.go（后面步骤会自动生成里面的代码）
3. 在api目录中仿照user.go和faq.go创建自己的api代码
4. 在app.go中添加你需要暴露的接口类，如： api.UserApi{}, api.FaqApi{},
5. 脚本生成build.bat或者单步生成：
    -  (1). 运行命令生成api的路由代码(接口路径会显示在控制台)：
        -  go run -tags=create_router app.go apiRoutes.go
    -  (2). 运行命令生成可执行文件：
        -  go build app.go apiRoutes.go
    -  或直接运行测试环境（加dev参数运行，无参数运行会读取conf/prod.conf请自行配置生产环境值）
        -  go run  app.go apiRoutes.go dev

###说明：
1. 配置文件在conf目录中： 端口、跨域、mongo、redis等
2. 可以将数据库相关的结构放在models/db_struct.go中，仅用于接口的结构放在models/api_struct中
3. 单元测试可以写在app_test.go中（有范例）运行： go test -count=1 即可
4. 支持yaag文档生成，在运行test前将init/global.go中的BuildApiDoc改为true，然后运行程序（需要添加dev参数才会生成，为防止生产环境被执行）
生成的文件是apiDoc.html和apiDoc.html.json.会包含测试中访问的接口和收发数据。

后续会在创建路由代码的时候添加生成接口文档的功能。
