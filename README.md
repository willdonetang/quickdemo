目录结构
├── conf 
│   └── app.ini       //配置相关
├── crontab
│   └── crontab.go    //定时任务(支持crotab模式)
├── docs              //文档(后续可以加mysql的集成)
│   └── docs.go       //swagger文档集成(自动生成不需要手动改)
├── models
|   ├── ****models.go  //具体的model,
│   └── models.go     //数据库层，于具体的数据库打交道,setup可以具体设置callback以及日志记录
├── pkg               //第三方以及常用中间件配置
|   ├── app  //具体的model
|   |   ├── form.go         //具体验签的中间件
|   |   └── response.go     //统一返回结构体 code:返回码、data：具体数据、msg:错误信息、errMsg:前端调试信息
│   ├── e     //code以及msg
|   |   ├── code.go         //具体的code码维护
|   |   └── msg.go         //code码对应的msg
│   ├── errors     //自定义错误
|   |   └── errors.go         //具体自定义的接口
│   ├── file     //文件以及文件夹的常用方法集合
|   |   └── file.go         //文件以及文件夹的常用方法集合
│   ├── gredis   //redis的集成
|   |   └── redis.go         //redis的集成实现
│   ├── logf   //日志库集成
|   |   ├── logf.go         //日志库集成
|   |   └── file.go         //日志库集成
│   ├── setting   //配置集成
|   |   └── setting.go         //配置集成
│   ├── util   //工具包
|   |   ├── const              //常用的const常量
|   ├──    md5.go         //md5加密工具
|   ├──    pagination.go         //分页工具
├── routers
|   └──router.go         //路由，以及路由分组(具体的路由，请在具体的文件实现如routers/api/v1/user.go)
├── request            //请求体(考虑是否可以跟response合并）
├── └──user.go(举例)            //user分组的请求体以及返回体
├── response            //返回体
├── runtime             //运行时(可忽略)
├── └──log         //运行日志
├── validation(可以考虑移到pkg里面)     //结构体验签以及自定义验签库
.gitignore
Dockerfile     //docker部署
main.go        //程序入口

本项目由go语言生成的web项目

其中框架为gin框架.github地址为:https://github.com/gin-gonic/gin

常用规范
1.前后端文档维护为swagger文档,其中swag由代码自动生成，不需要人工维护
2.前后端请求方式统一为get和post两种方式
get的应用场景为获取、查询、下载等操作
post为更新、增加、删除等改变数据操作

get请求，参数统一放置为url后缀，如www.*****.com/***/***?x=1&y=2&z=3
post请求，参数统一放置为postbody,raw下，并且设置为json格式，最好添加header：Content-Type：application/json
如
{
	"id":1,
	"key":"key",
	"value":"value"
}
如有其他请求方式，需前后端谨慎商议决定
3.后端人员需遵守指定的mysql规范(mysql规范，请参考文档)，并建立好相关的表结构
4.前后端统一返回体
{
	"code":200,
	"data":interface{},  //具体前端需要的数据
	"msg":"信息"，
	"errmsg":错误调试信息
}


5.开发流程
新增接口流程为:
评审需求-->给定接口完成时间-->编写swagger文档-->编写代码-->自测
6.数据库生成model代码，具体现有的模版引擎长期维护
具体的代码由后端开发维护



常用文件讲解
1.初始化
1.1配置初始化
目前配置文件基于ini来配置
然后将配置读取到内存中，
注意app.ini在在gitignore中，本地启动请用自己的配置文件
后续优化点:
目前是ini的配置，后续可以用viper，支持多种配置格式
目前配置是固定的，改配置只能重新启动程序，后续可以用viper或者air去动态获取，实现动态修改
1.2日志初始化
目前日志库为logrus,github地址:https://github.com/sirupsen/logrus日志初始化由两部分构成
1.2.1.常规日志，
输出方式为json(可以修改为text)，需要注意目前日志输出重定向两个，屏幕和文件，如果有一些定制的需求(如错误日志单独存放、日志kafaka通道）请使用logrus--hook
1.2.2mysql日志
输出方式为json(可以修改为text)，同样属于常规日志，也可以通过hook去定制化实现
1.2.3日志等级请调用logf的常用方法，显示等级由配置文件决定
1.3model的初始化
model由gorm库实现github地址:https://github.com/go-gorm/gorm，常见gorm用户请参考具体文档
目前gorm的代码可以由代码生成数据库(automirage)，也可以用sql生成代码
1.3.1后续gorm的sql语句，简单的由gorm语句生成，复杂语句请使用原生sql
1.4定时任务
采用定时任务库cron，github地址:https://github.com/robfig/cron
目前支持crontab形式以及字符串形式
1.5验签库
目前gin中采用最多的是validator库，github地址:https://github.com/go-playground/validator
目前已经在gin中间件已经用了validator库,只需要添加中间件即可如
api.GET("/myplace", Validation(&request.Myplace{}), My.place)
只需要在具体的Myplace的tag中去添加相应的标签(具体标签请参考具体文档)
type Myplace struct {
	FinishedTime string `form:"finishedTime" json:"finishedTime,omitempty" binding:"required"`
	AwaitTrain   int    `form:"awaitTrain" json:"awaitTrain,omitempty" binding:"required,oneof=1 2 3 4"`
}
如果要实现自定义验证器，请注册标签并使用
如:
自定义标签
v.RegisterValidation("timeValidated", timeValidated)

//时间格式校验
func timeValidated(fl validator.FieldLevel) bool {
	if timeString, ok := fl.Field().Interface().(string); ok {
		if timeString != "" {
			_, err := time.Parse("2006-01-02 15:04:05", timeString)
			return err == nil
		}
	}
	return true
}

2.路由
路由即使http请求的路由，需要主要的是常用中间件，并且使用context作为路由的上下文
中间件(token验签、字段验证、跨域、日志)
context，用context来链接一次请求的完整体，

3.统一返回体
请使用pkg的app模块