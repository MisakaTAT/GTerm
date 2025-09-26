<template>
  <NModal
    v-model:show="visible"
    close-on-esc
    :negative-text="$t('frontend.connModal.cancel')"
    :on-close="resetForm"
    :positive-text="$t('frontend.connModal.confirm')"
    :show-icon="false"
    :title="isEdit ? $t('frontend.connModal.editTitle') : $t('frontend.connModal.title')"
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
    @positive-click="handleConfirm"
  >
    <NTabs v-model:value="activeTab" animated type="line" placement="left">
      <NTabPane name="basic" :tab="$t('frontend.connModal.basicConfig')">
        <NScrollbar style="max-height: 70vh; padding-right: 12px">
          <NForm ref="formRef" :model="formValue" :rules="rules">
            <NFormItem path="label" :label="$t('frontend.connModal.label')">
              <NInput
                v-model:value="formValue.label"
                clearable
                :placeholder="$t('frontend.connModal.placeholder.label')"
                :allow-input="value => !/\s/.test(value)"
              />
            </NFormItem>

            <NFormItem path="groupID" :label="$t('frontend.connModal.group')">
              <NSelect
                v-model:value="formValue.groupID"
                :options="groupOptions"
                :placeholder="$t('frontend.connModal.placeholder.group')"
              />
            </NFormItem>

            <NFormItem path="connProtocol" :label="$t('frontend.connModal.connProtocol')">
              <NSelect
                v-model:value="formValue.connProtocol"
                :options="connProtocolOptions"
                :placeholder="$t('frontend.connModal.placeholder.connProtocol')"
              />
            </NFormItem>

            <template v-if="formValue.connProtocol === ConnProtocol.Serial">
              <NFormItem path="serialPort" :label="$t('frontend.connModal.serialPort')">
                <NSelect
                  v-model:value="formValue.serialPort"
                  :options="serialPortsOptions"
                  :placeholder="$t('frontend.connModal.placeholder.serialPort')"
                />
              </NFormItem>

              <NFormItem path="baudRate" :label="$t('frontend.connModal.baudRate')">
                <NSelect
                  v-model:value="formValue.baudRate"
                  :options="baudRateOptions"
                  :placeholder="$t('frontend.connModal.placeholder.baudRate')"
                />
              </NFormItem>

              <div class="form-row">
                <NFormItem path="dataBits" :label="$t('frontend.connModal.dataBits')" class="form-item">
                  <NSelect
                    v-model:value="formValue.dataBits"
                    :options="dataBitsOptions"
                    :placeholder="$t('frontend.connModal.placeholder.dataBits')"
                  />
                </NFormItem>

                <NFormItem path="stopBits" :label="$t('frontend.connModal.stopBits')" class="form-item">
                  <NSelect
                    v-model:value="formValue.stopBits"
                    :options="stopBitsOptions"
                    :placeholder="$t('frontend.connModal.placeholder.stopBits')"
                  />
                </NFormItem>

                <NFormItem path="parity" :label="$t('frontend.connModal.parity')" class="form-item">
                  <NSelect
                    v-model:value="formValue.parity"
                    :options="parityOptions"
                    :placeholder="$t('frontend.connModal.placeholder.parity')"
                  />
                </NFormItem>
              </div>
            </template>

            <template v-if="formValue.connProtocol === ConnProtocol.SSH">
              <div class="form-row">
                <NFormItem path="host" :label="$t('frontend.connModal.host')" class="form-item">
                  <NInput
                    v-model:value="formValue.host"
                    clearable
                    :placeholder="$t('frontend.connModal.placeholder.host')"
                    :allow-input="value => !/\s/.test(value)"
                  />
                </NFormItem>
                <NFormItem path="port" :label="$t('frontend.connModal.port')" class="port-input">
                  <NInputNumber
                    v-model:value="formValue.port"
                    :min="1"
                    :max="65535"
                    :show-button="false"
                    :placeholder="$t('frontend.connModal.placeholder.port')"
                  />
                </NFormItem>
              </div>

              <NFormItem :label="$t('frontend.connModal.authType')" path="credential.authMethod">
                <div class="auth-type-container">
                  <NButtonGroup>
                    <NButton
                      :type="
                        !formValue.useCommonCredential && formValue.credential?.authMethod === AuthMethod.Password
                          ? 'primary'
                          : 'default'
                      "
                      @click="handleCredentialTypeChange(CredentialType.Password)"
                    >
                      <template #icon>
                        <Icon icon="ph:password" />
                      </template>
                      {{ $t('frontend.connModal.password') }}
                    </NButton>
                    <NButton
                      :type="
                        !formValue.useCommonCredential && formValue.credential?.authMethod === AuthMethod.PrivateKey
                          ? 'primary'
                          : 'default'
                      "
                      @click="handleCredentialTypeChange(CredentialType.PrivateKey)"
                    >
                      <template #icon>
                        <Icon icon="ph:key" />
                      </template>
                      {{ $t('frontend.connModal.privateKey') }}
                    </NButton>
                    <NButton
                      :type="formValue.useCommonCredential ? 'primary' : 'default'"
                      @click="handleCredentialTypeChange(CredentialType.Common)"
                    >
                      <template #icon>
                        <Icon icon="ph:vault" />
                      </template>
                      {{ $t('frontend.connModal.commonCredentialLib') }}
                    </NButton>
                  </NButtonGroup>
                </div>
              </NFormItem>

              <template v-if="formValue.useCommonCredential">
                <NFormItem path="credentialID" :label="$t('frontend.connModal.credential')">
                  <NSelect
                    v-model:value="formValue.credentialID"
                    :options="credentialOptions"
                    clearable
                    :placeholder="$t('frontend.connModal.placeholder.selectCredential')"
                    @update:value="handleSelectCredential"
                  />
                </NFormItem>
              </template>

              <template v-else>
                <NFormItem path="credential.username" :label="$t('frontend.connModal.username')">
                  <NInput
                    v-model:value="formValue.credential!.username"
                    clearable
                    :placeholder="$t('frontend.connModal.placeholder.username')"
                    :allow-input="value => !/\s/.test(value)"
                  />
                </NFormItem>

                <template v-if="formValue.credential!.authMethod === AuthMethod.Password">
                  <NFormItem path="credential.password" :label="$t('frontend.connModal.password')">
                    <NInput
                      v-model:value="formValue.credential!.password"
                      type="password"
                      show-password-on="click"
                      clearable
                      :placeholder="$t('frontend.connModal.placeholder.password')"
                      :allow-input="value => !/\s/.test(value)"
                    />
                  </NFormItem>
                </template>

                <template v-if="formValue.credential!.authMethod === AuthMethod.PrivateKey">
                  <NFormItem path="credential.privateKey" :label="$t('frontend.connModal.privateKey')">
                    <NInput
                      v-model:value="formValue.credential!.privateKey"
                      type="textarea"
                      :autosize="{ minRows: 3, maxRows: 3 }"
                      clearable
                      :placeholder="$t('frontend.connModal.placeholder.privateKey')"
                      :allow-input="value => !/\s/.test(value)"
                    />
                  </NFormItem>
                  <NFormItem path="credential.passphrase" :label="$t('frontend.connModal.passphrase')">
                    <NInput
                      v-model:value="formValue.credential!.passphrase"
                      type="password"
                      show-password-on="click"
                      clearable
                      :placeholder="$t('frontend.connModal.placeholder.passphrase')"
                      :allow-input="value => !/\s/.test(value)"
                    />
                  </NFormItem>
                </template>
              </template>
            </template>
          </NForm>
        </NScrollbar>
      </NTabPane>

      <NTabPane name="personal" :tab="$t('frontend.connModal.personalConfig')">
        <NForm ref="formRef" :model="formValue" :rules="rules">
          <NFormItem path="theme" :label="$t('frontend.connModal.theme')">
            <div class="theme-section">
              <div class="theme-preview">
                <ThemePreview :theme="selectedTheme" />
              </div>
              <NSelect
                v-model:value="formValue.theme"
                :options="themeOptions"
                :placeholder="t('frontend.connModal.placeholder.selectTheme')"
                filterable
                @update:value="handleThemeChange"
              />
            </div>
          </NFormItem>
        </NForm>
      </NTabPane>

      <NTabPane name="advanced" :tab="$t('frontend.connModal.advancedConfig')">
        <NScrollbar style="max-height: 70vh; padding-right: 12px">
          <template v-if="formValue.connProtocol === ConnProtocol.SSH">
            <NForm ref="formRef" :model="formValue" :rules="rules">
              <NFormItem path="sshCharset" :label="$t('frontend.connModal.charset')">
                <NSelect
                  v-model:value="formValue.sshCharset"
                  :options="charsetOptions"
                  filterable
                  :placeholder="$t('frontend.connModal.placeholder.charset')"
                />
              </NFormItem>

              <NCollapse accordion>
                <NCollapseItem name="ciphers">
                  <template #header>
                    <div style="display: flex; justify-content: space-between; width: 100%; align-items: center">
                      <span>{{ $t('frontend.connModal.ciphers') }}</span>
                      <NButton text size="tiny" @click.stop="handleResetAlgorithms('ciphers')">
                        {{ $t('frontend.connModal.resetToDefault') }}
                      </NButton>
                    </div>
                  </template>
                  <div class="algorithm-container">
                    <NList bordered size="small" class="algorithm-list">
                      <NListItem v-for="(cipher, index) in formValue.sshCiphers" :key="index" class="compact-item">
                        <div class="algorithm-item">
                          <NText class="algorithm-name">{{ cipher }}</NText>
                          <div class="algorithm-actions">
                            <NButton
                              quaternary
                              circle
                              size="tiny"
                              :disabled="index === 0"
                              @click="handleMoveAlgorithm('ciphers', index, -1)"
                            >
                              <template #icon>
                                <Icon icon="ph:arrow-up" />
                              </template>
                            </NButton>
                            <NButton
                              quaternary
                              circle
                              size="tiny"
                              :disabled="index === formValue.sshCiphers.length - 1"
                              @click="handleMoveAlgorithm('ciphers', index, 1)"
                            >
                              <template #icon>
                                <Icon icon="ph:arrow-down" />
                              </template>
                            </NButton>
                          </div>
                        </div>
                      </NListItem>
                    </NList>
                  </div>
                </NCollapseItem>

                <NCollapseItem name="keyExchanges">
                  <template #header>
                    <div style="display: flex; justify-content: space-between; width: 100%; align-items: center">
                      <span>{{ $t('frontend.connModal.keyExchanges') }}</span>
                      <NButton text size="tiny" @click.stop="handleResetAlgorithms('keyExchanges')">
                        {{ $t('frontend.connModal.resetToDefault') }}
                      </NButton>
                    </div>
                  </template>
                  <div class="algorithm-container">
                    <NList bordered size="small" class="algorithm-list">
                      <NListItem v-for="(kex, index) in formValue.sshKeyExchanges" :key="index" class="compact-item">
                        <div class="algorithm-item">
                          <NText class="algorithm-name">{{ kex }}</NText>
                          <div class="algorithm-actions">
                            <NButton
                              quaternary
                              circle
                              size="tiny"
                              :disabled="index === 0"
                              @click="handleMoveAlgorithm('keyExchanges', index, -1)"
                            >
                              <template #icon>
                                <Icon icon="ph:arrow-up" />
                              </template>
                            </NButton>
                            <NButton
                              quaternary
                              circle
                              size="tiny"
                              :disabled="index === formValue.sshKeyExchanges.length - 1"
                              @click="handleMoveAlgorithm('keyExchanges', index, 1)"
                            >
                              <template #icon>
                                <Icon icon="ph:arrow-down" />
                              </template>
                            </NButton>
                          </div>
                        </div>
                      </NListItem>
                    </NList>
                  </div>
                </NCollapseItem>

                <NCollapseItem name="macs">
                  <template #header>
                    <div style="display: flex; justify-content: space-between; width: 100%; align-items: center">
                      <span>{{ $t('frontend.connModal.macs') }}</span>
                      <NButton text size="tiny" @click.stop="handleResetAlgorithms('macs')">
                        {{ $t('frontend.connModal.resetToDefault') }}
                      </NButton>
                    </div>
                  </template>
                  <div class="algorithm-container">
                    <NList bordered size="small" class="algorithm-list">
                      <NListItem v-for="(mac, index) in formValue.sshMACs" :key="index" class="compact-item">
                        <div class="algorithm-item">
                          <NText class="algorithm-name">{{ mac }}</NText>
                          <div class="algorithm-actions">
                            <NButton
                              quaternary
                              circle
                              size="tiny"
                              :disabled="index === 0"
                              @click="handleMoveAlgorithm('macs', index, -1)"
                            >
                              <template #icon>
                                <Icon icon="ph:arrow-up" />
                              </template>
                            </NButton>
                            <NButton
                              quaternary
                              circle
                              size="tiny"
                              :disabled="index === formValue.sshMACs.length - 1"
                              @click="handleMoveAlgorithm('macs', index, 1)"
                            >
                              <template #icon>
                                <Icon icon="ph:arrow-down" />
                              </template>
                            </NButton>
                          </div>
                        </div>
                      </NListItem>
                    </NList>
                  </div>
                </NCollapseItem>

                <NCollapseItem name="publicKeyAlgorithms">
                  <template #header>
                    <div style="display: flex; justify-content: space-between; width: 100%; align-items: center">
                      <span>{{ $t('frontend.connModal.publicKeyAlgorithms') }}</span>
                      <NButton text size="tiny" @click.stop="handleResetAlgorithms('publicKeyAlgorithms')">
                        {{ $t('frontend.connModal.resetToDefault') }}
                      </NButton>
                    </div>
                  </template>
                  <div class="algorithm-container">
                    <NList bordered size="small" class="algorithm-list">
                      <NListItem
                        v-for="(algo, index) in formValue.sshPublicKeyAlgorithms"
                        :key="index"
                        class="compact-item"
                      >
                        <div class="algorithm-item">
                          <NText class="algorithm-name">{{ algo }}</NText>
                          <div class="algorithm-actions">
                            <NButton
                              quaternary
                              circle
                              size="tiny"
                              :disabled="index === 0"
                              @click="handleMoveAlgorithm('publicKeyAlgorithms', index, -1)"
                            >
                              <template #icon>
                                <Icon icon="ph:arrow-up" />
                              </template>
                            </NButton>
                            <NButton
                              quaternary
                              circle
                              size="tiny"
                              :disabled="index === formValue.sshPublicKeyAlgorithms.length - 1"
                              @click="handleMoveAlgorithm('publicKeyAlgorithms', index, 1)"
                            >
                              <template #icon>
                                <Icon icon="ph:arrow-down" />
                              </template>
                            </NButton>
                          </div>
                        </div>
                      </NListItem>
                    </NList>
                  </div>
                </NCollapseItem>

                <NCollapseItem name="hostKeyAlgorithms">
                  <template #header>
                    <div style="display: flex; justify-content: space-between; width: 100%; align-items: center">
                      <span>{{ $t('frontend.connModal.hostKeyAlgorithms') }}</span>
                      <NButton text size="tiny" @click.stop="handleResetAlgorithms('hostKeyAlgorithms')">
                        {{ $t('frontend.connModal.resetToDefault') }}
                      </NButton>
                    </div>
                  </template>
                  <div class="algorithm-container">
                    <NList bordered size="small" class="algorithm-list">
                      <NListItem
                        v-for="(algo, index) in formValue.sshHostKeyAlgorithms"
                        :key="index"
                        class="compact-item"
                      >
                        <div class="algorithm-item">
                          <NText class="algorithm-name">{{ algo }}</NText>
                          <div class="algorithm-actions">
                            <NButton
                              quaternary
                              circle
                              size="tiny"
                              :disabled="index === 0"
                              @click="handleMoveAlgorithm('hostKeyAlgorithms', index, -1)"
                            >
                              <template #icon>
                                <Icon icon="ph:arrow-up" />
                              </template>
                            </NButton>
                            <NButton
                              quaternary
                              circle
                              size="tiny"
                              :disabled="index === formValue.sshHostKeyAlgorithms.length - 1"
                              @click="handleMoveAlgorithm('hostKeyAlgorithms', index, 1)"
                            >
                              <template #icon>
                                <Icon icon="ph:arrow-down" />
                              </template>
                            </NButton>
                          </div>
                        </div>
                      </NListItem>
                    </NList>
                  </div>
                </NCollapseItem>
              </NCollapse>
            </NForm>
          </template>

          <template v-else-if="formValue.connProtocol === ConnProtocol.Serial">
            <NEmpty size="small" :description="$t('frontend.connModal.developing')">
              <template #icon>
                <Icon icon="ph:code" />
              </template>
            </NEmpty>
          </template>

          <template v-else>
            <NSpace justify="center" align="center" style="height: 200px">
              <NText>{{ $t('frontend.connModal.selectProtocolFirst') }}</NText>
            </NSpace>
          </template>
        </NScrollbar>
      </NTabPane>
    </NTabs>
  </NModal>
</template>

<script lang="ts" setup>
import { Icon } from '@iconify/vue';
import type { Connection, Credential, Group } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/dal/model';
import { AuthMethod, ConnProtocol } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/enums';
import {
  CreateConnection,
  FindConnectionByID,
  UpdateConnection,
} from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/connectionsrv';
import { ListCredential } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/credentialsrv';
import { ListGroup } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/groupsrv';
import { SerialPorts } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/terminalsrv';
import type { FormInst, FormRules } from 'naive-ui';
import {
  NButton,
  NButtonGroup,
  NCollapse,
  NCollapseItem,
  NEmpty,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NList,
  NListItem,
  NModal,
  NScrollbar,
  NSelect,
  NSpace,
  NTabPane,
  NTabs,
  NText,
} from 'naive-ui';
import type { SelectMixedOption } from 'naive-ui/es/select/src/interface';
import { computed, onMounted, onUpdated, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { defaultTheme, loadTheme, themeOptions as xtermThemeOptions } from '@/themes/xtermjs';
import { useCall } from '@/utils/call';
import ThemePreview from '../components/ThemePreview.vue';

const props = defineProps<{
  show: boolean;
  isEdit: boolean;
  connectionId?: number;
}>();

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'success'): void;
}>();

const { t } = useI18n();
const formRef = ref<FormInst | null>(null);
const activeTab = ref('basic');
const { call } = useCall();

const visible = computed({
  get: () => props.show,
  set: value => {
    emit('update:show', value);
  },
});

const CredentialType = {
  Password: 0,
  PrivateKey: 1,
  Common: 2,
} as const;

const connProtocolOptions = [
  { label: ConnProtocol.SSH, value: ConnProtocol.SSH },
  { label: ConnProtocol.Serial, value: ConnProtocol.Serial },
];

const defaultCredential: Partial<Credential> = {
  label: '',
  username: '',
  password: '',
  privateKey: '',
  passphrase: '',
  isCommonCredential: false,
  authMethod: AuthMethod.Password,
};

const defaultConnection: Partial<Connection> = {
  label: '',
  host: '',
  port: 22,
  serialPort: undefined,
  connProtocol: undefined,
  credentialID: undefined,
  useCommonCredential: false,
  groupID: undefined,
  baudRate: 9600,
  dataBits: 8,
  stopBits: 0, // OneStopBit
  parity: 0, // NoParity
  theme: 'Default',
  sshCharset: 'UTF-8',
  sshCiphers: [
    'aes256-gcm@openssh.com',
    'aes128-gcm@openssh.com',
    'chacha20-poly1305@openssh.com',
    'aes256-ctr',
    'aes192-ctr',
    'aes128-ctr',
    'aes256-cbc',
    'aes192-cbc',
    'aes128-cbc',
    'blowfish-cbc',
    '3des-cbc',
  ],
  sshKeyExchanges: [
    'curve25519-sha256',
    'curve25519-sha256@libssh.org',
    'ecdh-sha2-nistp256',
    'ecdh-sha2-nistp384',
    'ecdh-sha2-nistp521',
    'diffie-hellman-group18-sha512',
    'diffie-hellman-group16-sha512',
    'diffie-hellman-group-exchange-sha1',
    'diffie-hellman-group-exchange-sha256',
    'diffie-hellman-group1-sha1',
    'diffie-hellman-group14-sha1',
    'diffie-hellman-group14-sha256',
  ],
  sshMACs: [
    'hmac-sha2-256-etm@openssh.com',
    'hmac-sha2-512-etm@openssh.com',
    'hmac-sha1-etm@openssh.com',
    'hmac-sha1-96-etm@openssh.com',
    'hmac-md5-etm@openssh.com',
    'hmac-md5-96-etm@openssh.com',
    'hmac-sha2-256',
    'hmac-sha2-512',
    'hmac-sha1',
    'hmac-sha1-96',
    'hmac-md5',
    'hmac-md5-96',
    'aead-poly1305',
    'aead-gcm',
  ],
  sshPublicKeyAlgorithms: [
    'ssh-ed25519-cert-v01@openssh.com',
    'ecdsa-sha2-nistp512-cert-v01@openssh.com',
    'ecdsa-sha2-nistp384-cert-v01@openssh.com',
    'ecdsa-sha2-nistp256-cert-v01@openssh.com',
    'rsa-sha2-512-cert-v01@openssh.com',
    'rsa-sha2-256-cert-v01@openssh.com',
    'ssh-rsa-cert-v01@openssh.com',
    'ssh-ed25519',
    'ecdsa-sha2-nistp256',
    'rsa-sha2-512',
    'rsa-sha2-256',
    'ssh-rsa-cert-v01@openssh.com',
    'ssh-dss-cert-v01@openssh.com',
    'ssh-ed25519',
    'ecdsa-sha2-nistp512',
    'ecdsa-sha2-nistp384',
    'ecdsa-sha2-nistp256',
    'rsa-sha2-512',
    'rsa-sha2-256',
    'ssh-rsa',
    'ssh-dss',
  ],
  sshHostKeyAlgorithms: [
    'ssh-ed25519',
    'ecdsa-sha2-nistp256',
    'ecdsa-sha2-nistp384',
    'ecdsa-sha2-nistp521',
    'rsa-sha2-512',
    'rsa-sha2-256',
    'ssh-rsa',
    'ssh-dss',
    'ssh-dss-cert-v01@openssh.com',
  ],
};

const tempCachedCredentials = ref<{
  password: Credential;
  privateKey: Credential;
}>({
  password: createCredentialObject(AuthMethod.Password),
  privateKey: createCredentialObject(AuthMethod.PrivateKey),
});

const groupOptions = ref<SelectMixedOption[]>([]);
const serialPortsOptions = ref<SelectMixedOption[]>([]);
const credentialOptions = ref<SelectMixedOption[]>([]);

const formValue = ref<Connection>(createConnectionObject());

function createCredentialObject(authMethod: AuthMethod): Credential {
  return {
    ...defaultCredential,
    authMethod,
  } as Credential;
}

function createConnectionObject(isCommon = false): Connection {
  const conn = {
    ...defaultConnection,
    useCommonCredential: isCommon,
  } as Connection;

  if (!isCommon) {
    conn.credential = createCredentialObject(AuthMethod.Password);
  }
  return conn;
}

function prepareConnectionForEdit(connection: Connection): Connection {
  const conn = { ...connection } as Connection;

  if (conn.credential && !conn.useCommonCredential) {
    if (conn.credential.authMethod === AuthMethod.Password) {
      tempCachedCredentials.value.password = { ...conn.credential } as Credential;
    } else if (conn.credential.authMethod === AuthMethod.PrivateKey) {
      tempCachedCredentials.value.privateKey = { ...conn.credential } as Credential;
    }
  }

  return conn;
}

const baudRateOptions = [
  { label: '9600', value: 9600 },
  { label: '19200', value: 19200 },
  { label: '38400', value: 38400 },
  { label: '57600', value: 57600 },
  { label: '115200', value: 115200 },
];

const parityOptions = [
  { label: 'None', value: 0 }, // NoParity
  { label: 'Odd', value: 1 }, // OddParity
  { label: 'Even', value: 2 }, // EvenParity
  { label: 'Mark', value: 3 }, // MarkParity
  { label: 'Space', value: 4 }, // SpaceParity
];

const dataBitsOptions = [
  { label: '5', value: 5 },
  { label: '6', value: 6 },
  { label: '7', value: 7 },
  { label: '8', value: 8 },
];

const stopBitsOptions = [
  { label: '1', value: 0 }, // OneStopBit
  { label: '1.5', value: 1 }, // OnePointFiveStopBits
  { label: '2', value: 2 }, // TwoStopBits
];

const rules = computed<FormRules>(() => ({
  label: {
    required: true,
    message: t('frontend.connModal.validation.labelRequired'),
    trigger: 'blur',
  },
  connProtocol: {
    required: true,
    message: t('frontend.connModal.validation.connProtocolRequired'),
    trigger: 'blur',
  },
  host: {
    required: formValue.value.connProtocol === ConnProtocol.SSH,
    message: t('frontend.connModal.validation.hostRequired'),
    trigger: 'blur',
  },
  port: {
    required: formValue.value.connProtocol === ConnProtocol.SSH,
    type: 'number',
    message: t('frontend.connModal.validation.portRequired'),
    trigger: ['blur', 'change'],
    validator: (rule, value) => {
      if (typeof value !== 'number' || value < 1 || value > 65535) {
        return new Error(t('frontend.connModal.validation.portInvalid'));
      }
    },
  },
  serialPort: {
    required: formValue.value.connProtocol === ConnProtocol.Serial,
    message: t('frontend.connModal.validation.serialPortRequired'),
    trigger: ['blur', 'change'],
  },
  baudRate: {
    required: formValue.value.connProtocol === ConnProtocol.Serial,
    type: 'number',
    message: t('frontend.connModal.validation.baudRateRequired'),
    trigger: ['blur', 'change'],
  },
  dataBits: {
    required: formValue.value.connProtocol === ConnProtocol.Serial,
    type: 'number',
    message: t('frontend.connModal.validation.dataBitsRequired'),
    trigger: ['blur', 'change'],
  },
  stopBits: {
    required: formValue.value.connProtocol === ConnProtocol.Serial,
    type: 'number',
    message: t('frontend.connModal.validation.stopBitsRequired'),
    trigger: ['blur', 'change'],
  },
  parity: {
    required: formValue.value.connProtocol === ConnProtocol.Serial,
    type: 'number',
    message: t('frontend.connModal.validation.parityRequired'),
    trigger: ['blur', 'change'],
  },
  'credential.username': {
    required: !formValue.value.useCommonCredential,
    message: t('frontend.connModal.validation.usernameRequired'),
    trigger: 'blur',
  },
  'credential.password': {
    required: !formValue.value.useCommonCredential && formValue.value.credential?.authMethod === AuthMethod.Password,
    message: t('frontend.connModal.validation.passwordRequired'),
    trigger: 'blur',
  },
  'credential.privateKey': {
    required: !formValue.value.useCommonCredential && formValue.value.credential?.authMethod === AuthMethod.PrivateKey,
    message: t('frontend.connModal.validation.privateKeyRequired'),
    trigger: 'blur',
  },
  credentialID: {
    required: formValue.value.useCommonCredential,
    message: t('frontend.connModal.validation.credentialRequired'),
    trigger: ['blur', 'change'],
    validator: (rule, value) => {
      if (formValue.value.useCommonCredential && !value) {
        return new Error(t('frontend.connModal.validation.credentialRequired'));
      }
      return true;
    },
  },
}));

const handleCredentialTypeChange = (credentialType: number) => {
  if (credentialType === CredentialType.Common) {
    formValue.value.useCommonCredential = true;
    formValue.value.credentialID = null;
    formValue.value.credential = null;
    return;
  }
  formValue.value.useCommonCredential = false;
  formValue.value.credentialID = null;
  formValue.value.credential =
    credentialType === CredentialType.Password
      ? tempCachedCredentials.value.password
      : tempCachedCredentials.value.privateKey;
  formValue.value.credential.authMethod =
    credentialType === CredentialType.Password ? AuthMethod.Password : AuthMethod.PrivateKey;
};

const handleSelectCredential = async (id: number) => {
  if (!id) {
    formValue.value.credentialID = null;
    return;
  }
  formValue.value.credentialID = id;
  formValue.value.useCommonCredential = true;
};

const fetchGroups = async () => {
  const result = await call(ListGroup);
  return result.data || [];
};

const fetchSerialPorts = async () => {
  const result = await call(SerialPorts);
  return result.data || [];
};

const fetchCredentials = async () => {
  const result = await call(ListCredential);
  return result.data || [];
};

const initOptions = async () => {
  const [groups, credentials, serialPorts] = await Promise.all([fetchGroups(), fetchCredentials(), fetchSerialPorts()]);
  groupOptions.value = groups.map((group: Group) => ({
    label: group.name,
    value: group.id,
  }));
  credentialOptions.value = credentials.map((credential: Credential) => ({
    label: credential.label,
    value: credential.id,
  }));
  serialPortsOptions.value = serialPorts.map((serialPort: string) => ({
    label: serialPort,
    value: serialPort,
  }));
};

const initModalData = async () => {
  await initOptions();

  if (props.connectionId && props.connectionId > 0 && props.isEdit) {
    const result = await call(FindConnectionByID, {
      args: [props.connectionId],
    });

    if (result.ok) {
      formValue.value = prepareConnectionForEdit(result.data);
    } else {
      emit('update:show', false);
    }
  } else {
    formValue.value = createConnectionObject();
  }
};

onUpdated(() => {
  if (props.show) {
    initModalData();
  }
});

onMounted(() => {
  if (props.show) {
    initModalData();
  }
});

const handleConfirm = async () => {
  try {
    await formRef.value?.validate();

    const backendFunc = props.isEdit ? UpdateConnection : CreateConnection;
    const result = await call(backendFunc, {
      args: [formValue.value],
    });

    if (result.ok) {
      emit('update:show', false);
      emit('success');
    }

    return result.ok;
  } catch {
    return false;
  }
};

const resetForm = () => {
  formValue.value = createConnectionObject();
  tempCachedCredentials.value = {
    password: createCredentialObject(AuthMethod.Password),
    privateKey: createCredentialObject(AuthMethod.PrivateKey),
  };
  activeTab.value = 'basic';
  emit('update:show', false);
};

const themeOptions = xtermThemeOptions;
const selectedTheme = ref(defaultTheme);

onMounted(async () => {
  selectedTheme.value = await loadTheme(formValue.value.theme);
});

const handleThemeChange = async (newTheme: string) => {
  selectedTheme.value = await loadTheme(newTheme);
};

// 默认算法列表
const defaultCiphers = [...(defaultConnection.sshCiphers || [])];
const defaultKeyExchanges = [...(defaultConnection.sshKeyExchanges || [])];
const defaultMACs = [...(defaultConnection.sshMACs || [])];
const defaultPublicKeyAlgorithms = [...(defaultConnection.sshPublicKeyAlgorithms || [])];
const defaultHostKeyAlgorithms = [...(defaultConnection.sshHostKeyAlgorithms || [])];

const charsetOptions = [
  { label: 'UTF-8', value: 'UTF-8' },
  { label: 'GBK', value: 'GBK' },
  { label: 'GB18030', value: 'GB18030' },
  { label: 'ISO-8859-1', value: 'ISO-8859-1' },
  { label: 'ISO-8859-15', value: 'ISO-8859-15' },
  { label: 'Windows-1252', value: 'Windows-1252' },
];

type AlgorithmType = 'ciphers' | 'keyExchanges' | 'macs' | 'publicKeyAlgorithms' | 'hostKeyAlgorithms';

const handleMoveAlgorithm = (type: AlgorithmType, index: number, direction: -1 | 1) => {
  const newIndex = index + direction;

  if (type === 'ciphers' && formValue.value.sshCiphers) {
    if (newIndex < 0 || newIndex >= formValue.value.sshCiphers.length) return;
    const temp = formValue.value.sshCiphers[index];
    formValue.value.sshCiphers[index] = formValue.value.sshCiphers[newIndex];
    formValue.value.sshCiphers[newIndex] = temp;
  } else if (type === 'keyExchanges' && formValue.value.sshKeyExchanges) {
    if (newIndex < 0 || newIndex >= formValue.value.sshKeyExchanges.length) return;
    const temp = formValue.value.sshKeyExchanges[index];
    formValue.value.sshKeyExchanges[index] = formValue.value.sshKeyExchanges[newIndex];
    formValue.value.sshKeyExchanges[newIndex] = temp;
  } else if (type === 'macs' && formValue.value.sshMACs) {
    if (newIndex < 0 || newIndex >= formValue.value.sshMACs.length) return;
    const temp = formValue.value.sshMACs[index];
    formValue.value.sshMACs[index] = formValue.value.sshMACs[newIndex];
    formValue.value.sshMACs[newIndex] = temp;
  } else if (type === 'publicKeyAlgorithms' && formValue.value.sshPublicKeyAlgorithms) {
    if (newIndex < 0 || newIndex >= formValue.value.sshPublicKeyAlgorithms.length) return;
    const temp = formValue.value.sshPublicKeyAlgorithms[index];
    formValue.value.sshPublicKeyAlgorithms[index] = formValue.value.sshPublicKeyAlgorithms[newIndex];
    formValue.value.sshPublicKeyAlgorithms[newIndex] = temp;
  } else if (type === 'hostKeyAlgorithms' && formValue.value.sshHostKeyAlgorithms) {
    if (newIndex < 0 || newIndex >= formValue.value.sshHostKeyAlgorithms.length) return;
    const temp = formValue.value.sshHostKeyAlgorithms[index];
    formValue.value.sshHostKeyAlgorithms[index] = formValue.value.sshHostKeyAlgorithms[newIndex];
    formValue.value.sshHostKeyAlgorithms[newIndex] = temp;
  }
};

const handleResetAlgorithms = (type: AlgorithmType) => {
  if (type === 'ciphers') {
    formValue.value.sshCiphers = [...defaultCiphers];
  } else if (type === 'keyExchanges') {
    formValue.value.sshKeyExchanges = [...defaultKeyExchanges];
  } else if (type === 'macs') {
    formValue.value.sshMACs = [...defaultMACs];
  } else if (type === 'publicKeyAlgorithms') {
    formValue.value.sshPublicKeyAlgorithms = [...defaultPublicKeyAlgorithms];
  } else if (type === 'hostKeyAlgorithms') {
    formValue.value.sshHostKeyAlgorithms = [...defaultHostKeyAlgorithms];
  }
};
</script>

<style lang="less" scoped>
.form-row {
  display: flex;
  align-items: center;
  width: 100%;
  gap: 8px;

  .form-item {
    flex: 1;
  }
}

.auth-type-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.port-input {
  width: 120px;
}

.tooltip-text {
  font-size: 12px;
  line-height: 1.5;
}

.theme-preview {
  width: 100%;
  height: 100%;
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid var(--border-color);
  margin-bottom: 8px;
}

.theme-section {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.algorithm-container {
  display: flex;
  flex-direction: column;
  width: 100%;

  .reset-text {
    text-align: right;
    font-size: 12px;
    color: var(--primary-color);
    cursor: pointer;
    margin-bottom: 4px;

    &:hover {
      text-decoration: underline;
    }
  }

  .algorithm-list {
    max-height: 160px;
    overflow-y: auto;

    :deep(.n-list-item) {
      padding: 2px 10px;
    }

    .compact-item {
      padding: 0;
    }
  }

  .algorithm-item {
    display: flex;
    align-items: center;
    width: 100%;

    .algorithm-name {
      flex: 1;
      font-size: 12px;
      padding-left: 8px;
    }

    .algorithm-actions {
      display: flex;
      gap: 2px;
    }
  }
}

:deep(.n-collapse) {
  .n-collapse-item {
    margin: 0px;

    .n-collapse-item__header {
      padding: 0px;
      height: 24px;
    }

    .n-collapse-item__content-inner {
      padding: 0px;
    }
  }
}
</style>
