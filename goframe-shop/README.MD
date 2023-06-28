# GoFrame Template For SingleRepo

Quick Start: 
- https://goframe.org/pages/viewpage.action?pageId=1114399

# 开发流程
GoFrame V2开发流程：
1. 设计表结构，初始化项目，修改配置文件
2. 使用gf gen dao 生成对应的dao/do/model
3. 编写api层：定义业务侧数据结构，提供对外接口的输入/输出数据结构
4. 编写model层，定义数据侧数据结构，提供对内的数据处理的输入/输出数据结构
5. 编写logic层，自动生成service层代码