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

+ 收藏部分一般建议用硬删除

# ORM链式操作-时间维护

> gdb模块支持对数据记录的写入、更新、删除时间自动填充，提高开发维护效率。为了便于时间字段名称、类型的统一维护，如果使用该特性，我们约定：
+ 字段应当允许被设置为`null` 
+ 字段的类型必须为时间类型，如data,datatime,timestamp。不支持数字类型字段，如int。、
+ 字段的名称不支持自定义设置，并且固定名称约定为：
    + created_at 用于保存记录的创建时间，仅写入一次
    + updated_at 用于保存记录的修改时间，每次记录变更时更新
    + deleted_at 用于保存记录的软删除特性，只有当记录删除时会写入一次

## 特性注意点
注意字段名称不区分大小写，也会忽略特殊字符。

注意名称字段可以通过配置文件进行自定义修改，并可以使用TimeMaintainDisabled配置完整关闭该特性。

## 特性的启用
当数据表中包含一个或多个字段时，该特性会自动启用。


# ORM链式操作-时间维护

1. Fields 用于指定需要操作的表字段,包括查询字段、写入字段、更新字段deng过滤。
2. FieldsEx 用于指定排除的字段。 ---> [Fields exclude]


## 空值过滤

当map/struct中存在空值如`nil`,`""`,`0`时，默认情况下，gdb会将其当作正常的输入参数，因此这些参数也会被更新到数据表。`OmitEmpty特性`可以在将数据写入到数据库之前过滤空值数据。
```go
func (m *Model) OmitEmpty() *Model
func (m *Model) OmitEmptyWhere() *Model
func (m *Model) OmitEmptyData() *Model 
```
批量写入/更新操作中OmitEmpty方法将会失效，因为在批量操作中，必须保证每个写入记录的字段是统一的。



omitempty标签

omitempty 和 omitNil方法


