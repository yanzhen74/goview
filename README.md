# goview
-----------

## 目录结构
-----------

- /private 这个目录存放的是前端打包构建的文件，主要包含模版文件，它并不对外访问，发布的时候，这个目录无需复制到服务器
- /private/router 前端路由
- /private/util 上传图片的脚本，时间转换，ajax请求脚本等辅助函数
- /public 所有的样式文件(css文件），前端脚本，网站小图标，上传的图片，打包生成的js文件都存放在这个目录
- /src go代码的存放目录
- /src/access 数据层go代码
- /src/cache 缓存层代码
- /src/controller 处理器代码
- /views go语言的视图模板

## 功能测试
-----------------
- 启动selenium服务器
    java -jar selenium-server-standalone-3.9.1.jar
- 在控制终端运行测试
    cd functional_test
    go test -v