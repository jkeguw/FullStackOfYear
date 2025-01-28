import { ref } from 'vue'
import { getReviews, addReview } from '@/api/review'
import type { Review } from '@/api/review'

export function useReview() {
  const reviews = ref<Review[]>([])
  const loading = ref(false)

  const fetchReviews = async (deviceId?: string) => {
    loading.value = true
    try {
      const res = await getReviews({ deviceId })
      reviews.value = res.data
    } finally {
      loading.value = false
    }
  }

  return {
    reviews,
    loading,
    fetchReviews
  }
}