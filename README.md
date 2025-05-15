## cactus
- CactusOps devops项目

## 目录结构
```
project/
├── config.go            # 配置初始化
└── main.go              # 应用入口 
├── config/
│   └── config.yaml      # 配置文件
├── api/                 # API层 - HTTP接口定义和路由
├── service/             # 业务逻辑层
├── types/               # 数据模型和DTO
├── sql/                 # 数据访问层
├── utils/               # 通用工具函数
├── pkg/                 # 可复用的公共包
├── scripts/             # 部署/构建脚本
├── test/                # 测试代码
├── go.mod
└── go.sum
```
## 功能