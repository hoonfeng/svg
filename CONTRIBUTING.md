# 贡献指南 / Contributing Guide

感谢您对SVG库项目的关注！我们欢迎所有形式的贡献，包括但不限于代码、文档、测试、问题报告和功能建议。

Thank you for your interest in the SVG library project! We welcome all forms of contributions, including but not limited to code, documentation, testing, issue reports, and feature suggestions.

## 🚀 快速开始 Quick Start

### 开发环境设置 Development Environment Setup

1. **安装Go** / **Install Go**
   - 确保安装Go 1.18或更高版本 / Ensure Go 1.18 or higher is installed
   - 验证安装：`go version` / Verify installation: `go version`

2. **克隆仓库** / **Clone Repository**
   ```bash
   git clone https://github.com/hoonfeng/svg.git
   cd svg
   ```

3. **安装依赖** / **Install Dependencies**
   ```bash
   go mod tidy
   ```

4. **运行测试** / **Run Tests**
   ```bash
   go test ./...
   ```

5. **运行示例** / **Run Examples**
   ```bash
   go run examples/animation_builder_demo.go
   ```

## 📝 贡献类型 Types of Contributions

### 🐛 错误报告 Bug Reports

在提交错误报告之前，请：
- 检查是否已有相关的issue
- 确保使用最新版本
- 提供详细的重现步骤

**错误报告模板：**
```markdown
**描述 / Description**
简要描述遇到的问题

**重现步骤 / Steps to Reproduce**
1. 执行...
2. 点击...
3. 查看错误...

**期望行为 / Expected Behavior**
描述您期望发生的情况

**实际行为 / Actual Behavior**
描述实际发生的情况

**环境信息 / Environment**
- OS: [例如 Windows 10, macOS 12.0, Ubuntu 20.04]
- Go版本: [例如 go1.19.3]
- 库版本: [例如 v1.0.0]

**附加信息 / Additional Context**
添加任何其他相关信息、截图等
```

### ✨ 功能请求 Feature Requests

在提交功能请求之前，请：
- 检查是否已有相关的issue或讨论
- 考虑该功能是否符合项目目标
- 提供详细的用例说明

**功能请求模板：**
```markdown
**功能描述 / Feature Description**
简要描述您希望添加的功能

**使用场景 / Use Case**
描述该功能的具体使用场景

**建议实现 / Suggested Implementation**
如果有想法，请描述可能的实现方式

**替代方案 / Alternatives**
描述您考虑过的其他解决方案
```

### 💻 代码贡献 Code Contributions

#### 开发流程 Development Workflow

1. **Fork仓库** / **Fork Repository**
   - 在GitHub上fork项目
   - 克隆您的fork到本地

2. **创建分支** / **Create Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # 或 / or
   git checkout -b fix/your-bug-fix
   ```

3. **开发和测试** / **Develop and Test**
   - 编写代码
   - 添加或更新测试
   - 确保所有测试通过
   - 运行代码格式化：`go fmt ./...`
   - 运行代码检查：`go vet ./...`

4. **提交更改** / **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

5. **推送分支** / **Push Branch**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **创建Pull Request** / **Create Pull Request**
   - 在GitHub上创建PR
   - 填写PR模板
   - 等待代码审查

#### 代码规范 Code Standards

**Go代码规范：**
- 遵循Go官方代码风格指南
- 使用`go fmt`格式化代码
- 使用`go vet`检查代码
- 变量和函数命名使用驼峰命名法
- 包名使用小写字母

**注释规范：**
- 所有公开的函数、类型、常量都必须有注释
- 注释应该同时包含中文和英文
- 复杂的算法或逻辑需要详细注释

**示例：**
```go
// Circle 创建一个圆形元素 / Creates a circle element
// x, y: 圆心坐标 / Center coordinates
// r: 半径 / Radius
func (s *SVG) Circle(x, y, r float64) *CircleElement {
    // 实现代码...
}
```

**测试规范：**
- 新功能必须包含单元测试
- 测试覆盖率应保持在80%以上
- 测试文件命名为`*_test.go`
- 使用表驱动测试模式

**示例测试：**
```go
func TestCircle(t *testing.T) {
    tests := []struct {
        name string
        x, y, r float64
        want string
    }{
        {"basic circle", 10, 20, 5, `<circle cx="10" cy="20" r="5"/>`},
        // 更多测试用例...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试实现...
        })
    }
}
```

### 📚 文档贡献 Documentation Contributions

我们欢迎以下类型的文档贡献：
- 修复文档中的错误或不准确信息
- 改进现有文档的清晰度
- 添加新的示例和教程
- 翻译文档到其他语言

**文档规范：**
- 使用Markdown格式
- 同时提供中文和英文版本
- 代码示例应该可以直接运行
- 包含适当的截图或图表

## 🔍 代码审查 Code Review

所有的代码贡献都需要经过代码审查。审查过程中我们会关注：

- **功能正确性** / **Functional Correctness**
  - 代码是否实现了预期功能
  - 是否处理了边界情况
  - 错误处理是否恰当

- **代码质量** / **Code Quality**
  - 代码是否清晰易读
  - 是否遵循项目的编码规范
  - 是否有适当的注释

- **性能** / **Performance**
  - 是否有性能问题
  - 内存使用是否合理
  - 算法复杂度是否可接受

- **测试** / **Testing**
  - 是否有足够的测试覆盖
  - 测试是否有意义
  - 是否测试了边界情况

## 📋 提交信息规范 Commit Message Convention

我们使用[Conventional Commits](https://www.conventionalcommits.org/)规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型 / Types:**
- `feat`: 新功能 / New feature
- `fix`: 错误修复 / Bug fix
- `docs`: 文档更新 / Documentation update
- `style`: 代码格式化 / Code formatting
- `refactor`: 代码重构 / Code refactoring
- `test`: 测试相关 / Test related
- `chore`: 构建或辅助工具变动 / Build or auxiliary tool changes

**示例 / Examples:**
```
feat(animation): add rotation animation support

fix(renderer): resolve memory leak in image rendering

docs: update API reference for text rendering

test(path): add unit tests for bezier curve parsing
```

## 🏷️ 发布流程 Release Process

项目使用语义化版本控制（Semantic Versioning）：

- **主版本号** / **Major**: 不兼容的API更改
- **次版本号** / **Minor**: 向后兼容的功能添加
- **修订号** / **Patch**: 向后兼容的错误修复

发布流程：
1. 更新版本号
2. 更新CHANGELOG.md
3. 创建发布标签
4. 发布到GitHub Releases

## 🤝 社区准则 Community Guidelines

我们致力于创建一个友好、包容的社区环境：

- **尊重他人** / **Respect Others**: 尊重不同的观点和经验水平
- **建设性反馈** / **Constructive Feedback**: 提供有帮助的、建设性的反馈
- **耐心帮助** / **Patient Help**: 耐心帮助新贡献者
- **专业交流** / **Professional Communication**: 保持专业和友好的交流方式

## 📞 获取帮助 Getting Help

如果您在贡献过程中遇到问题，可以通过以下方式获取帮助：

- 📧 创建GitHub Issue
- 💬 参与GitHub Discussions
- 📖 查看项目文档

## 🙏 致谢 Acknowledgments

感谢所有为项目做出贡献的开发者！您的贡献让这个项目变得更好。

Thank you to all the developers who have contributed to the project! Your contributions make this project better.

---

再次感谢您的贡献！🎉

Thank you again for your contribution! 🎉