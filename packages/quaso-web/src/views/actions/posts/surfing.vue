<template>
  <div class="max-w-[100vw]">
    <n-spin :show="reverting">
      <n-list bordered v-if="data.posts.length > 0">
        <n-list-item v-for="item in data.posts">
          <n-thing :title="data.related_authors[item.account_id].nickname">
            <template #description>
              <n-space size="small">
                <n-tag
                  v-if="data.related_authors[item.account_id].user_id === $principal.account.id"
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

            <div>{{ item.content }}</div>

            <n-space vertical class="mt-4">
              <n-card content-style="padding: 0" class="max-w-[800px]" v-for="img in item.attachments">
                <n-image
                  object-fit="cover"
                  class="post-image block"
                  :src="`/srv/subapps/quaso${img}`"
                />
              </n-card>
            </n-space>

            <div class="mt-2" v-if="item.belong_id != null">
              <n-alert class="post-reply-tips" :show-icon="false">
                This post is replying #{{ item.belong_id }}
              </n-alert>
            </div>

            <n-card size="small" class="mb-1 mt-2" content-style="padding: 8px" embedded>
              <div class="flex justify-around">
                <n-button quaternary size="small" @click="emits('reply', item)">
                  <template #icon>
                    <n-icon :component="ReplyRound" />
                  </template>
                  Reply
                </n-button>
                <n-button quaternary size="small" disabled>
                  <template #icon>
                    <n-icon :component="ShareRound" />
                  </template>
                  Share
                </n-button>
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
              </div>
            </n-card>
          </n-thing>
        </n-list-item>
      </n-list>
      <n-list bordered v-else>
        <n-list-item class="my-4">
          <n-empty description="There's no content for you." />
        </n-list-item>
      </n-list>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
import { http } from "@/utils/http"
import { ReplyRound, ShareRound, ThumbUpRound, ThumbDownRound } from "@vicons/material"
import { useMessage } from "naive-ui"
import { computed, onMounted, ref } from "vue"
import { usePrincipal } from "@/stores/principal"
import { useI18n } from "vue-i18n"

const emits = defineEmits(["reply"])

defineExpose({ fetch })

const { t } = useI18n()

const $message = useMessage()
const $principal = usePrincipal()

const reverting = ref(true)

const rawData = ref<any>({ posts: [], related_authors: {} })

const data = computed(() => rawData.value)

async function fetch() {
  try {
    reverting.value = true
    rawData.value = (await http.get("/api/posts")).data
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    reverting.value = false
  }
}

onMounted(() => {
  fetch()
})
</script>

<style>
.post-image img {
  max-width: 100%;
  max-height: 100%;
  display: block;
}

.post-reply-tips .n-alert-body {
  padding: 0 12px;
}
</style>