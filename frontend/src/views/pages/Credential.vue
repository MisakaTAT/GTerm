<template>
  <div class="credential-container">
    <div class="header">
      <div class="title">
        <span>{{ $t('frontend.credential.title') }}</span>
      </div>
      <div class="actions">
        <NInputGroup>
          <NInput :placeholder="$t('frontend.credential.search')" :allow-input="value => !/\s/.test(value)">
            <template #prefix>
              <Icon icon="ph:magnifying-glass" />
            </template>
          </NInput>
          <NButton type="primary" ghost @click="handleAdd">
            <template #icon>
              <Icon icon="ph:plus-bold" />
            </template>
            {{ $t('frontend.credential.add') }}
          </NButton>
        </NInputGroup>
      </div>
    </div>

    <NScrollbar class="content">
      <div class="content-wrapper">
        <NList v-if="creds?.length" hoverable clickable class="credential-list">
          <NListItem v-for="v in creds" :key="v.id">
            <div class="credential-item">
              <NThing>
                <template #avatar>
                  <div class="credential-type" :class="v.authMethod === AuthMethod.Password ? 'success' : 'warning'">
                    <Icon :icon="v.authMethod === AuthMethod.Password ? 'ph:password' : 'ph:key'" />
                  </div>
                </template>
                <template #header>
                  <div class="credential-header">
                    <span class="name">{{ v.label }}</span>
                  </div>
                </template>
                <template #description>
                  <div class="credential-info">
                    <span class="info-item">
                      <Icon icon="ph:user" />
                      {{ v.username }}
                    </span>
                    <span class="info-item">
                      <Icon icon="ph:clock" />
                      {{ formatTime(v.createdAt) }}
                    </span>
                  </div>
                </template>
              </NThing>
              <div class="credential-actions">
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <NButton circle text @click="handleCopy(v)">
                      <template #icon>
                        <Icon :icon="v.authMethod === AuthMethod.Password ? 'ph:copy' : 'ph:file-text'" />
                      </template>
                    </NButton>
                  </template>
                  {{
                    v.authMethod === AuthMethod.Password
                      ? $t('frontend.credential.actions.copyPassword')
                      : $t('frontend.credential.actions.viewKey')
                  }}
                </n-tooltip>
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <NButton circle text @click="handleEdit(v)">
                      <template #icon>
                        <Icon icon="ph:pencil-simple" />
                      </template>
                    </NButton>
                  </template>
                  {{ $t('frontend.credential.actions.edit') }}
                </n-tooltip>
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <NButton circle text type="error" @click="handleDelete(v)">
                      <template #icon>
                        <Icon icon="ph:trash" />
                      </template>
                    </NButton>
                  </template>
                  {{ $t('frontend.credential.actions.delete') }}
                </n-tooltip>
              </div>
            </div>
          </NListItem>
        </NList>
        <NEmpty v-else class="credential-empty" :description="$t('frontend.credential.empty')" />
      </div>
    </NScrollbar>

    <CredentialModal
      v-model:show="showModal"
      :is-edit="isEdit"
      :credential-id="credentialId"
      @success="handleSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import type { Credential } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/dal/model/models';
import { AuthMethod } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/enums/models';
import { DeleteCredential, ListCredential } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/credentialsrv';
import dayjs from 'dayjs';
import {
  NButton,
  NEmpty,
  NInput,
  NInputGroup,
  NList,
  NListItem,
  NScrollbar,
  NThing,
  useMessage,
  useThemeVars,
} from 'naive-ui';
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useCall } from '@/utils/call';
import CredentialModal from '@/views/modals/CredentialModal.vue';

const { t } = useI18n();
const message = useMessage();
const { call } = useCall();
const showModal = ref(false);
const isEdit = ref(false);
const credentialId = ref<number>(0);

const handleCopy = async (credential: Credential) => {
  if (credential.authMethod === AuthMethod.Password) {
    try {
      await navigator.clipboard.writeText(credential.password);
      message.success(t('frontend.credential.messages.passwordCopied'));
    } catch {
      message.error(t('frontend.credential.messages.copyFailed'));
    }
  }
};

const handleEdit = (credential: Credential) => {
  isEdit.value = true;
  credentialId.value = credential.id;
  showModal.value = true;
};

const creds = ref<Credential[]>();

const fetchCredentials = async () => {
  const result = await call(ListCredential);
  if (result.ok) {
    creds.value = result.data;
  }
  return result.data;
};

const handleDelete = async (credential: Credential) => {
  const result = await call(DeleteCredential, {
    args: [credential.id],
  });

  if (result.ok) {
    await fetchCredentials();
  }
};

const handleAdd = () => {
  isEdit.value = false;
  credentialId.value = 0;
  showModal.value = true;
};

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm');
};

const handleSuccess = () => {
  fetchCredentials();
};

onMounted(async () => {
  const [credsData] = await Promise.all([fetchCredentials()]);
  creds.value = credsData;
});

const themeVars = useThemeVars();
</script>

<style lang="less" scoped>
.credential-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
}

.header {
  position: sticky;
  top: 0;
  padding: 16px 24px;
  border-bottom: 1px solid v-bind('themeVars.borderColor');
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 1;

  .title {
    width: 100%;
    max-width: 800px;
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 16px;
    font-size: 18px;
    font-weight: 500;
    color: v-bind('themeVars.textColorBase');
  }

  .actions {
    width: 100%;
    max-width: 800px;

    :deep(.n-input-group) {
      width: 100%;
      display: flex;
      gap: 8px;

      .n-input {
        flex: 1;
      }
    }
  }
}

.content {
  flex: 1;
  min-height: 0;
}

.content-wrapper {
  padding: 16px;

  .credential-list {
    margin: 0 auto;
    width: 100%;
    max-width: 800px;
  }
}

.credential-type {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  transition: all 0.2s ease;

  &.success {
    background-color: rgba(63, 185, 80, 0.1);
    color: #3fb950;
  }

  &.warning {
    background-color: rgba(210, 153, 34, 0.1);
    color: #d29922;
  }
}

.credential-header {
  display: flex;
  align-items: center;
  gap: 8px;

  .name {
    font-weight: 500;
    color: v-bind('themeVars.textColorBase');
  }
}

.credential-info {
  display: flex;
  gap: 16px;
  margin-top: 4px;
  font-size: 13px;
  color: v-bind('themeVars.textColor3');

  .info-item {
    display: flex;
    align-items: center;
    gap: 4px;

    .iconify {
      font-size: 16px;
    }
  }
}

.credential-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 16px;

  :deep(.n-thing) {
    flex: 1;
    min-width: 0;
  }
}

.credential-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
  margin-right: -8px;

  :deep(.n-button) {
    width: 32px;
    height: 32px;

    .n-button__icon {
      font-size: 16px;
    }

    &:hover {
      background-color: v-bind('themeVars.hoverColor');
    }
  }
}

.credential-empty {
  margin: 120px auto;
  width: 100%;
  max-width: 800px;
}
</style>
