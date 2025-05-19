import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'

interface Stats {
  total_read: number
  reading_progress: number
  popular_author: string
}

export const useStatsStore = defineStore('stats', () => {
  const stats = ref<Stats>({ total_read: 0, reading_progress: 0, popular_author: 'N/A' })
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchStats = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await axios.get('http://localhost:8080/api/v1/stats')
      console.log('Fetched stats:', response.data) // Debug
      stats.value = response.data
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch stats'
      console.error('Stats fetch error:', err) // Debug
      stats.value = { total_read: 0, reading_progress: 0, popular_author: 'N/A' }
    } finally {
      loading.value = false
    }
  }

  return { stats, loading, error, fetchStats }
})
