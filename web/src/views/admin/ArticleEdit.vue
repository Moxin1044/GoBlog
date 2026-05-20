<template>
  <div class="article-edit">
    <div class="page-header mb-16">
      <h2>{{ isEdit ? $t('common.edit') : $t('common.create') }} {{ $t('article.title') }}</h2>
      <a-button @click="$router.back()">{{ $t('common.back') }}</a-button>
    </div>

    <a-form :model="formState" layout="vertical" @finish="handleSubmit">
      <a-row :gutter="24">
        <a-col :xs="24" :lg="16">
          <a-form-item :label="$t('article.title')" name="title" :rules="[{ required: true, message: $t('admin.titleRequired') }]">
            <a-input v-model:value="formState.title" size="large" :placeholder="$t('article.title')" />
          </a-form-item>

          <!-- Markdown Editor -->
          <a-form-item :label="$t('article.content')" name="content" :rules="[{ required: true, message: $t('admin.contentRequired') }]">
            <div class="md-editor">
              <!-- Toolbar -->
              <div class="md-toolbar">
                <a-space>
                  <a-button size="small" @click="insertMd('**', '**')" :title="$t('admin.bold')">
                    <BoldOutlined />
                  </a-button>
                  <a-button size="small" @click="insertMd('*', '*')" :title="$t('admin.italic')">
                    <ItalicOutlined />
                  </a-button>
                  <a-button size="small" @click="insertMd('## ', '')" :title="$t('admin.heading')">
                    <FontSizeOutlined />
                  </a-button>
                  <a-button size="small" @click="insertMd('[', '](url)')" :title="$t('admin.link')">
                    <LinkOutlined />
                  </a-button>
                  <a-button size="small" @click="insertMd('![alt](', ')')" :title="$t('admin.image')">
                    <PictureOutlined />
                  </a-button>
                  <a-button size="small" @click="insertMd('```\n', '\n```')" :title="$t('admin.codeBlock')">
                    <CodeOutlined />
                  </a-button>
                  <a-button size="small" @click="insertMd('| Header | Header |\n| --- | --- |\n| Cell | Cell |\n', '')" :title="$t('admin.table')">
                    <TableOutlined />
                  </a-button>
                  <a-button size="small" @click="insertMd('$$', '$$')" :title="$t('admin.formula')">
                    <FunctionOutlined />
                  </a-button>
                  <a-divider type="vertical" />
                  <a-radio-group v-model:value="editorMode" size="small" button-style="solid">
                    <a-radio-button value="split">{{ $t('admin.splitMode') }}</a-radio-button>
                    <a-radio-button value="edit">{{ $t('admin.editMode') }}</a-radio-button>
                  </a-radio-group>
                </a-space>
              </div>
              <!-- Editor Area -->
              <div class="md-content" :class="{ 'split-mode': editorMode === 'split', 'edit-mode': editorMode === 'edit' }">
                <div class="md-input">
                  <a-textarea
                    ref="textareaRef"
                    v-model:value="formState.content"
                    :rows="20"
                    :placeholder="$t('article.content')"
                    @paste="handlePaste"
                    class="md-textarea"
                  />
                </div>
                <div v-if="editorMode === 'split'" class="md-preview">
                  <div class="preview-content" v-html="renderedContent"></div>
                </div>
              </div>
            </div>
          </a-form-item>

          <a-form-item :label="$t('article.summary')" name="summary">
            <a-textarea v-model:value="formState.summary" :rows="3" :placeholder="$t('article.summary')" />
          </a-form-item>
        </a-col>

        <a-col :xs="24" :lg="8">
          <a-card :title="$t('admin.publishSettings')" size="small" class="mb-16">
            <a-form-item :label="$t('article.cover')" name="cover">
              <a-upload
                :before-upload="handleCoverUpload"
                :show-upload-list="false"
                accept="image/*"
              >
                <div v-if="formState.cover" class="cover-preview">
                  <img :src="formState.cover" alt="cover" />
                </div>
                <a-button v-else type="dashed" block>
                  <UploadOutlined /> {{ $t('article.cover') }}
                </a-button>
              </a-upload>
            </a-form-item>
            <a-form-item :label="$t('article.category')" name="category_id">
              <div class="select-with-add">
                <a-select v-model:value="formState.category_id" :placeholder="$t('article.category')" allow-clear style="flex: 1">
                  <a-select-option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</a-select-option>
                </a-select>
                <a-button type="text" @click="categoryModalVisible = true">
                  <PlusOutlined />
                </a-button>
              </div>
            </a-form-item>
            <a-form-item :label="$t('article.tags')" name="tag_ids">
              <div class="select-with-add">
                <a-select v-model:value="formState.tag_ids" mode="multiple" :placeholder="$t('article.tags')" style="flex: 1">
                  <a-select-option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</a-select-option>
                </a-select>
                <a-button type="text" @click="tagModalVisible = true">
                  <PlusOutlined />
                </a-button>
              </div>
            </a-form-item>
            <a-form-item :label="$t('article.publishTime')" name="published_at">
              <a-date-picker
                v-model:value="formState.published_at"
                show-time
                format="YYYY-MM-DD HH:mm:ss"
                :placeholder="$t('article.publishTime')"
                style="width: 100%"
              />
            </a-form-item>
            <a-form-item :label="$t('common.status')" name="status">
              <a-select v-model:value="formState.status">
                <a-select-option value="draft">{{ $t('article.draft') }}</a-select-option>
                <a-select-option value="published">{{ $t('article.published') }}</a-select-option>
              </a-select>
            </a-form-item>
          </a-card>

          <div class="mt-16">
            <a-space>
              <a-button type="primary" html-type="submit" :loading="saving">
                {{ formState.status === 'published' ? $t('article.published') : $t('common.save') }}
              </a-button>
              <a-button @click="$router.back()">{{ $t('common.cancel') }}</a-button>
            </a-space>
          </div>
        </a-col>
      </a-row>
    </a-form>

    <!-- Category Modal -->
    <a-modal
      v-model:open="categoryModalVisible"
      :title="$t('common.create') + $t('article.category')"
      @ok="handleCreateCategory"
      @cancel="categoryModalVisible = false"
    >
      <a-form layout="vertical">
        <a-form-item :label="$t('category.name')">
          <a-input v-model:value="newCategory.name" />
        </a-form-item>
        <a-form-item :label="$t('category.nameEn')">
          <a-input v-model:value="newCategory.nameEn" />
        </a-form-item>
        <a-form-item :label="$t('category.sort')">
          <a-input-number v-model:value="newCategory.sort" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Tag Modal -->
    <a-modal
      v-model:open="tagModalVisible"
      :title="$t('common.create') + $t('article.tag')"
      @ok="handleCreateTag"
      @cancel="tagModalVisible = false"
    >
      <a-form layout="vertical">
        <a-form-item :label="$t('tag.name')">
          <a-input v-model:value="newTag.name" />
        </a-form-item>
        <a-form-item :label="$t('tag.nameEn')">
          <a-input v-model:value="newTag.nameEn" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  UploadOutlined, BoldOutlined, ItalicOutlined, FontSizeOutlined,
  LinkOutlined, PictureOutlined, CodeOutlined, TableOutlined, FunctionOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { createArticle, adminGetArticle, updateArticle, adminGetCategories, adminGetTags, uploadImage, createCategory, createTag } from '@/api/admin'
import { useI18n } from 'vue-i18n'
import { marked } from 'marked'
import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const saving = ref(false)
const categories = ref<any[]>([])
const tags = ref<any[]>([])
const editorMode = ref<'split' | 'edit'>('split')
const textareaRef = ref<any>(null)

// Category modal
const categoryModalVisible = ref(false)
const newCategory = reactive({ name: '', nameEn: '', sort: 0 })

// Tag modal
const tagModalVisible = ref(false)
const newTag = reactive({ name: '', nameEn: '' })

const isEdit = computed(() => !!route.params.id)

const formState = reactive({
  title: '',
  content: '',
  summary: '',
  cover: '',
  category_id: undefined as number | undefined,
  tag_ids: [] as number[],
  status: 'draft',
  published_at: null as Dayjs | null,
})

const renderedContent = computed(() => {
  if (!formState.content) return ''
  return marked(formState.content)
})

function insertMd(before: string, after: string) {
  const el = textareaRef.value?.$el?.querySelector('textarea') as HTMLTextAreaElement | undefined
  if (!el) return
  const start = el.selectionStart
  const end = el.selectionEnd
  const selected = formState.content.substring(start, end)
  const replacement = before + selected + after
  formState.content = formState.content.substring(0, start) + replacement + formState.content.substring(end)
  // Restore cursor
  nextTick(() => {
    el.focus()
    el.setSelectionRange(start + before.length, start + before.length + selected.length)
  })
}

function nextTick(cb: () => void) {
  setTimeout(cb, 0)
}

async function handlePaste(e: ClipboardEvent) {
  const items = e.clipboardData?.items
  if (!items) return
  for (const item of items) {
    if (item.type.startsWith('image/')) {
      e.preventDefault()
      const file = item.getAsFile()
      if (!file) continue
      const formData = new FormData()
      formData.append('file', file)
      try {
        const res = await uploadImage(formData)
        const url = res.data?.url || ''
        const imgMd = `![image](${url})`
        formState.content += `\n${imgMd}\n`
      } catch { /* handled */ }
      break
    }
  }
}

async function fetchCategories() {
  try {
    const res = await adminGetCategories()
    categories.value = res.data?.list || res.data || []
  } catch { /* handled */ }
}

async function fetchTags() {
  try {
    const res = await adminGetTags()
    tags.value = res.data?.list || res.data || []
  } catch { /* handled */ }
}

async function fetchArticle() {
  if (!isEdit.value) return
  try {
    const id = Number(route.params.id)
    const res = await adminGetArticle(id)
    const data = res.data
    Object.assign(formState, {
      title: data.title,
      content: data.content,
      summary: data.summary,
      cover: data.cover,
      category_id: data.category_id,
      tag_ids: (data.tags || []).map((t: any) => t.id),
      status: data.status,
      published_at: data.published_at ? dayjs(data.published_at) : null,
    })
  } catch { /* handled */ }
}

async function handleCoverUpload(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  try {
    const res = await uploadImage(formData)
    formState.cover = res.data?.url || ''
  } catch { /* handled */ }
  return false
}

async function handleCreateCategory() {
  if (!newCategory.name) {
    message.warning(t('category.nameRequired'))
    return
  }
  try {
    const res = await createCategory(newCategory)
    categories.value.push(res.data)
    formState.category_id = res.data.id
    categoryModalVisible.value = false
    newCategory.name = ''
    newCategory.nameEn = ''
    newCategory.sort = 0
    message.success(t('common.success'))
  } catch { /* handled */ }
}

async function handleCreateTag() {
  if (!newTag.name) {
    message.warning(t('tag.nameRequired'))
    return
  }
  try {
    const res = await createTag(newTag)
    tags.value.push(res.data)
    formState.tag_ids.push(res.data.id)
    tagModalVisible.value = false
    newTag.name = ''
    newTag.nameEn = ''
    message.success(t('common.success'))
  } catch { /* handled */ }
}

async function handleSubmit() {
  saving.value = true
  try {
    const payload: Record<string, any> = { ...formState }
    if (formState.published_at) {
      payload.published_at = formState.published_at.format('YYYY-MM-DD HH:mm:ss')
    } else {
      delete payload.published_at
    }
    if (isEdit.value) {
      await updateArticle(Number(route.params.id), payload)
    } else {
      await createArticle(payload)
    }
    message.success(t('common.success'))
    router.push('/admin/articles')
  } catch { /* handled */ } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchCategories()
  fetchTags()
  fetchArticle()
})
</script>

<style scoped lang="less">
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  h2 {
    margin: 0;
  }
}

.select-with-add {
  display: flex;
  gap: 8px;
  align-items: center;
}

.cover-preview {
  img {
    width: 100%;
    max-height: 150px;
    object-fit: cover;
    border-radius: 6px;
  }
}

.md-editor {
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: 6px;
  overflow: hidden;

  .md-toolbar {
    padding: 8px 12px;
    border-bottom: 1px solid var(--border-color, #d9d9d9);
    background: var(--component-background, #fafafa);
  }

  .md-content {
    display: flex;
    min-height: 400px;

    &.split-mode {
      .md-input {
        width: 50%;
        border-right: 1px solid var(--border-color, #d9d9d9);
      }
      .md-preview {
        width: 50%;
      }
    }

    &.edit-mode {
      .md-input {
        width: 100%;
      }
      .md-preview {
        display: none;
      }
    }
  }

  .md-input {
    .md-textarea {
      border: none;
      border-radius: 0;
      box-shadow: none;
      resize: none;
      font-family: 'Courier New', Courier, monospace;
      font-size: 14px;
      line-height: 1.6;

      &:focus {
        box-shadow: none;
      }
    }
  }

  .md-preview {
    padding: 12px 16px;
    overflow-y: auto;
    background: var(--component-background, #fff);

    .preview-content {
      font-size: 14px;
      line-height: 1.8;

      :deep(h1), :deep(h2), :deep(h3), :deep(h4) {
        margin-top: 16px;
        margin-bottom: 8px;
      }

      :deep(pre) {
        background: var(--component-background, #f5f5f5);
        padding: 12px;
        border-radius: 4px;
        overflow-x: auto;
      }

      :deep(code) {
        background: var(--component-background, #f5f5f5);
        padding: 2px 4px;
        border-radius: 3px;
        font-size: 13px;
      }

      :deep(img) {
        max-width: 100%;
      }

      :deep(table) {
        border-collapse: collapse;
        width: 100%;

        th, td {
          border: 1px solid var(--border-color, #d9d9d9);
          padding: 8px 12px;
        }
      }
    }
  }
}
</style>
