<template>
  <div>
    <n-spin :show="reverting">
      <n-card v-if="layer === 'image'" content-style="padding: 0" class="post-media-container">
        <n-image
          object-fit="cover"
          class="post-media block"
          :src="`/srv/subapps/quaso${props.src}`"
        />
      </n-card>

      <n-card v-else-if="layer === 'video'" content-style="padding: 0" class="post-media-container">
        <video controls class="post-media block">
          <source :src="`/srv/subapps/quaso${props.src}`" :type="meta.mimetype">
        </video>
      </n-card>

      <n-card v-else-if="layer === 'audio'" content-style="padding: 0" class="post-media-container">
        <audio controls class="post-media block">
          <source :src="`/srv/subapps/quaso${props.src}`" :type="meta.mimetype">
        </audio>
      </n-card>

      <n-card v-else size="small">
        <a :href="`/srv/subapps/quaso${props.src}`" target="_blank">
          Download attachment {{ meta?.record?.name }}
        </a>
      </n-card>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
import { http } from "@/utils/http"
import { onMounted, ref } from "vue"
import { useMessage } from "naive-ui"
import { useI18n } from "vue-i18n"

const { t } = useI18n()

const props = defineProps<{ src: string }>()

const $message = useMessage()

const reverting = ref(false)

const layer = ref<"image" | "audio" | "video" | "download">("download")
const meta = ref<any>({})

async function fetch() {
  try {
    reverting.value = true

    const res = await http.get(`${props.src}/meta`)

    meta.value = res.data
    layer.value = res.data.mimetype.split("/")[0]
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e]))
  } finally {
    reverting.value = false
  }
}

onMounted(() => {
  fetch()
})
</script>

<style>
.post-media-container {
  max-width: 800px;

  display: flex;
  place-items: center;
}

.post-media, .post-media img {
  max-width: 100%;
  max-height: 100%;
  display: block;
}
</style>