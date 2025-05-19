<template>
  <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md flex flex-col justify-between">
    <div>
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ book.title }}</h3>
      <p class="text-gray-600 dark:text-gray-300">by {{ book.author }}</p>
      <p class="text-gray-500 dark:text-gray-400 capitalize">{{ book.status }}</p>
    </div>
    <div class="mt-4 flex justify-end space-x-2">
      <Dialog>
        <DialogTrigger as-child>
          <Button variant="outline">Edit</Button>
        </DialogTrigger>
        <DialogContent class="sm:max-w-md z-[60]">
          <DialogHeader>
            <DialogTitle>Edit Book</DialogTitle>
          </DialogHeader>
          <div class="space-y-4">
            <div>
              <Label for="title">Title</Label>
              <Input id="title" v-model="localBook.title" />
            </div>
            <div>
              <Label for="author">Author</Label>
              <Input id="author" v-model="localBook.author" />
            </div>
            <div>
              <Label for="status">Status</Label>
              <Select v-model="localBook.status">
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
            <Button @click="$emit('edit', localBook)">Update</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      <Button
        variant="destructive"
        class="bg-red-600 hover:bg-red-700 text-white font-bold"
        @click="$emit('delete', book.id)"
      >
        Delete
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { Button } from '../components/ui/button';
import { Dialog, DialogTrigger, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogClose } from '../components/ui/dialog';
import { Input } from '../components/ui/input';
import { Label } from '../components/ui/label';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '../components/ui/select';

const props = defineProps<{
  book: { id: string; title: string; author: string; status: string };
}>();
defineEmits(['edit', 'delete']);

const localBook = ref({ ...props.book });
</script>