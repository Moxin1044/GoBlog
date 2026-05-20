<template>
  <div class="register-page">
    <div class="register-card">
      <div class="register-header">
        <h2 class="register-title">{{ $t('auth.registerTitle') }}</h2>
        <p class="register-subtitle">{{ appStore.siteConfig.siteName }}</p>
      </div>
      <a-form
        :model="formState"
        :rules="rules"
        @finish="handleRegister"
        layout="vertical"
        size="large"
      >
        <a-form-item :label="$t('auth.username')" name="username">
          <a-input v-model:value="formState.username" :placeholder="$t('auth.username')">
            <template #prefix><UserOutlined class="input-icon" /></template>
          </a-input>
        </a-form-item>
        <a-form-item :label="$t('auth.email')" name="email">
          <a-input v-model:value="formState.email" :placeholder="$t('auth.email')">
            <template #prefix><MailOutlined class="input-icon" /></template>
          </a-input>
        </a-form-item>
        <a-form-item :label="$t('auth.phone')" name="phone">
          <a-input v-model:value="formState.phone" :placeholder="$t('auth.phone')">
            <template #prefix><PhoneOutlined class="input-icon" /></template>
          </a-input>
        </a-form-item>
        <a-form-item :label="$t('auth.verifyCode')" name="code">
          <div class="code-input">
            <a-input v-model:value="formState.code" :placeholder="$t('auth.verifyCode')" />
            <a-button :disabled="countdown > 0" @click="handleSendCode">
              {{ countdown > 0 ? `${countdown}s` : $t('auth.sendCode') }}
            </a-button>
          </div>
        </a-form-item>
        <a-form-item :label="$t('auth.password')" name="password">
          <a-input-password v-model:value="formState.password" :placeholder="$t('auth.password')">
            <template #prefix><LockOutlined class="input-icon" /></template>
          </a-input-password>
        </a-form-item>
        <a-form-item :label="$t('auth.confirmPassword')" name="confirmPassword">
          <a-input-password v-model:value="formState.confirmPassword" :placeholder="$t('auth.confirmPassword')">
            <template #prefix><LockOutlined class="input-icon" /></template>
          </a-input-password>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" block :loading="loading">
            {{ $t('nav.register') }}
          </a-button>
        </a-form-item>
        <div class="register-footer">
          {{ $t('auth.hasAccount') }}
          <router-link to="/login">{{ $t('nav.login') }}</router-link>
        </div>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { register, sendVerifyCode } from '@/api/auth'
import { useAppStore } from '@/stores/app'

const router = useRouter()
const { t } = useI18n()
const appStore = useAppStore()
const loading = ref(false)
const countdown = ref(0)

const formState = reactive({
  username: '',
  email: '',
  phone: '',
  code: '',
  password: '',
  confirmPassword: '',
})

const validateConfirm = async (_rule: any, value: string) => {
  if (value && value !== formState.password) {
    throw new Error(t('auth.confirmPasswordMismatch'))
  }
}

const rules = {
  username: [{ required: true, message: () => t('auth.usernameRequired') }],
  email: [
    { required: true, message: () => t('auth.emailRequired') },
    { type: 'email' as const, message: () => t('auth.emailInvalid') },
  ],
  phone: [
    { required: true, message: () => t('auth.phoneRequired') },
    { pattern: /^1[3-9]\d{9}$/, message: () => t('auth.phoneInvalid') },
  ],
  code: [{ required: true, message: () => t('auth.codeRequired') }],
  password: [
    { required: true, message: () => t('auth.passwordRequired') },
    { min: 6, message: () => t('auth.passwordMin') },
  ],
  confirmPassword: [
    { required: true, message: () => t('auth.confirmPasswordRequired') },
    { validator: validateConfirm },
  ],
}

async function handleSendCode() {
  if (!formState.email) {
    message.warning(t('auth.inputEmailFirst'))
    return
  }
  try {
    await sendVerifyCode(formState.email)
    message.success(t('auth.codeSent'))
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) clearInterval(timer)
    }, 1000)
  } catch {
    message.error(t('auth.codeSendFailed'))
  }
}

async function handleRegister() {
  loading.value = true
  try {
    await register(formState)
    message.success(t('auth.registerSuccess'))
    router.push('/login')
  } catch {
    // handled
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="less">
.register-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  padding: 24px;
}

.register-card {
  width: 420px;
  padding: 40px;
  background: var(--card-bg);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-title {
  font-size: 26px;
  font-weight: 700;
  margin-bottom: 4px;
}

.register-subtitle {
  color: var(--text-secondary);
  font-size: 14px;
}

.input-icon {
  color: var(--text-secondary);
}

.code-input {
  display: flex;
  gap: 8px;

  .ant-input {
    flex: 1;
  }
}

.register-footer {
  text-align: center;
  color: var(--text-secondary);
  a {
    margin-left: 4px;
    color: var(--primary-color);
  }
}

@media (max-width: 480px) {
  .register-card {
    width: 100%;
    padding: 24px;
  }
}
</style>
