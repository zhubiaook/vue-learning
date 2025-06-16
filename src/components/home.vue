<template>
  <el-container class="h-screen">
    <el-aside class="w-1/4 bg-gray-200"
      >Aside
      <el-scrollbar></el-scrollbar>
    </el-aside>
    <el-container>
      <el-header class="bg-gray-100">Header</el-header>
      <el-main class="bg-white">
        <el-space direction="vertical" size="large" alignment="start">
          <div v-for="message in messageList" :key="message.id">
            <div v-if="message.role === 'user'">
              <div class="bg-blue-200 p-2 rounded-lg justify-end">
                {{ message.content }}
              </div>
            </div>
            <div v-if="message.role === 'assistant'">
              <div class="bg-gray-200 p-2 rounded-lg">
                <div v-for="content in message.content" :key="content">
                  {{ content }}
                </div>
              </div>
            </div>
          </div>
        </el-space>
      </el-main>
      <el-footer class="bg-gray-100 p-2">
        <el-row :gutter="20">
          <el-col :span="20">
            <el-input
              class="w-full"
              v-model="query.input"
              placeholder="请输入内容"
              @keyup.enter.prevent="send2"
            />
          </el-col>
          <el-col :span="4">
            <el-button type="primary" @click="send2">发送</el-button>
          </el-col>
        </el-row>
      </el-footer>
    </el-container>
  </el-container>
</template>

<script setup>
import { reactive } from "vue";

const messageList = reactive([]);
const assistantContent = reactive({
  role: "assistant",
  content: ["思考中..."],
});
const userContent = reactive({
  role: "user",
  content: [""],
});

const query = reactive({
  input: "",
});

import { sendMessage } from "/src/api";
const send = async () => {
  let message = query.input;
  userContent.content = [message];
  messageList.push(userContent);
  query.input = "";
  messageList.push(assistantContent);
  // sleep 2 seconds, then for loop mock the assistant content
  await new Promise((resolve) => setTimeout(resolve, 2000));
  for (let i = 0; i < 10; i++) {
    assistantContent.content.push("message " + i);
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
};

import { sendSSEMessage } from "../api";

const send2 = () => {
  let message = query.input;
  query.input = "";
  messageList.push(userContent);
  messageList.push(assistantContent);

  let eventSource = sendSSEMessage(message);

  eventSource.addEventListener("chat", (event) => {
    const data = JSON.parse(event.data);
    if (data.content.includes("\n")) {
      const parts = data.content.split("\n");
      assistantContent.content[assistantContent.content.length - 1] += parts[0];
      for (let i = 1; i < parts.length; i++) {
        assistantContent.content.push(parts[i]);
      }
    } else {
      assistantContent.content[assistantContent.content.length - 1] +=
        data.content;
    }
  });
  eventSource.onerror = (error) => {
    console.error(error);
    eventSource.close();
  };
};
</script>
