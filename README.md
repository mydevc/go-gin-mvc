# go-gin-mvc
基于go gin 框架搭建的MVC架构的基础项目空架子
此项目集成了小型网站开发常用的功能：
1、基于redis存储的session;
2、基于redis存储的cache操作；
3、基于gorm的数据库操作；
4、基于beanstalk的队列服务；
5、类php laravel框架的数据验证；
6、csrf防跨站攻击；

其它注意事项
1、队列需要单启服务，与http独立
2、依赖包用govendor管理，命令：
govendor sync vendor/vendor.json 
