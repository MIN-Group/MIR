
# 0.快速部署（docker部署方案，一行代码一键部署）
此为快速部署方案，无需golang环境，无需gcc，甚至不需要代码，只需要docker和Dockerfile即可
如需手动部署请看后面的小节以及install.sh
- docker镜像构建
```bash
docker build -t dpmir .
```

- docker运行
```bash
 docker run -d  -p 13899:13899  --name dpmire dpmir        
```

- 启动mir
```bash
docker exec -it dpmire mir
```

- 运行mirc
```bash
docker exec -it dpmire mirc
```

-本地代码构建方案
将Dockerfile-local复制到与minlib和mir-go目录平级
```
docker build -t dpmir-local -f Dockerfile-local .
```

# 1.go环境配置
- 根据环境选择合适的go版本（树莓派4B是ARM，此处是AMD）
```bash
wget https://dl.google.com/go/go1.18.3.linux-amd64.tar.gz
```

- 解压：
```bash
sudo tar -zxvf go1.18.3.linux-amd64.tar.gz
```

- 移动：
```bash
sudo mv go /usr/local
```

- 创建库目录（“min”可以替换为自己的文件夹）：
```bash
mkdir /home/min/go
```

- 修改环境变量：

interactive-login终端，min用户加载
```bash
sudo vim /etc/profile
```
interactive-non-login终端，root用户加载：
```bash
sudo vim /root/.bashrc
```

- 两个文件添加的内容相同：
```bash
# go安装目录：
export GOROOT=/usr/local/go
# go工具链：
export PATH=$PATH:$GOROOT/bin
# go库目录：
export GOPATH=/home/min/go
# go配置文件目录：
export XDG_CONFIG_HOME=/home/min/.config
```

- 保存环境变量：
```bash
# min用户加载
source /etc/profile
# root用户加载
sudo -s
source /root/.bashrc
```

- 查看go配置：
```bash
go env
```

- 打开go mod：
```bash
go env -w GO111MODULE=on
```

- 修改go proxy：
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```
# 2.下载minlib
```bash
git clone http://git.sscfs.cn/pkusz-future-network-lab/common/minlib.git
# 切换分支
cd minlib
git checkout parallel-mir
go mod download
```
# 3.配置mir

- 下载mir-go
```bash
git clone http://git.sscfs.cn/pkusz-future-network-lab/mir/mir-go.git
# 切换分支
cd mir-go
git checkout parallel-mir
```

- 更新go mod
```bash
go mod tidy
```

- 创建本地文件夹
```bash
sudo mkdir /usr/local/etc/mir
```

- 传入配置文件
```bash
sudo cp mirconf.ini /usr/local/etc/mir/mirconf.ini
```

- 安装
```bash
# 不要用sudo安装
# 配置最后会要求输入默认身份密码
./install.sh
```

- 直接启动
```bash
sudo mir
```

-  安装成系统服务并启动
```bash
# 安装成系统服务
sudo mird install
# 启动程序
sudo mird start
# 终止程序
sudo mird stop
# 查看程序状态
sudo mird status
# 从系统服务中卸载 => 需要重新覆盖的时候先执行这个
sudo mird remove 
```

- 身份、密码修改
```bash
# 设置或者修改配置文件中的默认身份，对应的是配置文件：/usr/local/etc/mir/mirconf.ini 中的 DefaultId 配置项
# 接着调用 mirgen 设置或修改默认身份的密码
# 验证或者设置：
sudo mirgen
# 修改默认身份的密码：
sudo mirgen -rp
# 如果旧版的密码是明文，可以通过 -oldPasswdNoHash 参数兼容
sudo mirgen -rp -oldPasswdNoHash
```

- 终端日志输出位置 
   - Macos 
      - /usr/local/var/log/mird.err
      - /usr/local/var/log/mird.log
   - Linux => /var/log/mird.log

关于启动后服务的日志如何输出 => [https://blog.csdn.net/sinat_24092079/article/details/120676316](https://blog.csdn.net/sinat_24092079/article/details/120676316)

- mir使用
   - 使用mirc进入管理
   - 使用help查看mirc用法
