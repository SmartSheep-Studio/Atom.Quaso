<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-dialog-provider>
      <n-message-provider>
        <div class="w-full relative h-[100vh]">
          <n-layout has-sider position="absolute">
            <n-layout class="w-full h-full">
              <data-provider>
                <gatekeeper>
                  <router-view />
                </gatekeeper>
              </data-provider>
            </n-layout>
          </n-layout>
        </div>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script lang="ts" setup>
import DataProvider from "@/data-provider.vue";
import Gatekeeper from "@/components/global/gatekeeper.vue";

const themeOverrides = {
  common: {
    primaryColor: "#ca4d4dFF",
    primaryColorHover: "#DF5656FF",
    primaryColorPressed: "#C04747FF",
    primaryColorSuppl: "#A84141FF"
  }
}
</script>

<style>
.n-layout-header {
  height: 72px;
}

.n-layout-footer {
  padding: 24px;
}

.w-dialog {
  width: 520px;
}

a {
  color: #329fd2;
}
</style>
