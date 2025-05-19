export interface Book {
  id: number
  title: string
  author: string
  status: 'unread' | 'reading' | 'complete'
}

export interface BookFormData {
  title: string
  author: string
  status: 'unread' | 'reading' | 'complete'
}
