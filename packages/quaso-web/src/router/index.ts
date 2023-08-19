import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "layout.main",
      component: () => import("@/layouts/main-layout.vue"),
      children: [
        {
          path: "/",
          name: "plaza",
          component: () => import("@/views/plaza.vue"),
        },
      ]
    },
  ],
})

export default router
