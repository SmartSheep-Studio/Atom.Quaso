<template>
  <div>
    <div class="overflow-auto hide-scrollbar h-[100vh]">
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