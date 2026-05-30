<template>
  <div class="chat-panel-wrapper">
    <transition name="slide">
      <div v-if="visible" class="chat-panel">
        <div class="chat-header">
          <span class="chat-header-title">
            <RobotOutlined /> {{ $t('ai.chat') }}
          </span>
          <div class="chat-header-actions">
            <a-tooltip :title="$t('ai.newChat')">
              <a-button type="text" size="small" @click="startNewChat">
                <PlusOutlined />
              </a-button>
            </a-tooltip>
            <a-button type="text" size="small" @click="visible = false">
              <CloseOutlined />
            </a-button>
          </div>
        </div>

        <!-- No Config Warning -->
        <div v-if="!hasAIConfig" class="chat-no-config">
          <WarningOutlined class="no-config-icon" />
          <p>{{ $t('ai.noConfig') }}</p>
          <a-button type="primary" size="small" @click="goToConfig">
            {{ $t('ai.goToConfig') }}
          </a-button>
        </div>

        <!-- Article Context Badge -->
        <div v-if="articleContext && hasAIConfig" class="chat-context">
          <FileTextOutlined /> {{ $t('ai.contextArticle') }}: {{ articleContext.title }}
        </div>

        <!-- Messages -->
        <div v-if="hasAIConfig" class="chat-messages" ref="messagesRef">
          <div v-for="(msg, index) in messages" :key="index" :class="['chat-message', msg.role]">
            <div class="message-avatar">
              <a-avatar v-if="msg.role === 'user'" size="small">{{ userStore.username?.charAt(0) || 'U' }}</a-avatar>
              <a-avatar v-else size="small" :style="{ backgroundColor: '#1890ff' }">
                <template #icon><RobotOutlined /></template>
              </a-avatar>
            </div>
            <div class="message-content">
              <MarkdownRenderer v-if="msg.role === 'assistant'" :content="msg.content" :minimal="true" />
              <div v-else class="user-message-text">{{ msg.content }}</div>
            </div>
          </div>
          <div v-if="streaming" class="chat-message assistant">
            <div class="message-avatar">
              <a-avatar size="small" :style="{ backgroundColor: '#1890ff' }">
                <template #icon><RobotOutlined /></template>
              </a-avatar>
            </div>
            <div class="message-content">
              <a-spin size="small" />
            </div>
          </div>
        </div>

        <!-- Input -->
        <div v-if="hasAIConfig" class="chat-input">
          <a-textarea
            v-model:value="inputText"
            :placeholder="$t('ai.inputPlaceholder')"
            :auto-size="{ minRows: 1, maxRows: 4 }"
            @pressEnter="handleSend"
          />
          <a-button type="primary" :disabled="!inputText.trim() || streaming" @click="handleSend">
            <SendOutlined />
          </a-button>
        </div>
      </div>
    </transition>
    <a-tooltip v-if="!visible" :title="$t('ai.chat')">
      <a-button
        type="primary"
        shape="circle"
        size="large"
        class="chat-fab"
        @click="visible = true"
      >
        <MessageOutlined />
      </a-button>
    </a-tooltip>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  CloseOutlined, SendOutlined, MessageOutlined, RobotOutlined,
  PlusOutlined, WarningOutlined, FileTextOutlined,
} from '@ant-design/icons-vue'
import { streamChat } from '@/api/ai'
import { getAIConfig } from '@/api/user'
import { useUserStore } from '@/stores/user'
import MarkdownRenderer from './MarkdownRenderer.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const userStore = useUserStore()

const visible = ref(false)
const inputText = ref('')
const streaming = ref(false)
const messages = ref<{ role: string; content: string }[]>([])
const messagesRef = ref<HTMLElement>()
const aiConfigured = ref(false)

const hasAIConfig = computed(() => {
  return aiConfigured.value
})

const articleContext = computed(() => {
  if (route.name === 'ArticleDetail') {
    return { id: Number(route.params.id), title: document.querySelector('.article-title')?.textContent || '' }
  }
  return null
})

function scrollToBottom() {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  })
}

function startNewChat() {
  messages.value = []
}

function goToConfig() {
  router.push('/profile')
  visible.value = false
}

function handleSend(e?: any) {
  if (e?.shiftKey) return
  e?.preventDefault?.()
  const text = inputText.value.trim()
  if (!text || streaming.value) return

  messages.value.push({ role: 'user', content: text })
  inputText.value = ''
  streaming.value = true
  scrollToBottom()

  const assistantMsg = { role: 'assistant', content: '' }
  messages.value.push(assistantMsg)

  const data: { message: string; article_id?: number } = { message: text }
  if (articleContext.value?.id) {
    data.article_id = articleContext.value.id
  }

  streamChat(
    data,
    (content: string) => {
      assistantMsg.content += content
      scrollToBottom()
    },
    () => {
      streaming.value = false
      scrollToBottom()
    },
    (error: string) => {
      assistantMsg.content = `Error: ${error}`
      streaming.value = false
      scrollToBottom()
    }
  )
}

async function checkAIConfig() {
  if (!userStore.isLoggedIn()) return
  try {
    const res = await getAIConfig()
    const data = res.data
    aiConfigured.value = !!(data?.api_token && data?.model_name)
  } catch {
    aiConfigured.value = false
  }
}

onMounted(() => {
  checkAIConfig()
})
</script>

<style scoped lang="less">
.chat-panel-wrapper {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 1000;
}

.chat-fab {
  width: 56px;
  height: 56px;
  font-size: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: transform 0.2s;
  &:hover {
    transform: scale(1.1);
  }
}

.chat-panel {
  width: 400px;
  height: 560px;
  background: var(--bg-color);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--border-color);
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  background: var(--card-bg);
}

.chat-header-title {
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}

.chat-header-actions {
  display: flex;
  gap: 4px;
}

.chat-no-config {
  padding: 24px;
  text-align: center;
  color: var(--text-secondary);

  .no-config-icon {
    font-size: 32px;
    color: var(--warning-color, #faad14);
    margin-bottom: 8px;
  }

  p {
    margin-bottom: 12px;
  }
}

.chat-context {
  padding: 8px 16px;
  background: var(--blockquote-bg);
  font-size: 12px;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chat-message {
  display: flex;
  gap: 8px;

  &.user {
    flex-direction: row-reverse;

    .message-content {
      background: #1890ff;
      color: #fff;
      border-radius: 12px 12px 0 12px;
    }
  }

  &.assistant {
    .message-content {
      background: var(--card-bg);
      border: 1px solid var(--border-color);
      border-radius: 12px 12px 12px 0;
    }
  }
}

.message-content {
  max-width: 80%;
  padding: 8px 12px;
  font-size: 14px;
  line-height: 1.6;
  word-break: break-word;
}

.user-message-text {
  white-space: pre-wrap;
}

.chat-input {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
  background: var(--card-bg);
}

.slide-enter-active, .slide-leave-active {
  transition: all 0.3s ease;
}
.slide-enter-from, .slide-leave-to {
  transform: translateY(20px);
  opacity: 0;
}

@media (max-width: 480px) {
  .chat-panel {
    width: calc(100vw - 32px);
    height: 60vh;
    bottom: 16px;
    right: 16px;
  }
}
</style>
