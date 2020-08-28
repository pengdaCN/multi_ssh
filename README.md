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

--filter：通过命令行指定，选择需要执行的主机

```
格式如下：
字段!=/==值
如：
work="p" group="worker1"
表示选择扩展字段work="p" group="worker1"的主机
目前filter功能只支持扩展字段的处理，或许计划添加ip的选择
```

**hosts文件格式**

```
panda, 123456, lcoal.pengda.org:22; `group="work1" cpu="intel"`
# panda, 123456, lcoal.pengda.org:22
```

hosts文件有四个字段，分别由`,`分割，别是登录用户名，密码，主机位置，扩展信息，由`#`号开头为注释

扩展信息由反引号包裹，语法与golang的struct的tag一样

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

palybook：将lua脚本执行并运行

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

#### playbook介绍与使用

该功能通过gopher-lua库提供的lua虚拟机以及api实现

playbook模块目前提供3个由golang包装的函数，用于调用m_terminal包提供的方法

1. shell函数

   ```lua
   -- 函数头如下
   shell(id int, sudo bool, cmd string) out -> playbookResult
   --[[
   该函数需要传递3个参数，id参数用于确定需要执行的主机，sudo参数用于确定是否以sudo方式执行，cmd参数及真正要执行的命令，主要，在multi_ssh提供的公开执行一条命令的方式中，都在将命令进行预先的处理，使其在英语语系下执行
   ]]
   -- 调试示例如下
   local tab = shell(id, false, 'whoami')
   --[[
   其中tab是shell函数返回的table类型，在lua中，由于默认是glabol的，所有建议所有的值都设置为local
   ]]
   
   --[[
   关于golang包装函数返回table
   所有有glang包装的函数，如有返回值，table类型键值都是如下
   ]]
   tab = {
       u={ -- u键表示用户信息
           user='登录用户名',
           host='登录主机名',
       }
       msg='命令执行输出的结果，有stdout和stderr',
       errInfo='如命令执行失败，在golang中对错误的描述字符串',
       code='int 类型，是执行命令完后的状态码'
   }
   ```

2. copy函数

   ```lua
   -- copy函数头如下
   copy(id int, sudo, exists bool, src []string, dst string, attr map<lua table>) out -> playbookResult
   --[[
   id参数与shell一样，sudo参数可上传到服务器任意位置，exists文件夹不存在就创建，src，需要拷贝的一些文件，必须用数组，dst上传的目标位置，attr，需要设置的上传后的文件属性
   ]]
   -- 调用示例
   local = copy(id, false, false, {'/tmp/data.txt'}, '~', nil)
   -- 注意，目前copy函数的属性设置功能还没实现
   ```

3. script函数

   ```lua
   -- script函数头如下
   script(id int, sudo bool, script_path string, args string) out -> playbookResult
   --[[
   id参数与shell一样，sudo参数以sudo方式执行，script_path参数本地执行脚本的位置，args参数，脚本的参数
   ]]
   -- 调用示例
   local = script(id, false, '/tmp/1.sh', '')
   ```

关于需要执行的lua脚本说明

有multi_ssh执行的lua脚本，必须有一个exec函数，multi_ssh会自动调用exec函数，同事，exec函数头必须如下所示

```lua
-- exec函数头
function exec(id)
   shell(false, 'echo hello world') 
end
```

exec方法的id有multi_ssh调用时自动传入

playbook新增函数：

1. out 用于输出打印数据，格式有--format指定
2. extraInfo 用于获取当前执行的主机的扩展信息
3. hostInfo 用于获取主机的信息
4. outln 换行输出
5. setCode 设置返回状态码
6. setErrInfo 设置错误信息