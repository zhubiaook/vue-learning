import { createRouter, createWebHistory } from "vue-router";

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      component: () => import("../components/home.vue"),
    },
    {
      path: "/h7",
      component: () => import("../components/h7.vue"),
    },
    {
      path: "/ops",
      component: () => import("../components/operation.vue"),
    },
    {
      path: "/chat",
      component: () => import("../components/chat.vue"),
    },
  ],
});
