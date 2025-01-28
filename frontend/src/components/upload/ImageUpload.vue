<template>
  <el-upload
    v-model:file-list="fileList"
    action=""
    :http-request="uploadFile"
    :before-upload="beforeUpload"
    multiple
  >
    <el-button type="primary">Upload Image</el-button>
  </el-upload>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, UploadFile } from 'element-plus'
import { uploadFile } from '@/utils/upload'

const fileList = ref<UploadFile[]>([])

const beforeUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
  }
  return isImage
}
</script>