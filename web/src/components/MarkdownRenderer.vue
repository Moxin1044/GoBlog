<template>
  <div class="markdown-renderer" :class="{ 'minimal-mode': minimal }" v-html="renderedHtml"></div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import '@/assets/styles/markdown-dark.less'

const props = withDefaults(defineProps<{
  content: string
  minimal?: boolean
}>(), {
  minimal: false,
})

const emit = defineEmits<{
  (e: 'toc-generated', items: { id: string; text: string; level: number }[]): void
}>()

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight(str: string, lang: string): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(str, { language: lang, ignoreIllegals: true }).value}</code></pre>`
      } catch {
        // fallback
      }
    }
    return `<pre class="hljs"><code>${MarkdownIt.prototype.utils.escapeHtml(str)}</code></pre>`
  },
})

// Add heading IDs for TOC
let headingIndex = 0
const tocItems: { id: string; text: string; level: number }[] = []

md.core.ruler.push('heading_ids', (state) => {
  tocItems.length = 0
  headingIndex = 0
  for (const token of state.tokens) {
    if (token.type === 'heading_open') {
      const level = parseInt(token.tag.slice(1))
      // Find the inline token after heading_open
      const inlineToken = state.tokens[state.tokens.indexOf(token) + 1]
      if (inlineToken && inlineToken.type === 'inline') {
        const text = inlineToken.content
        const id = `heading-${headingIndex++}`
        token.attrSet('id', id)
        tocItems.push({ id, text, level })
      }
    }
  }
})

// Task list support
md.core.ruler.push('task_lists', (state) => {
  for (const token of state.tokens) {
    if (token.type === 'inline') {
      // Handle [ ] and [x] at the start of list items
      token.content = token.content.replace(/^\[(x|X| )\]\s*/, (match, check) => {
        return check.toLowerCase() === 'x' ? '✅ ' : '⬜ '
      })
    }
  }
})

const renderedHtml = computed(() => {
  const result = md.render(props.content || '')
  // Emit TOC items after render
  if (tocItems.length) {
    emit('toc-generated', [...tocItems])
  }
  return result
})
</script>

<style scoped lang="less">
.markdown-renderer {
  line-height: 1.8;
  font-size: 15px;
  word-break: break-word;

  :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
    margin-top: 1.5em;
    margin-bottom: 0.5em;
    font-weight: 600;
    scroll-margin-top: 80px;
  }

  :deep(h1) { font-size: 1.8em; }
  :deep(h2) {
    font-size: 1.5em;
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 0.3em;
  }
  :deep(h3) { font-size: 1.3em; }
  :deep(h4) { font-size: 1.15em; }

  :deep(p) {
    margin-bottom: 1em;
  }

  :deep(a) {
    color: var(--primary-color);
    text-decoration: none;
    &:hover {
      text-decoration: underline;
    }
  }

  :deep(img) {
    max-width: 100%;
    border-radius: 6px;
    margin: 8px 0;
  }

  :deep(blockquote) {
    margin: 1em 0;
    padding: 0.5em 1em;
    border-left: 4px solid var(--primary-color);
    background: var(--blockquote-bg);
    border-radius: 0 4px 4px 0;
  }

  :deep(code) {
    padding: 2px 6px;
    border-radius: 3px;
    font-size: 0.9em;
    background: var(--code-bg);
  }

  :deep(pre) {
    margin: 1em 0;
    border-radius: 8px;
    overflow-x: auto;
    border: 1px solid var(--border-color);

    code {
      padding: 16px;
      display: block;
      background: transparent;
      font-size: 14px;
      line-height: 1.6;
    }
  }

  :deep(table) {
    width: 100%;
    border-collapse: collapse;
    margin: 1em 0;
    display: block;
    overflow-x: auto;

    th, td {
      border: 1px solid var(--border-color);
      padding: 8px 12px;
      text-align: left;
    }

    th {
      background: var(--table-header-bg);
      font-weight: 600;
    }

    tr:nth-child(even) {
      background: var(--table-stripe-bg);
    }
  }

  :deep(ul), :deep(ol) {
    padding-left: 2em;
    margin-bottom: 1em;
  }

  :deep(li) {
    margin-bottom: 0.25em;
  }

  :deep(hr) {
    border: none;
    border-top: 1px solid var(--border-color);
    margin: 2em 0;
  }

  // Task list styles
  :deep(input[type="checkbox"]) {
    margin-right: 6px;
  }

  &.minimal-mode {
    font-size: 14px;
    line-height: 1.6;

    :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
      margin-top: 0.8em;
      margin-bottom: 0.3em;
    }

    :deep(p) {
      margin-bottom: 0.5em;
    }

    :deep(pre) {
      margin: 0.5em 0;
    }
  }
}
</style>
