import { ref } from 'vue'
import { getReviews, addReview, getReview } from '@/api/review'
import type { Review, ReviewListParams } from '@/api/review'

export function useReview() {
  const reviews = ref<Review[]>([])
  const loading = ref(false)

  const fetchReviews = async (params?: any) => {
    loading.value = true
    try {
      const res = await getReviews(params)
      if (res && res.data) {
        reviews.value = res.data.reviews || []
      }
      // 直接返回适当的格式，兼容ReviewListPage的处理
      return {
        data: reviews.value,
        total: res?.data?.total || 0
      }
    } catch (error) {
      console.error('获取评测列表失败', error)
      return {
        data: [],
        total: 0
      }
    } finally {
      loading.value = false
    }
  }

  // 创建评测
  const createReview = async (data: Omit<Review, 'id'>) => {
    loading.value = true
    try {
      const response = await addReview(data)
      return response.data.data
    } catch (error) {
      console.error('创建评测失败', error)
      return null
    } finally {
      loading.value = false
    }
  }

  // 获取单个评测详情
  const getReviewDetail = async (id: string) => {
    loading.value = true
    try {
      const response = await getReview(id)
      return response.data
    } catch (error) {
      console.error('获取评测详情失败', error)
      return null
    } finally {
      loading.value = false
    }
  }

  // 更新评测
  const updateReview = async (id: string, data: Partial<Review>) => {
    loading.value = true
    try {
      // 这里需要实现updateReview API
      // const response = await updateReview(id, data)
      // return response.data
      console.warn('updateReview API 未实现')
      return null
    } catch (error) {
      console.error('更新评测失败', error)
      return null
    } finally {
      loading.value = false
    }
  }

  return {
    reviews,
    loading,
    fetchReviews,
    createReview,
    updateReview,
    getReview: getReviewDetail
  }
}