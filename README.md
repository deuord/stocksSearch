# 全球股市情报面板

实时全球股市行情监控，支持美股、港股、A股、韩股、澳股、日股、新加坡多市场分类查看，配备3D地球可视化。

## 快速启动

```bash
# 后端 (端口 8080)
cd backend && go run main.go

# 前端 (端口 3000，自动代理API)
cd frontend && npm install && npm run dev
```

## 技术栈

- **后端**: Go + 标准库 net/http（MVC架构）
- **前端**: React 19 + TypeScript + Vite
- **3D**: Three.js + @react-three/fiber
- **数据**: iTick API

## 文档

- [架构文档](ARCHITECTURE.md)
- [需求流水线](CHANGELOG.md)
