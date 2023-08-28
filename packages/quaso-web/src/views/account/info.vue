<template>
  <div>
    <n-spin :show="reverting">
      <n-card class="rounded-none" style="border-bottom: 0">
        <n-page-header :title="account?.user?.nickname" @back="$router.back()">
          <template #avatar>
            <n-avatar color="transparent" :src="`/srv/subapps/quaso${account?.user?.avatar_url}`" />
          </template>

          <template #header>
            <n-breadcrumb>
              <n-breadcrumb-item @click="$router.push({ name: 'plaza' })">Plaza</n-breadcrumb-item>
              <n-breadcrumb-item>@{{ account?.user?.name }}</n-breadcrumb-item>
            </n-breadcrumb>
          </template>

          <template #subtitle>
            {{ account?.user?.description }}
          </template>

          <template #extra>
            <n-space>
              <div>
                <n-button secondary size="small" type="tertiary" :loading="submitting" @click="subscribe"
                  v-if="account?.is_subscribed">
                  <template #icon>
                    <n-icon :component="CheckRound" />
                  </template>

                  Subscribed
                </n-button>
                <n-button secondary size="small" type="primary" :loading="submitting" @click="subscribe" v-else>
                  <template #icon>
                    <n-icon :component="PlusRound" />
                  </template>

                  Subscribe
                </n-button>
              </div>
            </n-space>
          </template>
        </n-page-header>

        <n-divider class="inset-divider" />

        <n-grid :cols="4" item-responsive responsive="screen" :x-gap="8" :y-gap="8" class="text-center">
          <n-gi span="2 m:1">
            <div class="font-bold">Subscribers</div>
            <div class="font-mono text-xl">{{ account?.subscribers }}</div>
          </n-gi>
          <n-gi span="2 m:1">
            <div class="font-bold">Subscriptions</div>
            <div class="font-mono text-xl">{{ account?.subscriptions }}</div>
          </n-gi>
          <n-gi span="2 m:1">
            <div class="font-bold">Posts</div>
            <div class="font-mono text-xl">{{ account?.posts }}</div>
          </n-gi>
          <n-gi span="2 m:1">
            <div class="font-bold">Likes</div>
            <div class="font-mono text-xl">{{ account?.likes }}</div>
          </n-gi>
        </n-grid>
      </n-card>

      <div>
        <n-card size="small" class="rounded-none" style="border-bottom: 0">
          <div class="flex">
            <n-select v-model:value="postFilterOptions.type" :options="[
              { label: 'All', value: 'all' },
              { label: 'Text', value: 'text' },
              { label: 'Image', value: 'image' },
              { label: 'Audio', value: 'audio' },
              { label: 'Video', value: 'video' },
            ]" />
          </div>
        </n-card>

        <div>
          <n-list v-if="posts?.posts.length > 0" bordered hoverable class="rounded-none" style="--n-border-radius: 0">
            <n-list-item v-for="item in posts?.posts"
              @click="$router.push({ name: 'plaza.focus', params: { post: item.id } })">
              <post-player :post="item" @reply="replyPost(item)" @share="sharePost(item)" />
            </n-list-item>
          </n-list>
          <n-list bordered class="rounded-none" v-else>
            <n-list-item class="my-4">
              <n-empty description="There's no content for you." />
            </n-list-item>
          </n-list>

          <n-card size="small" class="rounded-none" style="border-top: 0">
            <div class="flex justify-center">
              <n-pagination v-model:page="postFilterOptions.page" :page-count="postPostPage" />
            </div>
          </n-card>
        </div>
      </div>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
import { CheckRound, PlusRound } from "@vicons/material"
import { computed, onMounted, reactive, ref, watch } from "vue"
import { http } from "@/utils/http"
import { useRoute } from "vue-router"
import { useI18n } from "vue-i18n"
import { useMessage } from "naive-ui"
import { usePosts } from "@/stores/posts"
import PostPlayer from "@/components/player/post-player.vue"

const emits = defineEmits(["post"])

const { t } = useI18n()

const $posts = usePosts()
const $route = useRoute()
const $message = useMessage()

const reverting = ref(true)

const account = ref<any>({})
const posts = ref<any>({ posts: [], total: 0 })

const postPostPage = computed(() => {
  return Math.ceil(posts.value.total / 5)
})
const postFilterOptions = reactive({
  type: "all",
  page: 1
})

async function fetch() {
  try {
    reverting.value = true
    account.value = (await http.get(`/api/accounts/${$route.params.account}`)).data
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    reverting.value = false
  }
}

async function fetchPosts() {
  try {
    reverting.value = true

    const res = await http.get(`/api/accounts/${$route.params.account}/posts`, {
      params: {
        type: postFilterOptions.type === "all" ? undefined : postFilterOptions.type,
        skip: (postFilterOptions.page - 1) * 5
      }
    })

    posts.value = res.data
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    reverting.value = false
  }
}

const submitting = ref(false)

async function subscribe() {
  try {
    submitting.value = true

    const res = await http.post(`/api/accounts/${$route.params.account}/subscribe`)
    await fetch()

    $message.success(res.status === 200 ? "Successfully subscribed" : "Successfully cancelled subscribe")
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    submitting.value = false
  }
}

watch(postFilterOptions, () => {
  fetchPosts()
})

watch($posts, (value) => {
  if (value.isReverting) {
    fetchPosts()
  }
}, { deep: true })

onMounted(() => {
  fetch()
  fetchPosts()
})

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
</script>

<style>
.inset-divider {
  margin: 18px -24px !important;
  width: calc(100% + 48px);
}
</style>
