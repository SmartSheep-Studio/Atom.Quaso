import { defineStore } from "pinia"
import { ref } from "vue"
import { http } from "@/utils/http"
import { useMessage } from "naive-ui"
import { useI18n } from "vue-i18n"

export const usePosts = defineStore("posts", () => {
  const isReverting = ref(true)

  const data = ref<any>({ posts: [], related_authors: {} })

  const { t } = useI18n()
  const $message = useMessage()

  async function fetch() {
    try {
      isReverting.value = true
      data.value = (await http.get("/api/posts")).data
    } catch (e: any) {
      $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
    } finally {
      isReverting.value = false
    }
  }

  return { isReverting, data, fetch }
})