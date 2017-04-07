# Private Vendor Tool
------
## 不成熟的小工具
我们自己的项目用一用的小工具，用于把私有库搞进vendor目录

vendor.toml文件示例:
```toml
[[dep]]
git="git项目地址"
branch="分支或者tag 默认是master"
path="库的路径 默认通过git项目地址生成"
[[dep]]
git="https://github.com/cocotyty/summer"
branch="master"
path="github.com/cocotyty/summer"
```

---
命令|含义
--- |---
pvendor|读取当前目录下的vendor.toml ,然后通过git把私有项目弄进去
---