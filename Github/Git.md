

# GIt

## GIt创建仓库

1. **先创建本地仓库**

```shell
#https://blog.csdn.net/qq_45890970/article/details/121381096
#创建本地文件夹(目录)
mkdir learn
#将该目录(文件夹)变成Git可以管理的仓库
git init learn
#随便创建一个文件.(检测之后的git push是否成功)
vi test.txt
#将文件提交到本地仓库,-m 后面的" "为提交的说明信息
git add test.txt
git commit -m "add test.txt"
```

2. 创建远程仓库

```
官网创建public仓库
```

3. 连接本地仓库到远程仓库-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

 

   连接远程库

```shell
git remote add origin git@github.com:hchcfly/Co-learing.git
#查看是否连接成功
git remote -v
origin  https://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@github.com/hchcfly/Co-learing (fetch)
origin  https://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@github.com/hchcfly/Co-learing (push)
```

生成自己的token(个人访问令牌)

```shell
#https://blog.csdn.net/Joy_Cheung666/article/details/119832970?utm_source=app&app_version=4.18.0&code=app_1562916241&uLinkId=usr1mkqgl919blen
#获取令牌后
git remote set-url origin https://<your_token>@github.com/<USERNAME>/<REPO>
#<your_token>:创建的token
#<USERNAME>:github用户名
#<REPO>:项目名称
git push
```

```shell
#-u:首次执行,Git不但会把本地的master分支内容推送的远程新的master分支，还会把本地的master分支和远程的master分支关联起来,之后用git push就行
git push -u origin master
```

## 问题

## 建仓库流程

```shell
https://blog.csdn.net/qq_45890970/article/details/121381096
```

## ssh问题

```shell
https://blog.csdn.net/qq_35495339/article/details/92847819?ops_request_misc=%257B%2522request%255Fid%2522%253A%2522163714056316780261957998%2522%252C%2522scm%2522%253A%252220140713.130102334..%2522%257D&request_id=163714056316780261957998&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~all~sobaiduend~default-1-92847819.pc_search_all_es&utm_term=gitthub%E7%94%9F%E6%88%90ssh%E5%AF%86%E9%92%A5&spm=1018.2226.3001.4187
```

## 本地提交的和远程仓库不能合并

```shell
#出现Compare & pull request绿色按钮
https://www.5axxw.com/questions/content/qzq3j7
```
## 每次git push git pull要输入密码

```shell
https://blog.csdn.net/love910809/article/details/124273642
```




