<template>
  <div class="system-config">
    <div class="page-header mb-16">
      <h2>{{ $t('admin.systemConfig') }}</h2>
    </div>

    <a-spin :spinning="loading">
      <a-tabs v-model:activeKey="activeTab">
        <!-- Site Basic Config -->
        <a-tab-pane key="basic" :tab="$t('admin.siteBasic')">
          <a-form :model="formState" layout="vertical" @finish="handleSaveBasic">
            <a-form-item :label="$t('site.siteName')" name="site_name">
              <a-input v-model:value="formState.site_name" />
            </a-form-item>
            <a-form-item :label="$t('site.siteLogo')" name="site_logo">
              <a-upload
                :before-upload="handleLogoUpload"
                :show-upload-list="false"
                accept="image/*"
              >
                <div v-if="formState.site_logo" class="logo-preview">
                  <img :src="formState.site_logo" alt="Logo" />
                </div>
                <a-button v-else type="dashed">
                  <UploadOutlined /> {{ $t('common.upload') }}
                </a-button>
              </a-upload>
            </a-form-item>
            <a-form-item :label="$t('site.copyright')" name="copyright">
              <a-input v-model:value="formState.copyright" />
            </a-form-item>
            <a-form-item :label="$t('site.icp')" name="icp">
              <a-input v-model:value="formState.icp" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" html-type="submit" :loading="saving">{{ $t('common.save') }}</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- Register Switch -->
        <a-tab-pane key="register" :tab="$t('site.registerSwitch')">
          <a-form :model="formState" layout="vertical" @finish="handleSaveRegister">
            <a-form-item :label="$t('site.registerSwitch')">
              <a-switch v-model:checked="formState.register_enabled" />
              <span class="ml-8">{{ formState.register_enabled ? $t('common.enable') : $t('common.disable') }}</span>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" html-type="submit" :loading="saving">{{ $t('common.save') }}</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- Notification Config -->
        <a-tab-pane key="notification" :tab="$t('admin.notificationConfig')">
          <a-form :model="formState" layout="vertical" @finish="handleSaveNotification">
            <a-card :title="$t('site.feishuConfig')" size="small" class="mb-16">
              <a-form-item label="Webhook URL">
                <a-input v-model:value="formState.feishu_webhook" />
              </a-form-item>
              <a-form-item :label="$t('user.feishuToken')">
                <a-input-password v-model:value="formState.feishu_token" />
              </a-form-item>
            </a-card>

            <a-card :title="$t('site.smtpConfig')" size="small" class="mb-16">
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item label="SMTP Host">
                    <a-input v-model:value="formState.smtp_host" />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item label="SMTP Port">
                    <a-input-number v-model:value="formState.smtp_port" style="width: 100%" />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item label="SMTP User">
                    <a-input v-model:value="formState.smtp_user" />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item label="SMTP Password">
                    <a-input-password v-model:value="formState.smtp_password" />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item :label="$t('admin.smtpSender')">
                    <a-input v-model:value="formState.smtp_sender" />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-card>

            <a-form-item>
              <a-button type="primary" html-type="submit" :loading="saving">{{ $t('common.save') }}</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { UploadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { getSystemConfig, updateSystemConfig, uploadImage } from '@/api/admin'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const activeTab = ref('basic')

const formState = reactive({
  site_name: 'GoBlog',
  site_logo: '',
  copyright: 'Copyright © 2024 GoBlog. All rights reserved.',
  icp: '',
  register_enabled: true,
  smtp_host: '',
  smtp_port: 465,
  smtp_user: '',
  smtp_password: '',
  smtp_sender: '',
  feishu_webhook: '',
  feishu_token: '',
})

async function fetchConfig() {
  loading.value = true
  try {
    const res = await getSystemConfig()
    const data = res.data || {}
    // 处理 register_enabled 可能是字符串的情况
    if (data.register_enabled !== undefined) {
      data.register_enabled = data.register_enabled === 'true' || data.register_enabled === true
    }
    Object.assign(formState, data)
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

async function handleLogoUpload(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  try {
    const res = await uploadImage(formData)
    formState.site_logo = res.data?.url || ''
  } catch { /* handled */ }
  return false
}

async function handleSaveBasic() {
  saving.value = true
  try {
    await updateSystemConfig({
      site_name: formState.site_name,
      site_logo: formState.site_logo,
      copyright: formState.copyright,
      icp: formState.icp,
    })
    message.success(t('common.success'))
  } catch { /* handled */ } finally {
    saving.value = false
  }
}

async function handleSaveRegister() {
  saving.value = true
  try {
    await updateSystemConfig({ register_enabled: formState.register_enabled })
    message.success(t('common.success'))
  } catch { /* handled */ } finally {
    saving.value = false
  }
}

async function handleSaveNotification() {
  saving.value = true
  try {
    await updateSystemConfig({
      feishu_webhook: formState.feishu_webhook,
      feishu_token: formState.feishu_token,
      smtp_host: formState.smtp_host,
      smtp_port: formState.smtp_port,
      smtp_user: formState.smtp_user,
      smtp_password: formState.smtp_password,
      smtp_sender: formState.smtp_sender,
    })
    message.success(t('common.success'))
  } catch { /* handled */ } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="less">
.page-header {
  h2 { margin: 0; }
}

.logo-preview {
  img {
    max-height: 60px;
    border-radius: 4px;
  }
}

.ml-8 {
  margin-left: 8px;
}
</style>
