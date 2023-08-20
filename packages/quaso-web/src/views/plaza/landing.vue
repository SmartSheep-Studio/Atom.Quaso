<template>
  <div>
    <div :class="isUnderShadow ? 'h-max' : 'h-full'" class="overflow-auto hide-scrollbar">
      <surfing ref="surfingInst" @reply="replyPost" @share="sharePost" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue"
import Surfing from "@/views/plaza/actions/surfing.vue"

const emits = defineEmits(["post"])

const surfingInst = ref<any | null>(null)

function replyPost(parent: any) {
  emits("post", {
    belong_to: parent
  })
}

function sharePost(parent: any) {
  emits("post", {
    content: `Shared a post [${parent.content.substring(0, 5)}...](/srv/subapps/quaso/plaza/${parent.id})`
  })
}

// Use for dynamic calculate height
const isUnderShadow = computed(() => {
  return (window as any).__POWERED_BY_WUJIE__ != null
})
</script>

<style>
.hide-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.hide-scrollbar::-webkit-scrollbar {
  display: none;
}
</style>