<template>
  <div>
    <router-view @post="triggerPost" :key="$route.fullPath" />

    <n-modal v-model:show="posting" class="w-[680px]" display-directive="show">
      <new-post ref="postingInst" @submit="afterPost" />
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

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { usePosts } from "@/stores/posts"
import { PostAddRound } from "@vicons/material"
import NewPost from "@/views/plaza/actions/post.vue"

const $posts = usePosts()

const posting = ref(false)

const postingInst = ref<any | null>(null)

function triggerPost(parent: any) {
  posting.value = true
  setTimeout(() => {
    postingInst.value.trigger(parent)
  }, 250)
}

function afterPost() {
  $posts.fetch()
  posting.value = false
}
</script>