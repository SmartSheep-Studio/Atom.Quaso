<template>
  <div class="lg:container px-2">
    <n-grid item-responsive responsive="screen" :x-gap="8" :y-gap="8">
      <n-gi span="0 m:6">
        <div>
          <div class="flex gap-2 items-center px-4 pt-4 pb-2">
            <div>
              <img :src="Logo" width="64" height="64" />
            </div>
            <div>
              <div class="text-lg">Quaso</div>
              <div class="mt-[-4px] text-xs">Know what's new in the world</div>
            </div>
          </div>

          <n-menu
            v-model:value="menuKey"
            :options="menuOptions"
          />
        </div>
      </n-gi>

      <n-gi span="24 m:12">
        <router-view />
      </n-gi>

      <n-gi span="0 m:6"></n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { Component, computed, h, Ref, ref, watch } from "vue"
import { MenuOption, NIcon } from "naive-ui"
import { RouterLink, useRoute } from "vue-router"
import { ExploreRound } from "@vicons/material"
import { usePrincipal } from "@/stores/principal"
import { useI18n } from "vue-i18n"
import Logo from "@/assets/icon.png"

const { t } = useI18n()

const $route = useRoute()
const $principal = usePrincipal()

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuKey = ref($route.name)
const menuOptions: Ref<MenuOption[]> = computed(() =>
  $principal.isSigned
    ? [
      {
        label: () => h(RouterLink, { to: { name: "plaza" } }, { default: () => t("nav.plaza") }),
        icon: renderIcon(ExploreRound),
        key: "plaza"
      }
    ]
    : []
)

watch($route, (v) => {
  menuKey.value = v.name
})
</script>