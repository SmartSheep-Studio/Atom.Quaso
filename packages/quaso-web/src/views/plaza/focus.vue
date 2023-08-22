<template>
  <div>
    <n-spin :show="reverting">
      <n-card class="rounded-t-none">
        <n-page-header
          :title="`Post #${post?.id}`"
          @back="$router.back()"
        >
          <template #header>
            <n-breadcrumb>
              <n-breadcrumb-item @click="$router.push({ name: 'plaza' })">Plaza</n-breadcrumb-item>
              <n-breadcrumb-item>Post #{{ post?.id }}</n-breadcrumb-item>
            </n-breadcrumb>
          </template>

          <template #subtitle>
            <div class="flex gap-1">
              <div>{{ new Date(post?.published_at).toLocaleString() }}</div>
              <div v-if="post?.is_edited">(Edited)</div>
            </div>
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

        <post-player
          v-if="!reverting"
          :post="post"
          @reply="replyPost(post)"
          @share="sharePost(post)"
        />
      </n-card>

      <div class="py-4 px-4">
        <div class="text-lg">Comments ({{ post?.comments?.length }})</div>
      </div>

      <n-list v-if="post?.comments?.length > 0" bordered hoverable>
        <n-list-item
          v-for="item in post?.comments"
          @click="$router.push({ name: 'plaza.focus', params: { post: item.id } })"
        >
          <post-player
            :post="item"
            @reply="replyPost(item)"
            @share="sharePost(item)"
          />
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
import PostPlayer from "@/components/player/post-player.vue"

const { t } = useI18n()

const emits = defineEmits(["post"])

const $posts = usePosts()
const $route = useRoute()
const $router = useRouter()
const $dialog = useDialog()
const $message = useMessage()
const $principal = usePrincipal()

const reverting = ref(true)

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