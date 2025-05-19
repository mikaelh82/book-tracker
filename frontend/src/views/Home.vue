<template>
  <div>
    <div class="mb-6 flex flex-col sm:flex-row gap-4">
      <div class="flex-1">
        <Input
          v-model="searchQuery"
          placeholder="Search by title..."
          class="w-full"
        />
      </div>
      <div class="flex gap-4">
        <Select v-model="filters.status">
          <SelectTrigger>
            <SelectValue placeholder="Filter by status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All</SelectItem>
            <SelectItem value="unread">Unread</SelectItem>
            <SelectItem value="reading">Reading</SelectItem>
            <SelectItem value="complete">Complete</SelectItem>
          </SelectContent>
        </Select>
        <Select v-model="filters.sortBy">
          <SelectTrigger>
            <SelectValue placeholder="Sort by" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="title">Title (A-Z)</SelectItem>
            <SelectItem value="author">Author (A-Z)</SelectItem>
          </SelectContent>
        </Select>
        <Dialog>
          <DialogTrigger as-child>
            <Button>Add Book</Button>
          </DialogTrigger>
          <DialogContent class="sm:max-w-md z-[60]">
            <DialogHeader>
              <DialogTitle>Add Book</DialogTitle>
            </DialogHeader>
            <div class="space-y-4">
              <div>
                <Label for="title">Title</Label>
                <Input id="title" v-model="form.title" />
              </div>
              <div>
                <Label for="author">Author</Label>
                <Input id="author" v-model="form.author" />
              </div>
              <div>
                <Label for="status">Status</Label>
                <Select v-model="form.status">
                  <SelectTrigger>
                    <SelectValue placeholder="Select status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="unread">Unread</SelectItem>
                    <SelectItem value="reading">Reading</SelectItem>
                    <SelectItem value="complete">Complete</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <DialogFooter class="sm:justify-start">
              <DialogClose as-child>
                <Button variant="outline">Cancel</Button>
              </DialogClose>
              <Button @click="saveBook(false)">Add</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>

    <div v-if="loading" class="text-center">
      <span class="text-gray-500">Loading...</span>
    </div>
    <div v-else-if="error" class="text-center text-red-500">
      {{ error }}
    </div>
    <div v-else-if="filteredBooks.length === 0" class="text-center">
      <span class="text-gray-500">No books found.</span>
    </div>
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <BookCard
        v-for="book in filteredBooks"
        :key="book.id"
        :book="book"
        @edit="saveBook(true, $event)"
        @delete="deleteBook(book.id)"
      />
    </div>

    <div v-if="hasMore" class="mt-4 text-center">
      <Button @click="loadMore">Load More</Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useBookStore } from '../stores/books.ts';
import BookCard from '../components/BookCard.vue';
import { Input } from '../components/ui/input';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '../components/ui/select';
import { Button } from '../components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogTrigger, DialogClose } from '../components/ui/dialog';
import { Label } from '../components/ui/label';
import { debounce } from 'lodash';

interface Book {
  id: string;
  title: string;
  author: string;
  status: string;
}

const bookStore = useBookStore();
const { books, loading, error, hasMore } = storeToRefs(bookStore);
const searchQuery = ref('');
const filters = ref({ status: 'all', sortBy: '' });
const form = ref({ id: '', title: '', author: '', status: '' });
const page = ref(1);

const filteredBooks = computed(() => {
  const bookList = books.value ?? [];
  console.log('Books in computed:', bookList);
  let result = bookList;
  if (searchQuery.value) {
    result = result.filter(book =>
      book.title.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
  }
  if (filters.value.status && filters.value.status !== 'all') {
    result = result.filter(book => book.status === filters.value.status);
  }
  if (filters.value.sortBy) {
    result = [...result].sort((a, b) =>
      a[filters.value.sortBy].localeCompare(b[filters.value.sortBy])
    );
  }
  return result;
});

const fetchBooks = () => {
  const status = filters.value.status === 'all' ? '' : filters.value.status;
  console.log('Fetching with params:', { status, limit: 6, offset: (page.value - 1) * 6 });
  bookStore.fetchBooks({ status, limit: 6, offset: (page.value - 1) * 6 });
};

const debouncedFetch = debounce(fetchBooks, 300);

watch([searchQuery, filters], () => {
  page.value = 1;
  fetchBooks();
}, { deep: true });

const loadMore = () => {
  page.value++;
  fetchBooks();
};

const saveBook = async (isEditing: boolean, book?: Book) => {
  console.log('Saving book:', form.value, 'isEditing:', isEditing); // Debug
  if (book) {
    form.value = { ...book };
  }
  if (!form.value.title || !form.value.author || !form.value.status) {
    alert('Please fill all fields');
    return;
  }
  try {
    if (isEditing) {
      await bookStore.updateBook(form.value);
    } else {
      await bookStore.createBook(form.value);
    }
    form.value = { id: '', title: '', author: '', status: '' };
    page.value = 1;
    fetchBooks();
  } catch (err: any) {
    alert(err.message);
  }
};

const deleteBook = async (id: string) => {
  if (confirm('Are you sure you want to delete this book?')) {
    try {
      await bookStore.deleteBook(id);
      page.value = 1;
      fetchBooks();
    } catch (err: any) {
      alert(err.message);
    }
  }
};

onMounted(() => {
  console.log('Home mounted, triggering fetch');
  fetchBooks();
});
</script>