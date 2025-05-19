<template>
  <div>
    <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">Statistics</h2>
    <div v-if="loading" class="text-center">
      <span class="text-gray-500">Loading...</span>
    </div>
    <div v-else-if="error" class="text-center text-red-500">
      {{ error }}
    </div>
    <div v-else class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Total Read</h3>
        <p class="text-2xl text-blue-500 dark:text-blue-400">{{ stats.total_read }}</p>
      </div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Reading Progress</h3>
        <p class="text-2xl text-blue-500 dark:text-blue-400">{{ stats.reading_progress }}%</p>
      </div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Popular Author</h3>
        <p class="text-2xl text-blue-500 dark:text-blue-400">{{ stats.popular_author }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useStatsStore } from '../stores/stats.ts';

const statsStore = useStatsStore();
const { stats, loading, error } = storeToRefs(statsStore);

onMounted(() => {
  console.log('Stats page mounted, fetching stats'); // Debug
  statsStore.fetchStats();
  console.log('Stats after fetch:', stats.value); // Debug
});
</script>