<template>
  <div :class="isUnderShadow ? 'max-h-max' : 'max-h-full'">
    <n-card size="small" class="rounded-none" style="border-bottom: 0">
      <div class="flex">
        <n-select
          v-model:value="$posts.filterOptions.type"
          :options="[
            { label: 'All', value: 'all' },
            { label: 'Text', value: 'text' },
            { label: 'Image', value: 'image' },
            { label: 'Audio', value: 'audio' },
            { label: 'Video', value: 'video' },
          ]"
        />
      </div>
    </n-card>

    <n-spin :show="$posts.isReverting">
      <n-list v-if="data.posts.length > 0" bordered hoverable class="rounded-none" style="--n-border-radius: 0">
        <n-list-item
          v-for="item in data.posts"
          @click="$router.push({ name: 'plaza.focus', params: { post: item.id } })"
        >
          <post-player
            :post="item"
            @reply="emits('reply', item)"
            @share="emits('share', item)"
          />
        </n-list-item>
      </n-list>
      <n-list bordered class="rounded-none" v-else>
        <n-list-item class="my-4">
          <n-empty description="There's no content for you." />
        </n-list-item>
      </n-list>

      <n-card size="small" class="rounded-none" style="border-top: 0">
        <div class="flex justify-center">
          <n-pagination v-model:page="$posts.filterOptions.page" :page-count="$posts.totalPage" />
        </div>
      </n-card>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted } from "vue"
import { usePosts } from "@/stores/posts"
import PostPlayer from "@/components/player/post-player.vue"

const emits = defineEmits(["reply", "share"])

defineExpose({ fetch })

const $posts = usePosts()

const data = computed(() => $posts.data)

onMounted(() => {
  $posts.fetch()
})

// Use for dynamic calculate height
const isUnderShadow = computed(() => {
  return (window as any).__POWERED_BY_WUJIE__ != null
})
</script>

<style>
.post-reply-tips .n-alert-body {
  padding: 12px;
}

.max-h-full {
  height: 100vh;
}

.max-h-max {
  height: calc(100vh - 72px);
}
</style>