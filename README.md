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
!=与==用于匹配相等满足，还是不相等满足
如：
work="p" group="worker1" IP=="192.168.0.1" USER=="cat" POER=="23"
建议用户在hosts文件扩展字段定义key时，key小写，在filter中，大写的key，一般有特殊用途

IP匹配使用
用于从hosts文件挑选合适ip的主机
规则：IP只能用于ipv4的匹配，不支持ipv6，目前支持ip范围，ip网络位，单个ip，多个ip范围等方式，同时，书写匹配ip是一定要书写完整，不能简写与省略，ip匹配中，使用-表示范围匹配，使用,来分割多个匹配以下是IP配置示例
IP=="192.168.0.3" # 表示匹配192.168.0.3 这个ip的主机
IP=="192.168.0.0/16" # 表示匹配192.168.0.0/16网段的所有主机
IP=="192.168,101.2,3,4,5" # 表示匹配192.168,192.101,中所有的.3.4.5的ip地址
IP=="192.168.101-180.2-30" # 表示匹配192.168.101到192.168.180,ip从2到30的所有ip地址，注意，ip范围匹配是使用的闭区间进行匹配的
IP=="192.168,172.2-60,58,90-95" # 混合匹配，规则同上面

PORT匹配使用
用于从hosts文件中挑选指定ssh 端口连接的主机
规则：PORT后只能指定一个端口，不能想IP指定多个，实例如下
PORT=="23" # 表示用于匹配ssh连接是23的端口主机

USER匹配使用
用于匹配hosts文件中的ssh登录用户名
规则：USER后只能指定一个用户名，实例如下
USER=="root" # 表示选择hosts文件所有由root用户登录的主机

除这些，其他的 所有key都是以hosts文件中扩展信息为匹配信息源，实例如下

如有hosts文件如下
ppp,13456,192.168.0.3:22; `group="worker" arch="intel"`

则可用如下key value匹配它

group=="worker" arch=="intel"

也可用!= 排出他

注意，一个--filter中可以有多个条件，但是他们需要使用空白字符作为分割，如下

'group!="worker" IP=="192.168.0.0/24"'

只有当所有条件符合的主机才会被选择执行
```

**hosts文件格式**

```
panda, 123456, lcoal.pengda.org:22; `group="work1" cpu="intel"`
# panda, 123456, lcoal.pengda.org:22
```

hosts文件有四个字段，分别由`,`分割，别是登录用户名，密码，主机位置，扩展信息，由`#`号开头为注释

扩展信息由反引号包裹，语法与golang的struct的tag一样，需要注意，当一项主机信息包含扩展信息是，在写完标准信息后一定要在期结尾处加上`;`分号，这是multi_ssh对于hosts文件语法的硬性要求

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

palybook：使用lua脚本，调用multi_ssh提供的方法，进行调用

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

使用multi_ssh执行的lua脚本，必须有一个exec函数，multi_ssh会自动调用exec函数，同事，exec函数头必须如下所示

```lua
-- exec函数头
function exec(term)
   	local r1 = term.shell({sudo=true, 'sudo whoami'})
    term.outln(r1.msg)
end
```

exec方法的term参数由multi_ssh调用时自动传入

有关exec函数的term参数介绍

> exec函数的term参数有multi_ssh 在执行对应主机时自动传入，其中传入lua虚拟机中的term是一个table，是对go函数的封装，在golang对应的一个go方法持有该lua thread操作的terminal的闭包函数，正是通过该中方式实现对不同主机的区分操作，对于lua而言，term对象就是所有对terminal操作的集合

**term对象所支持的方法**

1. shell

   ```lua
   function exec(term)
      	local r1 = term.shell({sudo=true, 'sudo whoami'})
       term.outln(r1.msg)
   end
   ```

   shell函数等同于shell子命令，用与执行一条shell命令

   shell函数接收一个table，参数table规则如下，键`sudo`的值是一个boolean值，用控制是否开始sudo功能，该功能需要登录被操作主机有执行sudo的权利，否则即使使用改参数也无效

   shell函数所接受的table参数的中的sudo参数可以不写，不写则为false

2. script

   ```lua
   function exec(term)
      	local r1 = term.script({sudo=true, './example.sh'})
       term.outln(r1.msg)
   end
   ```

   script函数等同于script子命令，用于在远端的主机上执行登录的一个脚本

   script函数同shell函数一样，也收一个table对象，其中sudo参数用于以sudo身份执行脚本，另一个非key-value的参数是需要在远端执行脚本的路径

3. copy

   ```lua
   function exec(term)
      	local r1 = term.copy({{'~/example.sh', '~/example.txt'} ,'~'})
       term.outln(r1.code)
   end
   ```

   copy函数等同于copy子命令，用于将本地文件传送到远端

   sopy函数接受一个table对象的参数，其中第一个非key-value形式的值，可以是string，可以是table类型，用于表示本地需要拷贝到远端的文件，当参数用string类型时，表示本地有一个文件需要传输到远端，主要，类型为string时，不支持多个本地文件的表示方法，当类型为table类型时，其中table的元素需要为string类型，表示需要从本地拷贝到远端的文件，其中第 二个位key-value的值为需要拷贝到 远端文件的路径，需要注意的是，这个两参数不能省略

4. extraInfo

   ```lua
   function exec(term)
      	local info = term.extraInfo()
       term.outln(type(info))
       for k, v in pairs(info) do
           print(k, v)
       end
   end
   ```

   extraInfo函数适用于在远端执行时动态的获取hosts文件配置的扩展信息参数

   extraInfo函数返回一个table类型，其中键和值都是string类型

5. out

   ```lua
   function exec(term)
       term.out('hello world')
   end
   ```

   out函数，使用multi_ssh统一输出，不进行换行

   out函数与outln函数本质一样，都是一样输出字符串到multi_ssh的执行结果中，注意，在multi_ssh不建议使用print函数，应为该函数由lua虚拟机进行自行输出，不受到multi_ssh控制，可能会打乱multi_ssh的输出样式，推荐使用out和outln进行输出打印

6. outln

   ```lua
   function exec(term)
       term.outln('hello world')
   end
   ```

   outln函数与out函数一样，唯一的区别是outln函数会进行换行输出

7. setCode与setErrinfo

   ```lua
   function exec(term)
       term.setCode(10)
       term.setErrInfo(1, '系统错误')
   end
   ```

   setCode函数用于设置lua脚本执行完时，执行结果的状态码，setErrInfo函数用于设置lua脚本执行完后状态码，与错误信息，注意，setCode与setErrInfo都可以设置状态，最终以最后一个执行的为准，同时，需要注意的是，setCode只接受一个整数类型的参数，setErrInfo则在setCode的基础上增加接受一个string类型参数作为错误的描述

playbook执行lua脚本是从lua脚本定义的exec函数开始的，multi_ssh的lua脚本推荐整个文件没有在函数之外的代码，建议过于复杂的逻辑抽离成单独的函数，最终在exec函数中调用

