<template>
  <div class="home-page">
    <!-- Banner / Carousel -->
    <div class="page-banner" v-if="featuredArticles.length">
      <a-carousel autoplay :dots="true" class="banner-carousel">
        <div v-for="article in featuredArticles" :key="article.id" class="banner-slide" @click="$router.push(`/article/${article.id}`)">
          <div class="banner-image-wrapper">
            <img v-if="article.cover" :src="article.cover" :alt="article.title" class="banner-image" />
            <div v-else class="banner-image banner-placeholder">
              <div class="placeholder-text">{{ appStore.siteConfig.siteName }}</div>
            </div>
          </div>
          <div class="banner-overlay">
            <h2 class="banner-article-title">{{ article.title }}</h2>
            <p class="banner-article-summary">{{ article.summary }}</p>
          </div>
        </div>
      </a-carousel>
    </div>
    <div class="page-banner banner-simple" v-else>
      <h1 class="banner-title">{{ appStore.siteConfig.siteName }}</h1>
      <p class="banner-desc">{{ $t('common.search') }}...</p>
    </div>

    <!-- Search Bar -->
    <div class="search-bar mb-16">
      <a-input-search
        v-model:value="searchText"
        :placeholder="$t('common.search')"
        size="large"
        enter-button
        allow-clear
        @search="handleSearch"
      />
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar mb-16">
      <div class="filter-left">
        <a-radio-group v-model:value="filterType" button-style="solid" @change="handleFilterChange">
          <a-radio-button value="all">{{ $t('common.all') }}</a-radio-button>
          <a-radio-button value="category">{{ $t('article.categoryFilter') }}</a-radio-button>
          <a-radio-button value="tag">{{ $t('article.tagFilter') }}</a-radio-button>
        </a-radio-group>
      </div>
    </div>

    <!-- Category Tags -->
    <div v-if="filterType === 'category'" class="filter-tags mb-16">
      <a-tag
        v-for="cat in categories"
        :key="cat.id"
        :color="selectedCategory === cat.id ? 'blue' : ''"
        class="filter-tag"
        @click="toggleCategory(cat.id)"
      >
        {{ cat.name }} ({{ cat.count || 0 }})
      </a-tag>
    </div>

    <!-- Tag Cloud -->
    <div v-if="filterType === 'tag'" class="filter-tags mb-16">
      <a-tag
        v-for="tag in tags"
        :key="tag.id"
        :color="selectedTag === tag.id ? 'blue' : getTagColor(tag.id)"
        class="filter-tag"
        @click="toggleTag(tag.id)"
      >
        {{ tag.name }}
      </a-tag>
    </div>

    <!-- Active Filters -->
    <div v-if="activeFilters.length" class="active-filters mb-16">
      <a-tag
        v-for="filter in activeFilters"
        :key="filter.key"
        closable
        @close="removeFilter(filter.key)"
        color="blue"
      >
        {{ filter.label }}
      </a-tag>
      <a-button type="link" size="small" @click="clearFilters">{{ $t('common.all') }}</a-button>
    </div>

    <!-- Search Result Hint -->
    <div v-if="searchKeyword" class="search-hint mb-16">
      {{ $t('article.searchResult') }}: "{{ searchKeyword }}"
      <a-button type="link" size="small" @click="clearSearch">{{ $t('common.cancel') }}</a-button>
    </div>

    <!-- Article Grid -->
    <a-spin :spinning="loading">
      <a-row :gutter="[16, 16]">
        <a-col v-for="article in articles" :key="article.id" :xs="24" :sm="12" :lg="8">
          <ArticleCard :article="article" />
        </a-col>
      </a-row>
      <a-empty v-if="!loading && !articles.length" :description="$t('common.noData')" class="mt-24" />
    </a-spin>

    <!-- Pagination -->
    <div class="pagination-wrapper mt-24" v-if="total > pageSize">
      <a-pagination
        v-model:current="currentPage"
        :total="total"
        :page-size="pageSize"
        show-quick-jumper
        :show-total="(total: number) => `${$t('common.total')} ${total} ${$t('common.items')}`"
        @change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getArticleList, getCategories, getTags } from '@/api/article'
import { useAppStore } from '@/stores/app'
import ArticleCard from '@/components/ArticleCard.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const appStore = useAppStore()

const articles = ref<any[]>([])
const featuredArticles = ref<any[]>([])
const categories = ref<any[]>([])
const tags = ref<any[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(9)
const total = ref(0)
const filterType = ref('all')
const selectedCategory = ref<number | null>(null)
const selectedTag = ref<number | null>(null)
const searchText = ref('')
const searchKeyword = ref('')

const tagColors = ['blue', 'green', 'orange', 'red', 'purple', 'cyan', 'magenta', 'geekblue', 'volcano', 'gold']

function getTagColor(id: number) {
  return tagColors[id % tagColors.length]
}

const activeFilters = computed(() => {
  const filters: { key: string; label: string }[] = []
  if (selectedCategory.value) {
    const cat = categories.value.find(c => c.id === selectedCategory.value)
    if (cat) filters.push({ key: 'category', label: cat.name })
  }
  if (selectedTag.value) {
    const tag = tags.value.find(t => t.id === selectedTag.value)
    if (tag) filters.push({ key: 'tag', label: tag.name })
  }
  return filters
})

async function fetchFeatured() {
  try {
    const res = await getArticleList({ page: 1, page_size: 5, is_top: true })
    featuredArticles.value = res.data?.list || []
    if (!featuredArticles.value.length) {
      const res2 = await getArticleList({ page: 1, page_size: 5 })
      featuredArticles.value = res2.data?.list || []
    }
  } catch {
    // silently fail
  }
}

async function fetchCategories() {
  try {
    const res = await getCategories()
    categories.value = Array.isArray(res.data) ? res.data : (res.data?.list || [])
  } catch {
    // silently fail
  }
}

async function fetchTags() {
  try {
    const res = await getTags()
    tags.value = Array.isArray(res.data) ? res.data : (res.data?.list || [])
  } catch {
    // silently fail
  }
}

async function fetchArticles() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    if (searchKeyword.value) params.search = searchKeyword.value
    if (selectedCategory.value) params.category_id = selectedCategory.value
    if (selectedTag.value) params.tag_id = selectedTag.value
    const res = await getArticleList(params)
    articles.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function handleFilterChange() {
  currentPage.value = 1
  selectedCategory.value = null
  selectedTag.value = null
  fetchArticles()
}

function toggleCategory(id: number) {
  selectedCategory.value = selectedCategory.value === id ? null : id
  currentPage.value = 1
  fetchArticles()
}

function toggleTag(id: number) {
  selectedTag.value = selectedTag.value === id ? null : id
  currentPage.value = 1
  fetchArticles()
}

function removeFilter(key: string) {
  if (key === 'category') selectedCategory.value = null
  if (key === 'tag') selectedTag.value = null
  currentPage.value = 1
  fetchArticles()
}

function clearFilters() {
  selectedCategory.value = null
  selectedTag.value = null
  currentPage.value = 1
  fetchArticles()
}

function handleSearch(value: string) {
  searchKeyword.value = value.trim()
  currentPage.value = 1
  fetchArticles()
}

function clearSearch() {
  searchText.value = ''
  searchKeyword.value = ''
  currentPage.value = 1
  fetchArticles()
}

watch(() => route.query, (query) => {
  const search = query.search as string
  const category = query.category as string
  const tag = query.tag as string
  if (search) {
    searchKeyword.value = search
    searchText.value = search
  }
  if (category) {
    selectedCategory.value = Number(category)
    filterType.value = 'category'
  }
  if (tag) {
    selectedTag.value = Number(tag)
    filterType.value = 'tag'
  }
  fetchArticles()
}, { immediate: false })

onMounted(() => {
  const search = route.query.search as string
  const category = route.query.category as string
  const tag = route.query.tag as string
  if (search) {
    searchKeyword.value = search
    searchText.value = search
  }
  if (category) {
    selectedCategory.value = Number(category)
    filterType.value = 'category'
  }
  if (tag) {
    selectedTag.value = Number(tag)
    filterType.value = 'tag'
  }
  fetchFeatured()
  fetchCategories()
  fetchTags()
  fetchArticles()
})
</script>

<style scoped lang="less">
.home-page {
  max-width: 100%;
}

.page-banner {
  margin-bottom: 24px;
  border-radius: 12px;
  overflow: hidden;
}

.banner-simple {
  text-align: center;
  padding: 48px 0 32px;
  background: linear-gradient(135deg, var(--primary-color) 0%, #36cfc9 100%);
  color: #fff;
}

.banner-title {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 8px;
}

.banner-desc {
  color: rgba(255, 255, 255, 0.85);
  font-size: 16px;
}

.banner-carousel {
  :deep(.slick-dots) {
    bottom: 12px;
    li button {
      background: rgba(255, 255, 255, 0.6);
      border-radius: 4px;
    }
    li.slick-active button {
      background: #fff;
    }
  }
}

.banner-slide {
  cursor: pointer;
  position: relative;
}

.banner-image-wrapper {
  height: 320px;
  overflow: hidden;
}

.banner-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.banner-placeholder {
  background: linear-gradient(135deg, var(--primary-color) 0%, #36cfc9 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.placeholder-text {
  color: #fff;
  font-size: 32px;
  font-weight: 700;
}

.banner-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 24px 32px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  color: #fff;
}

.banner-article-title {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 8px;
  color: #fff;
}

.banner-article-summary {
  font-size: 14px;
  opacity: 0.9;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.search-bar {
  max-width: 500px;
  margin-left: auto;
  margin-right: auto;
}

.filter-bar {
  display: flex;
  justify-content: center;
}

.filter-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: center;
}

.filter-tag {
  cursor: pointer;
  transition: transform 0.2s;
  &:hover {
    transform: scale(1.05);
  }
}

.active-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.search-hint {
  font-size: 14px;
  color: var(--text-secondary);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}

@media (max-width: 576px) {
  .banner-image-wrapper {
    height: 200px;
  }
  .banner-article-title {
    font-size: 18px;
  }
  .banner-overlay {
    padding: 16px;
  }
}
</style>
