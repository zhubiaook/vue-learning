<template>
  <div class="p-4">
    <ul class="bg-gray-100 p-4 rounded space-y-1">
      <li v-for="(line, index) in messages" :key="index">{{ line }}</li>
    </ul>
    <el-space>
      <el-input v-model="input" />
      <el-button @click="sendMessage">Send</el-button>
    </el-space>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref } from "vue";

const messages = ref<string[]>([""]);
let eventSource: EventSource | null = null;
const input = ref("");

import { sendSSEMessage } from "../api";

const sendMessage = () => {
  eventSource = sendSSEMessage(input.value);

  eventSource.addEventListener("chat", (event) => {
    try {
      const data = JSON.parse(event.data);
      const content = data.content;

      if (content.includes("\n")) {
        const parts = content.split("\n");
        messages.value[messages.value.length - 1] += parts[0];
        for (let i = 1; i < parts.length; i++) {
          messages.value.push(parts[i]);
        }
      } else {
        messages.value[messages.value.length - 1] += content;
      }
    } catch (err) {
      console.error("Invalid SSE data:", event.data);
    }
  });

  eventSource.onerror = (err) => {
    console.error("SSE error:", err);
    eventSource?.close();
  };
};

onBeforeUnmount(() => {
  eventSource?.close();
});
</script>

<style scoped></style>
