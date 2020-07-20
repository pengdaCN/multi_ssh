## multi_ssh
### 什么是multi_ssh
> multi_ssh是一个简单的ssh客户客户端工具，方便在运维数量不多时机器时，提供一个简单的ssh批量处理工具

### multi_ssh使用

**example**

```shell
## 连接单台主机执行命令
multi_ssh --line 'panda, 123456, local.panda.org:22' shell 'whoami'
```

在该实例中，multi_ssh 命令的flag`--line`作用是从命令行中传递需要被操作的主机，格式与hosts格式一样，`shell`子命令是用于重终端执行一行命令使用，在该示例中则代表执行`whoami`命令

在multi_ssh中，子命令代表执行具体的操作，根命令的flag用于选择需要被执行具体操作的主机

#### multi_ssh flag

--line：从命令行读一行的用户信息，格式与hosts文件格式一样

--hosts：从命令行读取hosts文件的位置，默认是`./hosts`

--format：定义输出的格式

**hosts文件格式**

```
panda, 123456, lcoal.pengda.org:22
# panda, 123456, lcoal.pengda.org:22
```

hosts文件有三个字段，分别由`,`分割，别是登录用户名，密码，主机位置，由`#`号开头为注释

**format格式定义**

```
输出格式以#{keyname}为要展示的一个属性
目前定义的key有
user：用于远程连接的用户名
host：远程链接的主机地址
msg：远程终端输出的信息
code：本次执行的返回值
err：本次执行的错误信息
如：#{msg}err:#{err}returncode:#{code}
```

#### 子命令

shell：执行一行shell命令，更多可以通过help获取帮助

script：将本地文件在远端上执行，更多可以通过help获取帮助

copy：将本地文件拷贝到远端，可使用`~`方式代表当前登录用户的家目录，更多可以通过help获取帮助

​	**选项：**

​	--sudo：将需要拷贝的文件放在被操作主机的任意位置

​	--exists：当拷贝的目录不存在会自动创建

#### 常见示例

```shell
# 在hosts.txt文件中的主机执行pt.sh脚本，使用exmine -v作为参数
multi_ssh --hosts hosts.txt script --sudo --args 'examine -v' pt.sh
# 执hosts.txt文件中的主机执行单条命令，自动输入sudo密码
multi_ssh --hosts hosts.txt shell --sudo 'sudo shutdown now'
# 从命令行中出入一条主机信息进行操作
multi_ssh --line 'panda, 123456, local.panda.org:22' shell 'you-get --version'
```