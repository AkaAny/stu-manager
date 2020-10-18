# stu-manager
A simple student and dormitory CRUD impled  by Go

## Getting Started
1.  起一个MySQL容器，命令如下：
```
docker run --name mysql-stu-manager -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 -d mysql
```
2.  直接在GoLand里编译就完事了，默认端口8080
请求示例：[导入PostMan](StuManager.postman_collection.json)

## 出了bug怎么办
直接Q我就完事了

## 代码使用
用之前在麻烦注释里标注一下（连session都没上，鉴权都没有的CRUD不会真的有人用吧）
