<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h2 class="login-title">{{ $t('auth.loginTitle') }}</h2>
        <p class="login-subtitle">{{ appStore.siteConfig.siteName }}</p>
      </div>
      <div class="login-type-switch">
        <a-radio-group v-model:value="loginType" button-style="solid" size="small">
          <a-radio-button value="user">{{ $t('auth.userLogin') }}</a-radio-button>
          <a-radio-button value="admin">{{ $t('auth.adminLogin') }}</a-radio-button>
        </a-radio-group>
      </div>
      <a-form
        :model="formState"
        :rules="rules"
        @finish="handleLogin"
        layout="vertical"
        size="large"
      >
        <a-form-item name="account">
          <a-input v-model:value="formState.account" :placeholder="$t('auth.accountPlaceholder')">
            <template #prefix><UserOutlined class="input-icon" /></template>
          </a-input>
        </a-form-item>
        <a-form-item name="password">
          <a-input-password v-model:value="formState.password" :placeholder="$t('auth.password')">
            <template #prefix><LockOutlined class="input-icon" /></template>
          </a-input-password>
        </a-form-item>
        <a-form-item v-if="loginType === 'user'">
          <div class="login-options">
            <a-checkbox v-model:checked="formState.remember">{{ $t('auth.autoLogin') }}</a-checkbox>
            <a class="forgot-link">{{ $t('auth.forgotPassword') }}</a>
          </div>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" block :loading="loading">
            {{ $t('nav.login') }}
          </a-button>
        </a-form-item>
        <div v-if="loginType === 'user'" class="login-footer">
          {{ $t('auth.noAccount') }}
          <router-link to="/register">{{ $t('nav.register') }}</router-link>
        </div>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { login } from '@/api/auth'
import { adminLogin } from '@/api/admin'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const userStore = useUserStore()
const appStore = useAppStore()
const loading = ref(false)
const loginType = ref<'user' | 'admin'>('user')

const formState = reactive({
  account: '',
  password: '',
  remember: true,
})

const rules = {
  account: [{ required: true, message: () => t('auth.usernameRequired') }],
  password: [{ required: true, message: () => t('auth.passwordRequired') }],
}

async function handleLogin() {
  loading.value = true
  try {
    if (loginType.value === 'admin') {
      const res = await adminLogin({ account: formState.account, password: formState.password })
      const data = res.data
      userStore.setLogin({
        token: data.token,
        username: data.user.username,
        role: data.user.role,
        userId: data.user.id,
      })
      message.success(t('auth.loginSuccess'))
      router.push('/admin')
    } else {
      const res = await login({ account: formState.account, password: formState.password })
      const data = res.data
      userStore.setLogin({
        token: data.token,
        username: data.user.username,
        role: data.user.role,
        userId: data.user.id,
      })
      message.success(t('auth.loginSuccess'))
      const redirect = (route.query.redirect as string) || '/'
      router.push(redirect)
    }
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="less">
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  padding: 24px;
}

.login-card {
  width: 420px;
  padding: 40px;
  background: var(--card-bg);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

.login-header {
  text-align: center;
  margin-bottom: 16px;
}

.login-type-switch {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
}

.login-title {
  font-size: 26px;
  font-weight: 700;
  margin-bottom: 4px;
}

.login-subtitle {
  color: var(--text-secondary);
  font-size: 14px;
}

.input-icon {
  color: var(--text-secondary);
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.forgot-link {
  font-size: 14px;
  color: var(--primary-color);
}

.login-footer {
  text-align: center;
  color: var(--text-secondary);
  a {
    margin-left: 4px;
    color: var(--primary-color);
  }
}

@media (max-width: 480px) {
  .login-card {
    width: 100%;
    padding: 24px;
  }
}
</style>
