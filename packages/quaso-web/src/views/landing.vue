<template>
  <div class="container">
    <div class="px-4 md:px-8 lg:px-24 xl:px-48">
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
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue"
import NewPost from "@/views/actions/new-post.vue"
import Surfing from "@/views/actions/surfing.vue"

const surfingInst = ref<any | null>(null)
</script>
