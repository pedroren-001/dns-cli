# dns-cli

基于腾讯云 DNSPod API 的命令行工具，用于管理域名解析记录（增、删、查）。

## 安装

```bash
go install github.com/pedroren-001/dns-cli@latest
```

或从源码构建：

```bash
git clone https://github.com/pedroren-001/dns-cli.git
cd dns-cli
go build -o dns-cli
```

## 使用前配置

需要设置腾讯云 API 密钥环境变量：

```bash
export TENCENTCLOUD_SECRET_ID=<your-secret-id>
export TENCENTCLOUD_SECRET_KEY=<your-secret-key>
```

密钥获取：https://console.cloud.tencent.com/cam/capi

## 用法

```bash
# 列出域名下所有解析
dns-cli list -d example.com

# 添加 A 记录
dns-cli add -d example.com --type A --sub www --value 1.2.3.4

# 添加 CNAME
dns-cli add -d example.com --type CNAME --sub blog --value cdn.example.com

# 删除记录（先 list 查 record-id）
dns-cli list -d example.com
dns-cli rm -d example.com --record-id 12345
```

支持的记录类型：A、AAAA、CNAME、MX、TXT 等。默认 TTL 600 秒，默认线路"默认"。

## 仓库镜像

本仓库同时维护在两个 Git 服务上，**一次 `git push` 双端同步**：

- **GitHub**（主发布，`go install` 用）: <https://github.com/pedroren-001/dns-cli>
- **Forgejo**（自建镜像）: `ssh://git@forge.renyanjie.cn:36322/pedroren/dns-cli.git`

### 本地配置双推

首次 clone 后一次性配好：

```bash
git clone https://github.com/pedroren-001/dns-cli.git
cd dns-cli

# 移除默认的单一 push URL，改成两个
git remote set-url --add --push origin ssh://git@forge.renyanjie.cn:36322/pedroren/dns-cli.git
git remote set-url --add --push origin git@github.com:pedroren-001/dns-cli.git
```

验证 —— `git remote -v` 应看到 3 行：

```
origin  https://github.com/pedroren-001/dns-cli.git (fetch)
origin  ssh://git@forge.renyanjie.cn:36322/pedroren/dns-cli.git (push)
origin  git@github.com:pedroren-001/dns-cli.git (push)
```

之后 `git push origin master` 会同时推 forge 和 GitHub。`git fetch` 仍从 GitHub 拉（保持公开可 clone）。

### 原理

- `git remote set-url --add --push` 给同一个 remote 添加**额外的 push URL**
- 一个 remote 允许**一个 fetch URL + 多个 push URL**
- Git push 时会**依次尝试**每个 push URL，任一失败则整体失败

### 前置：SSH key 已注册

双推依赖 SSH 认证。首次配置前确保：

```bash
ssh -T git@github.com     # 应返回 "Hi <username>! You've successfully authenticated..."
ssh -T git@forge.renyanjie.cn -p 36322    # 类似响应
```

如未注册，把 `~/.ssh/id_ed25519.pub` 分别登记到 GitHub 和 Forge 的 SSH keys 页面。

## License

MIT
