import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      redirect: { name: "plaza" }
    },
    {
      path: "/",
      name: "layout.main",
      component: () => import("@/layouts/main-layout.vue"),
      children: [
        {
          path: "/",
          name: "layout.plaza",
          component: () => import("@/layouts/plaza-layout.vue"),
          children: [
            {
              path: "/plaza",
              name: "plaza",
              component: () => import("@/views/plaza/landing.vue")
            },
            {
              path: "/plaza/:post",
              name: "plaza.focus",
              component: () => import("@/views/plaza/focus.vue")
            },
            {
              path: "/accounts/:account",
              name: "plaza.account.info",
              component: () => import("@/views/account/info.vue")
            }
          ]
        }
      ]
    }
  ]
})

export default router
