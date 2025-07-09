# 贡献指南 / Contributing Guide

感谢您对SVG库的关注！我们欢迎各种形式的贡献。

Thank you for your interest in the SVG library! We welcome all forms of contributions.

## 开发环境设置 / Development Environment Setup

### 前置要求 / Prerequisites

- Go 1.20 或更高版本 / Go 1.20 or higher
- Git
- 文本编辑器或IDE / Text editor or IDE

### 克隆仓库 / Clone Repository

```bash
git clone https://github.com/hoonfeng/svg.git
cd svg
```

### 安装依赖 / Install Dependencies

```bash
go mod download
```

### 运行测试 / Run Tests

```bash
go test ./...
```

## 贡献类型 / Types of Contributions

### 1. 错误报告 / Bug Reports

如果您发现了错误，请创建一个Issue并包含以下信息：

If you find a bug, please create an Issue with the following information:

- 错误描述 / Bug description
- 重现步骤 / Steps to reproduce
- 期望行为 / Expected behavior
- 实际行为 / Actual behavior
- 环境信息 / Environment information
- 代码示例 / Code example

**错误报告模板 / Bug Report Template:**

```markdown
## 错误描述 / Bug Description
[简要描述错误 / Brief description of the bug]

## 重现步骤 / Steps to Reproduce
1. [第一步 / First step]
2. [第二步 / Second step]
3. [看到错误 / See error]

## 期望行为 / Expected Behavior
[描述您期望发生的事情 / Describe what you expected to happen]

## 实际行为 / Actual Behavior
[描述实际发生的事情 / Describe what actually happened]

## 环境 / Environment
- Go版本 / Go version: [例如 go1.20]
- 操作系统 / OS: [例如 Windows 10, macOS 12, Ubuntu 20.04]
- SVG库版本 / SVG library version: [例如 v1.0.0]

## 代码示例 / Code Example
```go
// 您的代码示例 / Your code example
```
```

### 2. 功能请求 / Feature Requests

我们欢迎新功能的建议！请创建一个Issue并说明：

We welcome suggestions for new features! Please create an Issue and explain:

- 功能描述 / Feature description
- 使用场景 / Use cases
- 预期API设计 / Expected API design
- 实现建议 / Implementation suggestions

**功能请求模板 / Feature Request Template:**

```markdown
## 功能描述 / Feature Description
[清楚简洁地描述您想要的功能 / Clear and concise description of the feature]

## 动机 / Motivation
[解释为什么这个功能有用 / Explain why this feature would be useful]

## 详细设计 / Detailed Design
[描述功能应该如何工作 / Describe how the feature should work]

## API设计示例 / API Design Example
```go
// 展示预期的API使用方式 / Show expected API usage
```

## 替代方案 / Alternatives
[描述您考虑过的替代解决方案 / Describe alternative solutions you've considered]
```

### 3. 代码贡献 / Code Contributions

#### 开发流程 / Development Workflow

1. **Fork仓库 / Fork the repository**
2. **创建功能分支 / Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **进行更改 / Make your changes**
4. **添加测试 / Add tests**
5. **运行测试 / Run tests**
   ```bash
   go test ./...
   ```
6. **提交更改 / Commit changes**
   ```bash
   git commit -m "Add: your feature description"
   ```
7. **推送分支 / Push branch**
   ```bash
   git push origin feature/your-feature-name
   ```
8. **创建Pull Request / Create Pull Request**

#### 提交信息规范 / Commit Message Convention

我们使用以下格式的提交信息：

We use the following commit message format:

```
<type>: <description>

[optional body]

[optional footer]
```

**类型 / Types:**
- `Add`: 新功能 / New feature
- `Fix`: 错误修复 / Bug fix
- `Update`: 更新现有功能 / Update existing feature
- `Remove`: 删除功能 / Remove feature
- `Docs`: 文档更改 / Documentation changes
- `Style`: 代码格式更改 / Code style changes
- `Refactor`: 重构 / Refactoring
- `Test`: 测试相关 / Test related
- `Chore`: 构建过程或辅助工具的变动 / Build process or auxiliary tool changes

**示例 / Examples:**
```
Add: circle drawing functionality

Fix: incorrect path calculation in bezier curves

Docs: update API reference for text rendering
```

### 4. 文档贡献 / Documentation Contributions

文档改进包括：

Documentation improvements include:

- API文档 / API documentation
- 教程和示例 / Tutorials and examples
- README更新 / README updates
- 代码注释 / Code comments

## 代码审查流程 / Code Review Process

1. **自动检查 / Automated Checks**
   - 所有测试必须通过 / All tests must pass
   - 代码格式检查 / Code formatting checks
   - 静态分析 / Static analysis

2. **人工审查 / Manual Review**
   - 代码质量 / Code quality
   - 设计一致性 / Design consistency
   - 性能考虑 / Performance considerations
   - 文档完整性 / Documentation completeness

3. **反馈处理 / Feedback Handling**
   - 及时响应审查意见 / Respond to review comments promptly
   - 进行必要的修改 / Make necessary changes
   - 更新测试和文档 / Update tests and documentation

## 代码规范 / Code Standards

### Go代码规范 / Go Code Standards

- 遵循Go官方代码风格 / Follow official Go code style
- 使用`gofmt`格式化代码 / Use `gofmt` to format code
- 使用`golint`检查代码 / Use `golint` to check code
- 遵循Go命名约定 / Follow Go naming conventions

### 注释规范 / Comment Standards

- 公共API必须有文档注释 / Public APIs must have documentation comments
- 复杂逻辑需要解释性注释 / Complex logic needs explanatory comments
- 注释应该解释"为什么"而不是"什么" / Comments should explain "why" not "what"
- 支持中英文双语注释 / Support bilingual comments (Chinese/English)

### 测试规范 / Testing Standards

- 新功能必须包含测试 / New features must include tests
- 测试覆盖率应该保持在80%以上 / Test coverage should be above 80%
- 使用表驱动测试 / Use table-driven tests
- 包含边界条件测试 / Include edge case tests

### 文档规范 / Documentation Standards

- API文档使用godoc格式 / API documentation uses godoc format
- 示例代码必须可运行 / Example code must be runnable
- 文档应该包含使用场景 / Documentation should include use cases
- 支持中英文双语文档 / Support bilingual documentation

## 发布流程 / Release Process

1. **版本规划 / Version Planning**
   - 功能冻结 / Feature freeze
   - 测试和修复 / Testing and fixes
   - 文档更新 / Documentation updates

2. **版本发布 / Version Release**
   - 创建发布分支 / Create release branch
   - 更新版本号 / Update version number
   - 生成变更日志 / Generate changelog
   - 创建Git标签 / Create Git tag

3. **发布后 / Post-release**
   - 监控反馈 / Monitor feedback
   - 修复紧急问题 / Fix critical issues
   - 规划下一版本 / Plan next version

## 社区准则 / Community Guidelines

- 保持友善和尊重 / Be kind and respectful
- 欢迎新贡献者 / Welcome new contributors
- 提供建设性反馈 / Provide constructive feedback
- 遵循行为准则 / Follow code of conduct

## 获取帮助 / Getting Help

如果您需要帮助，可以：

If you need help, you can:

- 创建Issue提问 / Create an Issue to ask questions
- 查看现有文档 / Check existing documentation
- 参考示例代码 / Refer to example code
- 联系维护者 / Contact maintainers

## 许可证 / License

通过贡献代码，您同意您的贡献将在MIT许可证下授权。

By contributing code, you agree that your contributions will be licensed under the MIT License.

---

再次感谢您的贡献！🎉

Thank you again for your contributions! 🎉