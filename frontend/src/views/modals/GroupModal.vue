<template>
  <NModal
    v-model:show="visible"
    close-on-esc
    :negative-text="$t('frontend.groupModal.cancel')"
    :on-close="resetForm"
    :positive-text="$t('frontend.groupModal.confirm')"
    :show-icon="false"
    :title="isEdit ? $t('frontend.groupModal.editTitle') : $t('frontend.groupModal.title')"
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
    @positive-click="handleConfirm"
  >
    <NForm ref="formRef" :model="formValue" :rules="rules">
      <NFormItem path="name" :label="$t('frontend.groupModal.name')">
        <NInput
          v-model:value="formValue.name"
          clearable
          :placeholder="$t('frontend.groupModal.placeholder.name')"
          :allow-input="value => !/\s/.test(value)"
        />
      </NFormItem>
    </NForm>
  </NModal>
</template>

<script lang="ts" setup>
import type { Group } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/dal/model';
import { CreateGroup, UpdateGroup } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/groupsrv';
import type { FormInst, FormRules } from 'naive-ui';
import { NForm, NFormItem, NInput, NModal } from 'naive-ui';
import { computed, onMounted, onUpdated, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useCall } from '@/utils/call';

const props = defineProps<{
  show: boolean;
  isEdit: boolean;
  group?: Group;
}>();
const emit = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'success'): void;
}>();
const { t } = useI18n();
const formRef = ref<FormInst | null>(null);
const { call } = useCall();

const visible = computed({
  get: () => props.show,
  set: value => emit('update:show', value),
});

const defaultGroup: Partial<Group> = {
  name: '',
};

function createGroupObject(): Group {
  return { ...defaultGroup } as Group;
}

const formValue = ref<Group>(createGroupObject());

const rules: FormRules = {
  name: {
    required: true,
    message: t('frontend.groupModal.validation.nameRequired'),
    trigger: 'blur',
  },
};

const initModalData = () => {
  if (props.group) {
    formValue.value = { ...props.group } as Group;
  } else {
    formValue.value = createGroupObject();
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

    const backendFunc = props.isEdit ? UpdateGroup : CreateGroup;
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
  formValue.value = createGroupObject();
  emit('update:show', false);
};
</script>

<style lang="less" scoped>
:deep(.n-form-item .n-form-item-label) {
  font-size: 13px;
}

:deep(.n-input) {
  font-size: 13px;
}
</style>
