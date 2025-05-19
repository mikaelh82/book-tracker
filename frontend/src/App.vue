<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900 transition-colors">
    <Navbar />
    <main class="container mx-auto p-4">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import Navbar from './components/Navbar.vue';
</script>