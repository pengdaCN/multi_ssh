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

​	**选项：**

​	--sudo：将需要拷贝的文件放在被操作主机的任意位置

​	--exists：当拷贝的目录不存在会自动创建

palybook：使用lua脚本，调用multi_ssh提供的方法，进行调用

​	**选项**

​	--set-args：在multi_ssh 调用任何一个函数之前，进行将一些变量的值进行提前的设置，在set-args中，只能将变量设置为字符串，多个变量以`,`分隔

​	exmaple：

​	`--set-args 'name = "张三", age = "18"'`

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

使用multi_ssh执行的lua脚本，必须有一个exec函数，multi_ssh会自动调用exec函数，同时，exec函数头必须如下所示

```lua
-- exec函数头
function exec(term)
   	local r1 = term.shell({sudo=true, 'sudo whoami'})
    term.outln(r1.msg)
end
```

exec方法的term参数由multi_ssh调用时自动传入

有关exec函数的term参数介绍

> exec函数的term参数有multi_ssh 在执行对应主机时自动传入，其中传入lua虚拟机中的term是一个自读的userdata，是对go函数的封装，在golang对应的一个go方法持有该lua thread操作的terminal的闭包函数，正是通过该中方式实现对不同主机的区分操作，对于lua而言，term对象就是所有对terminal操作的集合

在playbook中除了主机执行入口函数exec，还有4各对应playbook周期的钩子函数，分别是BEGIN，OVER，EXEC_BEGIN，EXEC_OVER

，其中，BEGIN和OVER函数头如下，这两个函数分别在整个playbook开始执行之前调用和执行结束之后调用

```lua
function BEGIN()
    --statment
end

function OVER()
    --statemet
end
```

EXEC_BEGIN与EXEC_OVER函数分别是在exec方法执行之前和之后调用的，其函数头如下

```lua
function EXEC_BEGIN(term)
    --statment
end

function EXEC_OVER(term)
    --statment
end
```

关于exec函数，该函数是multi_ssh针对每个被操作的主机都需要执行的函数，对应一个gorouting，所有，在exec函数中的都是并发的操作，对于全局变量的读写一定要注意并发的安全性

##### 关于exec方法所接收的term对象和所提供的方法以及其说明

term对象是一个自读的userdata对象，是一组对当前主机操作的方法集合，其结构如下

```lua
term = {
    shell: function({sudo: bool, 'commands'}),
    script: function({sudo: bool, text: str, 'script_path'}),
    copy: function({sudo: bool, exists: bool, {'src1', 'src2'}| 'src', 'dst'}),
    context: function({text: string, filename: string),
    out: function(msg: str),
    outln: function(msg: str),
    extraInfo: function(),
    hostInfo: function(),
    setCode: function(code: int),
    setErrCode: function(code: int, errInfo: str),
    sleep: function(secends: int),
    hostInfo: {
        line: int,
        ip: str,
        port: str,
        user: str,
        extra: {
            key: val
        }
    }
    iota: int
    exit: function()
}
```

##### 关于提供的tools工具对象所提供的方法和说明

tools对象是multi_ssh实现的一组全局工具方法，通过使用自读的userdata注入到全局变量tools中，其结构如下

```lua
tools = {
    sleep: function(secends: int),
    setShareIotaMax: function(max: int),
    getShareIota: function() -> int,
    newWaitGroup: function() -> {
        add: function(i: int),
        done: function(),
        wait: function(),
    },
    newTokenBucket: function(max: int) -> {
    	get: function() -> int,
    },
    newMux: function() -> {
    	lock: function(),
       	unlock: function(),
        rLock: function(),
        rUnlock: function(),
  	},
   	newSafeTable: function() -> {
        append: function(val),
        set: function(key, val),
        get: function(key) -> val,
        len: function() -> int,
        rLock: function(),
        rUnlock: function(),
        into: function() -> {},
    },
    newOnce: function() -> {
    	Do: function(function()),
    },
	str: {
        split: function(src: str, option(sep, default=' ')) -> []str,
        hasPrefix: function(s: str) -> bool,
        hasSuffix: function(s: str) -> bool,
        trim: function(s: str) -> bool,
        replace: function(s: str, old: str, new: str, option(count: int, default=-1)) -> []str,
        contain: function(s: str, sub: str) -> bool
    },
    re: {
        match: function(s: str, re: str) -> bool,
        find: function(s: str, re: str, mode: option(mode: emnu('sub', 'sub_all', 'str', 'str_all'), default='sub_all')) -> emnu(str, []str),
        replace: function(s: str, re: str, new: str) -> str
        split: function(s: str, re: str, count: option(count: int, default=-1)) -> []str,
        splitSpace: function(s: str) -> []str,
    },
}
```

### 关于使用密钥进行登陆的配置

如需要使用密钥进行登陆，则需要在hosts文件中，该主机的扩展字段中设置`PRIKEY='pri_key_path'`，如下

```txt
root,,192.168.0.2:22; `PRIKEY="/home/user/id_rsa"`
root,123456,192.168.0.3:22; `PRIKEY="/home/user/id_rsa"`
```

