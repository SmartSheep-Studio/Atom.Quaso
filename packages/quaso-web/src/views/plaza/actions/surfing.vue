<template>
  <div class="max-w-[100vw]">
    <n-spin :show="$posts.isReverting">
      <n-list v-if="data.posts.length > 0" bordered hoverable class="rounded-none">
        <n-list-item
          v-for="item in data.posts"
          @click="$router.push({ name: 'plaza.focus', params: { post: item.id } })"
        >
          <n-thing :title="item.account.nickname">
            <template #description>
              <n-space size="small">
                <n-tag
                  v-if="item.account.user_id === $principal.account.id"
                  size="small"
                  type="warning"
                  :bordered="false"
                >
                  You
                </n-tag>
                <n-tag class="capitalize" size="small" type="primary" :bordered="false">
                  {{ item.type }}
                </n-tag>
                <n-tag
                  v-for="tag in item.tags"
                  size="small"
                  type="info"
                  :bordered="false"
                >
                  {{ tag }}
                </n-tag>
              </n-space>
            </template>

            <div>
              <vue-markdown :source="item.content" />
            </div>

            <n-space vertical class="mt-4">
              <n-card content-style="padding: 0" class="max-w-[800px]" v-for="img in item.attachments">
                <attachment-player :src="img" />
              </n-card>
            </n-space>

            <div class="mt-2" v-if="item.belong_id != null">
              <n-alert class="post-reply-tips" :show-icon="false">
                This post is replying
                <router-link :to="{ name: 'plaza.focus', params: { post: item.belong_id } }" @click.stop>
                  #{{ item.belong_id }}
                </router-link>
              </n-alert>
            </div>

            <n-card size="small" class="mb-1 mt-2" content-style="padding: 8px" embedded>
              <div class="flex justify-around">
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-button quaternary size="small" @click.stop="emits('reply', item)">
                      <template #icon>
                        <n-icon :component="ReplyRound" />
                      </template>
                      {{ item.comment_count }}
                    </n-button>
                  </template>
                  Reply
                </n-tooltip>
                <n-button quaternary size="small" disabled>
                  <template #icon>
                    <n-icon :component="ThumbUpRound" />
                  </template>
                  Like
                </n-button>
                <n-button quaternary size="small" disabled>
                  <template #icon>
                    <n-icon :component="ThumbDownRound" />
                  </template>
                  Dislike
                </n-button>
                <n-button quaternary size="small" @click.stop="emits('share', item)">
                  <template #icon>
                    <n-icon :component="ShareRound" />
                  </template>
                  Share
                </n-button>
              </div>
            </n-card>
          </n-thing>
        </n-list-item>
      </n-list>
      <n-list bordered class="rounded-none" v-else>
        <n-list-item class="my-4">
          <n-empty description="There's no content for you." />
        </n-list-item>
      </n-list>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted } from "vue"
import { ReplyRound, ShareRound, ThumbDownRound, ThumbUpRound } from "@vicons/material"
import { usePrincipal } from "@/stores/principal"
import { usePosts } from "@/stores/posts"
import VueMarkdown from 'vue-markdown-render'
import AttachmentPlayer from "@/components/player/attachment-player.vue"

const emits = defineEmits(["reply", "share"])

defineExpose({ fetch })

const $posts = usePosts()
const $principal = usePrincipal()

const data = computed(() => $posts.data)

onMounted(() => {
  $posts.fetch()
})
</script>

<style>
.post-reply-tips .n-alert-body {
  padding: 12px;
}
</style>