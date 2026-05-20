<template>
  <div class="article-detail">
    <a-spin :spinning="loading">
      <template v-if="article">
        <article class="article-content-wrapper">
          <h1 class="article-title">{{ article.title }}</h1>
          <div class="article-meta">
            <span><ClockCircleOutlined /> {{ formatDate(article.published_at) }}</span>
            <span><EyeOutlined /> {{ article.view_count || 0 }}</span>
            <span><LikeOutlined /> {{ article.like_count || 0 }}</span>
            <span><MessageOutlined /> {{ article.comment_count || 0 }}</span>
          </div>
          <div class="article-tags mb-16">
            <a-tag v-if="article.category" color="blue">{{ article.category.name || article.category }}</a-tag>
            <a-tag v-for="tag in (article.tags || [])" :key="tag.id">{{ tag.name }}</a-tag>
          </div>

          <!-- AI Summary -->
          <div class="ai-summary mb-16">
            <div class="ai-summary-header">
              <span><RobotOutlined /> {{ $t('article.aiSummary') }}</span>
              <a-button type="link" size="small" @click="refreshSummary" :loading="summaryLoading">
                <ReloadOutlined /> {{ $t('article.refreshSummary') }}
              </a-button>
            </div>
            <p v-if="article.ai_summary">{{ article.ai_summary }}</p>
            <p v-else class="ai-summary-empty">{{ $t('common.noData') }}</p>
          </div>

          <!-- Article Content + TOC Layout -->
          <div class="article-body">
            <div class="article-main">
              <MarkdownRenderer :content="article.content" @toc-generated="handleTocGenerated" />
            </div>
            <a-affix :offset-top="80" class="toc-affix" v-if="tocItems.length">
              <div class="toc-wrapper">
                <h4 class="toc-title">{{ $t('article.tableOfContents') }}</h4>
                <ul class="toc-list">
                  <li
                    v-for="item in tocItems"
                    :key="item.id"
                    :class="['toc-item', `toc-level-${item.level}`, { active: activeTocId === item.id }]"
                  >
                    <a @click.prevent="scrollToHeading(item.id)">{{ item.text }}</a>
                  </li>
                </ul>
              </div>
            </a-affix>
          </div>

          <!-- Like Button -->
          <div class="article-actions mt-24">
            <a-button
              :type="article.liked ? 'primary' : 'default'"
              :ghost="!article.liked"
              @click="handleLike"
              size="large"
            >
              <LikeFilled v-if="article.liked" />
              <LikeOutlined v-else />
              {{ article.liked ? $t('article.liked') : $t('article.like') }} {{ article.like_count || 0 }}
            </a-button>
          </div>
        </article>

        <a-divider />

        <!-- Comment Section -->
        <div class="comment-section">
          <h3>{{ $t('comment.content') }} ({{ comments.length }})</h3>

          <!-- Comment Input -->
          <div class="comment-form mb-16">
            <div v-if="!userStore.isLoggedIn()" class="comment-nickname mb-8">
              <a-input
                v-model:value="commentNickname"
                :placeholder="$t('comment.nicknamePlaceholder')"
                size="small"
                style="max-width: 200px"
              />
            </div>
            <a-textarea
              v-model:value="commentText"
              :placeholder="$t('comment.placeholder')"
              :rows="3"
            />
            <div class="comment-form-actions mt-8">
              <span v-if="replyTo" class="reply-hint">
                {{ $t('comment.replyTo') }} {{ replyTo.nickname || $t('comment.anonymous') }}
                <a-button type="link" size="small" @click="cancelReply">{{ $t('comment.cancelReply') }}</a-button>
              </span>
              <a-button type="primary" @click="submitCommentHandler" :disabled="!commentText.trim()" :loading="submitting">
                {{ $t('comment.submit') }}
              </a-button>
            </div>
          </div>

          <!-- Comment List -->
          <div class="comment-list">
            <div v-for="comment in topLevelComments" :key="comment.id" class="comment-item">
              <div class="comment-main">
                <a-avatar class="comment-avatar">{{ (comment.nickname || 'A').charAt(0).toUpperCase() }}</a-avatar>
                <div class="comment-body">
                  <div class="comment-header">
                    <span class="comment-author">{{ comment.nickname || $t('comment.anonymous') }}</span>
                    <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
                  </div>
                  <div class="comment-content">{{ comment.content }}</div>
                  <a-button type="link" size="small" @click="handleReply(comment)">
                    <MessageOutlined /> {{ $t('comment.reply') }}
                  </a-button>
                </div>
              </div>
              <!-- Nested Replies -->
              <div v-if="getReplies(comment.id).length" class="comment-replies">
                <div v-for="reply in getReplies(comment.id)" :key="reply.id" class="comment-item reply-item">
                  <a-avatar class="comment-avatar" size="small">{{ (reply.nickname || 'A').charAt(0).toUpperCase() }}</a-avatar>
                  <div class="comment-body">
                    <div class="comment-header">
                      <span class="comment-author">{{ reply.nickname || $t('comment.anonymous') }}</span>
                      <span v-if="reply.reply_to" class="reply-to">{{ $t('comment.replyTo') }} {{ reply.reply_to }}</span>
                      <span class="comment-time">{{ formatDate(reply.created_at) }}</span>
                    </div>
                    <div class="comment-content">{{ reply.content }}</div>
                    <a-button type="link" size="small" @click="handleReply(reply)">
                      <MessageOutlined /> {{ $t('comment.reply') }}
                    </a-button>
                  </div>
                </div>
              </div>
            </div>
            <a-empty v-if="!comments.length" :description="$t('common.noData')" />
          </div>
        </div>
      </template>
    </a-spin>
    <a-button class="back-btn" @click="$router.back()">
      <ArrowLeftOutlined /> {{ $t('common.back') }}
    </a-button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import {
  ClockCircleOutlined,
  EyeOutlined,
  LikeOutlined,
  LikeFilled,
  MessageOutlined,
  RobotOutlined,
  ReloadOutlined,
  ArrowLeftOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { getArticleDetail, likeArticle, getComments, submitComment, generateAISummary } from '@/api/article'
import { useUserStore } from '@/stores/user'
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'

const route = useRoute()
const { t } = useI18n()
const userStore = useUserStore()

const article = ref<any>(null)
const comments = ref<any[]>([])
const loading = ref(true)
const commentText = ref('')
const commentNickname = ref('')
const submitting = ref(false)
const summaryLoading = ref(false)
const replyTo = ref<any>(null)
const tocItems = ref<{ id: string; text: string; level: number }[]>([])
const activeTocId = ref('')

function formatDate(date: string) {
  return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : ''
}

const topLevelComments = computed(() => {
  return comments.value.filter(c => !c.parent_id)
})

function getReplies(parentId: number) {
  return comments.value.filter(c => c.parent_id === parentId)
}

async function fetchArticle() {
  loading.value = true
  try {
    const id = Number(route.params.id)
    const res = await getArticleDetail(id)
    article.value = res.data
  } catch {
    // handled
  } finally {
    loading.value = false
  }
}

async function fetchComments() {
  try {
    const id = Number(route.params.id)
    const res = await getComments(id)
    comments.value = res.data?.list || []
  } catch {
    // handled
  }
}

async function handleLike() {
  try {
    await likeArticle(article.value.id)
    if (article.value.liked) {
      article.value.like_count = Math.max(0, (article.value.like_count || 0) - 1)
      article.value.liked = false
    } else {
      article.value.like_count = (article.value.like_count || 0) + 1
      article.value.liked = true
    }
  } catch {
    // handled
  }
}

async function refreshSummary() {
  summaryLoading.value = true
  try {
    const id = Number(route.params.id)
    const res = await generateAISummary(id)
    article.value.ai_summary = res.data?.summary || ''
    message.success(t('common.success'))
  } catch {
    // handled
  } finally {
    summaryLoading.value = false
  }
}

function handleReply(comment: any) {
  replyTo.value = comment
  commentText.value = ''
}

function cancelReply() {
  replyTo.value = null
}

async function submitCommentHandler() {
  if (!commentText.value.trim()) return
  submitting.value = true
  try {
    const data: Record<string, any> = {
      content: commentText.value,
    }
    if (!userStore.isLoggedIn() && commentNickname.value.trim()) {
      data.nickname = commentNickname.value.trim()
    }
    if (replyTo.value) {
      data.parent_id = replyTo.value.id
      data.reply_to = replyTo.value.nickname || t('comment.anonymous')
    }
    await submitComment(article.value.id, data)
    message.success(t('comment.submitSuccess'))
    commentText.value = ''
    replyTo.value = null
    fetchComments()
  } catch {
    // handled
  } finally {
    submitting.value = false
  }
}

function handleTocGenerated(items: { id: string; text: string; level: number }[]) {
  tocItems.value = items
}

function scrollToHeading(id: string) {
  const el = document.getElementById(id)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    activeTocId.value = id
  }
}

function handleScroll() {
  for (let i = tocItems.value.length - 1; i >= 0; i--) {
    const el = document.getElementById(tocItems.value[i].id)
    if (el) {
      const rect = el.getBoundingClientRect()
      if (rect.top <= 100) {
        activeTocId.value = tocItems.value[i].id
        return
      }
    }
  }
  if (tocItems.value.length) {
    activeTocId.value = tocItems.value[0].id
  }
}

onMounted(() => {
  fetchArticle()
  fetchComments()
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped lang="less">
.article-detail {
  max-width: 100%;
}

.article-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 12px;
  line-height: 1.4;
}

.article-meta {
  display: flex;
  gap: 20px;
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 12px;
  flex-wrap: wrap;

  span {
    display: flex;
    align-items: center;
    gap: 4px;
  }
}

.article-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.ai-summary {
  padding: 16px;
  background: var(--blockquote-bg);
  border-radius: 8px;
  border-left: 4px solid #1890ff;
}

.ai-summary-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  margin-bottom: 8px;
}

.ai-summary-empty {
  color: var(--text-secondary);
  font-style: italic;
}

.article-body {
  display: flex;
  gap: 24px;
}

.article-main {
  flex: 1;
  min-width: 0;
}

.toc-affix {
  width: 220px;
  flex-shrink: 0;
}

.toc-wrapper {
  padding: 12px 16px;
  background: var(--card-bg);
  border-radius: 8px;
  border: 1px solid var(--border-color);
  max-height: calc(100vh - 120px);
  overflow-y: auto;
}

.toc-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 8px;
}

.toc-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.toc-item {
  a {
    color: var(--text-secondary);
    text-decoration: none;
    font-size: 13px;
    line-height: 1.8;
    display: block;
    transition: color 0.2s;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;

    &:hover {
      color: var(--primary-color);
    }
  }

  &.active a {
    color: var(--primary-color);
    font-weight: 600;
  }

  &.toc-level-3 { padding-left: 12px; }
  &.toc-level-4 { padding-left: 24px; }
  &.toc-level-5 { padding-left: 36px; }
}

.article-actions {
  display: flex;
  gap: 12px;
}

.comment-section {
  h3 {
    margin-bottom: 16px;
  }
}

.comment-form {
  display: flex;
  flex-direction: column;
}

.comment-nickname {
  display: flex;
}

.comment-form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.reply-hint {
  font-size: 13px;
  color: var(--text-secondary);
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.comment-item {
  display: flex;
  gap: 12px;
}

.comment-main {
  display: flex;
  gap: 12px;
  width: 100%;
}

.comment-avatar {
  flex-shrink: 0;
  background: var(--primary-color);
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.comment-author {
  font-weight: 600;
  font-size: 14px;
}

.comment-time {
  font-size: 12px;
  color: var(--text-secondary);
}

.reply-to {
  font-size: 12px;
  color: var(--primary-color);
}

.comment-content {
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 4px;
  word-break: break-word;
}

.comment-replies {
  margin-left: 44px;
  padding-left: 16px;
  border-left: 2px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 8px;
}

.reply-item {
  .comment-avatar {
    font-size: 12px;
  }
}

.back-btn {
  margin-top: 24px;
}

@media (max-width: 992px) {
  .article-body {
    flex-direction: column;
  }
  .toc-affix {
    width: 100%;
    order: -1;
  }
  .toc-wrapper {
    max-height: 200px;
  }
}
</style>
