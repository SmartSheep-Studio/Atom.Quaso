<template>
  <div class="w-[720px]">
    <n-card :title="payload.edit_to != null ? 'Edit Post' : 'New Post'" size="large">
      <n-form @submit.prevent="submit">
        <div>
          <n-input
            v-model:value="payload.content"
            type="textarea"
            placeholder="What's happened?"
          />
        </div>

        <n-alert
          v-if="payload.edit_to != null"
          class="post-reply-tips mt-2"
          :show-icon="false"
          closable
          @close="reset()"
        >
          You're editing post #{{ payload.edit_to?.id }}
        </n-alert>

        <n-alert
          v-if="payload.belong_to != null"
          class="post-reply-tips mt-2"
          :show-icon="false"
          closable
          @close="payload.belong_to = null"
        >
          You're replying post #{{ payload.belong_to?.id }}
        </n-alert>

        <div class="mt-2 flex gap-2">
          <n-select :options="options" v-model:value="payload.type" class="w-[160px]" />

          <n-dynamic-tags v-model:value="payload.tags" class="post-tags" />
        </div>

        <div class="flex justify-between mt-5">
          <n-space>
            <n-tooltip trigger="hover" placement="bottom">
              <template #trigger>
                <n-button circle class="rounded-[3px]" @click="popups.attachments = true">
                  <template #icon>
                    <n-icon :component="AttachFileRound" />
                  </template>
                </n-button>
              </template>
              Attachments
            </n-tooltip>
            <n-tooltip trigger="hover" placement="bottom">
              <template #trigger>
                <n-button circle class="rounded-[3px]" @click="popups.schedule = true">
                  <template #icon>
                    <n-icon :component="ScheduleRound" />
                  </template>
                </n-button>
              </template>
              Schedule
            </n-tooltip>
          </n-space>

          <n-button type="primary" :loading="submitting" attr-type="submit">Publish</n-button>
        </div>
      </n-form>
    </n-card>

    <n-modal v-model:show="popups.attachments" display-directive="show">
      <n-card class="w-[680px]" title="Attachments" size="huge" :bordered="false">
        <n-upload
          multiple
          directory-dnd
          action="/api/attachments"
          :max="5"
          :custom-request="attach"
          :show-remove-button="false"
        >
          <n-upload-dragger>
            <div style="margin-bottom: 12px">
              <n-icon size="48" :depth="3">
                <archive-round />
              </n-icon>
            </div>
            <n-text style="font-size: 16px">
              Click or drag files to this area to upload
            </n-text>
            <n-p depth="3" style="margin: 8px 0 0 0">
              Please do not upload sensitive data, such as your bank card number and password,
              credit card number expiration date and security code
            </n-p>
          </n-upload-dragger>
        </n-upload>
      </n-card>
    </n-modal>

    <n-modal v-model:show="popups.schedule" display-directive="show">
      <n-card class="w-[680px]" title="Schedule" size="huge" :bordered="false">
        <n-date-picker type="datetime" clearable v-model:value="payload.published_at" />
      </n-card>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { ArchiveRound, AttachFileRound, ScheduleRound } from "@vicons/material"
import { type UploadCustomRequestOptions, useMessage } from "naive-ui"
import { reactive, ref } from "vue"
import { useI18n } from "vue-i18n"
import { http } from "@/utils/http"

const emits = defineEmits(["submit"])

defineExpose({ trigger })

const { t } = useI18n()

const $message = useMessage()

const submitting = ref(false)

const popups = reactive({ attachments: false, schedule: false })

const options = [
  { label: "Text", value: "text" },
  { label: "Image", value: "image" },
  { label: "Audio", value: "audio" },
  { label: "Video", value: "video" }
]

const payload = reactive<any>({
  type: "text",
  content: "",
  tags: [],
  attachments: [],
  belong_to: null,
  edit_to: null,
  published_at: null
})

async function submit() {
  try {
    submitting.value = true

    const pd = {
      type: payload.type,
      content: payload.content,
      tags: payload.tags,
      attachments: payload.attachments,
      belong_to: payload.belong_to?.id ?? undefined,
      published_at: new Date(payload.published_at)
    }

    if (payload.edit_to == null) {
      await http.post("/api/posts", pd)
    } else {
      await http.put(`/api/posts/${payload.edit_to.id}`, pd)
    }

    reset()
    emits("submit")
    $message.success("Successfully published.")
  } catch (e: any) {
    $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
  } finally {
    submitting.value = false
  }
}

function attach({ file, data, headers, action, onFinish, onError, onProgress }: UploadCustomRequestOptions) {
  const formData = new FormData()
  if (data) {
    Object.keys(data).forEach((key) => {
      formData.append(
        key,
        data[key as keyof UploadCustomRequestOptions["data"]]
      )
    })
  }
  formData.append("file", file.file as File)

  http.post(action as string, formData, {
    headers: headers as Record<string, string>,
    onUploadProgress: (data: any) => {
      onProgress({ percent: Math.ceil(data.percent) })
    }
  })
    .then((res) => {
      payload.attachments.push(res.data.url)

      if (file.type?.startsWith("image")) {
        payload.type = "image"
      } else if (file.type?.startsWith("audio")) {
        payload.type = "audio"
      } else if (file.type?.startsWith("video")) {
        payload.type = "video"
      }

      $message.success(`Upload file "${file.name}" successfully.`)
      onFinish()
    })
    .catch((e) => {
      $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
      onError()
    })
}

function trigger(overrides: any) {
  payload.belong_to = overrides.belong_to
  payload.type = overrides.type ?? payload.type
  payload.content = overrides.content
  payload.tags = overrides.tags ?? payload.tags
  payload.edit_to = overrides.edit_to ?? null
  payload.published_at = overrides.published_at
  payload.attachments = overrides.attachments ?? payload.attachments
}

function reset() {
  payload.type = "text"
  payload.content = ""
  payload.tags = []
  payload.attachments = []
  payload.belong_to = null
  payload.edit_to = null
  payload.published_at = null
}
</script>

<style>
.post-tags, .post-tags .n-button, .post-tags .n-tag {
  height: 34px;
}

.post-tags .n-tag {
  gap: 4px;
  padding: 0 12px;
}

.post-reply-tips {
  height: 34px;
  display: flex;
  align-items: center;
}

.post-reply-tips .n-alert-body {
  padding: 0 12px;
}

.post-reply-tips .n-alert__close {
  margin: 6.5px 12px 0 0;
}
</style>