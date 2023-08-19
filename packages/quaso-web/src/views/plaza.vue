<template>
  <div>
    <div class="h-[100vh] overflow-auto hide-scrollbar">
      <surfing ref="surfingInst" @reply="replyPost" />
    </div>

    <n-modal v-model:show="posting" class="w-[680px]" display-directive="show">
      <new-post ref="postingInst" @submitted="surfingInst.fetch()" />
    </n-modal>

    <div class="fixed bottom-[24px] right-[24px]">
      <n-button circle type="primary" size="large" class="p-[24px] shadow-2xl" @click="posting = true">
        <template #icon>
          <n-icon :component="PostAddRound" />
        </template>
      </n-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue"
import { PostAddRound } from "@vicons/material"
import NewPost from "@/views/actions/posts/new.vue"
import Surfing from "@/views/actions/posts/surfing.vue"

const posting = ref(false)

const postingInst = ref<any | null>(null)
const surfingInst = ref<any | null>(null)

function replyPost(parent: any) {
  posting.value = true
  setTimeout(() => {
    postingInst.value.reply(parent)
  }, 250)
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