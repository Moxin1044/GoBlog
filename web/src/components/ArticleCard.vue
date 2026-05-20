<template>
  <a-card hoverable class="article-card" @click="$router.push(`/article/${article.id}`)">
    <template #cover>
      <div class="card-cover-wrapper">
        <img v-if="article.cover" :src="article.cover" :alt="article.title" class="card-cover" />
        <div v-else class="card-cover card-placeholder">
          <FileTextOutlined class="placeholder-icon" />
        </div>
      </div>
    </template>
    <div class="card-body">
      <h3 class="card-title">{{ article.title }}</h3>
      <p class="card-summary">{{ article.summary }}</p>
      <div class="card-tags">
        <a-tag v-if="article.category" color="blue" size="small">{{ article.category.name || article.category }}</a-tag>
        <a-tag v-for="tag in (article.tags || []).slice(0, 2)" :key="tag.id" size="small">{{ tag.name }}</a-tag>
      </div>
      <div class="card-meta-info">
        <span class="meta-item">
          <ClockCircleOutlined /> {{ formatDate(article.published_at || article.created_at) }}
        </span>
        <span class="meta-item">
          <EyeOutlined /> {{ article.view_count || 0 }}
        </span>
        <span class="meta-item">
          <LikeOutlined /> {{ article.like_count || 0 }}
        </span>
      </div>
    </div>
  </a-card>
</template>

<script setup lang="ts">
import { ClockCircleOutlined, EyeOutlined, LikeOutlined, FileTextOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'

defineProps<{
  article: Record<string, any>
}>()

function formatDate(date: string) {
  return date ? dayjs(date).format('YYYY-MM-DD') : ''
}
</script>

<style scoped lang="less">
.article-card {
  cursor: pointer;
  border-radius: 12px;
  transition: transform 0.25s ease, box-shadow 0.25s ease;
  overflow: hidden;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  }

  :deep(.ant-card-cover) {
    border-top-left-radius: 12px;
    border-top-right-radius: 12px;
  }

  :deep(.ant-card-body) {
    padding: 16px;
  }
}

.card-cover-wrapper {
  height: 180px;
  overflow: hidden;
}

.card-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;

  .article-card:hover & {
    transform: scale(1.05);
  }
}

.card-placeholder {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.placeholder-icon {
  font-size: 48px;
  color: rgba(255, 255, 255, 0.6);
}

.card-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  line-height: 1.4;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.card-summary {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.card-meta-info {
  display: flex;
  gap: 12px;
  color: var(--text-secondary);
  font-size: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 3px;
}
</style>
