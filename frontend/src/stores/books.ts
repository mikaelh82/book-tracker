import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'

interface Book {
  id: string
  title: string
  author: string
  status: string
}

export const useBookStore = defineStore('book', () => {
  const books = ref<Book[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const hasMore = ref(true)

  const fetchBooks = async ({ status = '', limit = 10, offset = 0 }) => {
    loading.value = true
    error.value = null
    try {
      const response = await axios.get('http://localhost:8080/api/v1/books', {
        params: { status, limit, offset },
      })
      console.log('Raw API response:', response.data) // Debug raw response
      const newBooks: Book[] = Array.isArray(response.data) ? response.data : []
      console.log('Parsed books:', newBooks) // Debug parsed books
      books.value = offset === 0 ? newBooks : [...books.value, ...newBooks]
      console.log('Updated books.value:', books.value) // Debug final state
      console.log('Store state:', {
        books: books.value,
        loading: loading.value,
        error: error.value,
        hasMore: hasMore.value,
      }) // Debug store
      hasMore.value = newBooks.length === limit
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch books'
      hasMore.value = false
      console.error('Fetch error:', err) // Debug
    } finally {
      loading.value = false
    }
  }

  const createBook = async (book: Book) => {
    try {
      const response = await axios.post('http://localhost:8080/api/v1/books', book)
      books.value = [...books.value, response.data]
    } catch (err: any) {
      throw new Error(err.response?.data?.error || 'Failed to create book')
    }
  }

  const updateBook = async (book: Book) => {
    try {
      await axios.put(`http://localhost:8080/api/v1/books/${book.id}`, book)
      books.value = books.value.map((b) => (b.id === book.id ? book : b))
    } catch (err: any) {
      throw new Error(err.response?.data?.error || 'Failed to update book')
    }
  }

  const deleteBook = async (id: string) => {
    try {
      await axios.delete(`http://localhost:8080/api/v1/books/${id}`)
      books.value = books.value.filter((b) => b.id !== id)
    } catch (err: any) {
      throw new Error(err.response?.data?.error || 'Failed to delete book')
    }
  }

  const clearBooks = () => {
    books.value = []
    hasMore.value = true
  }

  return {
    books,
    loading,
    error,
    hasMore,
    fetchBooks,
    createBook,
    updateBook,
    deleteBook,
    clearBooks,
  }
})
