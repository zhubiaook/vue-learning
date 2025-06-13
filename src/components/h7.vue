<!-- The exported code uses Tailwind CSS. Install Tailwind CSS in your dev environment to ensure all styles work. -->
<template>
  <div class="app-container bg-white min-h-screen flex">
    <!-- 左侧导航栏 -->
    <div
      class="sidebar relative bg-gray-100 text-gray-800"
      :class="{ 'sidebar-collapsed': isSidebarCollapsed }"
      :style="{ width: sidebarWidth + 'px' }"
    >
      <div class="p-4 font-bold text-lg border-b border-gray-200">Octopus</div>
      <!-- 新建聊天按钮 -->
      <div class="p-3">
        <button
          class="w-full bg-blue-600 text-white py-2 px-4 rounded-lg flex items-center justify-center !rounded-button whitespace-nowrap cursor-pointer"
          @click="startNewChat"
        >
          <el-icon class="mr-2"><Plus /></el-icon>
          <span>新建聊天</span>
        </button>
      </div>
      <!-- 导航菜单 -->
      <div class="nav-menu">
        <div class="nav-group">
          <div
            class="nav-group-header cursor-pointer flex items-center justify-between p-2 hover:bg-gray-200"
            @click="toggleServers"
          >
            <span>Servers</span>
            <el-icon :class="{ 'transform rotate-180': !serversCollapsed }"
              ><ArrowDown
            /></el-icon>
          </div>
          <div v-show="!serversCollapsed" class="nav-group-content">
            <div
              class="nav-item px-4 py-2 flex justify-between items-center hover:bg-gray-200 cursor-pointer"
            >
              <span>testServer</span>
              <span class="text-xs bg-gray-300 px-1 rounded">1.8.9</span>
            </div>
          </div>
        </div>
        <div class="nav-group">
          <div
            class="nav-group-header cursor-pointer flex items-center justify-between p-2 hover:bg-gray-200"
            @click="toggleTools"
          >
            <span>Tools</span>
            <el-icon :class="{ 'transform rotate-180': !toolsCollapsed }"
              ><ArrowDown
            /></el-icon>
          </div>
        </div>
      </div>
      <!-- 聊天历史 -->
      <div class="chat-history mt-4 border-t border-gray-200 pt-2">
        <div class="px-3 py-2 text-sm font-medium text-gray-500">聊天历史</div>
        <div
          v-for="(chat, index) in chatHistory"
          :key="index"
          class="chat-item flex justify-between items-center p-3 hover:bg-gray-200 cursor-pointer"
          @click="selectChat(index)"
        >
          <span class="text-gray-700 truncate">{{ chat.title }}</span>
          <span class="text-gray-500 text-xs">{{ chat.count }}</span>
        </div>
      </div>
      <!-- 侧边栏调整手柄 -->
      <div
        class="resize-handle absolute top-0 right-0 w-1 h-full cursor-col-resize bg-blue-500 opacity-0 hover:opacity-50"
        @mousedown="startResizing"
      ></div>
      <!-- 折叠按钮 -->
      <div
        class="collapse-btn absolute bottom-4 right-2 cursor-pointer text-gray-600 hover:text-gray-900"
        @click="toggleSidebar"
      >
        <el-icon v-if="!isSidebarCollapsed"><Fold /></el-icon>
        <el-icon v-else><Expand /></el-icon>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="main-content flex-1 flex flex-col">
      <!-- 聊天内容区域 -->
      <div class="chat-content flex-1 flex flex-col p-6 bg-white">
        <div class="messages flex-1 overflow-y-auto">
          <div
            v-if="currentMessages.length === 0"
            class="flex items-center justify-center h-full"
          >
            <div class="text-gray-500 text-center">
              <img
                :src="welcomeImage"
                alt="Welcome"
                class="w-64 h-64 mx-auto mb-4 object-cover"
              />
              <h2 class="text-2xl mb-4">欢迎使用 MCP 智能助手</h2>
              <p>开始一个新的对话吧</p>
            </div>
          </div>
          <div v-else>
            <div
              v-for="(message, index) in currentMessages"
              :key="index"
              class="message mb-6"
            >
              <div
                class="flex"
                :class="
                  message.sender === 'user' ? 'justify-end' : 'justify-start'
                "
              >
                <div
                  class="message-content max-w-3xl relative"
                  :class="message.sender === 'user' ? 'order-1' : 'order-2'"
                >
                  <div
                    class="message-bubble p-3 rounded-lg"
                    :class="
                      message.sender === 'user'
                        ? 'bg-blue-600 text-white'
                        : 'bg-white text-gray-800'
                    "
                  >
                    <div
                      v-if="message.sender === 'bot'"
                      v-html="renderMarkdown(message.content)"
                      class="markdown-content"
                    ></div>
                    <div v-else>{{ message.content }}</div>
                    <button
                      v-if="message.sender === 'bot'"
                      @click="copyMessage(message.content)"
                      class="absolute top-2 right-2 text-gray-500 hover:text-gray-700 p-1 rounded"
                    >
                      <el-icon><Document /></el-icon>
                    </button>
                  </div>
                </div>
                <div
                  class="message-avatar mx-2"
                  :class="message.sender === 'user' ? 'order-2' : 'order-1'"
                >
                  <div
                    class="w-8 h-8 rounded-full flex items-center justify-center"
                    :class="
                      message.sender === 'bot'
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-300'
                    "
                  >
                    <el-icon v-if="message.sender === 'bot'"
                      ><ChatDotRound
                    /></el-icon>
                    <el-icon v-else><User /></el-icon>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- 输入区域 -->
        <div class="chat-input-container mt-4 border-t border-gray-200 pt-4">
          <div class="relative">
            <textarea
              v-model="messageInput"
              class="w-full bg-gray-100 text-gray-800 rounded-lg p-4 pr-16 resize-none focus:outline-none focus:ring-1 focus:ring-blue-500 border-none"
              placeholder="输入您的问题..."
              rows="2"
              @keydown.enter.prevent="sendMessage"
            ></textarea>
            <div class="absolute right-2 bottom-2 flex items-center">
              <div class="text-xs text-gray-500 mr-2">Shift+Enter 换行</div>
              <button
                @click="sendMessage"
                class="bg-blue-600 text-white p-2 rounded-lg !rounded-button whitespace-nowrap cursor-pointer"
                :disabled="!messageInput.trim()"
              >
                <el-icon><Position /></el-icon>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from "vue";
import { ElMessage } from "element-plus";
import {
  ArrowDown,
  Fold,
  Expand,
  ChatDotRound,
  User,
  Position,
  Plus,
  Document,
} from "@element-plus/icons-vue";
import { marked } from "marked";
import hljs from "highlight.js";
import "highlight.js/styles/github-dark.css";

// 侧边栏状态
const sidebarWidth = ref(240);
const isSidebarCollapsed = ref(false);
const minSidebarWidth = 150;
const maxSidebarWidth = 400;
const isResizing = ref(false);
// 菜单折叠状态
const serversCollapsed = ref(false);
const toolsCollapsed = ref(true);
const resourcesCollapsed = ref(true);
const promptsCollapsed = ref(true);
// 聊天相关
const messageInput = ref("");
const currentMessages = ref<
  { sender: string; content: string; timestamp: number }[]
>([]);
const chatHistory = ref([
  { title: "Good ones", count: 0 },
  { title: "Is there?", count: 0 },
  { title: "Hello chat", count: 0 },
  { title: "Hello there", count: 0 },
]);
// 欢迎图片
const welcomeImage =
  "https://readdy.ai/api/search-image?query=A%20modern%2C%20clean%2C%20minimalist%20illustration%20of%20a%20friendly%20AI%20assistant%20or%20chatbot%2C%20with%20soft%20blue%20and%20white%20colors%2C%20simple%20geometric%20shapes%2C%20and%20a%20welcoming%20design.%20The%20image%20should%20have%20a%20light%20background%20with%20subtle%20patterns%2C%20perfect%20for%20a%20chat%20application%20welcome%20screen&width=512&height=512&seq=welcome1&orientation=squarish";

// 只需加 as any 绕过类型检查
marked.setOptions({
  highlight: function (code: string, lang: string) {
    if (lang && hljs.getLanguage(lang)) {
      return hljs.highlight(code, { language: lang }).value;
    }
    return hljs.highlightAuto(code).value;
  },
  breaks: true,
} as any);

// 顶层定义，暴露给模板
const renderMarkdown = (text: string) => {
  return marked.parse(text);
};

// 侧边栏调整函数
const startResizing = (e: MouseEvent) => {
  isResizing.value = true;
  document.addEventListener("mousemove", handleMouseMove);
  document.addEventListener("mouseup", stopResizing);
  e.preventDefault();
};
const handleMouseMove = (e: MouseEvent) => {
  if (!isResizing.value) return;
  const newWidth = e.clientX;
  if (newWidth >= minSidebarWidth && newWidth <= maxSidebarWidth) {
    sidebarWidth.value = newWidth;
  }
};
const stopResizing = () => {
  isResizing.value = false;
  document.removeEventListener("mousemove", handleMouseMove);
  document.removeEventListener("mouseup", stopResizing);
};
// 折叠/展开侧边栏
const toggleSidebar = () => {
  isSidebarCollapsed.value = !isSidebarCollapsed.value;
};
// 折叠/展开菜单组
const toggleServers = () => {
  serversCollapsed.value = !serversCollapsed.value;
};
const toggleTools = () => {
  toolsCollapsed.value = !toolsCollapsed.value;
};
const toggleResources = () => {
  resourcesCollapsed.value = !resourcesCollapsed.value;
};
const togglePrompts = () => {
  promptsCollapsed.value = !promptsCollapsed.value;
};
// 开始新聊天
const startNewChat = () => {
  currentMessages.value = [];
  messageInput.value = "";
};
// 选择聊天
const selectChat = (index: number) => {
  // 模拟选择历史聊天
  startNewChat();
  const selectedChat = chatHistory.value[index];
  // 模拟加载历史消息
  currentMessages.value = [
    {
      sender: "user",
      content: selectedChat.title,
      timestamp: Date.now() - 1000,
    },
    {
      sender: "bot",
      content: `这是关于 "${selectedChat.title}" 的历史对话。`,
      timestamp: Date.now(),
    },
  ];
};
// 发送消息
const sendMessage = () => {
  if (!messageInput.value.trim()) return;
  // 添加用户消息
  currentMessages.value.push({
    sender: "user",
    content: messageInput.value,
    timestamp: Date.now(),
  });
  const userMessage = messageInput.value;
  messageInput.value = "";
  // 使用模拟SSE响应
  simulateSSEResponse(userMessage);
};
// 模拟SSE流式响应
const simulateSSEResponse = (userMessage: string) => {
  // 创建一个初始的空消息
  const botMessageIndex = currentMessages.value.length;
  currentMessages.value.push({
    sender: "bot",
    content: "",
    timestamp: Date.now(),
  });
  // 根据用户输入生成一些示例响应
  let fullResponse = "";
  // 默认响应
  fullResponse = "正在思考...";
  // 根据用户输入匹配不同的响应
  if (userMessage.toLowerCase().includes("v3")) {
    fullResponse = `# V3 版本更新说明
    ## 主要功能更新
    1. 全新的用户界面
    - 更直观的操作体验
    - 优化的视觉设计
    - 响应式布局适配
    2. 性能优化
    - 更快的响应速度
    - 更低的资源占用
    - 更好的兼容性
    3. 新增特性
    \`\`\`javascript
    // 示例代码
    const features = {
    ui: 'Modern Design',
    performance: 'Optimized',
    compatibility: 'Enhanced'
    };
    function checkVersion() {
    return 'V3.0.0';
    }
    \`\`\`
    4. 技术栈升级
    - Vue 3 + TypeScript
    - Tailwind CSS
    - Element Plus
    5. 开发者工具
    - 更完善的API文档
    - 更强大的调试功能
    - 更多的开发示例
    ## 即将推出
    - 更多的主题选项
    - 插件系统
    - 自定义工作流
    > 欢迎体验新版本，如有问题请随时反馈！`;
  } else if (userMessage.toLowerCase().includes("hello")) {
    fullResponse =
      "# 你好！\n\n很高兴与你交流。我是一个AI助手，可以帮助你回答问题、提供信息或者进行有趣的对话。\n\n## 我能做什么？\n\n- 回答问题\n- 提供信息\n- 编写代码\n- 创意写作\n\n```javascript\n// 这是一个简单的JavaScript示例\nfunction sayHello(name) {\n  return `Hello, ${name}!`;\n}\n\nconsole.log(sayHello('World'));\n```";
  } else if (
    userMessage.toLowerCase().includes("code") ||
    userMessage.toLowerCase().includes("编程")
  ) {
    fullResponse =
      "# 编程示例\n\n以下是几种常见编程语言的Hello World示例：\n\n## Python\n```python\nprint('Hello, World!')\n```\n\n## JavaScript\n```javascript\nconsole.log('Hello, World!');\n```\n\n## Java\n```java\npublic class HelloWorld {\n  public static void main(String[] args) {\n    System.out.println(\"Hello, World!\");\n  }\n}\n```";
  } else {
    fullResponse =
      "感谢你的消息！我是MCP的AI助手，很高兴能帮助你。请告诉我你需要什么帮助，我会尽力提供支持。\n\n如果你需要了解特定功能或有任何问题，请随时告诉我。";
  }

  // 确保响应不为空
  if (!fullResponse.trim()) {
    fullResponse =
      "我没有完全理解你的问题。请尝试用不同的方式描述，我会尽力帮助你。";
  }
  // 模拟流式响应
  let currentIndex = 0;
  const streamInterval = setInterval(() => {
    if (currentIndex < fullResponse.length) {
      // 每次添加1-3个字符，模拟打字效果
      const chunkSize = Math.floor(Math.random() * 3) + 1;
      const nextChunk = fullResponse.substring(
        currentIndex,
        currentIndex + chunkSize
      );
      currentIndex += chunkSize;
      // 更新当前消息内容
      currentMessages.value[botMessageIndex].content = fullResponse.substring(
        0,
        currentIndex
      );
    } else {
      clearInterval(streamInterval);
    }
  }, 30);
};
// 复制消息内容
const copyMessage = (content: string) => {
  navigator.clipboard
    .writeText(content)
    .then(() => {
      ElMessage({
        message: "复制成功",
        type: "success",
        duration: 2000,
      });
    })
    .catch(() => {
      ElMessage({
        message: "复制失败",
        type: "error",
        duration: 2000,
      });
    });
};
// 组件挂载和卸载时的事件处理
onMounted(() => {
  window.addEventListener("resize", handleWindowResize);
});
onUnmounted(() => {
  window.removeEventListener("resize", handleWindowResize);
});
const handleWindowResize = () => {
  // 在小屏幕下自动折叠侧边栏
  if (window.innerWidth < 768 && !isSidebarCollapsed.value) {
    isSidebarCollapsed.value = true;
  }
};
</script>
<style scoped>
.app-container {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  background-color: #f7f7f8;
}
.sidebar {
  transition: width 0.3s ease;
  height: 100vh;
  overflow-y: auto;
  z-index: 10;
}
.sidebar-collapsed {
  width: 50px !important;
}
.sidebar-collapsed .nav-group-header span,
.sidebar-collapsed .nav-group-content,
.sidebar-collapsed .nav-group-header el-icon,
.sidebar-collapsed .chat-history,
.sidebar-collapsed button span {
  display: none;
}
.resize-handle {
  transition: opacity 0.2s;
}
.nav-group-header {
  transition: background-color 0.2s;
}
.chat-content {
  height: calc(100vh - 120px);
  overflow-y: auto;
}
.message-bubble {
  max-width: 100%;
  word-wrap: break-word;
}
/* 自定义样式，避免与组件库冲突 */
.mcp-chat-ui .markdown-content {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
    Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
}
/* 确保输入框没有默认边框 */
textarea:focus {
  outline: none;
}
/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
::-webkit-scrollbar-track {
  background: #f1f1f1;
}
::-webkit-scrollbar-thumb {
  background: #ccc;
  border-radius: 3px;
}
::-webkit-scrollbar-thumb:hover {
  background: #aaa;
}
/* 确保markdown内容正确显示 */
:deep(.markdown-content) {
  color: inherit;
}
:deep(.markdown-content h1) {
  font-size: 1.8rem;
  margin-top: 1rem;
  margin-bottom: 1rem;
  font-weight: 600;
}
:deep(.markdown-content h2) {
  font-size: 1.5rem;
  margin-top: 0.8rem;
  margin-bottom: 0.8rem;
  font-weight: 600;
}
:deep(.markdown-content p) {
  margin-bottom: 1rem;
}
:deep(.markdown-content pre) {
  background-color: #282c34;
  padding: 1rem;
  border-radius: 4px;
  overflow-x: auto;
  margin: 1rem 0;
  color: #abb2bf;
}
:deep(.markdown-content code) {
  font-family: "Courier New", Courier, monospace;
  font-size: 0.9rem;
}
:deep(.markdown-content ul),
:deep(.markdown-content ol) {
  margin-left: 1.5rem;
  margin-bottom: 1rem;
}
:deep(.markdown-content li) {
  margin-bottom: 0.5rem;
}
</style>
