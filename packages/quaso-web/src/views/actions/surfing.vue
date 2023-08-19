<template>
  <div>
    <n-spin :show="reverting">
      <n-list bordered v-if="data.posts.length > 0" hoverable clickable>
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
                  class="post-image"
                  :src="`/srv/subapps/quaso${img}`"
                />
              </n-card>
            </n-space>
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
import { useMessage } from "naive-ui"
import { computed, onMounted, ref } from "vue"
import { usePrincipal } from "@/stores/principal"
import { useI18n } from "vue-i18n"

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
}
</style>