<template>
  <n-card size="small" content-style="padding: 8px" embedded>
    <div class="flex justify-around">
      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button quaternary size="small" @click.stop="emits('reply')">
            <template #icon>
              <n-icon :component="ReplyRound" />
            </template>
            {{ post.comment_count }}
          </n-button>
        </template>
        Reply
      </n-tooltip>
      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button
            quaternary
            size="small"
            :type="props.post.is_liked ? 'primary' : 'empty'"
            @click.stop="emits('like')"
          >
            <template #icon>
              <n-icon :component="ThumbUpRound" />
            </template>
            {{ props.post.like_count }}
          </n-button>
        </template>
        Like
      </n-tooltip>
      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button
            quaternary
            size="small"
            :type="props.post.is_disliked ? 'primary' : 'empty'"
            @click="emits('dislike')"
          >
            <template #icon>
              <n-icon :component="ThumbDownRound" />
            </template>
            {{ props.post.dislike_count }}
          </n-button>
        </template>
        Dislike
      </n-tooltip>
      <n-button quaternary size="small" @click="emits('share')">
        <template #icon>
          <n-icon :component="ShareRound" />
        </template>
        Share
      </n-button>
    </div>
  </n-card>
</template>

<script setup lang="ts">
import { ReplyRound, ShareRound, ThumbDownRound, ThumbUpRound } from "@vicons/material"

const props = defineProps<{ post: any }>()
const emits = defineEmits(["reply", "like", "dislike", "share"])
</script>