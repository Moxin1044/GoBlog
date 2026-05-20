<template>
  <div class="profile-page">
    <a-card :title="$t('nav.profile')" class="profile-card">
      <a-tabs v-model:activeKey="activeTab">
        <!-- Basic Info Tab -->
        <a-tab-pane key="info" :tab="$t('user.nickname')">
          <a-form :model="profileForm" layout="vertical" @finish="handleUpdateProfile">
            <a-form-item :label="$t('user.avatar')">
              <div class="avatar-upload-wrapper">
                <a-upload
                  :before-upload="handleAvatarUpload"
                  :show-upload-list="false"
                  accept="image/*"
                >
                  <a-avatar :size="80" :src="profileForm.avatar" class="avatar-upload">
                    {{ profileForm.nickname?.charAt(0) || 'U' }}
                  </a-avatar>
                </a-upload>
                <div class="avatar-hint">{{ $t('common.upload') }}</div>
              </div>
            </a-form-item>
            <a-form-item :label="$t('user.nickname')" name="nickname">
              <a-input v-model:value="profileForm.nickname" />
            </a-form-item>
            <a-form-item :label="$t('auth.email')" name="email">
              <a-input v-model:value="profileForm.email" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" html-type="submit" :loading="saving">{{ $t('common.save') }}</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- Change Password Tab -->
        <a-tab-pane key="password" :tab="$t('user.changePassword')">
          <a-form :model="passwordForm" layout="vertical" @finish="handleChangePassword">
            <a-form-item :label="$t('user.oldPassword')" name="old_password">
              <a-input-password v-model:value="passwordForm.old_password" />
            </a-form-item>
            <a-form-item :label="$t('user.newPassword')" name="new_password">
              <a-input-password v-model:value="passwordForm.new_password" />
            </a-form-item>
            <a-form-item :label="$t('user.confirmPassword')" name="confirm_password">
              <a-input-password v-model:value="passwordForm.confirm_password" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" html-type="submit" :loading="changingPwd">{{ $t('common.save') }}</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- Subscription Tab -->
        <a-tab-pane key="subscription" :tab="$t('user.subscription')">
          <a-form :model="subscriptionForm" layout="vertical" @finish="handleUpdateSubscription">
            <a-form-item :label="$t('user.bindEmail')" name="email">
              <a-input v-model:value="subscriptionForm.email" />
            </a-form-item>
            <a-form-item :label="$t('user.feishuToken')" name="feishu_token">
              <a-input-password v-model:value="subscriptionForm.feishu_token" />
            </a-form-item>
            <a-form-item :label="$t('user.notifyChannel')">
              <a-checkbox-group v-model:value="subscriptionForm.notify_channels">
                <a-checkbox value="email">{{ $t('user.notifyEmail') }}</a-checkbox>
                <a-checkbox value="feishu">{{ $t('user.notifyFeishu') }}</a-checkbox>
              </a-checkbox-group>
            </a-form-item>
            <a-form-item :label="$t('user.categorySubscription')">
              <a-checkbox-group v-model:value="subscriptionForm.subscribed_categories" class="category-checkbox-group">
                <a-row :gutter="[8, 8]">
                  <a-col v-for="cat in allCategories" :key="cat.id" :span="8">
                    <a-checkbox :value="cat.id">{{ cat.name }}</a-checkbox>
                  </a-col>
                </a-row>
              </a-checkbox-group>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" html-type="submit" :loading="savingSubscription">{{ $t('common.save') }}</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- AI Config Tab -->
        <a-tab-pane key="ai" :tab="$t('user.aiConfig')">
          <a-form :model="aiForm" layout="vertical" @finish="handleUpdateAIConfig">
            <a-form-item :label="$t('user.apiToken')" name="api_token">
              <a-input-password v-model:value="aiForm.api_token" />
            </a-form-item>
            <a-form-item :label="$t('user.apiUrl')" name="api_url">
              <a-input v-model:value="aiForm.api_url" placeholder="https://api.openai.com/v1" />
            </a-form-item>
            <a-form-item :label="$t('user.model')">
              <div class="model-select-wrapper">
                <a-radio-group v-model:value="isCustomModel" style="margin-bottom: 8px">
                  <a-radio :value="false">{{ $t('user.selectFromList') }}</a-radio>
                  <a-radio :value="true">{{ $t('user.customModel') }}</a-radio>
                </a-radio-group>
                <a-select
                  v-if="!isCustomModel"
                  v-model:value="aiForm.model_id"
                  :placeholder="$t('user.selectModel')"
                  allow-clear
                  show-search
                  :filter-option="filterModelOption"
                  style="width: 100%"
                >
                  <a-select-option v-for="m in availableModels" :key="m.id" :value="m.id">
                    {{ m.provider ? `[${m.provider}] ` : '' }}{{ m.name }}
                  </a-select-option>
                </a-select>
                <a-input v-else :placeholder="$t('user.model')" />
              </div>
            </a-form-item>
            <a-form-item :label="$t('user.temperature')" name="temperature">
              <div class="slider-wrapper">
                <a-slider v-model:value="aiForm.temperature" :min="0" :max="2" :step="0.1" />
                <span class="slider-value">{{ aiForm.temperature }}</span>
              </div>
            </a-form-item>
            <a-form-item :label="$t('user.maxContext')" name="max_context">
              <a-input-number v-model:value="aiForm.max_context" :min="1" :max="128" style="width: 200px" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" html-type="submit" :loading="savingAI">{{ $t('common.save') }}</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import {
  getUserInfo, updateUserInfo, changePassword, uploadAvatar,
  getAIConfig, updateAIConfig, getSubscription, updateSubscription,
  getAvailableModels,
} from '@/api/user'
import { getCategories } from '@/api/article'

const { t } = useI18n()

const activeTab = ref('info')
const saving = ref(false)
const changingPwd = ref(false)
const savingAI = ref(false)
const savingSubscription = ref(false)
const allCategories = ref<any[]>([])
const availableModels = ref<any[]>([])

const profileForm = reactive({
  nickname: '',
  email: '',
  avatar: '',
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

const subscriptionForm = reactive({
  email: '',
  feishu_token: '',
  notify_channels: [] as string[],
  subscribed_categories: [] as number[],
})

const aiForm = reactive({
  api_token: '',
  api_url: '',
  model_id: 0,
  custom_model: '',
  temperature: 0.7,
  max_context: 10,
})

const isCustomModel = ref(false)

function filterModelOption(input: string, option: any) {
  return option.label?.toLowerCase().includes(input.toLowerCase())
}

async function fetchUserInfo() {
  try {
    const res = await getUserInfo()
    Object.assign(profileForm, res.data)
  } catch {
    // handled
  }
}

async function fetchAvailableModels() {
  try {
    const res = await getAvailableModels()
    availableModels.value = res.data || []
  } catch {
    // handled
  }
}

async function fetchAIConfig() {
  try {
    const res = await getAIConfig()
    Object.assign(aiForm, {
      api_token: res.data?.api_token || '',
      api_url: res.data?.api_url || '',
      model_id: res.data?.model_id || 0,
      temperature: res.data?.temperature || 0.7,
      max_context: res.data?.max_context || 10,
    })
    // 如果有 model_id 且存在于 availableModels 中，则不是自定义模型
    const hasSelectedModel = availableModels.value.some((m: any) => m.id === aiForm.model_id)
    isCustomModel.value = !hasSelectedModel && !!res.data?.model_id
    if (!isCustomModel.value && !hasSelectedModel) {
      aiForm.model_id = 0
    }
  } catch {
    // handled
  }
}

async function fetchSubscription() {
  try {
    const res = await getSubscription()
    Object.assign(subscriptionForm, res.data)
  } catch {
    // handled
  }
}

async function fetchCategories() {
  try {
    const res = await getCategories()
    allCategories.value = res.data?.list || []
  } catch {
    // handled
  }
}

async function handleUpdateProfile() {
  saving.value = true
  try {
    await updateUserInfo(profileForm)
    message.success(t('user.profileSaved'))
  } catch {
    // handled
  } finally {
    saving.value = false
  }
}

async function handleAvatarUpload(file: File) {
  const formData = new FormData()
  formData.append('avatar', file)
  try {
    const res = await uploadAvatar(formData)
    profileForm.avatar = res.data?.url || ''
    message.success(t('user.avatarUploadSuccess'))
  } catch {
    // handled
  }
  return false
}

async function handleChangePassword() {
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    message.error(t('user.passwordMismatch'))
    return
  }
  changingPwd.value = true
  try {
    await changePassword(passwordForm)
    message.success(t('user.passwordChanged'))
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch {
    // handled
  } finally {
    changingPwd.value = false
  }
}

async function handleUpdateSubscription() {
  savingSubscription.value = true
  try {
    await updateSubscription(subscriptionForm)
    message.success(t('user.subscriptionSaved'))
  } catch {
    // handled
  } finally {
    savingSubscription.value = false
  }
}

async function handleUpdateAIConfig() {
  savingAI.value = true
  try {
    const submitData: any = {
      api_token: aiForm.api_token,
      api_url: aiForm.api_url,
      model_id: isCustomModel.value ? 0 : aiForm.model_id,
      temperature: aiForm.temperature,
      max_context: aiForm.max_context,
    }
    await updateAIConfig(submitData)
    message.success(t('user.aiConfigSaved'))
  } catch {
    // handled
  } finally {
    savingAI.value = false
  }
}

onMounted(async () => {
  await fetchAvailableModels()
  fetchUserInfo()
  fetchAIConfig()
  fetchSubscription()
  fetchCategories()
})
</script>

<style scoped lang="less">
.profile-page {
  max-width: 700px;
  margin: 0 auto;
}

.profile-card {
  border-radius: 12px;
}

.avatar-upload-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.avatar-upload {
  cursor: pointer;
  transition: opacity 0.2s;
  &:hover {
    opacity: 0.8;
  }
}

.avatar-hint {
  font-size: 12px;
  color: var(--text-secondary);
}

.model-select-wrapper {
  width: 100%;
}

.slider-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;

  .ant-slider {
    flex: 1;
  }
}

.slider-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--primary-color);
  min-width: 32px;
  text-align: center;
}

.category-checkbox-group {
  width: 100%;
}
</style>
