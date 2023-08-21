<template>
  <div>
    <n-spin :show="reverting">
      <n-card class="rounded-t-none">
        <n-page-header
          :title="post?.account?.nickname"
          @back="$router.back()"
        >
          <template #header>
            <n-breadcrumb>
              <n-breadcrumb-item @click="$router.push({ name: 'plaza' })">Plaza</n-breadcrumb-item>
              <n-breadcrumb-item>Post #{{ post?.id }}</n-breadcrumb-item>
            </n-breadcrumb>
          </template>

          <template #subtitle>
            <n-space size="small">
              <n-tag
                v-if="post?.account?.user_id === $principal.account.id"
                size="small"
                type="warning"
                :bordered="false"
              >
                You
              </n-tag>
              <n-tag class="capitalize" size="small" type="primary" :bordered="false">
                {{ post?.type }}
              </n-tag>
              <n-tag
                v-for="tag in post?.tags"
                size="small"
                type="info"
                :bordered="false"
              >
                {{ tag }}
              </n-tag>
            </n-space>
          </template>

          <template #extra>
            <n-dropdown :options="dropdownOptions" placement="bottom-start" @select="dropdownHandler">
              <n-button :bordered="false" style="padding: 0 4px">
                ···
              </n-button>
            </n-dropdown>
          </template>
        </n-page-header>

        <n-divider class="inset-divider" />

        <n-thing>
          <div>
            <vue-markdown :source="post?.content ?? ''" />
          </div>

          <n-space vertical class="mt-4">
            <n-card content-style="padding: 0" class="max-w-[800px]" v-for="img in post?.attachments">
              <attachment-player :src="img" />
            </n-card>
          </n-space>

          <div class="mt-2" v-if="post?.belong_id != null">
            <n-alert class="post-reply-tips" :show-icon="false">
              This post is replying
              <router-link :to="{ name: 'plaza.focus', params: { post: post?.belong_id } }" @click.stop>
                #{{ post?.belong_id }}
              </router-link>
            </n-alert>
          </div>

          <n-card size="small" class="mb-1 mt-2" content-style="padding: 8px" embedded>
            <div class="flex justify-around">
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-button quaternary size="small" @click.stop="replyPost(post)">
                    <template #icon>
                      <n-icon :component="ReplyRound" />
                    </template>
                    {{ post?.comment_count }}
                  </n-button>
                </template>
                Reply
              </n-tooltip>
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-button
                    quaternary
                    size="small"
                    :type="post?.is_liked ? 'primary' : 'empty'"
                    @click.stop="like(post)"
                  >
                    <template #icon>
                      <n-icon :component="ThumbUpRound" />
                    </template>
                    {{ post?.like_count }}
                  </n-button>
                </template>
                Like
              </n-tooltip>
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-button
                    quaternary
                    size="small"
                    :type="post?.is_disliked ? 'primary' : 'empty'"
                    @click.stop="dislike(post)"
                  >
                    <template #icon>
                      <n-icon :component="ThumbDownRound" />
                    </template>
                    {{ post?.dislike_count }}
                  </n-button>
                </template>
                Dislike
              </n-tooltip>
              <n-button quaternary size="small" @click.stop="sharePost(post)">
                <template #icon>
                  <n-icon :component="ShareRound" />
                </template>
                Share
              </n-button>
            </div>
          </n-card>
        </n-thing>
      </n-card>

      <div class="py-4 px-4">
        <div class="text-lg">Comments ({{ post?.comments?.length }})</div>
      </div>

      <n-list v-if="post?.comments?.length > 0" bordered hoverable>
        <n-list-item
          v-for="item in post?.comments"
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
              <attachment-player v-for="img in item.attachments" :src="img" />
            </n-space>

            <n-card size="small" class="mb-1 mt-2" content-style="padding: 8px" embedded>
              <div class="flex justify-around">
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-button quaternary size="small" @click.stop="replyPost(item)">
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
                <n-button quaternary size="small" @click.stop="sharePost(item)">
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
import { computed, onMounted, ref, watch } from "vue"
import { http } from "@/utils/http"
import { useRoute, useRouter } from "vue-router"
import { useDialog, useMessage } from "naive-ui"
import { useI18n } from "vue-i18n"
import { usePosts } from "@/stores/posts"
import { usePrincipal } from "@/stores/principal"
import { ReplyRound, ShareRound, ThumbDownRound, ThumbUpRound } from "@vicons/material"
import VueMarkdown from "vue-markdown-render"
import AttachmentPlayer from "@/components/player/attachment-player.vue"

const { t } = useI18n()

const emits = defineEmits(["post"])

const $posts = usePosts()
const $route = useRoute()
const $router = useRouter()
const $dialog = useDialog()
const $message = useMessage()
const $principal = usePrincipal()

const reverting = ref(false)

const post = ref<any>({})

const dropdownOptions = computed(() => {
  const items: any[] = [
    { label: "Report", key: "report", disabled: true }
  ]

  if (post.value?.account?.user_id === $principal.account.id) {
    items.push(
      { label: "Edit", key: "edit" },
      { label: "Delete", key: "delete", props: { style: { color: "#de576d" } } }
    )
  }

  return items
})

function dropdownHandler(key: string) {
  switch (key) {
    case "edit":
      emits("post", {
        ...post.value,
        edit_to: post.value
      })
      break
    case "delete":
      $dialog.warning({
        title: "Are you confirm?",
        content: "Are you sure you want to delete this post? This operation cannot be undo and will processed immediately.",
        positiveText: "Confirm",
        negativeText: "Cancel",
        onPositiveClick: async () => {
          await http.delete(`/api/posts/${$route.params.post}`)
          $message.success("Successfully deleted this post.")
          await $router.push({ name: "plaza" })
        }
      })
      break
  }
}

async function fetch() {
  try {
    reverting.value = true
    post.value = (await http.get(`/api/posts/${$route.params.post}`)).data
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    reverting.value = false
  }
}

async function like(item: any) {
  try {
    $posts.isReverting = true
    const res = await http.post(`/api/posts/${item.id}/like`)
    await fetch()
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
    await fetch()
    $message.success(res.status === 200 ? "Successfully disliked" : "Successfully cancelled dislike")
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    $posts.isReverting = false
  }
}

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

onMounted(() => {
  fetch()
})

watch($posts, (value) => {
  if (value.isReverting) {
    fetch()
  }
}, { deep: true })
</script>

<style>
.post-reply-tips .n-alert-body {
  padding: 0 12px;
}

.post-reply-tips .n-alert__close {
  margin: 6.5px 12px 0 0;
}

.inset-divider {
  margin: 18px -24px !important;
  width: calc(100% + 48px);
}
</style>