<template>
  <n-thing>
    <template #avatar>
      <n-avatar color="transparent" :src="`/srv/subapps/quaso${props.post.author.avatar_url}`" />
    </template>

    <template #header>
      <div>
        <div>{{ props.post.author.nickname }}</div>
        <div class="text-xs text-gray-600 mb-1">{{ props.post.author.description }}</div>
      </div>
    </template>

    <n-space size="small" class="mt-[-12px]">
      <n-tag
        v-if="props.post.author.id === $principal.account.id"
        type="warning"
        size="small"
        :bordered="false"
      >
        You
      </n-tag>
      <n-tag class="capitalize" size="small" type="primary" :bordered="false">
        {{ props.post.type }}
      </n-tag>
      <n-tag
        v-for="tag in props.post.tags"
        size="small"
        type="info"
        :bordered="false"
      >
        {{ tag }}
      </n-tag>
    </n-space>

    <div>
      <vue-markdown :source="props.post.content" />
    </div>

    <n-space vertical class="mt-4">
      <attachment-player v-for="img in props.post.attachments" :src="img" :key="img" @click.stop />
    </n-space>

    <post-widget
      class="mb-1 mt-2"
      :post="props.post"
      @reply="emits('reply')"
      @like="like"
      @dislike="dislike"
      @share="emits('share')"
      @click.stop
    />
  </n-thing>
</template>
<script setup lang="ts">
import AttachmentPlayer from "@/components/player/attachment-player.vue"
import PostWidget from "@/components/player/post-widget.vue"
import VueMarkdown from "vue-markdown-render"
import { useMessage } from "naive-ui"
import { usePrincipal } from "@/stores/principal"
import { usePosts } from "@/stores/posts"
import { http } from "@/utils/http"
import { useI18n } from "vue-i18n"

const { t } = useI18n()

const $posts = usePosts()
const $principal = usePrincipal()
const $message = useMessage()

const props = defineProps<{ post: any }>()
const emits = defineEmits(["reply", "like", "dislike", "share"])

async function like() {
  try {
    $posts.isReverting = true
    const res = await http.post(`/api/posts/${props.post.id}/like`)
    await $posts.fetch()
    $message.success(res.status === 200 ? "Successfully liked" : "Successfully cancelled like")
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    $posts.isReverting = false
    emits("like")
  }
}

async function dislike() {
  try {
    $posts.isReverting = true
    const res = await http.post(`/api/posts/${props.post.id}/dislike`)
    await $posts.fetch()
    $message.success(res.status === 200 ? "Successfully disliked" : "Successfully cancelled dislike")
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    $posts.isReverting = false
    emits("dislike")
  }
}
</script>