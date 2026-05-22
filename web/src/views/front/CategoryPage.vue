<template>
  <div class="category-page">
    <a-page-header :title="category?.name || $t('article.category')" class="page-header" />
    <div class="article-list">
      <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
    </div>
    <a-pagination
      v-model:current="page"
      :total="total"
      :page-size="pageSize"
      :show-total="(total: number) => $t('common.totalCount', { total })"
      @change="fetchArticles"
      class="pagination"
    />
    <a-empty v-if="articles.length === 0 && loading === false" :description="$t('common.noData')" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import ArticleCard from '@/components/ArticleCard.vue'
import { getArticleList, getCategories } from '@/api/article'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const articles = ref<any[]>([])
const categories = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const category = computed(() => {
  const id = Number(route.params.id)
  return categories.value.find(c => c.id === id)
})

const fetchArticles = async () => {
  loading.value = true
  try {
    const res = await getArticleList({
      page: page.value,
      page_size: pageSize.value,
      category_id: route.params.id
    })
    articles.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    message.error(t('common.error'))
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const res = await getCategories()
    categories.value = res.data || []
  } catch {
    message.error(t('common.error'))
  }
}

onMounted(() => {
  fetchCategories()
  fetchArticles()
})
</script>

<style scoped lang="less">
.category-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;

  .page-header {
    margin-bottom: 24px;
    background: #fff;
    border-radius: 8px;
  }

  .article-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 24px;
  }

  .pagination {
    margin-top: 32px;
    display: flex;
    justify-content: center;
  }
}
</style>
