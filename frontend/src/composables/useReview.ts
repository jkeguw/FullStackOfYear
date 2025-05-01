import { ref } from 'vue';
import { getReviews, addReview, getReview } from '@/api/review';
import type { Review, ReviewListParams } from '@/api/review';

export function useReview() {
  const reviews = ref<Review[]>([]);
  const loading = ref(false);

  const fetchReviews = async (params?: ReviewListParams) => {
    loading.value = true;
    try {
      const res = await getReviews(params);
      
      // Verify response structure matches expected format from API
      if (res && res.code === 0 && res.data) {
        // Check if the response has the expected structure from design doc
        const reviewData = res.data;
        
        if (reviewData.reviews && Array.isArray(reviewData.reviews)) {
          reviews.value = reviewData.reviews;
        } else {
          console.warn('API response missing reviews array or has invalid format');
          reviews.value = [];
        }
        
        // Return in format compatible with ReviewListPage
        return {
          data: reviews.value,
          total: reviewData.total || 0
        };
      } else {
        console.warn('Invalid API response format', res);
        return {
          data: [],
          total: 0
        };
      }
    } catch (error) {
      console.error('Failed to fetch review list', error);
      return {
        data: [],
        total: 0
      };
    } finally {
      loading.value = false;
    }
  };

  // Create review
  const createReview = async (data: Omit<Review, 'id'>) => {
    loading.value = true;
    try {
      const response = await addReview(data);
      if (response.data && response.data.code === 0 && response.data.data) {
        return response.data.data;
      } else {
        console.warn('Invalid API response format when creating review', response);
        return null;
      }
    } catch (error) {
      console.error('Failed to create review', error);
      return null;
    } finally {
      loading.value = false;
    }
  };

  // Get single review details
  const getReviewDetail = async (id: string) => {
    loading.value = true;
    try {
      const response = await getReview(id);
      if (response && response.code === 0 && response.data) {
        return response.data;
      } else {
        console.warn('Invalid API response format when getting review details', response);
        return null;
      }
    } catch (error) {
      console.error('Failed to get review details', error);
      return null;
    } finally {
      loading.value = false;
    }
  };

  // Update review
  const updateReview = async (id: string, data: Partial<Review>) => {
    loading.value = true;
    try {
      // Implementation of updateReview API needed
      // const response = await updateReview(id, data)
      // return response.data
      console.warn('updateReview API not implemented yet');
      return null;
    } catch (error) {
      console.error('Failed to update review', error);
      return null;
    } finally {
      loading.value = false;
    }
  };

  return {
    reviews,
    loading,
    fetchReviews,
    createReview,
    updateReview,
    getReview: getReviewDetail
  };
}
