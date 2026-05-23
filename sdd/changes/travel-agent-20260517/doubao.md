# 火山引擎-豆包接入文档
## 基本信息
API Key: 0dd386ea-afd4-45f1-a7da-efb3dc2e6df7
可用模型：
- doubao-seed-1-8-251228
- doubao-seed-2-0-lite-260428
- deepseek-v3-2-251201

调用地址：https://ark.cn-beijing.volces.com/api/v3/chat/completions

## 调用示例
以curl命令为例：
```shell
curl https://ark.cn-beijing.volces.com/api/v3/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ARK_API_KEY" \
  -d $'{
    "messages": [
        {
            "content": "You are a helpful assistant.",
            "role": "system"
        },
        {
            "content": "hello",
            "role": "user"
        }
    ],
    "model": "doubao-seed-2-0-lite-260215",
    "stream": true
}'
```
### 请求参数说明
                                                    
| 参数名称          | 父节点     | 参数类型 | 是否必传 | 说明                                                                                                                      |
|-------------------|------------|----------|:---------|:--------------------------------------------------------------------------------------------------------------------------|
| model             | /          | string   | 是       | 调用的模型 ID （Model ID）                                                                                                  |
| messages          | /          | object[] | 是       | 消息列表，不同模型支持不同类型的消息，如文本、图片、视频、音频等。                                                              |
| role              | messages   | string   | 是       | 发送消息的角色：支持system、user、assistant                                                                                  |
| content           | messages   | string   | 是       | 消息内容                                                                                                                  |
| reasoning_content | messages   | string   | 否       | 模型消息中思维链内容，仅doubao-seed-1.8、deepseek-v3.2、doubao-seed-2.0支持该字段                                            |
| tool_calls        | messages   | object[] | 否       | 模型消息中工具调用部分。                                                                                                   |
| function          | tool_calls | object   | 是       | 模型返回的需调用的函数信息。如果有tool_calls字段则function字段必填                                                         |
| name              | function   | string   | 是       | 需调用的函数的名称。                                                                                                       |
| arguments         | function   | string   | 是       | 需调用的函数的入参，JSON 格式。                                                                                             |
| id                | tool_calls | string   | 是       | 需调用的工具的 ID，由模型生成。                                                                                             |
| type              | tool_calls | string   | 是       | 消息类型，当前仅支持function                                                                                               |
| tool_call_id      | messages   | string   | 是       | 模型生成的需调用工具请求时，生成的ID。在程序调用工具的返回需要附上同一 ID，来关联工具结构与模型请求。避免多工具调用时混淆信息。 |
| thinking          | /          | object   | 否       | 控制模型是否开启深度思考模式。                                                                                             |
| type              | thinking   | string   | 是       | 取值范围：enabled， disabled，auto                                                                                           |
| stream            | /          | boolean  | 否       | 响应内容是否流式返回，默认false                                                                                            |

详细信息参考：https://www.volcengine.com/docs/82379/1494384?lang=zh






