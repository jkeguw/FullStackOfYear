import { ref } from 'vue'
import { getDevices, addDevice } from '@/api/device'
import type { Device } from '@/api/device'

export function useDevice() {
  const devices = ref<Device[]>([])
  const loading = ref(false)

  const fetchDevices = async () => {
    loading.value = true
    try {
      const res = await getDevices()
      devices.value = res.data
    } finally {
      loading.value = false
    }
  }

  return {
    devices,
    loading,
    fetchDevices
  }
}