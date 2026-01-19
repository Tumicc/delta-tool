# 贡献指南 (Contributing Guide)

感谢您对 Delta Tool 项目的关注！我们欢迎各种形式的贡献。

## 🤝 如何贡献

### 报告问题
- 在提交 Issue 前，请先搜索是否已有相同问题
- 提供详细的问题描述和复现步骤
- 附上相关的错误日志或截图

### 提出建议
- 清楚描述建议的功能或改进
- 说明该功能的使用场景和价值
- 如果可能，提供实现思路或参考示例

### 提交代码

#### 开发流程

1. **Fork 项目**
   ```bash
   # 在 GitHub 上 Fork 本项目
   ```

2. **克隆到本地**
   ```bash
   git clone https://github.com/your-username/delta-tool.git
   cd delta-tool
   ```

3. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或修复问题
   git checkout -b fix/your-bug-fix
   ```

4. **安装依赖**
   ```bash
   # 安装 Go 依赖
   go mod download

   # 安装前端依赖
   cd frontend && npm install
   ```

5. **开发与测试**
   ```bash
   # 启动开发服务器
   wails dev
   ```

6. **提交更改**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

7. **推送到 Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

8. **创建 Pull Request**
   - 在 GitHub 上创建 PR
   - 填写 PR 模板
   - 等待代码审查

## 📋 代码规范

### Go 代码规范
- 遵循 [Effective Go](https://golang.org/doc/effective_go) 指南
- 使用 `gofmt` 格式化代码
- 导出函数添加注释
- 错误处理要完整

### Vue.js 代码规范
- 遵循 [Vue.js 风格指南](https://vuejs.org/style-guide/)
- 组件名使用 PascalCase
- 使用 Composition API
- Props 定义要包含类型和默认值

### 提交信息规范
使用约定式提交格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型：
- `feat`: 新功能
- `fix`: 修复 Bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具更新

示例：
```
feat(excel): add support for new data source

- Implement parsing logic for new Excel format
- Add unit tests for new parser
- Update documentation

Closes #123
```

## 🧪 测试

在提交 PR 前，请确保：
- 应用能够成功构建
- 新功能有相应的测试
- 现有测试仍然通过

```bash
# 运行 Go 测试
go test ./...

# 构建测试
wails build
```

## 📁 项目结构

```
delta-tool/
├── app/              # 后端业务逻辑
├── cmd/              # 应用入口
├── frontend/         # Vue.js 前端
├── data/             # 数据文件
├── scripts/          # 构建脚本
└── docs/             # 文档
```

### 添加新功能

#### 后端 (Go)
1. 在 `app/` 目录下添加或修改文件
2. 导出需要在前端调用的函数
3. 运行 `wails dev` 生成绑定

#### 前端 (Vue.js)
1. 在 `frontend/src/` 下添加组件
2. 使用 `wailsjs/go/` 中的绑定调用后端
3. 保持组件简洁，复用性强

## 📝 文档

- 代码注释使用清晰的语言
- 复杂逻辑添加说明
- 更新相关文档（README, CONTRIBUTING.md）
- API 变更需要更新文档

## 🎯 优先任务

查看 [Issues](https://github.com/yourusername/delta-tool/issues) 标记为 `good first issue` 或 `help wanted` 的问题。

## 💬 交流

- 提交 Issue 讨论问题
- 参与 Discussions 交流想法
- 遵循行为准则，保持友善

## 📄 许可

提交代码即表示您同意将代码以 [MIT License](LICENSE) 授权。

---

感谢您的贡献！🎉
