<template>
  <n-modal
    v-model:show="visible"
    close-on-esc
    :negative-text="$t('connModal.cancel')"
    :on-close="resetForm"
    :positive-text="$t('connModal.confirm')"
    :show-icon="false"
    :title="isEdit ? $t('connModal.editTitle') : $t('connModal.title')"
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
    @positive-click="handleConfirm"
  >
    <n-tabs animated type="line" placement="left" v-model:value="activeTab">
      <n-tab-pane name="basic" :tab="$t('connModal.basicConfig')">
        <n-scrollbar style="max-height: 70vh; padding-right: 12px">
          <n-form ref="formRef" :model="formValue" :rules="rules">
            <n-form-item path="label" :label="$t('connModal.label')">
              <n-input
                v-model:value="formValue.label"
                clearable
                :placeholder="$t('connModal.placeholder.label')"
                :allow-input="value => !/\s/.test(value)"
              />
            </n-form-item>

            <n-form-item path="groupID" :label="$t('connModal.group')">
              <n-select
                v-model:value="formValue.groupID"
                :options="groupOptions"
                :placeholder="$t('connModal.placeholder.group')"
              />
            </n-form-item>

            <n-form-item path="connProtocol" :label="$t('connModal.connProtocol')">
              <n-select
                v-model:value="formValue.connProtocol"
                :options="connProtocolOptions"
                :placeholder="$t('connModal.placeholder.connProtocol')"
              />
            </n-form-item>

            <template v-if="formValue.connProtocol === ConnProtocol.SERIAL">
              <n-form-item path="serialPort" :label="$t('connModal.serialPort')">
                <n-select
                  v-model:value="formValue.serialPort"
                  :options="serialPortsOptions"
                  :placeholder="$t('connModal.placeholder.serialPort')"
                />
              </n-form-item>

              <n-form-item path="baudRate" :label="$t('connModal.baudRate')">
                <n-select
                  v-model:value="formValue.baudRate"
                  :options="baudRateOptions"
                  :placeholder="$t('connModal.placeholder.baudRate')"
                />
              </n-form-item>

              <div class="form-row">
                <n-form-item path="dataBits" :label="$t('connModal.dataBits')" class="form-item">
                  <n-select
                    v-model:value="formValue.dataBits"
                    :options="dataBitsOptions"
                    :placeholder="$t('connModal.placeholder.dataBits')"
                  />
                </n-form-item>

                <n-form-item path="stopBits" :label="$t('connModal.stopBits')" class="form-item">
                  <n-select
                    v-model:value="formValue.stopBits"
                    :options="stopBitsOptions"
                    :placeholder="$t('connModal.placeholder.stopBits')"
                  />
                </n-form-item>

                <n-form-item path="parity" :label="$t('connModal.parity')" class="form-item">
                  <n-select
                    v-model:value="formValue.parity"
                    :options="parityOptions"
                    :placeholder="$t('connModal.placeholder.parity')"
                  />
                </n-form-item>
              </div>
            </template>

            <template v-if="formValue.connProtocol === ConnProtocol.SSH">
              <div class="form-row">
                <n-form-item path="host" :label="$t('connModal.host')" class="form-item">
                  <n-input
                    v-model:value="formValue.host"
                    clearable
                    :placeholder="$t('connModal.placeholder.host')"
                    :allow-input="value => !/\s/.test(value)"
                  />
                </n-form-item>
                <n-form-item path="port" :label="$t('connModal.port')" class="port-input">
                  <n-input-number
                    v-model:value="formValue.port"
                    :min="1"
                    :max="65535"
                    :show-button="false"
                    :placeholder="$t('connModal.placeholder.port')"
                  />
                </n-form-item>
              </div>

              <n-form-item :label="$t('connModal.authType')" path="credential.authMethod">
                <div class="auth-type-container">
                  <n-button-group>
                    <n-button
                      :type="
                        !formValue.useCommonCredential && formValue.credential?.authMethod === AuthMethod.PASSWORD
                          ? 'primary'
                          : 'default'
                      "
                      @click="handleCredentialTypeChange(CredentialType.Password)"
                    >
                      <template #icon>
                        <Icon icon="ph:password" />
                      </template>
                      {{ $t('connModal.password') }}
                    </n-button>
                    <n-button
                      :type="
                        !formValue.useCommonCredential && formValue.credential?.authMethod === AuthMethod.PRIVATEKEY
                          ? 'primary'
                          : 'default'
                      "
                      @click="handleCredentialTypeChange(CredentialType.PrivateKey)"
                    >
                      <template #icon>
                        <Icon icon="ph:key" />
                      </template>
                      {{ $t('connModal.privateKey') }}
                    </n-button>
                    <n-button
                      :type="formValue.useCommonCredential ? 'primary' : 'default'"
                      @click="handleCredentialTypeChange(CredentialType.Common)"
                    >
                      <template #icon>
                        <Icon icon="ph:vault" />
                      </template>
                      {{ $t('connModal.commonCredentialLib') }}
                    </n-button>
                  </n-button-group>
                </div>
              </n-form-item>

              <template v-if="formValue.useCommonCredential">
                <n-form-item path="credentialID" :label="$t('connModal.credential')">
                  <n-select
                    v-model:value="formValue.credentialID"
                    :options="credentialOptions"
                    clearable
                    :placeholder="$t('connModal.placeholder.selectCredential')"
                    @update:value="handleSelectCredential"
                  />
                </n-form-item>
              </template>

              <template v-else>
                <n-form-item path="credential.username" :label="$t('connModal.username')">
                  <n-input
                    v-model:value="formValue.credential!.username"
                    clearable
                    :placeholder="$t('connModal.placeholder.username')"
                    :allow-input="value => !/\s/.test(value)"
                  />
                </n-form-item>

                <template v-if="formValue.credential!.authMethod === AuthMethod.PASSWORD">
                  <n-form-item path="credential.password" :label="$t('connModal.password')">
                    <n-input
                      v-model:value="formValue.credential!.password"
                      type="password"
                      show-password-on="click"
                      clearable
                      :placeholder="$t('connModal.placeholder.password')"
                      :allow-input="value => !/\s/.test(value)"
                    />
                  </n-form-item>
                </template>

                <template v-if="formValue.credential!.authMethod === AuthMethod.PRIVATEKEY">
                  <n-form-item path="credential.privateKey" :label="$t('connModal.privateKey')">
                    <n-input
                      v-model:value="formValue.credential!.privateKey"
                      type="textarea"
                      :autosize="{ minRows: 3, maxRows: 3 }"
                      clearable
                      :placeholder="$t('connModal.placeholder.privateKey')"
                      :allow-input="value => !/\s/.test(value)"
                    />
                  </n-form-item>
                  <n-form-item path="credential.passphrase" :label="$t('connModal.passphrase')">
                    <n-input
                      v-model:value="formValue.credential!.passphrase"
                      type="password"
                      show-password-on="click"
                      clearable
                      :placeholder="$t('connModal.placeholder.passphrase')"
                      :allow-input="value => !/\s/.test(value)"
                    />
                  </n-form-item>
                </template>
              </template>
            </template>
          </n-form>
        </n-scrollbar>
      </n-tab-pane>

      <n-tab-pane name="personal" :tab="$t('connModal.personalConfig')">
        <n-form ref="formRef" :model="formValue" :rules="rules">
          <n-form-item path="theme" :label="$t('connModal.theme')">
            <div class="theme-section">
              <div class="theme-preview">
                <theme-preview :theme="selectedTheme" />
              </div>
              <n-select
                v-model:value="formValue.theme"
                :options="themeOptions"
                :placeholder="t('connModal.placeholder.selectTheme')"
                filterable
                @update:value="handleThemeChange"
              />
            </div>
          </n-form-item>
        </n-form>
      </n-tab-pane>

      <n-tab-pane name="advanced" :tab="$t('connModal.advancedConfig')">
        <n-empty size="small" :description="$t('connModal.developing')">
          <template #icon>
            <Icon icon="ph:code" />
          </template>
        </n-empty>
      </n-tab-pane>
    </n-tabs>
  </n-modal>
</template>

<script lang="ts" setup>
import { enums, model } from '@wailsApp/go/models';
import { ListGroup } from '@wailsApp/go/services/GroupSrv';
import { CreateConnection, UpdateConnection } from '@wailsApp/go/services/ConnectionSrv';
import { Icon } from '@iconify/vue';
import ThemePreview from '../components/ThemePreview.vue';
import { themeOptions as xtermThemeOptions, loadTheme, defaultTheme } from '@/themes/xtermjs';

import {
  FormInst,
  FormRules,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NButton,
  NButtonGroup,
  NSelect,
  NTabPane,
  NTabs,
  useMessage,
  NEmpty,
  NScrollbar,
} from 'naive-ui';
import { SelectMixedOption } from 'naive-ui/es/select/src/interface';
import { ref, computed, onMounted, onUpdated } from 'vue';
import { useI18n } from 'vue-i18n';
import { SerialPorts } from '@wailsApp/go/services/TerminalSrv';
import { ListCredential } from '@wailsApp/go/services/CredentialSrv';

const { AuthMethod, ConnProtocol } = enums;

const { t } = useI18n();
const formRef = ref<FormInst | null>(null);
const activeTab = ref('basic');
const message = useMessage();

const props = defineProps<{
  show: boolean;
  isEdit: boolean;
  connection?: model.Connection;
}>();

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'success'): void;
}>();

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
  { label: ConnProtocol.SERIAL, value: ConnProtocol.SERIAL },
];

const defaultCredential: Partial<model.Credential> = {
  label: '',
  username: '',
  password: '',
  privateKey: '',
  passphrase: '',
  isCommonCredential: false,
  authMethod: AuthMethod.PASSWORD,
};

const defaultConnection: Partial<model.Connection> = {
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
};

const tempCachedCredentials = ref<{
  password: model.Credential;
  privateKey: model.Credential;
}>({
  password: createCredentialObject(AuthMethod.PASSWORD),
  privateKey: createCredentialObject(AuthMethod.PRIVATEKEY),
});

const groupOptions = ref<SelectMixedOption[]>([]);
const serialPortsOptions = ref<SelectMixedOption[]>([]);
const credentialOptions = ref<SelectMixedOption[]>([]);

const formValue = ref<model.Connection>(createConnectionObject());

function createCredentialObject(authMethod: enums.AuthMethod): model.Credential {
  return {
    ...defaultCredential,
    authMethod,
  } as model.Credential;
}

function createConnectionObject(isCommon = false): model.Connection {
  const conn = {
    ...defaultConnection,
    useCommonCredential: isCommon,
  } as model.Connection;

  if (!isCommon) {
    conn.credential = createCredentialObject(AuthMethod.PASSWORD);
  }
  return conn;
}

function prepareConnectionForEdit(connection: model.Connection): model.Connection {
  const conn = { ...connection } as model.Connection;

  if (conn.credential && !conn.useCommonCredential) {
    if (conn.credential.authMethod === AuthMethod.PASSWORD) {
      tempCachedCredentials.value.password = { ...conn.credential } as model.Credential;
    } else if (conn.credential.authMethod === AuthMethod.PRIVATEKEY) {
      tempCachedCredentials.value.privateKey = { ...conn.credential } as model.Credential;
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
    message: t('connModal.validation.labelRequired'),
    trigger: 'blur',
  },
  connProtocol: {
    required: true,
    message: t('connModal.validation.connProtocolRequired'),
    trigger: 'blur',
  },
  host: {
    required: formValue.value.connProtocol === ConnProtocol.SSH,
    message: t('connModal.validation.hostRequired'),
    trigger: 'blur',
  },
  port: {
    required: formValue.value.connProtocol === ConnProtocol.SSH,
    type: 'number',
    message: t('connModal.validation.portRequired'),
    trigger: ['blur', 'change'],
    validator: (rule, value) => {
      if (typeof value !== 'number' || value < 1 || value > 65535) {
        return new Error(t('connModal.validation.portInvalid'));
      }
    },
  },
  serialPort: {
    required: formValue.value.connProtocol === ConnProtocol.SERIAL,
    message: t('connModal.validation.serialPortRequired'),
    trigger: ['blur', 'change'],
  },
  baudRate: {
    required: formValue.value.connProtocol === ConnProtocol.SERIAL,
    type: 'number',
    message: t('connModal.validation.baudRateRequired'),
    trigger: ['blur', 'change'],
  },
  dataBits: {
    required: formValue.value.connProtocol === ConnProtocol.SERIAL,
    type: 'number',
    message: t('connModal.validation.dataBitsRequired'),
    trigger: ['blur', 'change'],
  },
  stopBits: {
    required: formValue.value.connProtocol === ConnProtocol.SERIAL,
    type: 'number',
    message: t('connModal.validation.stopBitsRequired'),
    trigger: ['blur', 'change'],
  },
  parity: {
    required: formValue.value.connProtocol === ConnProtocol.SERIAL,
    type: 'number',
    message: t('connModal.validation.parityRequired'),
    trigger: ['blur', 'change'],
  },
  'credential.username': {
    required: !formValue.value.useCommonCredential,
    message: t('connModal.validation.usernameRequired'),
    trigger: 'blur',
  },
  'credential.password': {
    required: !formValue.value.useCommonCredential && formValue.value.credential?.authMethod === AuthMethod.PASSWORD,
    message: t('connModal.validation.passwordRequired'),
    trigger: 'blur',
  },
  'credential.privateKey': {
    required: !formValue.value.useCommonCredential && formValue.value.credential?.authMethod === AuthMethod.PRIVATEKEY,
    message: t('connModal.validation.privateKeyRequired'),
    trigger: 'blur',
  },
  credentialID: {
    required: formValue.value.useCommonCredential,
    message: t('connModal.validation.credentialRequired'),
    trigger: ['blur', 'change'],
    validator: (rule, value) => {
      if (formValue.value.useCommonCredential && !value) {
        return new Error(t('connModal.validation.credentialRequired'));
      }
      return true;
    },
  },
}));

const handleCredentialTypeChange = (credentialType: number) => {
  if (credentialType === CredentialType.Common) {
    formValue.value.useCommonCredential = true;
    formValue.value.credentialID = undefined;
    formValue.value.credential = undefined;
    return;
  }
  formValue.value.useCommonCredential = false;
  formValue.value.credentialID = undefined;
  formValue.value.credential =
    credentialType === CredentialType.Password
      ? tempCachedCredentials.value.password
      : tempCachedCredentials.value.privateKey;
  formValue.value.credential.authMethod =
    credentialType === CredentialType.Password ? AuthMethod.PASSWORD : AuthMethod.PRIVATEKEY;
};

const handleSelectCredential = async (id: number) => {
  if (!id) {
    formValue.value.credentialID = undefined;
    return;
  }
  formValue.value.credentialID = id;
  formValue.value.useCommonCredential = true;
};

const fetchGroups = async () => {
  const resp = await ListGroup();
  if (!resp.ok) {
    message.error(resp.msg);
    return [];
  }
  return resp.data;
};

const fetchSerialPorts = async () => {
  const resp = await SerialPorts();
  if (!resp.ok) {
    message.error(resp.msg);
    return [];
  }
  return resp.data;
};

const fetchCredentials = async () => {
  const resp = await ListCredential();
  if (!resp.ok) {
    message.error(resp.msg);
    return [];
  }
  return resp.data || [];
};

const initOptions = async () => {
  const [groups, credentials, serialPorts] = await Promise.all([fetchGroups(), fetchCredentials(), fetchSerialPorts()]);
  groupOptions.value = groups.map((group: model.Group) => ({
    label: group.name,
    value: group.id,
  }));
  credentialOptions.value = credentials.map((credential: model.Credential) => ({
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

  if (props.connection) {
    formValue.value = prepareConnectionForEdit(props.connection);
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
    const resp = props.isEdit ? await UpdateConnection(formValue.value) : await CreateConnection(formValue.value);

    if (!resp.ok) {
      message.error(resp.msg);
      return false;
    }

    message.success(props.isEdit ? t('message.updateSuccess') : t('message.createSuccess'));
    emit('update:show', false);
    emit('success');
  } catch (errors) {
    console.error('表单验证错误:', errors);
    return false;
  }
};

const resetForm = () => {
  formValue.value = createConnectionObject();
  tempCachedCredentials.value = {
    password: createCredentialObject(AuthMethod.PASSWORD),
    privateKey: createCredentialObject(AuthMethod.PRIVATEKEY),
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
</style>
