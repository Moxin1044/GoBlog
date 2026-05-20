<template>
  <div class="app-sidebar">
    <!-- Categories -->
    <a-card class="sidebar-card" :title="$t('nav.categories')" size="small">
      <div class="category-list">
        <a-tag
          v-for="cat in categories"
          :key="cat.id"
          class="category-tag"
          :color="selectedCategory === cat.id ? 'blue' : ''"
          @click="handleCategoryClick(cat)"
        >
          {{ cat.name }} ({{ cat.count || 0 }})
        </a-tag>
        <a-empty v-if="!categories.length" :description="$t('common.noData')" :image="simpleImage" />
      </div>
    </a-card>

    <!-- Tag Cloud -->
    <a-card class="sidebar-card" :title="$t('nav.tags')" size="small">
      <div class="tag-cloud">
        <a-tag
          v-for="tag in tags"
          :key="tag.id"
          :color="getTagColor(tag.id)"
          :style="{ fontSize: getTagSize(tag.count || 0) }"
          class="cloud-tag"
          @click="handleTagClick(tag)"
        >
          {{ tag.name }}
        </a-tag>
        <a-empty v-if="!tags.length" :description="$t('common.noData')" :image="simpleImage" />
      </div>
    </a-card>

    <!-- Hot Articles -->
    <a-card class="sidebar-card" :title="`🔥 ${$t('article.viewCount')}`" size="small">
      <a-list :data-source="hotArticles" size="small" :split="true">
        <template #renderItem="{ item, index }">
          <a-list-item class="hot-article-item">
            <router-link :to="`/article/${item.id}`" class="hot-article-link">
              <span :class="['hot-rank', { 'top-three': index < 3 }]">{{ index + 1 }}</span>
              {{ item.title }}
            </router-link>
          </a-list-item>
        </template>
      </a-list>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Empty } from 'ant-design-vue'
import { getCategories, getTags, getArticleList } from '@/api/article'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

const categories = ref<any[]>([])
const tags = ref<any[]>([])
const hotArticles = ref<any[]>([])
const selectedCategory = ref<number | null>(null)

const tagColors = ['blue', 'green', 'orange', 'red', 'purple', 'cyan', 'magenta', 'geekblue', 'volcano', 'gold']

function getTagColor(id: number) {
  return tagColors[id % tagColors.length]
}

function getTagSize(count: number) {
  const minSize = 12
  const maxSize = 18
  const maxCount = Math.max(...tags.value.map(t => t.count || 0), 1)
  const ratio = count / maxCount
  return `${minSize + ratio * (maxSize - minSize)}px`
}

function handleCategoryClick(cat: any) {
  selectedCategory.value = selectedCategory.value === cat.id ? null : cat.id
  router.push(selectedCategory.value ? `/?category=${cat.id}` : '/')
}

function handleTagClick(tag: any) {
  router.push(`/?tag=${tag.id}`)
}

onMounted(async () => {
  try {
    const [catRes, tagRes, articleRes] = await Promise.all([
      getCategories(),
      getTags(),
      getArticleList({ page: 1, page_size: 10, sort: 'view_count' }),
    ])
    categories.value = catRes.data?.list || []
    tags.value = tagRes.data?.list || []
    hotArticles.value = articleRes.data?.list || []
  } catch {
    // silently fail
  }
})
</script>

<style scoped lang="less">
.app-sidebar {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.sidebar-card {
  border-radius: 8px;
}

.category-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.category-tag {
  cursor: pointer;
  transition: all 0.2s;
  &:hover {
    transform: scale(1.05);
    opacity: 0.85;
  }
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.cloud-tag {
  cursor: pointer;
  transition: all 0.2s;
  &:hover {
    transform: scale(1.1);
  }
}

.hot-article-item {
  padding: 4px 0 !important;
}

.hot-article-link {
  color: var(--text-color);
  text-decoration: none;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;

  &:hover {
    color: var(--primary-color);
  }
}

.hot-rank {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  background: var(--border-color);
  color: var(--text-secondary);
  flex-shrink: 0;

  &.top-three {
    background: var(--primary-color);
    color: #fff;
  }
}
</style>
