import { defineStore } from "pinia"
import { computed, reactive, ref, watch } from "vue"
import { http } from "@/utils/http"
import { useMessage } from "naive-ui"
import { useI18n } from "vue-i18n"

export const usePosts = defineStore("posts", () => {
  const isReverting = ref(true)

  const filterOptions = reactive({
    type: 'all',
    page: 1
  })

  const data = ref<any>({ total: 0, posts: [] })

  const totalPage = computed(() => {
    return Math.ceil(data.value.total / 5)
  })

  const { t } = useI18n()
  const $message = useMessage()

  async function fetch() {
    try {
      isReverting.value = true

      const res = await http.get("/api/posts", {
        params: {
          type: filterOptions.type === 'all' ? undefined : filterOptions.type,
          skip: (filterOptions.page - 1) * 5
        }
      })

      data.value = res.data
    } catch (e: any) {
      $message.error(t("common.feedback.unknown-error", [e.response.body ?? e.message]))
    } finally {
      isReverting.value = false
    }
  }

  watch(filterOptions, () => {
    fetch()
  })

  return { isReverting, filterOptions, totalPage, data, fetch }
})