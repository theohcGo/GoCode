# GoFrame Template For SingleRepo

Quick Start: 
- https://goframe.org/pages/viewpage.action?pageId=1114399

# 开发流程

dao   : 数据访问。和底层数据库交互，仅仅包含最基础的CRUD方法 gf gen dao
api   : 对外接口
model : 结构模型。数据结构管理模块。管理数据实体对象。


GoFrame V2开发流程：
1. 设计表结构，初始化项目，修改配置文件
2. 使用gf gen dao 生成对应的dao/do/model
3. 编写api层：定义业务侧数据结构，提供对外接口的输入/输出数据结构
4. 编写model层，定义数据侧数据结构，提供对内的数据处理的输入/输出数据结构
5. 编写logic层，自动生成service层代码
6. 在service层代码生成RegisterXX()方法后，在对应的logic模块注册服务
7. 编写controller层，接受/解析用户输入的 


# 轮播图的删除
git branch

//创建新的分支
git checkout -b l03_轮播图删除_链式操作_软删除_时间维护

// TODO:
//为什么有delete_time字段时就是软删除?
unscope
...
https://goframe.org/pages/viewpage.action?pageId=1114139#ORM%E9%93%BE%E5%BC%8F%E6%93%8D%E4%BD%9C%E6%97%B6%E9%97%B4%E7%BB%B4%E6%8A%A4-Unscoped%E5%BF%BD%E7%95%A5%E6%97%B6%E9%97%B4%E7%89%B9%E6%80%A7
...

收藏部分一般建议用硬删除


# ORM链式操作的字段过滤

1. Fields用于指定需要操作的表字段。
2. FieldsEx用于指定排除的字段。 ---> [Fields exclude]


## 空值过滤

omitempty标签

omitempty 和 omitNil方法


