<template>
  <div class="docs-page">
    <header class="page-header">
      <div class="logo" @click="$router.push('/')">
        <el-icon><Connection /></el-icon>
        <span>Token Hub</span>
      </div>
      <nav class="nav">
        <el-button @click="$router.push('/')">首页</el-button>
        <el-button @click="$router.push('/models')">模型列表</el-button>
        <el-button type="primary">接入说明</el-button>
        <template v-if="userStore.isAuthenticated">
          <el-button @click="$router.push('/user/dashboard')">控制台</el-button>
        </template>
        <template v-else>
          <el-button @click="$router.push('/login')">登录</el-button>
        </template>
      </nav>
    </header>

    <main class="main-content">
      <div class="page-title">
        <h1>API接入说明</h1>
        <p>Token Hub 提供与 OpenAI API 兼容的接口，您可以无缝迁移现有代码</p>
      </div>

      <div class="docs-container">
        <aside class="sidebar">
          <el-menu :default-active="activeSection" class="sidebar-menu">
            <el-menu-item index="quickstart">快速开始</el-menu-item>
            <el-menu-item index="authentication">认证方式</el-menu-item>
            <el-menu-item index="chat-completion">聊天对话</el-menu-item>
            <el-menu-item index="completion">文本补全</el-menu-item>
            <el-menu-item index="models">获取模型列表</el-menu-item>
            <el-menu-item index="error-handling">错误处理</el-menu-item>
          </el-menu>
        </aside>

        <div class="content">
          <section id="quickstart" class="doc-section">
            <h2>快速开始</h2>
            <p>使用 Token Hub API 非常简单，只需要以下几个步骤：</p>
            <ol>
              <li>注册账号并登录</li>
              <li>在控制台创建 API 密钥</li>
              <li>充值获取积分</li>
              <li>使用 API 密钥调用接口</li>
            </ol>
            
            <h3>API 基础地址</h3>
            <div class="code-block">
              <code>https://api.your-domain.com/api/v1</code>
            </div>
            <p>或者在开发环境中使用：</p>
            <div class="code-block">
              <code>http://localhost:8080/api/v1</code>
            </div>
          </section>

          <section id="authentication" class="doc-section">
            <h2>认证方式</h2>
            <p>所有 API 请求都需要在 HTTP 请求头中携带 API 密钥进行认证。</p>
            
            <h3>请求头格式</h3>
            <div class="code-block">
              <pre>
Authorization: Bearer your_api_key
              </pre>
            </div>

            <h3>示例代码</h3>
            <div class="code-block">
              <pre>
import openai

openai.api_key = "your_api_key"
openai.api_base = "http://localhost:8080/api/v1"

response = openai.ChatCompletion.create(
    model="gpt-3.5-turbo",
    messages=[{"role": "user", "content": "Hello"}]
)
print(response)
              </pre>
            </div>

            <div class="notice">
              <el-icon><InfoFilled /></el-icon>
              <span>API 密钥可以在用户控制台中创建和管理，请妥善保管，不要泄露给他人。</span>
            </div>
          </section>

          <section id="chat-completion" class="doc-section">
            <h2>聊天对话接口</h2>
            <p>用于进行多轮对话，支持上下文记忆。</p>
            
            <h3>请求参数</h3>
            <el-table :data="chatParams" style="width: 100%" border>
              <el-table-column prop="name" label="参数名" width="150" />
              <el-table-column prop="type" label="类型" width="100" />
              <el-table-column prop="required" label="必填" width="80" />
              <el-table-column prop="desc" label="说明" />
            </el-table>

            <h3>请求示例</h3>
            <div class="code-block">
              <pre>
curl http://localhost:8080/api/v1/chat/completions \
  -H "Authorization: Bearer your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "system", "content": "You are a helpful assistant."},
      {"role": "user", "content": "Hello!"}
    ],
    "temperature": 0.7,
    "max_tokens": 1000
  }'
              </pre>
            </div>

            <h3>响应示例</h3>
            <div class="code-block">
              <pre>
{
  "id": "chatcmpl-xxxx",
  "object": "chat.completion",
  "created": 1700000000,
  "model": "gpt-3.5-turbo",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello! How can I help you today?"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 20,
    "completion_tokens": 15,
    "total_tokens": 35
  }
}
              </pre>
            </div>
          </section>

          <section id="completion" class="doc-section">
            <h2>文本补全接口</h2>
            <p>用于传统的文本补全任务。</p>
            
            <h3>请求示例</h3>
            <div class="code-block">
              <pre>
curl http://localhost:8080/api/v1/completions \
  -H "Authorization: Bearer your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo-instruct",
    "prompt": "Once upon a time",
    "max_tokens": 100,
    "temperature": 0.7
  }'
              </pre>
            </div>
          </section>

          <section id="models" class="doc-section">
            <h2>获取模型列表</h2>
            <p>获取当前可用的所有模型列表。</p>
            
            <h3>请求示例</h3>
            <div class="code-block">
              <pre>
curl http://localhost:8080/api/v1/models \
  -H "Authorization: Bearer your_api_key"
              </pre>
            </div>

            <h3>响应示例</h3>
            <div class="code-block">
              <pre>
{
  "object": "list",
  "data": [
    {
      "id": "gpt-3.5-turbo",
      "object": "model",
      "created": 1700000000,
      "owned_by": "openai",
      "provider": "OpenAI",
      "description": "GPT-3.5 Turbo 模型"
    },
    {
      "id": "claude-3-opus",
      "object": "model",
      "created": 1700000000,
      "owned_by": "anthropic",
      "provider": "Anthropic",
      "description": "Claude 3 Opus 模型"
    }
  ]
}
              </pre>
            </div>
          </section>

          <section id="error-handling" class="doc-section">
            <h2>错误处理</h2>
            <p>API 请求失败时会返回相应的错误信息。</p>
            
            <h3>错误响应格式</h3>
            <div class="code-block">
              <pre>
{
  "code": 400,
  "message": "错误描述",
  "data": null
}
              </pre>
            </div>

            <h3>常见错误码</h3>
            <el-table :data="errorCodes" style="width: 100%" border>
              <el-table-column prop="code" label="状态码" width="100" />
              <el-table-column prop="desc" label="说明" />
            </el-table>
          </section>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
const activeSection = ref('quickstart')

const chatParams = [
  { name: 'model', type: 'string', required: '是', desc: '模型名称，如 gpt-3.5-turbo' },
  { name: 'messages', type: 'array', required: '是', desc: '消息列表，包含 role 和 content' },
  { name: 'temperature', type: 'number', required: '否', desc: '采样温度，0-2 之间' },
  { name: 'max_tokens', type: 'integer', required: '否', desc: '最大生成 token 数' },
  { name: 'stream', type: 'boolean', required: '否', desc: '是否流式响应' },
]

const errorCodes = [
  { code: 400, desc: '请求参数错误' },
  { code: 401, desc: '认证失败，API 密钥无效' },
  { code: 403, desc: '权限不足' },
  { code: 404, desc: '资源不存在' },
  { code: 429, desc: '请求过于频繁' },
  { code: 500, desc: '服务器内部错误' },
]
</script>

<style scoped lang="scss">
.docs-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.page-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  height: 70px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1000;

  .logo {
    display: flex;
    align-items: center;
    font-size: 20px;
    font-weight: bold;
    color: #409eff;
    cursor: pointer;

    .el-icon {
      font-size: 24px;
      margin-right: 8px;
    }
  }

  .nav {
    .el-button {
      margin-left: 10px;
    }
  }
}

.main-content {
  padding: 100px 40px 60px;

  .page-title {
    text-align: center;
    margin-bottom: 40px;

    h1 {
      font-size: 36px;
      margin-bottom: 10px;
    }

    p {
      color: #666;
      font-size: 16px;
    }
  }
}

.docs-container {
  display: flex;
  max-width: 1400px;
  margin: 0 auto;
  gap: 40px;
}

.sidebar {
  width: 250px;
  flex-shrink: 0;

  .sidebar-menu {
    position: sticky;
    top: 100px;
    border-radius: 8px;
  }
}

.content {
  flex: 1;
  min-width: 0;

  .doc-section {
    background: #fff;
    border-radius: 12px;
    padding: 30px;
    margin-bottom: 30px;

    h2 {
      font-size: 24px;
      margin-bottom: 20px;
      padding-bottom: 10px;
      border-bottom: 2px solid #409eff;
    }

    h3 {
      font-size: 18px;
      margin: 20px 0 15px;
      color: #303133;
    }

    p {
      color: #606266;
      line-height: 1.8;
      margin-bottom: 15px;
    }

    ol {
      color: #606266;
      line-height: 2;
      padding-left: 20px;
      margin-bottom: 20px;
    }

    .code-block {
      background: #1e1e1e;
      border-radius: 8px;
      padding: 16px;
      margin: 15px 0;
      overflow-x: auto;

      pre, code {
        color: #d4d4d4;
        font-family: 'Consolas', 'Monaco', monospace;
        font-size: 14px;
        line-height: 1.6;
        white-space: pre-wrap;
        word-wrap: break-word;
      }
    }

    .notice {
      display: flex;
      align-items: flex-start;
      gap: 10px;
      padding: 15px;
      background: #ecf5ff;
      border-left: 4px solid #409eff;
      border-radius: 4px;
      margin: 20px 0;
      color: #409eff;

      .el-icon {
        font-size: 20px;
        flex-shrink: 0;
      }
    }
  }
}
</style>
