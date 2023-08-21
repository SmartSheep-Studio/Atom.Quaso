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
              <attachment-player v-for="img in item.attachments" :src="img" @click.prevent />
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
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-button
                      quaternary
                      size="small"
                      :type="item.is_liked ? 'primary' : 'empty'"
                      @click.stop="like(item)"
                    >
                      <template #icon>
                        <n-icon :component="ThumbUpRound" />
                      </template>
                      {{ item.like_count }}
                    </n-button>
                  </template>
                  Like
                </n-tooltip>
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-button
                      quaternary
                      size="small"
                      :type="item.is_disliked ? 'primary' : 'empty'"
                      @click.stop="dislike(item)"
                    >
                      <template #icon>
                        <n-icon :component="ThumbDownRound" />
                      </template>
                      {{ item.dislike_count }}
                    </n-button>
                  </template>
                  Dislike
                </n-tooltip>
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
import { ReplyRound, ShareRound, ThumbDownRound, ThumbUpRound } from "@vicons/material"
import { useMessage } from "naive-ui"
import { usePrincipal } from "@/stores/principal"
import { usePosts } from "@/stores/posts"
import { http } from "@/utils/http"
import VueMarkdown from "vue-markdown-render"
import AttachmentPlayer from "@/components/player/attachment-player.vue"
import { useI18n } from "vue-i18n"

const { t } = useI18n()

const emits = defineEmits(["reply", "share"])

defineExpose({ fetch })

const $posts = usePosts()
const $principal = usePrincipal()
const $message = useMessage()

const data = computed(() => $posts.data)

async function like(item: any) {
  try {
    $posts.isReverting = true
    const res = await http.post(`/api/posts/${item.id}/like`)
    await $posts.fetch()
    $message.success(res.status === 200 ? "Successfully liked" : "Successfully cancelled like")
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    $posts.isReverting = false
  }
}

async function dislike(item: any) {
  try {
    $posts.isReverting = true
    const res = await http.post(`/api/posts/${item.id}/dislike`)
    await $posts.fetch()
    $message.success(res.status === 200 ? "Successfully disliked" : "Successfully cancelled dislike")
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    $posts.isReverting = false
  }
}

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