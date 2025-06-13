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
                {{ message.content }}
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
              @keyup.enter.prevent="send"
            />
          </el-col>
          <el-col :span="4">
            <el-button type="primary" @click="send">发送</el-button>
          </el-col>
        </el-row>
      </el-footer>
    </el-container>
  </el-container>
</template>

<script setup>
import { reactive } from "vue";

const messageList = reactive([]);

const query = reactive({
  input: "",
});

import { sendMessage } from "/src/api";
const send = async () => {
  let message = query.input;
  messageList.push(
    {
      role: "user",
      content: message,
    },
    {
      role: "assistant",
      content: "思考中...",
    }
  );
  query.input = "";
  const result = await sendMessage(message);
  messageList.pop();
  messageList.push({
    role: "assistant",
    content: result.data,
  });
};
</script>
