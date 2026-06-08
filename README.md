# bloom-tool

用于分布式系统中的布隆过滤器工具，基于 gRPC 实现。

## 功能特性

- 布隆过滤器：支持元素添加和存在性检查
- 键值对存储：支持键值对的存储和读取
- 多地址连接：支持多个服务地址，自动选择可用节点
- 连接池管理：自动处理 gRPC 连接状态

## 安装

```bash
go get github.com/ou8zz/bloom-tool
```

## 使用示例

### 初始化 gRPC 客户端

```go
import "github.com/ou8zz/bloom-tool/service"

func main() {
    addrs := []string{"localhost:50051", "192.168.1.100:50051"}
    service.InitGRpcClient(addrs)
    defer service.CloseGRpcClient()
}
```

### 布隆过滤器操作

```go
// 添加元素
service.Add("item_key")

// 检查元素是否存在
exists := service.Exists("item_key")
```

### 键值对存储操作

```go
// 存储键值对
service.StoreIP("key", "value")

// 加载值
value := service.LoadIP("key")
```

## 协议定义

服务使用 Protocol Buffers 定义接口，位于 `proto/bloom.proto`。

### BloomService

```protobuf
service BloomService {
  // 布隆过滤器操作
  rpc Add (ItemRequest) returns (ItemResponse);
  rpc Exists (ItemRequest) returns (ItemResponse);

  // 键值对存储
  rpc Store (MapRequest) returns (MapResponse);
  rpc Load (MapRequest) returns (MapResponse);
}
```

## 依赖

- Go 1.26+
- google.golang.org/grpc v1.80.0
- github.com/golang/protobuf v1.5.4
