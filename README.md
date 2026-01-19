# 🎮 三角洲行动改枪码小助手

![Built with Wails](https://img.shields.io/badge/Built%20with-Wails-ff69b4.svg)
![Backend Go](https://img.shields.io/badge/Backend-Go-00ADD8.svg)
![Frontend Vue.js](https://img.shields.io/badge/Frontend-Vue.js-4FC08D.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

> **一个懒人为了少翻几个网页做的工具** - 快速查抄高手的改枪码。

## 💡 为什么要做这个？

玩三角洲行动的时候，每次想玩个不常玩的枪，或者武器改动后，就要：

1. 打开浏览器
2. 翻收藏夹/搜视频找心仪的码
3. 复制那一串该死的 21 位改枪码
4. 被固排催墨迹

**太麻烦了！** 😤

作为一个追求极致（懒）的玩家，我决定做一个工具，把主流的"轮椅"配装都整合在一起，一键复制，直接回到游戏继续得吃。

## ✨ 它能做什么

### 🎯 已实现功能

- **🔫 抄作业** - 刀仔和武器大师的配装都在这，不用到处找
- **🎮 双模式** - 烽火地带和全面战场分开，不搞混
- **📋 一键复制** - 点一下就把那串 神秘代码复制到剪贴板
- **🚀 秒开** - 基本上打开就能用。

### 📦 技术特点（给懂技术的看）

- Go + Wails + Vue.js，跨平台，Windows/Mac/Linux 都能用
- 本地 JSON 缓存，速度飞快
- 单文件部署，不用安装，下载就能用
- 代码结构清晰，想自己加功能也方便

## 🚀 怎么用

### 下载使用（推荐）

等我把构建好的版本放出来，下载解压就能用。

### 从源码运行（开发者）

如果你是个爱折腾的人：

```bash
# 1. 安装依赖
go mod download
cd frontend && npm install

# 2. 启动开发模式
wails dev

# 3. 或者构建一个可执行文件
wails build
```

## 📁 项目长这样

```
delta-tool/
├── app/              # 后端核心代码
│   ├── excel.go      # 解析 Excel（刀仔/武器大师的数据格式）
│   ├── cache.go      # 缓存管理（把 Excel 转成 JSON）
│   └── app.go        # 主应用逻辑
├── frontend/         # 前端界面（Vue.js）
├── data/             # 数据文件
│   ├── weapon_codes.json  # 缓存的武器数据
│   └── *.xlsx        # Excel 源文件（不包含在发布版中）
└── cmd/              # 程序入口
```

## 🎮 数据来源

[![抖音](https://img.shields.io/badge/刀仔（三角洲枪匠之王）-black?style=for-the-badge&logo=douyin&logoColor=white)](https://v.douyin.com/xsjMEZDNbbY/)

两个源的配装风格不太一样，建议都看看，找到适合你的。

## 🔧 开发者指南

### 更新数据缓存

如果你改了 Excel 文件，或者想重新生成缓存：

```bash
go run cmd/main.go generate-cache
```

这会读取 `data/` 下的 Excel 文件，生成 `weapon_codes.json`。

### 添加新的数据源

如果你想添加新的配装来源（比如某个 UP 主的 Excel）：

1. 在 `app/excel.go` 里加解析逻辑
2. 更新 `app/cache.go` 的版本号
3. 前端加个选择按钮

### 构建不同平台

```bash
# Windows
wails build -platform windows/amd64

# Mac Intel
wails build -platform darwin/amd64

# Mac Apple Silicon
wails build -platform darwin/arm64

# Linux
wails build -platform linux/amd64
```

## 📝 数据格式

每个武器配装长这样：

```json
{
  "id": "1",
  "mode": "烽火地带",
  "name": "M4A1",
  "tier": "T0",
  "price": 85,
  "build": "标准改装",
  "code": "6XXXXXXXXXXXXXXXXXXXX",
  "range": 52,
  "update_time": "2025-01",
  "source": "刀仔"
}
```

字段说明：

- `code` 就是你在游戏里输入的那串 21 位代码
- `tier` 是强度等级（T0 最强）
- `price` 是改装价格（单位：万）
- `range` 是有效射程（米）

## 🔍 常见问题

### 缓存文件在哪？

- Windows: `%APPDATA%/delta-tool/weapon_codes.json`
- macOS: `~/Library/Application Support/delta-tool/weapon_codes.json`

### 找不到某个枪？

可能数据源里没有，或者我解析错了。可以提个 Issue，我会加上。

### 想要 XX 功能？

目前核心功能就是"查抄作业复制代码"，其他功能（收藏、搜索筛选、统计分析）看心情/时间再加。

## 🤝 贡献

欢迎 PR！如果你想加功能或者修 Bug，看 [CONTRIBUTING.md](CONTRIBUTING.md)。

不过这个项目本来就是个懒人工具，别搞太复杂。🙃

## 📄 开源协议

MIT License - 爱咋用咋用，别告我就行。

## 📮 找我

有问题就提 [Issue](https://github.com/yourusername/delta-tool/issues)，或者直接看代码，写得很简单。

---

<div align="center">
  <strong>🎮 祝你把把得吃，把把 MVP！</strong>
  <br><br>
  <sub>（虽然工具做得简陋，但赢游戏才是重点，对吧？）</sub>
  <br><br>
  Made with ❤️ (and laziness) by Nolan
</div>
