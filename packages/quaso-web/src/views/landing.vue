<template>
  <div>
    <div class="pt-4">
      <new-post @submitted="surfingInst.fetch()" />
    </div>

    <div class="pt-4">
      <surfing ref="surfingInst" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue"
import NewPost from "@/views/actions/new-post.vue"
import Surfing from "@/views/actions/surfing.vue"

const surfingInst = ref<any | null>(null)
</script>
