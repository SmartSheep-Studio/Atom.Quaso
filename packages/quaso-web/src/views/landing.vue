<template>
  <div class="container">
    <n-grid item-responsive responsive="screen" :x-gap="8" :y-gap="8">
      <n-gi span="0 m:4 l:6"></n-gi>

      <n-gi span="24 m:16 l:12">
        <div class="pt-12 px-8 flex items-center">
          <div class="mt-1 pl-4">
            <div class="text-4xl font-bold">Quaso</div>
            <div class="text-lg">Know what's new in the world.</div>
          </div>
        </div>

        <div class="pt-10 px-8">
          <new-post @submitted="surfingInst.fetch()" />
        </div>

        <div class="pt-4 px-8">
          <surfing ref="surfingInst" />
        </div>
      </n-gi>

      <n-gi span="0 m:4 l:6"></n-gi>
    </n-grid>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue"
import NewPost from "@/views/actions/new-post.vue"
import Surfing from "@/views/actions/surfing.vue"

const surfingInst = ref<any | null>(null)
</script>
